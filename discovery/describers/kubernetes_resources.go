package describers

import (
	"context"
	"encoding/json" // Required for marshaling the output JSON
	"errors"        // Required for errors.Is
	"fmt"
	"github.com/opengovern/og-describer-kubernetes/discovery/provider"
	"log"
	"os"
	"os/signal" // Required for signal handling

	// Required for Retry-After parsing
	"strings" // Required for error aggregation and kind mapping
	"sync"    // For caching discovery results safely
	"syscall" // Required for signal types (SIGINT, SIGTERM)
	"time"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta" // Required for RESTMapper
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema" // Required for GVR, GVK

	// For UID type
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/discovery"               // Import discovery client
	"k8s.io/client-go/discovery/cached/memory" // For caching discovery results
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"       // Import rest package for Config type
	"k8s.io/client-go/restmapper" // Import RESTMapper
	"k8s.io/client-go/tools/clientcmd"
)

// --- Configurable Constants ---
const (
	defaultQPS           float32       = 50
	defaultBurst         int           = 100
	defaultLimit         int64         = 5000
	retryInitialInterval time.Duration = 1 * time.Second
	retryFactor          float64       = 2.0
	retrySteps           int           = 5
	retryJitter          float64       = 0.1
	retryMaxInterval     time.Duration = 15 * time.Second
	apiCallTimeout       time.Duration = 30 * time.Second
	idleTimeout          time.Duration = 3 * time.Minute
	neverExceedTimeout   time.Duration = 30 * time.Minute
)

// Global discovery information
var restMapper meta.RESTMapper
var discoveryOnce sync.Once

// --- Application Arguments Structure ---
type AppArgs struct {
	KubeconfigPath  string
	ResourceType    string
	Limit           int64
	QPS             float64
	Burst           int
	StreamMode      bool
	IncludeStatus   bool
	IncludeMetadata bool
}

// --- Output Structures ---

// Represents the data for a single listed Kubernetes object (map with snake_case keys)
type K8sObjectData map[string]interface{}

// Represents the collection of items for a specific resource type in non-stream mode
type ResourceTypeResult struct {
	ResourceTable string                                   `json:"resource_table"`
	TotalCount    int                                      `json:"total_count"`
	Items         []provider.KubernetesResourceDescription `json:"items"`
}

// Represents the final summary in non-stream mode
type ListSummary struct {
	Status              string         `json:"status"`
	TotalItemsProcessed int            `json:"total_items_processed"`
	ResourceTableCounts map[string]int `json:"resource_table_counts,omitempty"`
	Reason              string         `json:"reason,omitempty"`
	Error               string         `json:"error,omitempty"`
}

// Top-level structure for non-stream mode output
type NonStreamOutput struct {
	Results     map[string]ResourceTypeResult `json:"results"` // Map of resource_table -> results
	ListSummary ListSummary                   `json:"list_summary"`
}

// --- Kind to Resource Table Mapping ---
var kindToResourceTableMap = map[string]string{
	"clusterrole":              "k8_cluster_role",
	"clusterrolebinding":       "k8_cluster_role_binding",
	"configmap":                "k8_config_map",
	"cronjob":                  "k8_cronjob",
	"customresourcedefinition": "k8_custom_resource_definition",
	"daemonset":                "k8_daemonset",
	"deployment":               "k8_deployment",
	"endpointslice":            "k8_endpoint_slice",
	"endpoints":                "k8_endpoints",
	"event":                    "k8_event",
	"horizontalpodautoscaler":  "k8_horizontal_pod_autoscaler",
	"ingress":                  "k8_ingress",
	"job":                      "k8_job",
	"limitrange":               "k8_limit_range",
	"namespace":                "k8_namespace",
	"networkpolicy":            "k8_network_policy",
	"node":                     "k8_node",
	"persistentvolume":         "k8_persistent_volume",
	"persistentvolumeclaim":    "k8_persistent_volume_claim",
	"pod":                      "k8_pod",
	"poddisruptionbudget":      "k8_pod_disruption_budget",
	"podtemplate":              "k8_pod_template",
	"replicaset":               "k8_replicaset",
	"replicationcontroller":    "k8_replication_controller",
	"resourcequota":            "k8_resource_quota",
	"role":                     "k8_role",
	"rolebinding":              "k8_role_binding",
	"secret":                   "k8_secret",
	"service":                  "k8_service",
	"serviceaccount":           "k8_service_account",
	"statefulset":              "k8_stateful_set",
	"storageclass":             "k8_storage_class",
}

// getResourceTable determines the resource table name based on the object's Kind.
func getResourceTable(kind string) string {
	lowerKind := strings.ToLower(kind)
	if ref, ok := kindToResourceTableMap[lowerKind]; ok {
		return ref
	}
	return "k8_custom_resource"
}

// --- Main Execution ---

func GetKubernetesResources(kubeConfig string) ([]provider.KubernetesResourceDescription, error) {
	limit := defaultLimit
	qps := float64(defaultQPS)
	burst := defaultBurst
	stream := false
	includeStatus := false
	includeMetadata := false

	kubeconfigPath := "./kubeconfig.yaml"

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

	appArgs := AppArgs{
		KubeconfigPath: kubeconfigPath, ResourceType: "", Limit: limit,
		QPS: qps, Burst: burst, StreamMode: stream, IncludeStatus: includeStatus,
		IncludeMetadata: includeMetadata,
	}

	resources, err := Execute(appArgs)
	return resources, err
}

// --- Execute Function (Main Application Logic) ---

func Execute(args AppArgs) ([]provider.KubernetesResourceDescription, error) {
	config, err := clientcmd.BuildConfigFromFlags("", args.KubeconfigPath)
	if err != nil {
		log.Printf("Error building kubeconfig: %v", err)
		if !args.StreamMode {
			failStatus := map[string]interface{}{"results": map[string]interface{}{}, "list_summary": ListSummary{Status: "failed", Error: fmt.Sprintf("kubeconfig error: %v", err)}}
			failJSON, _ := json.MarshalIndent(failStatus, "", "  ")
			fmt.Println(string(failJSON))
		}
		return nil, err
	}
	config.QPS = float32(args.QPS)
	config.Burst = args.Burst
	log.Printf("Using client rate limiting: QPS=%.2f, Burst=%d", config.QPS, config.Burst)

	mainCtx, rootCancel := context.WithCancel(context.Background())
	defer rootCancel()
	timedCtx, timedCancel := context.WithTimeout(mainCtx, neverExceedTimeout)
	defer timedCancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() { sig := <-sigChan; log.Printf("Received signal: %v. Shutting down...", sig); rootCancel() }()

	logMsg := fmt.Sprintf("Starting lister with hard timeout limit: %v, idle timeout: %v. Stream mode: %v. Include Status: %v. Include Metadata: %v.",
		neverExceedTimeout, idleTimeout, args.StreamMode, args.IncludeStatus, args.IncludeMetadata)
	isListAll := args.ResourceType == ""
	if isListAll {
		logMsg += " Listing all resource types."
	} else {
		logMsg += fmt.Sprintf(" Listing resource type '%s'.", args.ResourceType)
	}
	logMsg += " Press Ctrl+C to interrupt."
	log.Println(logMsg)

	// --- Run the appropriate Lister Logic ---
	var itemsProcessed int
	var runErr error
	// resultsMap stores the structured results grouped by resource_table for non-stream mode
	resultsMap := make(map[string]ResourceTypeResult) // Changed type here

	if isListAll {
		itemsProcessed, _, runErr = ListAllTypes(timedCtx, config, args, resultsMap)
	} else {
		itemsProcessed, runErr = RunLister(timedCtx, config, args, resultsMap)
		// Populate resourceTableCounts for the summary if successful
		if runErr == nil || !errors.Is(runErr, context.Canceled) && !errors.Is(runErr, context.DeadlineExceeded) {
			gvr, _, findErr := findResourceGVR(timedCtx, args.ResourceType)
			if findErr == nil {
				kind := gvr.Resource
				// Try getting kind from results if available (more accurate)
				// Need to access the correct map key (resourceTable)
				resourceTable := getResourceTable(kind)
				if resResult, ok := resultsMap[resourceTable]; ok && len(resResult.Items) > 0 {
					if resResult.Items[0].Kind != "" {
						kind = resResult.Items[0].Kind
					}
					// Re-calculate resourceTable based on potentially more accurate kind
					resourceTable = getResourceTable(kind)
				}
			} else {
				log.Printf("Warning: Could not re-find GVR for summary count for type '%s'", args.ResourceType)
			}
		}
	}
	completionStatus := ListSummary{TotalItemsProcessed: itemsProcessed}

	if runErr != nil {
		ctxErr := context.Cause(timedCtx)
		if strings.HasPrefix(runErr.Error(), "partial failure listing types:") {
			completionStatus.Status = "partial_failure"
			completionStatus.Error = runErr.Error()
			log.Printf("Lister finished with partial failures: %v", runErr)
			err = runErr
		} else if errors.Is(ctxErr, context.DeadlineExceeded) {
			completionStatus.Status = "interrupted"
			completionStatus.Reason = fmt.Sprintf("hard timeout limit reached (%v)", neverExceedTimeout)
			log.Printf("Lister stopped due to hard timeout: %v", neverExceedTimeout)
		} else if errors.Is(ctxErr, context.Canceled) {
			completionStatus.Status = "interrupted"
			completionStatus.Reason = "operation cancelled (possibly by signal)"
			log.Printf("Lister stopped due to cancellation.")
		} else if runErr.Error() == fmt.Sprintf("operation timed out due to inactivity after %v", idleTimeout) {
			completionStatus.Status = "interrupted"
			completionStatus.Reason = "idle timeout reached"
			log.Printf("Lister stopped due to idle timeout: %v", idleTimeout)
		} else {
			completionStatus.Status = "failed"
			completionStatus.Error = runErr.Error()
			log.Printf("Lister failed with error: %v", runErr)
			err = runErr
		}
	} else {
		completionStatus.Status = "completed"
		log.Println("Lister finished successfully.")
	}

	var results []provider.KubernetesResourceDescription
	for _, res := range resultsMap {
		results = append(results, res.Items...)
	}

	return results, err
}

// --- Core Lister Logic for Single Type ---

// RunLister lists a *specific* resource type.
// If not streaming, it populates the resultsMap. Otherwise, prints items to stdout.
// Returns total items processed and error.
func RunLister(ctx context.Context, config *rest.Config, args AppArgs, resultsMap map[string]ResourceTypeResult) (int, error) { // Changed resultsMap type
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return 0, fmt.Errorf("error creating dynamic client: %w", err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return 0, fmt.Errorf("error creating standard clientset: %w", err)
	}
	initErr := initDiscovery(clientset.Discovery())
	if initErr != nil {
		return 0, initErr
	}

	gvr, isNamespaced, err := findResourceGVR(ctx, args.ResourceType)
	if err != nil {
		log.Printf("Error finding resource type '%s': %v", args.ResourceType, err)
		return 0, err
	}

	// handleList now returns the buffered items if not streaming
	itemsData, itemsProcessed, err := handleList(ctx, dynamicClient, gvr, isNamespaced, args)

	// If not streaming and successful, populate the resultsMap grouped by resource_table
	if !args.StreamMode && err == nil {
		if itemsProcessed > 0 {
			kind := gvr.Resource // Use resource name as fallback kind
			// Try to get kind from first item for more accuracy
			if len(itemsData) > 0 {
				if itemsData[0].Kind != "" {
					kind = itemsData[0].Kind
				}
			}
			resourceTable := getResourceTable(kind)
			// Create the ResourceTypeResult struct and add it to the map
			resultsMap[resourceTable] = ResourceTypeResult{
				ResourceTable: resourceTable,
				TotalCount:    itemsProcessed,
				Items:         itemsData,
			}
		}
	}

	return itemsProcessed, err
}

// --- Core Lister Logic for All Types ---

// ListAllTypes discovers and lists all listable resource types.
// If not streaming, it populates the resultsMap. Otherwise, prints items to stdout.
// Returns total items processed across all types, a map of counts aggregated by resource_table,
// and an aggregated error if any type failed.
func ListAllTypes(ctx context.Context, config *rest.Config, args AppArgs, resultsMap map[string]ResourceTypeResult) (int, map[string]int, error) { // Changed resultsMap type
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return 0, nil, fmt.Errorf("error creating dynamic client: %w", err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return 0, nil, fmt.Errorf("error creating standard clientset: %w", err)
	}
	initErr := initDiscovery(clientset.Discovery())
	if initErr != nil {
		return 0, nil, initErr
	}

	discoveryClient := memory.NewMemCacheClient(clientset.Discovery())
	resList, err := discoveryClient.ServerPreferredResources()
	if err != nil {
		if apierrors.IsNotFound(err) {
			log.Printf("Warning: ServerPreferredResources API not found, discovery might be incomplete.")
			return 0, make(map[string]int), nil
		}
		if discovery.IsGroupDiscoveryFailedError(err) {
			log.Printf("Warning: Partial discovery failure: %v", err)
		} else {
			return 0, nil, fmt.Errorf("failed to discover server resources: %w", err)
		}
	}

	totalItemsOverall := 0
	var failedTypes []string
	var lastError error
	resourceTableCounts := make(map[string]int)

	for _, list := range resList {
		select {
		case <-ctx.Done():
			err = fmt.Errorf("operation timed out or cancelled during discovery: %w", context.Cause(ctx))
			goto EndListAll
		default:
		}
		if len(list.APIResources) == 0 {
			continue
		}
		gv, err := schema.ParseGroupVersion(list.GroupVersion)
		if err != nil {
			log.Printf("Skipping resource list for invalid GroupVersion '%s': %v", list.GroupVersion, err)
			continue
		}

		for _, res := range list.APIResources {
			select {
			case <-ctx.Done():
				err = fmt.Errorf("operation timed out or cancelled during discovery: %w", context.Cause(ctx))
				goto EndListAll
			default:
			}
			listable := false
			for _, verb := range res.Verbs {
				if verb == "list" {
					listable = true
					break
				}
			}
			if !listable {
				continue
			}

			gvr := gv.WithResource(res.Name)
			isNamespaced := res.Namespaced
			gvrString := gvr.String()

			// handleList now returns buffered items if not streaming
			itemsData, itemsProcessed, listErr := handleList(ctx, dynamicClient, gvr, isNamespaced, args)

			// Determine resource table and aggregate counts
			kindForTableRef := res.Kind
			if kindForTableRef == "" {
				kindForTableRef = gvr.Resource
				log.Printf("[%s] Warning: Kind missing in discovery for resource, using resource name '%s' for resource_table lookup.", gvrString, kindForTableRef)
			}
			resourceTable := getResourceTable(kindForTableRef)
			resourceTableCounts[resourceTable] += itemsProcessed

			totalItemsOverall += itemsProcessed

			// If not streaming and successful, populate the resultsMap grouped by resource_table
			if !args.StreamMode && listErr == nil {
				if itemsProcessed > 0 {
					// Ensure the key exists before appending
					if _, ok := resultsMap[resourceTable]; !ok {
						// Create the entry if it doesn't exist
						resultsMap[resourceTable] = ResourceTypeResult{
							ResourceTable: resourceTable,
							TotalCount:    itemsProcessed, // Set initial count
							Items:         itemsData,      // Assign initial items
						}
					} else {
						// Append items and update count if entry already exists
						existingResult := resultsMap[resourceTable]
						existingResult.Items = append(existingResult.Items, itemsData...)
						existingResult.TotalCount += itemsProcessed // Should already be correct from aggregation, but recalculate for safety
						resultsMap[resourceTable] = existingResult
					}
				} else if _, ok := resultsMap[resourceTable]; !ok && listErr == nil {
					// Ensure an entry exists even if count is 0 for successful lists
					resultsMap[resourceTable] = ResourceTypeResult{
						ResourceTable: resourceTable,
						TotalCount:    0,
						Items:         []provider.KubernetesResourceDescription{},
					}
				}
			}

			// Handle errors for this specific type
			if listErr != nil {
				log.Printf("Error listing resource type %s: %v", gvrString, listErr)
				failedTypes = append(failedTypes, gvrString)
				lastError = listErr
				ctxErr := context.Cause(ctx)
				if ctxErr != nil && (errors.Is(listErr, ctxErr) || errors.Is(listErr, context.DeadlineExceeded) || errors.Is(listErr, context.Canceled)) {
					err = listErr
					goto EndListAll
				}
				if listErr.Error() == fmt.Sprintf("operation timed out due to inactivity after %v", idleTimeout) {
					err = listErr
					goto EndListAll
				}
			}
		}
	}

EndListAll:
	if len(failedTypes) > 0 {
		if err != nil && (errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) || err.Error() == fmt.Sprintf("operation timed out due to inactivity after %v", idleTimeout)) {
			return totalItemsOverall, resourceTableCounts, err
		}
		return totalItemsOverall, resourceTableCounts, fmt.Errorf("partial failure listing types: encountered errors for %s (last error: %w)", strings.Join(failedTypes, ", "), lastError)
	}
	if err != nil {
		return totalItemsOverall, resourceTableCounts, err
	}

	log.Printf("ListAllTypes finished successfully.")
	return totalItemsOverall, resourceTableCounts, nil
}

// initDiscovery initializes the RESTMapper safely using sync.Once
func initDiscovery(discoveryClient discovery.DiscoveryInterface) error {
	var initErr error
	discoveryOnce.Do(func() {
		cachedDiscoveryClient := memory.NewMemCacheClient(discoveryClient)
		restMapper = restmapper.NewDeferredDiscoveryRESTMapper(cachedDiscoveryClient)
	})
	if restMapper == nil {
		return fmt.Errorf("RESTMapper failed to initialize")
	}
	return initErr
}

// --- Retry Logic (Enhanced) ---
func executeWithRetry(ctx context.Context, operationName string, operation func(context.Context) error) error {
	var lastErr error
	backoff := wait.Backoff{
		Duration: retryInitialInterval, Factor: retryFactor, Steps: retrySteps, Jitter: retryJitter, Cap: retryMaxInterval,
	}
	currentRetry := 0

	err := wait.ExponentialBackoffWithContext(ctx, backoff, func(ctx context.Context) (bool, error) {
		currentRetry++
		opErr := operation(ctx)
		if opErr == nil {
			return true, nil
		} // Success

		if ctx.Err() != nil {
			return false, ctx.Err()
		} // Context cancelled/timed out

		var retryAfter time.Duration
		isRetriable := false
		_, suggestsDelay := apierrors.SuggestsClientDelay(opErr)

		if apierrors.IsTooManyRequests(opErr) {
			isRetriable = true
			if statusErr, ok := opErr.(apierrors.APIStatus); ok {
				if details := statusErr.Status().Details; details != nil && details.RetryAfterSeconds > 0 {
					retryAfter = time.Duration(details.RetryAfterSeconds) * time.Second
					log.Printf("Warning: Operation '%s' rate limited (429). Server suggests Retry-After: %v", operationName, retryAfter)
				} else {
					log.Printf("Warning: Operation '%s' rate limited (429). No Retry-After suggested.", operationName)
				}
			} else {
				log.Printf("Warning: Operation '%s' rate limited (429), but couldn't parse status details.", operationName)
			}
		} else if apierrors.IsServerTimeout(opErr) || apierrors.IsInternalError(opErr) || apierrors.IsServiceUnavailable(opErr) || suggestsDelay {
			isRetriable = true
		}

		if isRetriable {
			log.Printf("Warning: Retriable error during %s (attempt %d/%d), retrying... Error: %v", operationName, currentRetry, backoff.Steps+1, opErr)
			lastErr = opErr
			if retryAfter > 0 {
				log.Printf("Info: Respecting server Retry-After suggestion (%v) if longer than backoff.", retryAfter)
			}
			return false, nil // Continue retrying
		}

		log.Printf("Non-retriable error during %s: %v", operationName, opErr)
		return false, opErr // Stop retrying
	})

	// Handle final errors from retry mechanism
	if err == wait.ErrWaitTimeout {
		log.Printf("Error: Operation '%s' timed out after %d retries (retry mechanism timeout). Last error: %v", operationName, backoff.Steps+1, lastErr)
		ctxErr := context.Cause(ctx)
		if ctxErr != nil {
			return fmt.Errorf("operation %s stopped due to context during retries: %w", operationName, ctxErr)
		}
		if lastErr != nil {
			return fmt.Errorf("operation %s timed out after retries: %w", operationName, lastErr)
		}
		return fmt.Errorf("operation %s timed out after retries", operationName)
	} else if err != nil {
		ctxErr := context.Cause(ctx)
		if ctxErr != nil && (errors.Is(err, ctxErr) || errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled)) {
			return err // Return context error directly
		}
		log.Printf("Error: Operation '%s' failed. Final error: %v", operationName, err)
		return fmt.Errorf("operation %s failed: %w", operationName, err)
	}

	return nil // Success
}

// --- Dynamic Resource Discovery ---
func findResourceGVR(ctx context.Context, resourceType string) (gvr schema.GroupVersionResource, namespaced bool, err error) {
	if restMapper == nil {
		return schema.GroupVersionResource{}, false, fmt.Errorf("RESTMapper not initialized before finding GVR")
	}
	gk := schema.GroupKind{Kind: resourceType}
	mapping, err := restMapper.RESTMapping(gk, "")
	if err == nil {
		gvr = mapping.Resource
		namespaced = mapping.Scope.Name() == meta.RESTScopeNameNamespace
		return gvr, namespaced, nil
	}
	kindMappingErr := err
	gvk, err2 := restMapper.KindFor(schema.GroupVersionResource{Resource: resourceType})
	if err2 == nil {
		mapping, err3 := restMapper.RESTMapping(gvk.GroupKind(), gvk.Version)
		if err3 == nil {
			gvr = mapping.Resource
			namespaced = mapping.Scope.Name() == meta.RESTScopeNameNamespace
			return gvr, namespaced, nil
		}
		err = fmt.Errorf("mapping GVK %s failed: %w", gvk.String(), err3)
	} else {
		err = fmt.Errorf("failed to find resource type '%s': kind mapping error: %v, resource mapping error: %v", resourceType, kindMappingErr, err2)
	}
	return schema.GroupVersionResource{}, false, err
}

// --- List Handler ---
// Fetches resources page by page for a given GVR.
// If streaming, prints item JSON to stdout immediately.
// If not streaming, buffers item data.
// Returns ([]K8sObjectData, int, error) -> (buffered items, count, error).
func handleList(ctx context.Context, dynamicClient dynamic.Interface, gvr schema.GroupVersionResource, isNamespaced bool, args AppArgs) ([]provider.KubernetesResourceDescription, int, error) {
	logPrefix := fmt.Sprintf("[%s] ", gvr.String())
	logLimit := args.Limit
	if logLimit <= 0 {
		logLimit = -1
	}

	var listErr error
	continueToken := ""
	totalListed := 0
	lastProgressTime := time.Now()
	var itemsDataBuffer []provider.KubernetesResourceDescription

	// Initialize buffer only if not streaming
	if !args.StreamMode {
		itemsDataBuffer = make([]provider.KubernetesResourceDescription, 0)
	}

	for { // Pagination loop
		select {
		case <-ctx.Done():
			log.Printf("%sOperation stopped: Main context deadline exceeded or cancelled: %v", logPrefix, context.Cause(ctx))
			return itemsDataBuffer, totalListed, fmt.Errorf("operation timed out or cancelled: %w", context.Cause(ctx))
		default: // Continue
		}
		if time.Since(lastProgressTime) > idleTimeout {
			log.Printf("%sOperation stopped: Idle timeout exceeded (%v).", logPrefix, idleTimeout)
			return itemsDataBuffer, totalListed, fmt.Errorf("operation timed out due to inactivity after %v", idleTimeout)
		}

		listOptions := metav1.ListOptions{Continue: continueToken}
		if args.Limit > 0 {
			listOptions.Limit = args.Limit
		}
		var list *unstructured.UnstructuredList

		listFunc := func(c context.Context) error {
			var opErr error
			apiCtx, cancel := context.WithTimeout(c, apiCallTimeout)
			defer cancel()
			resourceInterface := dynamicClient.Resource(gvr)
			if isNamespaced {
				list, opErr = resourceInterface.Namespace("").List(apiCtx, listOptions)
			} else {
				list, opErr = resourceInterface.List(apiCtx, listOptions)
			}
			return opErr
		}

		listErr = executeWithRetry(ctx, fmt.Sprintf("%sList (page)", logPrefix), listFunc)
		if listErr != nil {
			if ctxErr := context.Cause(ctx); ctxErr != nil && (errors.Is(listErr, ctxErr) || errors.Is(listErr, context.DeadlineExceeded) || errors.Is(listErr, context.Canceled)) {
				return itemsDataBuffer, totalListed, listErr // Context error
			}
			log.Printf("%sPersistent error listing page: %v", logPrefix, listErr)
			return itemsDataBuffer, totalListed, listErr // Persistent error
		}

		if list != nil {
			lastProgressTime = time.Now()
			if len(list.Items) > 0 && (totalListed == 0 || totalListed%(5*int(defaultLimit)) < len(list.Items)) {
				log.Printf("%sProcessing page (processed %d items so far)...", logPrefix, totalListed)
			}

			for i := range list.Items {
				item := &list.Items[i]
				kind := item.GetKind()
				if kind == "" {
					kind = gvr.Resource
					log.Printf("%sWarning: Kind missing for item %s, using resource name '%s' for resource_table lookup.", logPrefix, item.GetName(), kind)
				}
				resourceTable := getResourceTable(kind)
				lowerKind := strings.ToLower(kind)

				// Prepare the main output data map with snake_case keys
				outputData := provider.KubernetesResourceDescription{
					Kind:              lowerKind,
					ObjectName:        item.GetName(),
					Namespace:         item.GetNamespace(),
					UID:               fmt.Sprintf("%s", item.GetUID()),
					CreationTimestamp: item.GetCreationTimestamp().Format(time.RFC3339),
					ResourceVersion:   item.GetResourceVersion(),
					ResourceTable:     resourceTable,
					ApiVersion:        item.GetAPIVersion(),
				}

				// Conditionally include metadata
				if args.IncludeMetadata {
					metadataObj := map[string]interface{}{
						"labels":      item.GetLabels(),
						"annotations": item.GetAnnotations(),
					}
					if metadataObj["labels"] == nil {
						metadataObj["labels"] = make(map[string]string)
					}
					if metadataObj["annotations"] == nil {
						metadataObj["annotations"] = make(map[string]string)
					}
					outputData.Metadata = metadataObj
				}

				// Conditionally include status
				if args.IncludeStatus {
					statusVal := item.Object["status"]
					if statusMap, ok := statusVal.(map[string]interface{}); !ok || statusMap == nil {
						outputData.Status = make(map[string]interface{})
					} else {
						outputData.Status = statusVal
					}
				}

				// Decide whether to stream or buffer
				if args.StreamMode {
					jsonData, marshalErr := json.MarshalIndent(outputData, "", "  ")
					if marshalErr != nil {
						log.Printf("%sError marshaling item %s (%s): %v", logPrefix, item.GetName(), item.GetNamespace(), marshalErr)
						continue
					}
					fmt.Println(string(jsonData)) // Print directly
				} else {
					itemsDataBuffer = append(itemsDataBuffer, outputData) // Buffer if not streaming
				}
				totalListed++
			} // End item processing loop

			continueToken = list.GetContinue()
			if args.Limit <= 0 || continueToken == "" {
				break
			}
		} else {
			log.Printf("%sWarning: List result was nil after successful API call.", logPrefix)
			break
		}
	} // End of pagination loop

	log.Printf("%sSuccessfully listed %d items for this type.", logPrefix, totalListed)
	return itemsDataBuffer, totalListed, nil // Return buffer (nil if streaming), count, and success
}
