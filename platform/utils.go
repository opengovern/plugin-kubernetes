// Description:
// This program inspects a Kubernetes Kubeconfig file. It can perform two main actions:
// 1. Get Cluster Info: Validates the Kubeconfig, connects to the cluster specified
//    by the current context, and returns details like endpoint, auth method,
//    TLS status, and server version.
// 2. Health Check: Validates the Kubeconfig and checks if the credentials have
//    sufficient read-only permissions (via 'list' verb checks) for a predefined
//    set of common Kubernetes resources across the cluster using SelfSubjectAccessReview.
//
// The program always outputs results in JSON format to standard output.
// Logging Behaviour:
// - By default, no logs are sent to standard error.
// - Set the LOG_LEVEL environment variable to enable logging:
//   - LOG_LEVEL=error: Logs only critical errors.
//   - LOG_LEVEL=warn: Logs warnings and errors.
//   - LOG_LEVEL=info: Logs info, warnings, and errors (shows major steps).
//   - LOG_LEVEL=debug: Logs all levels (verbose, for troubleshooting).
// Logs are written to standard error in JSON format.
//
// Exit codes indicate success (0), expected failure (1, e.g., validation error,
// permission denied, connection timeout), or internal/setup error (2+).
//
// Dependencies:
// - go.uber.org/zap
// - golang.org/x/xerrors (or use standard library error wrapping)
// - k8s.io/api
// - k8s.io/apimachinery
// - k8s.io/client-go
//
// How to Build:
// 1. Ensure Go is installed (>= 1.18 recommended).
// 2. Save this code as `main.go` (or your preferred name) in a directory.
// 3. Open a terminal in that directory.
// 4. Initialize Go module (if not already): go mod init <your_module_name> (e.g., go mod init kubechecker)
// 5. Add dependencies:
//    go get go.uber.org/zap
//    go get golang.org/x/xerrors
//    go get k8s.io/api/authorization/v1
//    go get k8s.io/apimachinery/pkg/apis/meta/v1
//    go get k8s.io/client-go/kubernetes
//    go get k8s.io/client-go/rest
//    go get k8s.io/client-go/tools/clientcmd
//    go get k8s.io/client-go/tools/clientcmd/api
// 6. Tidy dependencies: go mod tidy
// 7. Build the executable: go build -o kubechecker .
//
// How to Run:
//   # Get Cluster Information (no logging by default):
//   ./kubechecker --kubeconfig /path/to/your/kubeconfig.yaml
//
//   # Get Cluster Information (with info logging):
//   LOG_LEVEL=info ./kubechecker --kubeconfig /path/to/your/kubeconfig.yaml
//
//   # Run Read Permission Health Check (with debug logging):
//   LOG_LEVEL=debug ./kubechecker --kubeconfig /path/to/your/kubeconfig.yaml --health-check
//
//   # See all flags:
//   ./kubechecker --help

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	// Zap logger
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	// Kubernetes client-go libraries
	authorizationv1 "k8s.io/api/authorization/v1" // For SelfSubjectAccessReview
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest" // Added for rest.Config type
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"

	// For better error wrapping (can be replaced with standard library %w)
	"golang.org/x/xerrors"
)

// --- Configuration Constants ---
const (
	// MaxRetries defines how many times to retry fetching server version on transient failures.
	MaxRetries = 2
	// RetryDelay defines the duration to wait between retries.
	RetryDelay = 1 * time.Second
	// RequestTimeout defines the maximum time allowed for individual Kubernetes API requests.
	RequestTimeout = 8 * time.Second
	// EnvLogLevel is the environment variable used to control logging level.
	EnvLogLevel = "LOG_LEVEL"
)

// --- Output Structures (with snake_case JSON tags) ---

// ClusterInfo holds the extracted information about the Kubernetes cluster on success.
// NOTE: The 'error' field has been removed from the successful output.
type ClusterInfo struct {
	AuthMethod            string `json:"auth_method"`
	ContextName           string `json:"context_name"`
	Endpoint              string `json:"endpoint"`
	TLSServerVerification bool   `json:"tls_server_verification"`
	ServerVersion         string `json:"server_version,omitempty"` // Omit if not retrieved
}

// ErrorInfo holds information about a generic failure (validation, connection, internal).
type ErrorInfo struct {
	Error   bool   `json:"error"` // Always true for error JSON
	Message string `json:"message"`
	Details string `json:"details,omitempty"` // Technical details from underlying error
}

// HealthPermission represents a specific permission check failure during the health check.
type HealthPermission struct {
	Group    string `json:"group"`
	Version  string `json:"version"`
	Resource string `json:"resource"`
	Verb     string `json:"verb"`
	Scope    string `json:"scope"` // e.g., "cluster", "namespace"
	Reason   string `json:"reason,omitempty"`
}

// HealthStatus represents the overall outcome of the health check.
type HealthStatus struct {
	Healthy            bool               `json:"healthy"`
	Message            string             `json:"message"`
	MissingPermissions []HealthPermission `json:"missing_permissions,omitempty"` // Included only if not healthy
}

// --- Resource Definition for Health Check ---

// ResourceToCheck defines the GVR and scope for a health check permission.
type ResourceToCheck struct {
	Group    string
	Version  string
	Resource string
	Scope    string // Friendly scope name for reporting ("cluster" or "namespace")
	Friendly string // User-friendly name (e.g., "Pod", "Deployment")
}

// resourcesToVerify defines the list of resource types and their corresponding API groups/versions
// for which 'list' permission will be checked cluster-wide during the health check.
var resourcesToVerify = []ResourceToCheck{
	// Core API Group ("")
	{Group: "", Version: "v1", Resource: "configmaps", Scope: "namespace", Friendly: "ConfigMap"},
	{Group: "", Version: "v1", Resource: "endpoints", Scope: "namespace", Friendly: "Endpoints"},
	{Group: "", Version: "v1", Resource: "events", Scope: "namespace", Friendly: "Event"},
	{Group: "", Version: "v1", Resource: "limitranges", Scope: "namespace", Friendly: "LimitRange"},
	{Group: "", Version: "v1", Resource: "namespaces", Scope: "cluster", Friendly: "Namespace"},
	{Group: "", Version: "v1", Resource: "nodes", Scope: "cluster", Friendly: "Node"},
	{Group: "", Version: "v1", Resource: "persistentvolumeclaims", Scope: "namespace", Friendly: "PersistentVolumeClaim"},
	{Group: "", Version: "v1", Resource: "persistentvolumes", Scope: "cluster", Friendly: "PersistentVolume"},
	{Group: "", Version: "v1", Resource: "pods", Scope: "namespace", Friendly: "Pod"},
	{Group: "", Version: "v1", Resource: "podtemplates", Scope: "namespace", Friendly: "PodTemplate"},
	{Group: "", Version: "v1", Resource: "replicationcontrollers", Scope: "namespace", Friendly: "ReplicationController"},
	{Group: "", Version: "v1", Resource: "resourcequotas", Scope: "namespace", Friendly: "ResourceQuota"},
	{Group: "", Version: "v1", Resource: "secrets", Scope: "namespace", Friendly: "Secret"},
	{Group: "", Version: "v1", Resource: "serviceaccounts", Scope: "namespace", Friendly: "ServiceAccount"},
	{Group: "", Version: "v1", Resource: "services", Scope: "namespace", Friendly: "Service"},
	// Apps API Group ("apps")
	{Group: "apps", Version: "v1", Resource: "daemonsets", Scope: "namespace", Friendly: "DaemonSet"},
	{Group: "apps", Version: "v1", Resource: "deployments", Scope: "namespace", Friendly: "Deployment"},
	{Group: "apps", Version: "v1", Resource: "replicasets", Scope: "namespace", Friendly: "ReplicaSet"},
	{Group: "apps", Version: "v1", Resource: "statefulsets", Scope: "namespace", Friendly: "StatefulSet"},
	// Autoscaling API Group ("autoscaling")
	{Group: "autoscaling", Version: "v1", Resource: "horizontalpodautoscalers", Scope: "namespace", Friendly: "HorizontalPodAutoscaler"},
	// Batch API Group ("batch")
	{Group: "batch", Version: "v1", Resource: "cronjobs", Scope: "namespace", Friendly: "CronJob"},
	{Group: "batch", Version: "v1", Resource: "jobs", Scope: "namespace", Friendly: "Job"},
	// Discovery API Group ("discovery.k8s.io")
	{Group: "discovery.k8s.io", Version: "v1", Resource: "endpointslices", Scope: "namespace", Friendly: "EndpointSlice"},
	// Networking API Group ("networking.k8s.io")
	{Group: "networking.k8s.io", Version: "v1", Resource: "ingresses", Scope: "namespace", Friendly: "Ingress"},
	{Group: "networking.k8s.io", Version: "v1", Resource: "networkpolicies", Scope: "namespace", Friendly: "NetworkPolicy"},
	// Policy API Group ("policy")
	{Group: "policy", Version: "v1", Resource: "poddisruptionbudgets", Scope: "namespace", Friendly: "PodDisruptionBudget"},
	// RBAC API Group ("rbac.authorization.k8s.io")
	{Group: "rbac.authorization.k8s.io", Version: "v1", Resource: "clusterroles", Scope: "cluster", Friendly: "ClusterRole"},
	{Group: "rbac.authorization.k8s.io", Version: "v1", Resource: "clusterrolebindings", Scope: "cluster", Friendly: "ClusterRoleBinding"},
	{Group: "rbac.authorization.k8s.io", Version: "v1", Resource: "roles", Scope: "namespace", Friendly: "Role"},
	{Group: "rbac.authorization.k8s.io", Version: "v1", Resource: "rolebindings", Scope: "namespace", Friendly: "RoleBinding"},
	// Storage API Group ("storage.k8s.io")
	{Group: "storage.k8s.io", Version: "v1", Resource: "storageclasses", Scope: "cluster", Friendly: "StorageClass"},
	// API Extensions Group ("apiextensions.k8s.io")
	{Group: "apiextensions.k8s.io", Version: "v1", Resource: "customresourcedefinitions", Scope: "cluster", Friendly: "CustomResourceDefinition"},
}

// --- Helper Function for Consistent Error JSON ---

// createErrorJSON creates a standardized JSON string representation of an ErrorInfo struct.
// It logs an internal error if marshalling the ErrorInfo itself fails.
// Handles nil detailErr gracefully.
func createErrorJSON(message string, detailErr error, logger *zap.Logger) string {
	errInfo := ErrorInfo{
		Error:   true, // Indicate this JSON represents an error
		Message: message,
	}
	if detailErr != nil {
		errInfo.Details = fmt.Sprint(detailErr)
	}

	jsonBytes, err := json.Marshal(errInfo)
	if err != nil {
		if logger != nil {
			// Use a direct print to stderr in case logger setup failed or is disabled
			fmt.Fprintf(os.Stderr, `{"level":"error","timestamp":"%s","caller":"main/createErrorJSON","msg":"Internal alert: Failed to marshal ErrorInfo struct","original_message":"%s","error":"%v"}\n`,
				time.Now().Format(time.RFC3339), message, err)
		}
		return `{"error":true,"message":"Internal processing error while formatting error response"}`
	}
	return string(jsonBytes)
}

// --- Kubeconfig Validation Function ---

// ValidateKubeconfig reads, parses, and performs structural validation on a kubeconfig file.
// It returns the parsed Kubeconfig structure (api.Config), the REST config, and nil error on success.
// On failure, it logs the error (if logger is enabled) and returns nil configs along with the error.
func ValidateKubeconfig(ctx context.Context, kubeconfigPath string, logger *zap.Logger) (*api.Config, *rest.Config, error) {
	l := logger.With(zap.String("kubeconfig_path", kubeconfigPath))

	// 1. Read File
	kubeconfigBytes, err := os.ReadFile(kubeconfigPath)
	if err != nil {
		wrappedErr := xerrors.Errorf("failed to read kubeconfig file '%s': %w", kubeconfigPath, err)
		l.Error("Cannot read input file", zap.Error(wrappedErr)) // Logged at Error level
		return nil, nil, wrappedErr
	}

	// 2. Parse Structure (Lean Check)
	configAPI, err := clientcmd.Load(kubeconfigBytes)
	if err != nil {
		wrappedErr := xerrors.Errorf("failed to parse kubeconfig data from '%s': %w", kubeconfigPath, err)
		l.Error("Cannot parse kubeconfig data (is it valid YAML/JSON?)", zap.Error(wrappedErr)) // Logged at Error level
		return nil, nil, wrappedErr
	}

	// 3. Deeper Validation
	if configAPI.CurrentContext == "" {
		errMsg := "kubeconfig validation failed: current-context is not set"
		err := xerrors.New(errMsg)
		l.Error(errMsg) // Logged at Error level
		return nil, nil, err
	}
	l = l.With(zap.String("current_context", configAPI.CurrentContext))

	contextInfo, contextExists := configAPI.Contexts[configAPI.CurrentContext]
	if !contextExists {
		errMsg := fmt.Sprintf("current-context '%s' not found in contexts map", configAPI.CurrentContext)
		err := xerrors.New(errMsg)
		l.Error("Kubeconfig validation failed", zap.String("reason", errMsg))
		return nil, nil, err
	}

	if contextInfo.Cluster == "" {
		errMsg := fmt.Sprintf("cluster not defined for context '%s'", configAPI.CurrentContext)
		err := xerrors.New(errMsg)
		l.Error("Kubeconfig validation failed", zap.String("reason", errMsg))
		return nil, nil, err
	}
	l = l.With(zap.String("cluster_name", contextInfo.Cluster))

	_, clusterExists := configAPI.Clusters[contextInfo.Cluster]
	if !clusterExists {
		errMsg := fmt.Sprintf("cluster '%s' (referenced by context '%s') not found in clusters map", contextInfo.Cluster, configAPI.CurrentContext)
		err := xerrors.New(errMsg)
		l.Error("Kubeconfig validation failed", zap.String("reason", errMsg))
		return nil, nil, err
	}

	authInfoNameLog := "<anonymous>"
	if contextInfo.AuthInfo != "" {
		authInfoNameLog = contextInfo.AuthInfo
		_, authInfoExists := configAPI.AuthInfos[contextInfo.AuthInfo]
		if !authInfoExists {
			errMsg := fmt.Sprintf("authinfo (user) '%s' (referenced by context '%s') not found in users map", contextInfo.AuthInfo, configAPI.CurrentContext)
			err := xerrors.New(errMsg)
			l.Error("Kubeconfig validation failed", zap.String("reason", errMsg))
			return nil, nil, err
		}
	} else {
		l.Warn("AuthInfo (user) not defined for context. Cluster access might be anonymous.") // Logged at Warn level
	}
	l = l.With(zap.String("auth_info_name", authInfoNameLog))

	// 4. Build REST Config
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	loadingRules.ExplicitPath = kubeconfigPath
	configOverrides := &clientcmd.ConfigOverrides{}
	clientConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)

	restConfig, err := clientConfig.ClientConfig()
	if err != nil {
		wrappedErr := xerrors.Errorf("failed to build REST client config from kubeconfig: %w", err)
		l.Error("Cannot build client configuration", zap.Error(wrappedErr)) // Logged at Error level
		return nil, nil, wrappedErr
	}

	l.Info("Kubeconfig validation successful") // Logged at Info level
	return configAPI, restConfig, nil
}

// --- Cluster Information Retrieval Function ---

type KubeConfigInfo struct {
	AuthMethod            string `json:"auth_method"`
	ContextName           string `json:"context_name"`
	Endpoint              string `json:"endpoint"`
	TLSServerVerification bool   `json:"tls_server_verification"`
	ServerVersion         string `json:"server_version"`
}

// GetClusterInfo connects to the Kubernetes cluster using the validated configuration,
// retrieves metadata including the server version, and determines the authentication method used.
// It *always* returns a JSON string: ClusterInfo on success, or ErrorInfo on failure.
func GetClusterInfo(ctx context.Context, configAPI *api.Config, restConfig *rest.Config, logger *zap.Logger) string {
	fields := []zap.Field{
		zap.String("current_context", configAPI.CurrentContext),
	}
	if contextInfo, ok := configAPI.Contexts[configAPI.CurrentContext]; ok {
		fields = append(fields, zap.String("cluster_name", contextInfo.Cluster))
		authInfoName := contextInfo.AuthInfo
		if authInfoName == "" {
			authInfoName = "<anonymous>"
		}
		fields = append(fields, zap.String("auth_info_name", authInfoName))
	}
	l := logger.With(fields...)

	// 1. Extract Basic Info
	clusterName := configAPI.Contexts[configAPI.CurrentContext].Cluster
	clusterInfo := configAPI.Clusters[clusterName]

	output := ClusterInfo{
		ContextName:           configAPI.CurrentContext,
		Endpoint:              clusterInfo.Server,
		TLSServerVerification: !restConfig.Insecure,
		// NOTE: 'Error' field removed from successful output struct
	}

	authMethod := "Unknown"
	authInfoName := configAPI.Contexts[configAPI.CurrentContext].AuthInfo
	if authInfoName != "" {
		userInfo := configAPI.AuthInfos[authInfoName]
		if userInfo.Exec != nil {
			authMethod = "Exec Plugin"
		} else if userInfo.AuthProvider != nil {
			authMethod = fmt.Sprintf("Auth Provider (%s)", userInfo.AuthProvider.Name)
		} else if userInfo.ClientCertificate != "" || len(userInfo.ClientCertificateData) > 0 {
			authMethod = "Client Certificate"
		} else if userInfo.Token != "" || userInfo.TokenFile != "" {
			authMethod = "Token"
		} else if userInfo.Username != "" || userInfo.Password != "" {
			authMethod = "Basic Auth (Username/Password)"
		} else if userInfo.Impersonate != "" {
			authMethod = "Impersonation"
		}
	} else {
		authMethod = "None/Anonymous"
	}
	output.AuthMethod = authMethod
	l = l.With(zap.String("determined_auth_method", authMethod))

	// 2. Create Clientset
	clientRestConfig := *restConfig
	clientRestConfig.Timeout = RequestTimeout

	clientset, err := kubernetes.NewForConfig(&clientRestConfig)
	if err != nil {
		wrappedErr := xerrors.Errorf("failed to create Kubernetes clientset from REST config: %w", err)
		l.Error("Cannot create Kubernetes clientset", zap.Error(wrappedErr))
		return createErrorJSON("Internal error: Failed to create Kubernetes client", err, logger)
	}

	// 3. Get Server Version (with Retries)
	var serverVersionStr string
	var lastErr error
	var contextCancelled bool

	l.Info("Attempting to connect and retrieve server version...") // Logged at Info level
	for attempt := 0; attempt <= MaxRetries; attempt++ {
		select {
		case <-ctx.Done():
			lastErr = ctx.Err()
			l.Warn("Context cancelled before fetching server version attempt.", zap.Error(lastErr)) // Logged at Warn level
			contextCancelled = true
		default:
		}
		if contextCancelled {
			break
		}

		if attempt > 0 {
			l.Info(fmt.Sprintf("Retrying server version fetch (attempt %d/%d)", attempt, MaxRetries), zap.Duration("delay", RetryDelay)) // Logged at Info level
			select {
			case <-time.After(RetryDelay):
			case <-ctx.Done():
				lastErr = ctx.Err()
				l.Warn("Context cancelled during retry delay.", zap.Error(lastErr)) // Logged at Warn level
				contextCancelled = true
			}
			if contextCancelled {
				break
			}
		}

		// *** THE ACTUAL CALL - NO CONTEXT ARGUMENT ***
		versionInfo, err := clientset.Discovery().ServerVersion()

		if err != nil {
			lastErr = err
			if xerrors.Is(err, context.Canceled) { // Note: May not happen often without direct ctx pass
				l.Warn("Context cancelled during server version fetch (detected via error).", zap.Error(err)) // Logged at Warn level
				contextCancelled = true
				break
			}
			l.Warn("Error fetching server version", zap.Int("attempt", attempt), zap.Error(err)) // Logged at Warn level
			if attempt == MaxRetries {
				l.Error("Max retries reached for server version fetch.", zap.Int("max_retries", MaxRetries)) // Logged at Error level
			}
			continue
		}

		serverVersionStr = versionInfo.GitVersion
		output.ServerVersion = serverVersionStr
		lastErr = nil
		l.Info("Successfully retrieved server version", zap.String("server_version", serverVersionStr)) // Logged at Info level
		break
	}

	if lastErr != nil {
		errMsg := "Failed to connect to cluster or get server version after retries"
		if contextCancelled || xerrors.Is(lastErr, context.Canceled) {
			errMsg = "Operation cancelled while trying to get server version"
		}
		l.Error(errMsg, zap.Int("retries_attempted", MaxRetries), zap.Error(lastErr)) // Logged at Error level
		return createErrorJSON(errMsg, lastErr, logger)
	}

	// 4. Marshal Success JSON
	successJSONBytes, err := json.Marshal(output)
	if err != nil {
		wrappedErr := xerrors.Errorf("internal error: failed to marshal successful cluster info: %w", err)
		l.Error("Failed to marshal success response", zap.Error(wrappedErr)) // Logged at Error level
		return createErrorJSON("Internal error: Failed to create success JSON response", err, logger)
	}

	// 5. Return Success JSON
	return string(successJSONBytes)
}

// --- Health Check Function ---

// VerifyHealth checks if the credentials provided via the restConfig have 'list' permissions.
// It *always* returns a JSON string: HealthStatus indicating success or detailing failures.
func VerifyHealth(ctx context.Context, restConfig *rest.Config, logger *zap.Logger) string {
	l := logger

	// 1. Create Clientset
	healthCheckRestConfig := *restConfig
	healthCheckRestConfig.Timeout = RequestTimeout
	clientset, err := kubernetes.NewForConfig(&healthCheckRestConfig)
	if err != nil {
		wrappedErr := xerrors.Errorf("failed to create Kubernetes clientset for health check: %w", err)
		l.Error("Cannot create clientset for health check", zap.Error(wrappedErr)) // Logged at Error level
		return createErrorJSON("Failed to create Kubernetes client for health check", err, logger)
	}

	l.Info("Starting health check: Verifying 'list' permissions for core resources...") // Logged at Info level

	missingPermissions := []HealthPermission{}
	totalChecks := len(resourcesToVerify)
	failedChecks := 0
	contextCancelled := false // Flag for cancellation

	// 2. Perform Permission Checks
	for i, resourceCheck := range resourcesToVerify {
		select {
		case <-ctx.Done():
			l.Warn("Context cancelled during health check loop, aborting remaining checks.", zap.Error(ctx.Err())) // Logged at Warn level
			missingPermissions = append(missingPermissions, HealthPermission{
				Reason: fmt.Sprintf("Health check aborted due to context cancellation: %v", ctx.Err()),
			})
			contextCancelled = true
		default:
		}
		if contextCancelled {
			break // Exit outer loop
		}

		resourceAttributes := &authorizationv1.ResourceAttributes{
			Group: resourceCheck.Group, Version: resourceCheck.Version, Resource: resourceCheck.Resource,
			Verb: "list", Namespace: "",
		}
		ssar := &authorizationv1.SelfSubjectAccessReview{
			Spec: authorizationv1.SelfSubjectAccessReviewSpec{ResourceAttributes: resourceAttributes},
		}

		checkLogger := l.With(
			zap.String("check_resource", resourceCheck.Friendly),
			zap.String("group", resourceCheck.Group), zap.String("version", resourceCheck.Version),
			zap.String("resource", resourceCheck.Resource), zap.String("verb", "list"),
		)
		checkLogger.Debug("Performing SelfSubjectAccessReview check") // Logged at Debug level

		reqCtx, cancel := context.WithTimeout(ctx, RequestTimeout)
		review, err := clientset.AuthorizationV1().SelfSubjectAccessReviews().Create(reqCtx, ssar, metav1.CreateOptions{})
		cancel()

		if err != nil {
			failedChecks++
			checkLogger.Error("Failed to perform SelfSubjectAccessReview API call", zap.Error(err)) // Logged at Error level
			missingPermissions = append(missingPermissions, HealthPermission{
				Group: resourceCheck.Group, Version: resourceCheck.Version, Resource: resourceCheck.Resource,
				Verb: "list", Scope: resourceCheck.Scope, Reason: fmt.Sprintf("API call failed: %s", err.Error()),
			})
			if xerrors.Is(err, context.Canceled) || xerrors.Is(err, context.DeadlineExceeded) {
				checkLogger.Warn("Context cancelled or deadline exceeded during SSAR check, aborting remaining checks.", zap.Error(err)) // Logged at Warn level
				contextCancelled = true
				break // Exit outer loop
			}
			continue // Move to next check on other API errors
		}

		if !review.Status.Allowed {
			failedChecks++
			reason := review.Status.Reason
			if reason == "" {
				reason = "Permission denied by RBAC policy (no specific reason provided)"
			}
			checkLogger.Warn("Permission check failed: 'list' denied", zap.String("reason", reason)) // Logged at Warn level
			missingPermissions = append(missingPermissions, HealthPermission{
				Group: resourceCheck.Group, Version: resourceCheck.Version, Resource: resourceCheck.Resource,
				Verb: "list", Scope: resourceCheck.Scope, Reason: reason,
			})
		} else {
			checkLogger.Debug("Permission check succeeded: 'list' allowed") // Logged at Debug level
		}

		if (i+1)%10 == 0 || (i+1) == totalChecks {
			l.Info(fmt.Sprintf("Health check progress: %d/%d checks completed", i+1, totalChecks), zap.Int("failed_count", failedChecks)) // Logged at Info level
		}
	} // End loop

	// 3. Format Result JSON
	var healthResult HealthStatus
	finalMessage := ""
	if len(missingPermissions) == 0 {
		finalMessage = "Health check completed: All required read permissions verified successfully."
		l.Info(finalMessage) // Logged at Info level
		healthResult = HealthStatus{
			Healthy: true,
			Message: finalMessage,
		}
	} else {
		if contextCancelled {
			finalMessage = fmt.Sprintf("Health check aborted after %d checks: %d required 'list' permission checks failed or were skipped.", failedChecks+len(missingPermissions)-1, len(missingPermissions))
		} else {
			finalMessage = fmt.Sprintf("Health check completed: %d/%d required 'list' permission checks failed.", len(missingPermissions), totalChecks)
		}
		l.Warn(finalMessage, zap.Int("missing_count", len(missingPermissions))) // Logged at Warn level
		healthResult = HealthStatus{
			Healthy:            false,
			Message:            finalMessage,
			MissingPermissions: missingPermissions,
		}
	}

	resultJSONBytes, err := json.Marshal(healthResult)
	if err != nil {
		wrappedErr := xerrors.Errorf("internal error: failed to marshal health status result: %w", err)
		l.Error("Failed to marshal health status response", zap.Error(wrappedErr)) // Logged at Error level
		return createErrorJSON("Internal error: Failed to create health status JSON response", err, logger)
	}

	return string(resultJSONBytes)
}

// --- Main Function ---

func DoHealthcheck(kubeConfig string) (bool, error) {
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		zapcore.Lock(os.Stderr),
		zapcore.DebugLevel,
	)
	logger := zap.New(core, zap.AddCallerSkip(0)) // Add caller info if needed
	defer func() {
		_ = logger.Sync() // Flush logs before exit
	}()

	tmpDir := os.TempDir()

	// Create the full path for the kubeconfig file
	kubeconfigPath := filepath.Join(tmpDir, "kubeconfig.yaml")

	// Write the string into the file with 0600 permissions
	err := os.WriteFile(kubeconfigPath, []byte(kubeConfig), 0600)
	if err != nil {
		return false, err
	}

	// --- Basic File Existence Check ---
	if _, err := os.Stat(kubeconfigPath); os.IsNotExist(err) {
		return false, err
	}

	// --- Create Root Context ---
	ctx := context.Background()

	// --- Validate Kubeconfig (Common Step) ---
	logger.Info("Validating kubeconfig file...", zap.String("path", kubeconfigPath)) // Logged at Info level
	_, restConfig, err := ValidateKubeconfig(ctx, kubeconfigPath, logger)
	if err != nil {
		return false, err
	}

	// --- Execute Action Based on Flag ---
	var resultJSON string
	resultJSON = VerifyHealth(ctx, restConfig, logger)
	logger.Info("Fetching cluster information...", zap.String("kubeconfig_path", kubeconfigPath)) // Logged at Info level

	// --- Determine Exit Code Based on Result ---
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(resultJSON), &result); err == nil {
		if isHealthy, ok := result["healthy"].(bool); !ok || !isHealthy {
			return false, nil
		}
	} else {
		return false, err
	}
	return true, nil
}

func DoDiscovery(kubeConfig string) (map[string]string, error) {
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		zapcore.Lock(os.Stderr),
		zapcore.DebugLevel,
	)
	logger := zap.New(core, zap.AddCallerSkip(0)) // Add caller info if needed
	defer func() {
		_ = logger.Sync() // Flush logs before exit
	}()

	tmpDir := os.TempDir()

	// Create the full path for the kubeconfig file
	kubeconfigPath := filepath.Join(tmpDir, "kubeconfig.yaml")

	// Write the string into the file with 0600 permissions
	err := os.WriteFile(kubeconfigPath, []byte(kubeConfig), 0600)
	if err != nil {
		return nil, err
	}

	// --- Basic File Existence Check ---
	if _, err := os.Stat(kubeconfigPath); os.IsNotExist(err) {
		return nil, err
	}

	// --- Create Root Context ---
	ctx := context.Background()

	// --- Validate Kubeconfig (Common Step) ---
	logger.Info("Validating kubeconfig file...", zap.String("path", kubeconfigPath)) // Logged at Info level
	configAPI, restConfig, err := ValidateKubeconfig(ctx, kubeconfigPath, logger)
	if err != nil {
		return nil, err
	}

	// --- Execute Action Based on Flag ---
	var resultJSON string

	logger.Info("Fetching cluster information...", zap.String("kubeconfig_path", kubeconfigPath)) // Logged at Info level
	resultJSON = GetClusterInfo(ctx, configAPI, restConfig, logger)

	// --- Determine Exit Code Based on Result ---
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(resultJSON), &result); err != nil {
		return nil, err
	}

	return convertToStringMap(result), nil
}

func convertToStringMap(input map[string]interface{}) map[string]string {
	result := make(map[string]string)
	for k, v := range input {
		result[k] = fmt.Sprintf("%v", v)
	}
	return result
}
