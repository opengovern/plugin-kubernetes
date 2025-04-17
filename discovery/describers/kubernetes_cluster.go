package describers

import (
	"context"
	"encoding/json"
	"fmt"
	model "github.com/opengovern/og-describer-kubernetes/discovery/provider"
	"os"
	"time"

	// Zap logger
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

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

func DoDiscovery(kubeConfig string) (*model.KubernetesClusterDescription, error) {
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

	kubeconfigPath := "./kubeconfig.yaml"

	fmt.Println("kubeConfig", kubeConfig)

	// Write the string into the file with 0600 permissions
	err := os.WriteFile(kubeconfigPath, []byte(kubeConfig), 0600)
	if err != nil {
		return nil, err
	}

	fmt.Println("kubeconfigPath", kubeconfigPath)
	time.Sleep(2 * time.Second)

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
	var result model.KubernetesClusterDescription
	if err := json.Unmarshal([]byte(resultJSON), &result); err != nil {
		return nil, err
	}

	err = os.Remove(kubeconfigPath)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
