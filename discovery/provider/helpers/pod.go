package helpers

import (
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// --- Pod ---
type Pod struct {
	TypeMeta
	ObjectMeta
	Spec   PodSpec
	Status PodStatus // Define PodStatus and conversion later if needed
}

func ConvertPod(p *corev1.Pod) Pod {
	return Pod{
		TypeMeta:   ConvertTypeMeta(p.TypeMeta),
		ObjectMeta: ConvertObjectMeta(&p.ObjectMeta),
		Spec:       ConvertPodSpec(p.Spec),
		// Status:    ConvertPodStatus(p.Status), // Add if needed
	}
}

// --- PodSpec ---
type PodSpec struct {
	Volumes                       []Volume
	InitContainers                []Container
	Containers                    []Container
	EphemeralContainers           []EphemeralContainer
	RestartPolicy                 string // corev1.RestartPolicy
	TerminationGracePeriodSeconds *int64
	ActiveDeadlineSeconds         *int64
	DNSPolicy                     string // corev1.DNSPolicy
	NodeSelector                  map[string]string
	ServiceAccountName            string
	DeprecatedServiceAccount      string // Field is deprecated
	AutomountServiceAccountToken  *bool
	NodeName                      string
	HostNetwork                   bool
	HostPID                       bool
	HostIPC                       bool
	ShareProcessNamespace         *bool
	SecurityContext               *PodSecurityContext
	ImagePullSecrets              []LocalObjectReference
	Hostname                      string
	Subdomain                     string
	Affinity                      *Affinity
	SchedulerName                 string
	Tolerations                   []Toleration
	HostAliases                   []HostAlias
	PriorityClassName             string
	Priority                      *int32
	DNSConfig                     *PodDNSConfig
	ReadinessGates                []PodReadinessGate
	RuntimeClassName              *string
	EnableServiceLinks            *bool
	PreemptionPolicy              *string                      // corev1.PreemptionPolicy
	Overhead                      map[string]resource.Quantity // corev1.ResourceList
	TopologySpreadConstraints     []TopologySpreadConstraint
	SetHostnameAsFQDN             *bool
	OS                            *PodOS
	HostUsers                     *bool
	SchedulingGates               []PodSchedulingGate
	ResourceClaims                []PodResourceClaim
}

func ConvertRestartPolicy(rp corev1.RestartPolicy) string {
	return string(rp)
}

func ConvertDNSPolicy(p corev1.DNSPolicy) string {
	return string(p)
}

func ConvertPreemptionPolicy(pp *corev1.PreemptionPolicy) *string {
	if pp == nil {
		return nil
	}
	s := string(*pp)
	return &s
}

func ConvertResourceList(rl corev1.ResourceList) map[string]resource.Quantity {
	if rl == nil {
		return nil
	}
	result := make(map[string]resource.Quantity)
	for k, v := range rl {
		result[string(k)] = v
	}
	return result
}

func ConvertPodSpec(spec corev1.PodSpec) PodSpec {
	return PodSpec{
		Volumes:                       ConvertVolumes(spec.Volumes),
		InitContainers:                ConvertContainers(spec.InitContainers),
		Containers:                    ConvertContainers(spec.Containers),
		EphemeralContainers:           ConvertEphemeralContainers(spec.EphemeralContainers),
		RestartPolicy:                 ConvertRestartPolicy(spec.RestartPolicy),
		TerminationGracePeriodSeconds: spec.TerminationGracePeriodSeconds,
		ActiveDeadlineSeconds:         spec.ActiveDeadlineSeconds,
		DNSPolicy:                     ConvertDNSPolicy(spec.DNSPolicy),
		NodeSelector:                  spec.NodeSelector,
		ServiceAccountName:            spec.ServiceAccountName,
		DeprecatedServiceAccount:      spec.DeprecatedServiceAccount,
		AutomountServiceAccountToken:  spec.AutomountServiceAccountToken,
		NodeName:                      spec.NodeName,
		HostNetwork:                   spec.HostNetwork,
		HostPID:                       spec.HostPID,
		HostIPC:                       spec.HostIPC,
		ShareProcessNamespace:         spec.ShareProcessNamespace,
		SecurityContext:               ConvertPodSecurityContext(spec.SecurityContext),
		ImagePullSecrets:              ConvertLocalObjectReferences(spec.ImagePullSecrets),
		Hostname:                      spec.Hostname,
		Subdomain:                     spec.Subdomain,
		Affinity:                      ConvertAffinity(spec.Affinity),
		SchedulerName:                 spec.SchedulerName,
		Tolerations:                   ConvertTolerations(spec.Tolerations),
		HostAliases:                   ConvertHostAliases(spec.HostAliases),
		PriorityClassName:             spec.PriorityClassName,
		Priority:                      spec.Priority,
		DNSConfig:                     ConvertPodDNSConfig(spec.DNSConfig),
		ReadinessGates:                ConvertPodReadinessGates(spec.ReadinessGates),
		RuntimeClassName:              spec.RuntimeClassName,
		EnableServiceLinks:            spec.EnableServiceLinks,
		PreemptionPolicy:              ConvertPreemptionPolicy(spec.PreemptionPolicy),
		Overhead:                      ConvertResourceList(spec.Overhead),
		TopologySpreadConstraints:     ConvertTopologySpreadConstraints(spec.TopologySpreadConstraints),
		SetHostnameAsFQDN:             spec.SetHostnameAsFQDN,
		OS:                            ConvertPodOS(spec.OS),
		HostUsers:                     spec.HostUsers,
		SchedulingGates:               ConvertPodSchedulingGates(spec.SchedulingGates),
		ResourceClaims:                ConvertPodResourceClaims(spec.ResourceClaims),
	}
}

// --- PodTemplateSpec ---
type PodTemplateSpec struct {
	ObjectMeta // Embedded, needs conversion function
	Spec       PodSpec
}

// Note: ConvertObjectMeta is in model_helpers.go
func ConvertPodTemplateSpec(spec corev1.PodTemplateSpec) PodTemplateSpec {
	return PodTemplateSpec{
		ObjectMeta: ConvertObjectMeta(&spec.ObjectMeta),
		Spec:       ConvertPodSpec(spec.Spec),
	}
}

// --- Container ---
type Container struct {
	Name                     string
	Image                    string
	Command                  []string
	Args                     []string
	WorkingDir               string
	Ports                    []ContainerPort
	EnvFrom                  []EnvFromSource
	Env                      []EnvVar
	Resources                ResourceRequirements
	VolumeMounts             []VolumeMount
	VolumeDevices            []VolumeDevice
	LivenessProbe            *Probe
	ReadinessProbe           *Probe
	StartupProbe             *Probe
	Lifecycle                *Lifecycle
	TerminationMessagePath   string
	TerminationMessagePolicy string // corev1.TerminationMessagePolicy
	ImagePullPolicy          string // corev1.PullPolicy
	SecurityContext          *SecurityContext
	Stdin                    bool
	StdinOnce                bool
	TTY                      bool
	ResizePolicy             []ContainerResizePolicy
}

type ContainerResizePolicy struct {
	Resource      string // corev1.ResourceName
	RestartPolicy string // corev1.ContainerResizeRestartPolicy
}

func ConvertContainerResizePolicy(p corev1.ContainerResizePolicy) ContainerResizePolicy {
	return ContainerResizePolicy{
		Resource:      string(p.ResourceName), // ResourceName is the field in corev1
		RestartPolicy: string(p.RestartPolicy),
	}
}

func ConvertContainerResizePolicies(policies []corev1.ContainerResizePolicy) []ContainerResizePolicy {
	if policies == nil {
		return nil
	}
	result := make([]ContainerResizePolicy, len(policies))
	for i, p := range policies {
		result[i] = ConvertContainerResizePolicy(p)
	}
	return result
}

func ConvertTerminationMessagePolicy(p corev1.TerminationMessagePolicy) string {
	return string(p)
}

// ConvertPullPolicy is defined in volume.go, reuse or redefine if needed
// func ConvertPullPolicy(p corev1.PullPolicy) string { return string(p) }

func ConvertContainer(c corev1.Container) Container {
	return Container{
		Name:                     c.Name,
		Image:                    c.Image,
		Command:                  c.Command,
		Args:                     c.Args,
		WorkingDir:               c.WorkingDir,
		Ports:                    ConvertContainerPorts(c.Ports),
		EnvFrom:                  ConvertEnvFromSources(c.EnvFrom),
		Env:                      ConvertEnvVars(c.Env),
		Resources:                ConvertResourceRequirements(c.Resources), // Defined in volume.go
		VolumeMounts:             ConvertVolumeMounts(c.VolumeMounts),
		VolumeDevices:            ConvertVolumeDevices(c.VolumeDevices),
		LivenessProbe:            ConvertProbe(c.LivenessProbe),
		ReadinessProbe:           ConvertProbe(c.ReadinessProbe),
		StartupProbe:             ConvertProbe(c.StartupProbe),
		Lifecycle:                ConvertLifecycle(c.Lifecycle),
		TerminationMessagePath:   c.TerminationMessagePath,
		TerminationMessagePolicy: ConvertTerminationMessagePolicy(c.TerminationMessagePolicy),
		ImagePullPolicy:          ConvertPullPolicy(&c.ImagePullPolicy),
		SecurityContext:          ConvertSecurityContext(c.SecurityContext),
		Stdin:                    c.Stdin,
		StdinOnce:                c.StdinOnce,
		TTY:                      c.TTY,
		ResizePolicy:             ConvertContainerResizePolicies(c.ResizePolicy),
	}
}

func ConvertContainers(containers []corev1.Container) []Container {
	if containers == nil {
		return nil
	}
	result := make([]Container, len(containers))
	for i, c := range containers {
		result[i] = ConvertContainer(c)
	}
	return result
}

// --- EphemeralContainer ---
type EphemeralContainer struct {
	EphemeralContainerCommon
	TargetContainerName string
}

type EphemeralContainerCommon struct {
	Name                     string
	Image                    string
	Command                  []string
	Args                     []string
	WorkingDir               string
	Ports                    []ContainerPort
	EnvFrom                  []EnvFromSource
	Env                      []EnvVar
	Resources                ResourceRequirements // Defined in volume.go
	VolumeMounts             []VolumeMount
	VolumeDevices            []VolumeDevice
	LivenessProbe            *Probe
	ReadinessProbe           *Probe
	StartupProbe             *Probe
	Lifecycle                *Lifecycle
	TerminationMessagePath   string
	TerminationMessagePolicy string // corev1.TerminationMessagePolicy
	ImagePullPolicy          string // corev1.PullPolicy
	SecurityContext          *SecurityContext
	Stdin                    bool
	StdinOnce                bool
	TTY                      bool
	ResizePolicy             []ContainerResizePolicy
}

func ConvertEphemeralContainerCommon(ecc corev1.EphemeralContainerCommon) EphemeralContainerCommon {
	return EphemeralContainerCommon{
		Name:                     ecc.Name,
		Image:                    ecc.Image,
		Command:                  ecc.Command,
		Args:                     ecc.Args,
		WorkingDir:               ecc.WorkingDir,
		Ports:                    ConvertContainerPorts(ecc.Ports),
		EnvFrom:                  ConvertEnvFromSources(ecc.EnvFrom),
		Env:                      ConvertEnvVars(ecc.Env),
		Resources:                ConvertResourceRequirements(ecc.Resources), // Defined in volume.go
		VolumeMounts:             ConvertVolumeMounts(ecc.VolumeMounts),
		VolumeDevices:            ConvertVolumeDevices(ecc.VolumeDevices),
		LivenessProbe:            ConvertProbe(ecc.LivenessProbe),
		ReadinessProbe:           ConvertProbe(ecc.ReadinessProbe),
		StartupProbe:             ConvertProbe(ecc.StartupProbe),
		Lifecycle:                ConvertLifecycle(ecc.Lifecycle),
		TerminationMessagePath:   ecc.TerminationMessagePath,
		TerminationMessagePolicy: ConvertTerminationMessagePolicy(ecc.TerminationMessagePolicy),
		ImagePullPolicy:          ConvertPullPolicy(&ecc.ImagePullPolicy),
		SecurityContext:          ConvertSecurityContext(ecc.SecurityContext),
		Stdin:                    ecc.Stdin,
		StdinOnce:                ecc.StdinOnce,
		TTY:                      ecc.TTY,
		ResizePolicy:             ConvertContainerResizePolicies(ecc.ResizePolicy),
	}
}

func ConvertEphemeralContainer(ec corev1.EphemeralContainer) EphemeralContainer {
	return EphemeralContainer{
		EphemeralContainerCommon: ConvertEphemeralContainerCommon(ec.EphemeralContainerCommon),
		TargetContainerName:      ec.TargetContainerName,
	}
}

func ConvertEphemeralContainers(containers []corev1.EphemeralContainer) []EphemeralContainer {
	if containers == nil {
		return nil
	}
	result := make([]EphemeralContainer, len(containers))
	for i, c := range containers {
		result[i] = ConvertEphemeralContainer(c)
	}
	return result
}

// --- ContainerPort ---
type ContainerPort struct {
	Name          string
	HostPort      int32
	ContainerPort int32
	Protocol      string // corev1.Protocol
	HostIP        string
}

func ConvertProtocol(p corev1.Protocol) string {
	return string(p)
}

func ConvertContainerPort(p corev1.ContainerPort) ContainerPort {
	return ContainerPort{
		Name:          p.Name,
		HostPort:      p.HostPort,
		ContainerPort: p.ContainerPort,
		Protocol:      ConvertProtocol(p.Protocol),
		HostIP:        p.HostIP,
	}
}

func ConvertContainerPorts(ports []corev1.ContainerPort) []ContainerPort {
	if ports == nil {
		return nil
	}
	result := make([]ContainerPort, len(ports))
	for i, p := range ports {
		result[i] = ConvertContainerPort(p)
	}
	return result
}

// --- EnvVar ---
type EnvVar struct {
	Name      string
	Value     string
	ValueFrom *EnvVarSource
}

type EnvVarSource struct {
	FieldRef         *ObjectFieldSelector   // Defined in volume.go
	ResourceFieldRef *ResourceFieldSelector // Defined in volume.go
	ConfigMapKeyRef  *ConfigMapKeySelector
	SecretKeyRef     *SecretKeySelector
}

type ConfigMapKeySelector struct {
	LocalObjectReference // Name embedded
	Key                  string
	Optional             *bool
}

type SecretKeySelector struct {
	LocalObjectReference // Name embedded
	Key                  string
	Optional             *bool
}

func ConvertConfigMapKeySelector(s *corev1.ConfigMapKeySelector) *ConfigMapKeySelector {
	if s == nil {
		return nil
	}
	return &ConfigMapKeySelector{
		LocalObjectReference: ConvertLocalObjectReference(s.LocalObjectReference),
		Key:                  s.Key,
		Optional:             s.Optional,
	}
}

func ConvertSecretKeySelector(s *corev1.SecretKeySelector) *SecretKeySelector {
	if s == nil {
		return nil
	}
	return &SecretKeySelector{
		LocalObjectReference: ConvertLocalObjectReference(s.LocalObjectReference),
		Key:                  s.Key,
		Optional:             s.Optional,
	}
}

func ConvertEnvVarSource(s *corev1.EnvVarSource) *EnvVarSource {
	if s == nil {
		return nil
	}
	return &EnvVarSource{
		FieldRef:         ConvertObjectFieldSelector(s.FieldRef),           // From volume.go
		ResourceFieldRef: ConvertResourceFieldSelector(s.ResourceFieldRef), // From volume.go
		ConfigMapKeyRef:  ConvertConfigMapKeySelector(s.ConfigMapKeyRef),
		SecretKeyRef:     ConvertSecretKeySelector(s.SecretKeyRef),
	}
}

func ConvertEnvVar(e corev1.EnvVar) EnvVar {
	return EnvVar{
		Name:      e.Name,
		Value:     e.Value,
		ValueFrom: ConvertEnvVarSource(e.ValueFrom),
	}
}

func ConvertEnvVars(envs []corev1.EnvVar) []EnvVar {
	if envs == nil {
		return nil
	}
	result := make([]EnvVar, len(envs))
	for i, e := range envs {
		result[i] = ConvertEnvVar(e)
	}
	return result
}

// --- EnvFromSource ---
type EnvFromSource struct {
	Prefix       string
	ConfigMapRef *ConfigMapEnvSource
	SecretRef    *SecretEnvSource
}

type ConfigMapEnvSource struct {
	LocalObjectReference // Name embedded
	Optional             *bool
}

type SecretEnvSource struct {
	LocalObjectReference // Name embedded
	Optional             *bool
}

func ConvertConfigMapEnvSource(s *corev1.ConfigMapEnvSource) *ConfigMapEnvSource {
	if s == nil {
		return nil
	}
	return &ConfigMapEnvSource{
		LocalObjectReference: ConvertLocalObjectReference(s.LocalObjectReference),
		Optional:             s.Optional,
	}
}

func ConvertSecretEnvSource(s *corev1.SecretEnvSource) *SecretEnvSource {
	if s == nil {
		return nil
	}
	return &SecretEnvSource{
		LocalObjectReference: ConvertLocalObjectReference(s.LocalObjectReference),
		Optional:             s.Optional,
	}
}

func ConvertEnvFromSource(s corev1.EnvFromSource) EnvFromSource {
	return EnvFromSource{
		Prefix:       s.Prefix,
		ConfigMapRef: ConvertConfigMapEnvSource(s.ConfigMapRef),
		SecretRef:    ConvertSecretEnvSource(s.SecretRef),
	}
}

func ConvertEnvFromSources(sources []corev1.EnvFromSource) []EnvFromSource {
	if sources == nil {
		return nil
	}
	result := make([]EnvFromSource, len(sources))
	for i, s := range sources {
		result[i] = ConvertEnvFromSource(s)
	}
	return result
}

// --- VolumeMount ---
type VolumeMount struct {
	Name             string
	ReadOnly         bool
	MountPath        string
	SubPath          string
	MountPropagation *string // corev1.MountPropagationMode
	SubPathExpr      string
}

func ConvertMountPropagationMode(m *corev1.MountPropagationMode) *string {
	if m == nil {
		return nil
	}
	s := string(*m)
	return &s
}

func ConvertVolumeMount(vm corev1.VolumeMount) VolumeMount {
	return VolumeMount{
		Name:             vm.Name,
		ReadOnly:         vm.ReadOnly,
		MountPath:        vm.MountPath,
		SubPath:          vm.SubPath,
		MountPropagation: ConvertMountPropagationMode(vm.MountPropagation),
		SubPathExpr:      vm.SubPathExpr,
	}
}

func ConvertVolumeMounts(mounts []corev1.VolumeMount) []VolumeMount {
	if mounts == nil {
		return nil
	}
	result := make([]VolumeMount, len(mounts))
	for i, vm := range mounts {
		result[i] = ConvertVolumeMount(vm)
	}
	return result
}

// --- VolumeDevice ---
type VolumeDevice struct {
	Name       string
	DevicePath string
}

func ConvertVolumeDevice(vd corev1.VolumeDevice) VolumeDevice {
	return VolumeDevice{
		Name:       vd.Name,
		DevicePath: vd.DevicePath,
	}
}

func ConvertVolumeDevices(devices []corev1.VolumeDevice) []VolumeDevice {
	if devices == nil {
		return nil
	}
	result := make([]VolumeDevice, len(devices))
	for i, vd := range devices {
		result[i] = ConvertVolumeDevice(vd)
	}
	return result
}

// --- Probe ---
type Probe struct {
	ProbeHandler
	InitialDelaySeconds           int32
	TimeoutSeconds                int32
	PeriodSeconds                 int32
	SuccessThreshold              int32
	FailureThreshold              int32
	TerminationGracePeriodSeconds *int64
}

type ProbeHandler struct {
	Exec      *ExecAction
	HTTPGet   *HTTPGetAction
	TCPSocket *TCPSocketAction
	GRPC      *GRPCAction // Added in later k8s versions
}

type ExecAction struct {
	Command []string
}

type HTTPGetAction struct {
	Path        string
	Port        int // Should be IntOrString, simplified
	Host        string
	Scheme      string // corev1.URIScheme
	HTTPHeaders []HTTPHeader
}

type TCPSocketAction struct {
	Port int // Should be IntOrString, simplified
	Host string
}

type HTTPHeader struct {
	Name  string
	Value string
}

type GRPCAction struct {
	Port    int32
	Service *string
}

func ConvertExecAction(a *corev1.ExecAction) *ExecAction {
	if a == nil {
		return nil
	}
	return &ExecAction{
		Command: a.Command,
	}
}

func ConvertIntOrString(ios *intstr.IntOrString) int {
	if ios == nil {
		return 0 // Default value
	}
	// Corrected: Use IntVal field from intstr.IntOrString
	if ios.Type == intstr.Int {
		return int(ios.IntVal)
	}
	// Decide how to handle string type, maybe parse or return default
	// Returning 0 for simplicity if it's a string
	return 0
}

func ConvertURIScheme(s corev1.URIScheme) string {
	return string(s)
}

func ConvertHTTPHeader(h corev1.HTTPHeader) HTTPHeader {
	return HTTPHeader{
		Name:  h.Name,
		Value: h.Value,
	}
}

func ConvertHTTPHeaders(headers []corev1.HTTPHeader) []HTTPHeader {
	if headers == nil {
		return nil
	}
	result := make([]HTTPHeader, len(headers))
	for i, h := range headers {
		result[i] = ConvertHTTPHeader(h)
	}
	return result
}

func ConvertHTTPGetAction(a *corev1.HTTPGetAction) *HTTPGetAction {
	if a == nil {
		return nil
	}
	// Warning: a.Port is IntOrString, simplifying to int
	return &HTTPGetAction{
		Path:        a.Path,
		Port:        ConvertIntOrString(&a.Port), // Simplified conversion
		Host:        a.Host,
		Scheme:      ConvertURIScheme(a.Scheme),
		HTTPHeaders: ConvertHTTPHeaders(a.HTTPHeaders),
	}
}

func ConvertTCPSocketAction(a *corev1.TCPSocketAction) *TCPSocketAction {
	if a == nil {
		return nil
	}
	// Warning: a.Port is IntOrString, simplifying to int
	return &TCPSocketAction{
		Port: ConvertIntOrString(&a.Port), // Simplified conversion
		Host: a.Host,
	}
}

func ConvertGRPCAction(a *corev1.GRPCAction) *GRPCAction {
	if a == nil {
		return nil
	}
	return &GRPCAction{
		Port:    a.Port,
		Service: a.Service,
	}
}

func ConvertProbeHandler(h corev1.ProbeHandler) ProbeHandler {
	return ProbeHandler{
		Exec:      ConvertExecAction(h.Exec),
		HTTPGet:   ConvertHTTPGetAction(h.HTTPGet),
		TCPSocket: ConvertTCPSocketAction(h.TCPSocket),
		GRPC:      ConvertGRPCAction(h.GRPC),
	}
}

func ConvertProbe(p *corev1.Probe) *Probe {
	if p == nil {
		return nil
	}
	return &Probe{
		ProbeHandler:                  ConvertProbeHandler(p.ProbeHandler),
		InitialDelaySeconds:           p.InitialDelaySeconds,
		TimeoutSeconds:                p.TimeoutSeconds,
		PeriodSeconds:                 p.PeriodSeconds,
		SuccessThreshold:              p.SuccessThreshold,
		FailureThreshold:              p.FailureThreshold,
		TerminationGracePeriodSeconds: p.TerminationGracePeriodSeconds,
	}
}

// --- Lifecycle ---
type Lifecycle struct {
	PostStart *LifecycleHandler
	PreStop   *LifecycleHandler
}

type LifecycleHandler struct {
	Exec      *ExecAction
	HTTPGet   *HTTPGetAction
	TCPSocket *TCPSocketAction
	// Sleep added in recent k8s
	Sleep *SleepAction
}

type SleepAction struct {
	Seconds int64
}

func ConvertSleepAction(a *corev1.SleepAction) *SleepAction {
	if a == nil {
		return nil
	}
	return &SleepAction{
		Seconds: a.Seconds,
	}
}

func ConvertLifecycleHandler(h *corev1.LifecycleHandler) *LifecycleHandler {
	if h == nil {
		return nil
	}
	return &LifecycleHandler{
		Exec:      ConvertExecAction(h.Exec),
		HTTPGet:   ConvertHTTPGetAction(h.HTTPGet),
		TCPSocket: ConvertTCPSocketAction(h.TCPSocket),
		Sleep:     ConvertSleepAction(h.Sleep),
	}
}

func ConvertLifecycle(l *corev1.Lifecycle) *Lifecycle {
	if l == nil {
		return nil
	}
	return &Lifecycle{
		PostStart: ConvertLifecycleHandler(l.PostStart),
		PreStop:   ConvertLifecycleHandler(l.PreStop),
	}
}

// --- SecurityContext ---
type SecurityContext struct {
	Capabilities             *Capabilities
	Privileged               *bool
	SELinuxOptions           *SELinuxOptions
	WindowsOptions           *WindowsSecurityContextOptions
	RunAsUser                *int64
	RunAsGroup               *int64
	RunAsNonRoot             *bool
	ReadOnlyRootFilesystem   *bool
	AllowPrivilegeEscalation *bool
	ProcMount                *string // corev1.ProcMountType
	SeccompProfile           *SeccompProfile
}

type Capabilities struct {
	Add  []string // corev1.Capability
	Drop []string // corev1.Capability
}

type SELinuxOptions struct {
	User  string
	Role  string
	Type  string
	Level string
}

type WindowsSecurityContextOptions struct {
	GMSACredentialSpecName *string
	GMSACredentialSpec     *string
	RunAsUserName          *string
	HostProcess            *bool
}

type SeccompProfile struct {
	Type             string // corev1.SeccompProfileType
	LocalhostProfile *string
}

func ConvertCapability(c corev1.Capability) string {
	return string(c)
}

func ConvertCapabilitiesList(caps []corev1.Capability) []string {
	if caps == nil {
		return nil
	}
	result := make([]string, len(caps))
	for i, c := range caps {
		result[i] = ConvertCapability(c)
	}
	return result
}

func ConvertCapabilities(c *corev1.Capabilities) *Capabilities {
	if c == nil {
		return nil
	}
	return &Capabilities{
		Add:  ConvertCapabilitiesList(c.Add),
		Drop: ConvertCapabilitiesList(c.Drop),
	}
}

func ConvertSELinuxOptions(o *corev1.SELinuxOptions) *SELinuxOptions {
	if o == nil {
		return nil
	}
	return &SELinuxOptions{
		User:  o.User,
		Role:  o.Role,
		Type:  o.Type,
		Level: o.Level,
	}
}

func ConvertWindowsSecurityContextOptions(o *corev1.WindowsSecurityContextOptions) *WindowsSecurityContextOptions {
	if o == nil {
		return nil
	}
	return &WindowsSecurityContextOptions{
		GMSACredentialSpecName: o.GMSACredentialSpecName,
		GMSACredentialSpec:     o.GMSACredentialSpec,
		RunAsUserName:          o.RunAsUserName,
		HostProcess:            o.HostProcess,
	}
}

func ConvertProcMountType(p *corev1.ProcMountType) *string {
	if p == nil {
		return nil
	}
	s := string(*p)
	return &s
}

func ConvertSeccompProfileType(p corev1.SeccompProfileType) string {
	return string(p)
}

func ConvertSeccompProfile(p *corev1.SeccompProfile) *SeccompProfile {
	if p == nil {
		return nil
	}
	return &SeccompProfile{
		Type:             ConvertSeccompProfileType(p.Type),
		LocalhostProfile: p.LocalhostProfile,
	}
}

func ConvertSecurityContext(sc *corev1.SecurityContext) *SecurityContext {
	if sc == nil {
		return nil
	}
	return &SecurityContext{
		Capabilities:             ConvertCapabilities(sc.Capabilities),
		Privileged:               sc.Privileged,
		SELinuxOptions:           ConvertSELinuxOptions(sc.SELinuxOptions),
		WindowsOptions:           ConvertWindowsSecurityContextOptions(sc.WindowsOptions),
		RunAsUser:                sc.RunAsUser,
		RunAsGroup:               sc.RunAsGroup,
		RunAsNonRoot:             sc.RunAsNonRoot,
		ReadOnlyRootFilesystem:   sc.ReadOnlyRootFilesystem,
		AllowPrivilegeEscalation: sc.AllowPrivilegeEscalation,
		ProcMount:                ConvertProcMountType(sc.ProcMount),
		SeccompProfile:           ConvertSeccompProfile(sc.SeccompProfile),
	}
}

// --- PodSecurityContext ---
type PodSecurityContext struct {
	SELinuxOptions      *SELinuxOptions
	WindowsOptions      *WindowsSecurityContextOptions
	RunAsUser           *int64
	RunAsGroup          *int64
	RunAsNonRoot        *bool
	SupplementalGroups  []int64
	FSGroup             *int64
	Sysctls             []Sysctl
	FSGroupChangePolicy *string // corev1.PodFSGroupChangePolicy
	SeccompProfile      *SeccompProfile
}

type Sysctl struct {
	Name  string
	Value string
}

func ConvertSysctl(s corev1.Sysctl) Sysctl {
	return Sysctl{
		Name:  s.Name,
		Value: s.Value,
	}
}

func ConvertSysctls(sysctls []corev1.Sysctl) []Sysctl {
	if sysctls == nil {
		return nil
	}
	result := make([]Sysctl, len(sysctls))
	for i, s := range sysctls {
		result[i] = ConvertSysctl(s)
	}
	return result
}

func ConvertPodFSGroupChangePolicy(p *corev1.PodFSGroupChangePolicy) *string {
	if p == nil {
		return nil
	}
	s := string(*p)
	return &s
}

func ConvertPodSecurityContext(psc *corev1.PodSecurityContext) *PodSecurityContext {
	if psc == nil {
		return nil
	}
	return &PodSecurityContext{
		SELinuxOptions:      ConvertSELinuxOptions(psc.SELinuxOptions),
		WindowsOptions:      ConvertWindowsSecurityContextOptions(psc.WindowsOptions),
		RunAsUser:           psc.RunAsUser,
		RunAsGroup:          psc.RunAsGroup,
		RunAsNonRoot:        psc.RunAsNonRoot,
		SupplementalGroups:  psc.SupplementalGroups,
		FSGroup:             psc.FSGroup,
		Sysctls:             ConvertSysctls(psc.Sysctls),
		FSGroupChangePolicy: ConvertPodFSGroupChangePolicy(psc.FSGroupChangePolicy),
		SeccompProfile:      ConvertSeccompProfile(psc.SeccompProfile),
	}
}

// --- Affinity ---
type Affinity struct {
	NodeAffinity    *NodeAffinity
	PodAffinity     *PodAffinity
	PodAntiAffinity *PodAntiAffinity
}

type NodeAffinity struct {
	RequiredDuringSchedulingIgnoredDuringExecution  *NodeSelector
	PreferredDuringSchedulingIgnoredDuringExecution []PreferredSchedulingTerm
}

type PodAffinity struct {
	RequiredDuringSchedulingIgnoredDuringExecution  []PodAffinityTerm
	PreferredDuringSchedulingIgnoredDuringExecution []WeightedPodAffinityTerm
}

type PodAntiAffinity struct {
	RequiredDuringSchedulingIgnoredDuringExecution  []PodAffinityTerm
	PreferredDuringSchedulingIgnoredDuringExecution []WeightedPodAffinityTerm
}

type PreferredSchedulingTerm struct {
	Weight     int32
	Preference NodeSelectorTerm
}

type WeightedPodAffinityTerm struct {
	Weight          int32
	PodAffinityTerm PodAffinityTerm
}

type PodAffinityTerm struct {
	LabelSelector     *LabelSelector
	Namespaces        []string
	TopologyKey       string
	NamespaceSelector *LabelSelector
	MatchLabelKeys    []string
	MismatchLabelKeys []string
}

type NodeSelector struct {
	NodeSelectorTerms []NodeSelectorTerm
}

type NodeSelectorTerm struct {
	MatchExpressions []NodeSelectorRequirement
	MatchFields      []NodeSelectorRequirement
}

type NodeSelectorRequirement struct {
	Key      string
	Operator string
	Values   []string
}

func ConvertNodeSelectorOperator(op corev1.NodeSelectorOperator) string {
	return string(op)
}

func ConvertNodeSelectorRequirement(r corev1.NodeSelectorRequirement) NodeSelectorRequirement {
	return NodeSelectorRequirement{
		Key:      r.Key,
		Operator: ConvertNodeSelectorOperator(r.Operator),
		Values:   r.Values,
	}
}

func ConvertNodeSelectorRequirements(reqs []corev1.NodeSelectorRequirement) []NodeSelectorRequirement {
	if reqs == nil {
		return nil
	}
	result := make([]NodeSelectorRequirement, len(reqs))
	for i, r := range reqs {
		result[i] = ConvertNodeSelectorRequirement(r)
	}
	return result
}

func ConvertNodeSelectorTerm(t corev1.NodeSelectorTerm) NodeSelectorTerm {
	return NodeSelectorTerm{
		MatchExpressions: ConvertNodeSelectorRequirements(t.MatchExpressions),
		MatchFields:      ConvertNodeSelectorRequirements(t.MatchFields),
	}
}

func ConvertNodeSelectorTerms(terms []corev1.NodeSelectorTerm) []NodeSelectorTerm {
	if terms == nil {
		return nil
	}
	result := make([]NodeSelectorTerm, len(terms))
	for i, t := range terms {
		result[i] = ConvertNodeSelectorTerm(t)
	}
	return result
}

func ConvertNodeSelector(s *corev1.NodeSelector) *NodeSelector {
	if s == nil {
		return nil
	}
	return &NodeSelector{
		NodeSelectorTerms: ConvertNodeSelectorTerms(s.NodeSelectorTerms),
	}
}

func ConvertPreferredSchedulingTerm(t corev1.PreferredSchedulingTerm) PreferredSchedulingTerm {
	return PreferredSchedulingTerm{
		Weight:     t.Weight,
		Preference: ConvertNodeSelectorTerm(t.Preference),
	}
}

func ConvertPreferredSchedulingTerms(terms []corev1.PreferredSchedulingTerm) []PreferredSchedulingTerm {
	if terms == nil {
		return nil
	}
	result := make([]PreferredSchedulingTerm, len(terms))
	for i, t := range terms {
		result[i] = ConvertPreferredSchedulingTerm(t)
	}
	return result
}

func ConvertNodeAffinity(na *corev1.NodeAffinity) *NodeAffinity {
	if na == nil {
		return nil
	}
	return &NodeAffinity{
		RequiredDuringSchedulingIgnoredDuringExecution:  ConvertNodeSelector(na.RequiredDuringSchedulingIgnoredDuringExecution),
		PreferredDuringSchedulingIgnoredDuringExecution: ConvertPreferredSchedulingTerms(na.PreferredDuringSchedulingIgnoredDuringExecution),
	}
}

func ConvertPodAffinityTerm(t corev1.PodAffinityTerm) PodAffinityTerm {
	return PodAffinityTerm{
		LabelSelector:     ConvertLabelSelector(t.LabelSelector),
		Namespaces:        t.Namespaces,
		TopologyKey:       t.TopologyKey,
		NamespaceSelector: ConvertLabelSelector(t.NamespaceSelector),
		MatchLabelKeys:    t.MatchLabelKeys,
		MismatchLabelKeys: t.MismatchLabelKeys,
	}
}

func ConvertPodAffinityTerms(terms []corev1.PodAffinityTerm) []PodAffinityTerm {
	if terms == nil {
		return nil
	}
	result := make([]PodAffinityTerm, len(terms))
	for i, t := range terms {
		result[i] = ConvertPodAffinityTerm(t)
	}
	return result
}

func ConvertWeightedPodAffinityTerm(t corev1.WeightedPodAffinityTerm) WeightedPodAffinityTerm {
	return WeightedPodAffinityTerm{
		Weight:          t.Weight,
		PodAffinityTerm: ConvertPodAffinityTerm(t.PodAffinityTerm),
	}
}

func ConvertWeightedPodAffinityTerms(terms []corev1.WeightedPodAffinityTerm) []WeightedPodAffinityTerm {
	if terms == nil {
		return nil
	}
	result := make([]WeightedPodAffinityTerm, len(terms))
	for i, t := range terms {
		result[i] = ConvertWeightedPodAffinityTerm(t)
	}
	return result
}

func ConvertPodAffinity(pa *corev1.PodAffinity) *PodAffinity {
	if pa == nil {
		return nil
	}
	return &PodAffinity{
		RequiredDuringSchedulingIgnoredDuringExecution:  ConvertPodAffinityTerms(pa.RequiredDuringSchedulingIgnoredDuringExecution),
		PreferredDuringSchedulingIgnoredDuringExecution: ConvertWeightedPodAffinityTerms(pa.PreferredDuringSchedulingIgnoredDuringExecution),
	}
}

func ConvertPodAntiAffinity(paa *corev1.PodAntiAffinity) *PodAntiAffinity {
	if paa == nil {
		return nil
	}
	return &PodAntiAffinity{
		RequiredDuringSchedulingIgnoredDuringExecution:  ConvertPodAffinityTerms(paa.RequiredDuringSchedulingIgnoredDuringExecution),
		PreferredDuringSchedulingIgnoredDuringExecution: ConvertWeightedPodAffinityTerms(paa.PreferredDuringSchedulingIgnoredDuringExecution),
	}
}

func ConvertAffinity(a *corev1.Affinity) *Affinity {
	if a == nil {
		return nil
	}
	return &Affinity{
		NodeAffinity:    ConvertNodeAffinity(a.NodeAffinity),
		PodAffinity:     ConvertPodAffinity(a.PodAffinity),
		PodAntiAffinity: ConvertPodAntiAffinity(a.PodAntiAffinity),
	}
}

// --- Toleration ---
type Toleration struct {
	Key               string
	Operator          string // corev1.TolerationOperator
	Value             string
	Effect            string // corev1.TaintEffect
	TolerationSeconds *int64
}

func ConvertTolerationOperator(op corev1.TolerationOperator) string {
	return string(op)
}

func ConvertToleration(t corev1.Toleration) Toleration {
	return Toleration{
		Key:               t.Key,
		Operator:          ConvertTolerationOperator(t.Operator),
		Value:             t.Value,
		Effect:            ConvertTaintEffect(t.Effect),
		TolerationSeconds: t.TolerationSeconds,
	}
}

func ConvertTolerations(tolerations []corev1.Toleration) []Toleration {
	if tolerations == nil {
		return nil
	}
	result := make([]Toleration, len(tolerations))
	for i, t := range tolerations {
		result[i] = ConvertToleration(t)
	}
	return result
}

// --- HostAlias ---
type HostAlias struct {
	IP        string
	Hostnames []string
}

func ConvertHostAlias(ha corev1.HostAlias) HostAlias {
	return HostAlias{
		IP:        ha.IP,
		Hostnames: ha.Hostnames,
	}
}

func ConvertHostAliases(aliases []corev1.HostAlias) []HostAlias {
	if aliases == nil {
		return nil
	}
	result := make([]HostAlias, len(aliases))
	for i, ha := range aliases {
		result[i] = ConvertHostAlias(ha)
	}
	return result
}

// --- PodDNSConfig ---
type PodDNSConfig struct {
	Nameservers []string
	Searches    []string
	Options     []PodDNSConfigOption
}

type PodDNSConfigOption struct {
	Name  string
	Value *string
}

func ConvertPodDNSConfigOption(o corev1.PodDNSConfigOption) PodDNSConfigOption {
	return PodDNSConfigOption{
		Name:  o.Name,
		Value: o.Value,
	}
}

func ConvertPodDNSConfigOptions(options []corev1.PodDNSConfigOption) []PodDNSConfigOption {
	if options == nil {
		return nil
	}
	result := make([]PodDNSConfigOption, len(options))
	for i, o := range options {
		result[i] = ConvertPodDNSConfigOption(o)
	}
	return result
}

func ConvertPodDNSConfig(c *corev1.PodDNSConfig) *PodDNSConfig {
	if c == nil {
		return nil
	}
	return &PodDNSConfig{
		Nameservers: c.Nameservers,
		Searches:    c.Searches,
		Options:     ConvertPodDNSConfigOptions(c.Options),
	}
}

// --- PodReadinessGate ---
type PodReadinessGate struct {
	ConditionType string // corev1.PodConditionType
}

func ConvertPodConditionType(ct corev1.PodConditionType) string {
	return string(ct)
}

func ConvertPodReadinessGate(g corev1.PodReadinessGate) PodReadinessGate {
	return PodReadinessGate{
		ConditionType: ConvertPodConditionType(g.ConditionType),
	}
}

func ConvertPodReadinessGates(gates []corev1.PodReadinessGate) []PodReadinessGate {
	if gates == nil {
		return nil
	}
	result := make([]PodReadinessGate, len(gates))
	for i, g := range gates {
		result[i] = ConvertPodReadinessGate(g)
	}
	return result
}

// --- TopologySpreadConstraint ---
type TopologySpreadConstraint struct {
	MaxSkew            int32
	TopologyKey        string
	WhenUnsatisfiable  string
	LabelSelector      *LabelSelector
	MinDomains         *int32
	NodeAffinityPolicy *string
	NodeTaintsPolicy   *string
	MatchLabelKeys     []string
}

func ConvertUnsatisfiableConstraintAction(a corev1.UnsatisfiableConstraintAction) string {
	return string(a)
}

func ConvertNodeInclusionPolicy(p *corev1.NodeInclusionPolicy) *string {
	if p == nil {
		return nil
	}
	s := string(*p)
	return &s
}

func ConvertTopologySpreadConstraint(c corev1.TopologySpreadConstraint) TopologySpreadConstraint {
	return TopologySpreadConstraint{
		MaxSkew:            c.MaxSkew,
		TopologyKey:        c.TopologyKey,
		WhenUnsatisfiable:  ConvertUnsatisfiableConstraintAction(c.WhenUnsatisfiable),
		LabelSelector:      ConvertLabelSelector(c.LabelSelector),
		MinDomains:         c.MinDomains,
		NodeAffinityPolicy: ConvertNodeInclusionPolicy(c.NodeAffinityPolicy),
		NodeTaintsPolicy:   ConvertNodeInclusionPolicy(c.NodeTaintsPolicy),
		MatchLabelKeys:     c.MatchLabelKeys,
	}
}

func ConvertTopologySpreadConstraints(constraints []corev1.TopologySpreadConstraint) []TopologySpreadConstraint {
	if constraints == nil {
		return nil
	}
	result := make([]TopologySpreadConstraint, len(constraints))
	for i, c := range constraints {
		result[i] = ConvertTopologySpreadConstraint(c)
	}
	return result
}

// --- PodOS ---
type PodOS struct {
	Name string // corev1.OSName
}

func ConvertOSName(n corev1.OSName) string {
	return string(n)
}

func ConvertPodOS(os *corev1.PodOS) *PodOS {
	if os == nil {
		return nil
	}
	return &PodOS{
		Name: ConvertOSName(os.Name),
	}
}

// --- PodSchedulingGate ---
type PodSchedulingGate struct {
	Name string
}

func ConvertPodSchedulingGate(g corev1.PodSchedulingGate) PodSchedulingGate {
	return PodSchedulingGate{
		Name: g.Name,
	}
}

func ConvertPodSchedulingGates(gates []corev1.PodSchedulingGate) []PodSchedulingGate {
	if gates == nil {
		return nil
	}
	result := make([]PodSchedulingGate, len(gates))
	for i, g := range gates {
		result[i] = ConvertPodSchedulingGate(g)
	}
	return result
}

// --- PodResourceClaim ---
type PodResourceClaim struct {
	Name string
	// Removed Source field, fields are direct now
	ResourceClaimName         *string
	ResourceClaimTemplateName *string
}

// Removed ClaimSource struct and ConvertClaimSource function
/*
type ClaimSource struct {
	ResourceClaimName         *string
	ResourceClaimTemplateName *string
}

func ConvertClaimSource(s corev1.ClaimSource) ClaimSource {
	return ClaimSource{
		ResourceClaimName:         s.ResourceClaimName,
		ResourceClaimTemplateName: s.ResourceClaimTemplateName,
	}
}
*/

func ConvertPodResourceClaim(prc corev1.PodResourceClaim) PodResourceClaim {
	return PodResourceClaim{
		Name: prc.Name,
		// Access fields directly from corev1.PodResourceClaim
		ResourceClaimName:         prc.ResourceClaimName,
		ResourceClaimTemplateName: prc.ResourceClaimTemplateName,
	}
}

func ConvertPodResourceClaims(claims []corev1.PodResourceClaim) []PodResourceClaim {
	if claims == nil {
		return nil
	}
	result := make([]PodResourceClaim, len(claims))
	for i, c := range claims {
		result[i] = ConvertPodResourceClaim(c)
	}
	return result
}

// --- PodStatus (Basic Definition - expand if needed) ---
type PodStatus struct {
	Phase      string // corev1.PodPhase
	Conditions []PodCondition
	Message    string
	Reason     string
	// Add other fields like HostIP, PodIP, StartTime, ContainerStatuses etc. if required
}

type PodCondition struct {
	Type               string     // corev1.PodConditionType
	Status             string     // corev1.ConditionStatus
	LastProbeTime      *time.Time // Changed from metav1.Time
	LastTransitionTime *time.Time // Changed from metav1.Time
	Reason             string
	Message            string
}

// Add ConvertPodStatus, ConvertPodCondition etc. if a full status conversion is needed

// Helper for ImagePullSecrets (used in PodSpec)
func ConvertLocalObjectReferences(refs []corev1.LocalObjectReference) []LocalObjectReference {
	if refs == nil {
		return nil
	}
	result := make([]LocalObjectReference, len(refs))
	for i, r := range refs {
		result[i] = ConvertLocalObjectReference(r)
	}
	return result
}

// NOTE: ResourceRequirements, ConvertResourceRequirements, LocalObjectReference, ConvertLocalObjectReference defined in volume.go
// NOTE: IntOrString type from "k8s.io/apimachinery/pkg/util/intstr" is simplified to int in Probes
// NOTE: Types from metav1 (LabelSelector, Time, ObjectMeta, TypeMeta) are used directly or converted via functions in model_helpers.go
