package helpers

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	// metav1 is likely needed for some volume types, adding it back
)

type Volume struct {
	Name string
	VolumeSource
}

type HostPathVolumeSource struct {
	Path string
	Type *string // Should be *corev1.HostPathType ideally, but simplifying
}

func ConvertHostPathType(hpt *corev1.HostPathType) *string {
	if hpt == nil {
		return nil
	}
	s := string(*hpt)
	return &s
}

func ConvertHostPathVolumeSource(s *corev1.HostPathVolumeSource) *HostPathVolumeSource {
	if s == nil {
		return nil
	}
	return &HostPathVolumeSource{
		Path: s.Path,
		Type: ConvertHostPathType(s.Type),
	}
}

// --- Added definitions for all VolumeSource types ---

type GCEPersistentDiskVolumeSource struct {
	PDName    string
	FSType    string
	Partition int32
	ReadOnly  bool
}

func ConvertGCEPersistentDiskVolumeSource(s *corev1.GCEPersistentDiskVolumeSource) *GCEPersistentDiskVolumeSource {
	if s == nil {
		return nil
	}
	return &GCEPersistentDiskVolumeSource{
		PDName:    s.PDName,
		FSType:    s.FSType,
		Partition: s.Partition,
		ReadOnly:  s.ReadOnly,
	}
}

type AWSElasticBlockStoreVolumeSource struct {
	VolumeID  string
	FSType    string
	Partition int32
	ReadOnly  bool
}

func ConvertAWSElasticBlockStoreVolumeSource(s *corev1.AWSElasticBlockStoreVolumeSource) *AWSElasticBlockStoreVolumeSource {
	if s == nil {
		return nil
	}
	return &AWSElasticBlockStoreVolumeSource{
		VolumeID:  s.VolumeID,
		FSType:    s.FSType,
		Partition: s.Partition,
		ReadOnly:  s.ReadOnly,
	}
}

type GitRepoVolumeSource struct {
	Repository string
	Revision   string
	Directory  string
}

func ConvertGitRepoVolumeSource(s *corev1.GitRepoVolumeSource) *GitRepoVolumeSource {
	if s == nil {
		return nil
	}
	return &GitRepoVolumeSource{
		Repository: s.Repository,
		Revision:   s.Revision,
		Directory:  s.Directory,
	}
}

type ISCSIVolumeSource struct {
	TargetPortal      string
	IQN               string
	Lun               int32
	ISCSIInterface    string
	FSType            string
	ReadOnly          bool
	Portals           []string
	DiscoveryCHAPAuth bool
	SessionCHAPAuth   bool
	SecretRef         *LocalObjectReference
	InitiatorName     *string
}

func ConvertISCSIVolumeSource(s *corev1.ISCSIVolumeSource) *ISCSIVolumeSource {
	if s == nil {
		return nil
	}
	var secretRef *LocalObjectReference
	if s.SecretRef != nil {
		ref := ConvertLocalObjectReference(*s.SecretRef)
		secretRef = &ref
	}
	return &ISCSIVolumeSource{
		TargetPortal:      s.TargetPortal,
		IQN:               s.IQN,
		Lun:               s.Lun,
		ISCSIInterface:    s.ISCSIInterface,
		FSType:            s.FSType,
		ReadOnly:          s.ReadOnly,
		Portals:           s.Portals,
		DiscoveryCHAPAuth: s.DiscoveryCHAPAuth,
		SessionCHAPAuth:   s.SessionCHAPAuth,
		SecretRef:         secretRef,
		InitiatorName:     s.InitiatorName,
	}
}

type GlusterfsVolumeSource struct {
	EndpointsName string
	Path          string
	ReadOnly      bool
}

func ConvertGlusterfsVolumeSource(s *corev1.GlusterfsVolumeSource) *GlusterfsVolumeSource {
	if s == nil {
		return nil
	}
	return &GlusterfsVolumeSource{
		EndpointsName: s.EndpointsName,
		Path:          s.Path,
		ReadOnly:      s.ReadOnly,
	}
}

type RBDVolumeSource struct {
	CephMonitors []string
	RBDImage     string
	FSType       string
	RBDPool      string
	RadosUser    string
	Keyring      string
	SecretRef    *LocalObjectReference
	ReadOnly     bool
}

func ConvertRBDVolumeSource(s *corev1.RBDVolumeSource) *RBDVolumeSource {
	if s == nil {
		return nil
	}
	var secretRef *LocalObjectReference
	if s.SecretRef != nil {
		ref := ConvertLocalObjectReference(*s.SecretRef)
		secretRef = &ref
	}
	return &RBDVolumeSource{
		CephMonitors: s.CephMonitors,
		RBDImage:     s.RBDImage,
		FSType:       s.FSType,
		RBDPool:      s.RBDPool,
		RadosUser:    s.RadosUser,
		Keyring:      s.Keyring,
		SecretRef:    secretRef,
		ReadOnly:     s.ReadOnly,
	}
}

type FlexVolumeSource struct {
	Driver    string
	FSType    string
	SecretRef *LocalObjectReference
	ReadOnly  bool
	Options   map[string]string
}

func ConvertFlexVolumeSource(s *corev1.FlexVolumeSource) *FlexVolumeSource {
	if s == nil {
		return nil
	}
	var secretRef *LocalObjectReference
	if s.SecretRef != nil {
		ref := ConvertLocalObjectReference(*s.SecretRef)
		secretRef = &ref
	}
	return &FlexVolumeSource{
		Driver:    s.Driver,
		FSType:    s.FSType,
		SecretRef: secretRef,
		ReadOnly:  s.ReadOnly,
		Options:   s.Options,
	}
}

type CinderVolumeSource struct {
	VolumeID  string
	FSType    string
	ReadOnly  bool
	SecretRef *LocalObjectReference
}

func ConvertCinderVolumeSource(s *corev1.CinderVolumeSource) *CinderVolumeSource {
	if s == nil {
		return nil
	}
	var secretRef *LocalObjectReference
	if s.SecretRef != nil {
		ref := ConvertLocalObjectReference(*s.SecretRef)
		secretRef = &ref
	}
	return &CinderVolumeSource{
		VolumeID:  s.VolumeID,
		FSType:    s.FSType,
		ReadOnly:  s.ReadOnly,
		SecretRef: secretRef,
	}
}

type CephFSVolumeSource struct {
	Monitors   []string
	Path       string
	User       string
	SecretFile string
	SecretRef  *LocalObjectReference
	ReadOnly   bool
}

func ConvertCephFSVolumeSource(s *corev1.CephFSVolumeSource) *CephFSVolumeSource {
	if s == nil {
		return nil
	}
	var secretRef *LocalObjectReference
	if s.SecretRef != nil {
		ref := ConvertLocalObjectReference(*s.SecretRef)
		secretRef = &ref
	}
	return &CephFSVolumeSource{
		Monitors:   s.Monitors,
		Path:       s.Path,
		User:       s.User,
		SecretFile: s.SecretFile,
		SecretRef:  secretRef,
		ReadOnly:   s.ReadOnly,
	}
}

type FlockerVolumeSource struct {
	DatasetName string
	DatasetUUID string
}

func ConvertFlockerVolumeSource(s *corev1.FlockerVolumeSource) *FlockerVolumeSource {
	if s == nil {
		return nil
	}
	return &FlockerVolumeSource{
		DatasetName: s.DatasetName,
		DatasetUUID: s.DatasetUUID,
	}
}

type FCVolumeSource struct {
	TargetWWNs []string
	Lun        *int32
	FSType     string
	ReadOnly   bool
	WWIDs      []string
}

func ConvertFCVolumeSource(s *corev1.FCVolumeSource) *FCVolumeSource {
	if s == nil {
		return nil
	}
	return &FCVolumeSource{
		TargetWWNs: s.TargetWWNs,
		Lun:        s.Lun,
		FSType:     s.FSType,
		ReadOnly:   s.ReadOnly,
		WWIDs:      s.WWIDs,
	}
}

type AzureFileVolumeSource struct {
	SecretName string
	ShareName  string
	ReadOnly   bool
}

func ConvertAzureFileVolumeSource(s *corev1.AzureFileVolumeSource) *AzureFileVolumeSource {
	if s == nil {
		return nil
	}
	return &AzureFileVolumeSource{
		SecretName: s.SecretName,
		ShareName:  s.ShareName,
		ReadOnly:   s.ReadOnly,
	}
}

type VsphereVirtualDiskVolumeSource struct {
	VolumePath        string
	FSType            string
	StoragePolicyName string
	StoragePolicyID   string
}

func ConvertVsphereVirtualDiskVolumeSource(s *corev1.VsphereVirtualDiskVolumeSource) *VsphereVirtualDiskVolumeSource {
	if s == nil {
		return nil
	}
	return &VsphereVirtualDiskVolumeSource{
		VolumePath:        s.VolumePath,
		FSType:            s.FSType,
		StoragePolicyName: s.StoragePolicyName,
		StoragePolicyID:   s.StoragePolicyID,
	}
}

type QuobyteVolumeSource struct {
	Registry string
	Volume   string
	ReadOnly bool
	User     string
	Group    string
	Tenant   string
}

func ConvertQuobyteVolumeSource(s *corev1.QuobyteVolumeSource) *QuobyteVolumeSource {
	if s == nil {
		return nil
	}
	return &QuobyteVolumeSource{
		Registry: s.Registry,
		Volume:   s.Volume,
		ReadOnly: s.ReadOnly,
		User:     s.User,
		Group:    s.Group,
		Tenant:   s.Tenant,
	}
}

type AzureDiskVolumeSource struct {
	DiskName    string
	DataDiskURI string
	CachingMode *string // corev1.AzureDataDiskCachingMode
	FSType      *string
	ReadOnly    *bool
	Kind        *string // corev1.AzureDataDiskKind
}

func ConvertAzureDiskCachingMode(mode *corev1.AzureDataDiskCachingMode) *string {
	if mode == nil {
		return nil
	}
	s := string(*mode)
	return &s
}

func ConvertAzureDataDiskKind(kind *corev1.AzureDataDiskKind) *string {
	if kind == nil {
		return nil
	}
	s := string(*kind)
	return &s
}

func ConvertAzureDiskVolumeSource(s *corev1.AzureDiskVolumeSource) *AzureDiskVolumeSource {
	if s == nil {
		return nil
	}
	return &AzureDiskVolumeSource{
		DiskName:    s.DiskName,
		DataDiskURI: s.DataDiskURI,
		CachingMode: ConvertAzureDiskCachingMode(s.CachingMode),
		FSType:      s.FSType,
		ReadOnly:    s.ReadOnly,
		Kind:        ConvertAzureDataDiskKind(s.Kind),
	}
}

type PhotonPersistentDiskVolumeSource struct {
	PdID   string
	FSType string
}

func ConvertPhotonPersistentDiskVolumeSource(s *corev1.PhotonPersistentDiskVolumeSource) *PhotonPersistentDiskVolumeSource {
	if s == nil {
		return nil
	}
	return &PhotonPersistentDiskVolumeSource{
		PdID:   s.PdID,
		FSType: s.FSType,
	}
}

type PortworxVolumeSource struct {
	VolumeID string
	FSType   string
	ReadOnly bool
}

func ConvertPortworxVolumeSource(s *corev1.PortworxVolumeSource) *PortworxVolumeSource {
	if s == nil {
		return nil
	}
	return &PortworxVolumeSource{
		VolumeID: s.VolumeID,
		FSType:   s.FSType,
		ReadOnly: s.ReadOnly,
	}
}

type ScaleIOVolumeSource struct {
	Gateway          string
	System           string
	SecretRef        *LocalObjectReference
	SSLEnabled       bool
	ProtectionDomain string
	StoragePool      string
	StorageMode      string
	VolumeName       string
	FSType           string
	ReadOnly         bool
}

func ConvertScaleIOVolumeSource(s *corev1.ScaleIOVolumeSource) *ScaleIOVolumeSource {
	if s == nil {
		return nil
	}
	var secretRef *LocalObjectReference
	if s.SecretRef != nil {
		ref := ConvertLocalObjectReference(*s.SecretRef)
		secretRef = &ref
	}
	return &ScaleIOVolumeSource{
		Gateway:          s.Gateway,
		System:           s.System,
		SecretRef:        secretRef,
		SSLEnabled:       s.SSLEnabled,
		ProtectionDomain: s.ProtectionDomain,
		StoragePool:      s.StoragePool,
		StorageMode:      s.StorageMode,
		VolumeName:       s.VolumeName,
		FSType:           s.FSType,
		ReadOnly:         s.ReadOnly,
	}
}

type StorageOSVolumeSource struct {
	VolumeName      string
	VolumeNamespace string
	FSType          string
	ReadOnly        bool
	SecretRef       *LocalObjectReference
}

func ConvertStorageOSVolumeSource(s *corev1.StorageOSVolumeSource) *StorageOSVolumeSource {
	if s == nil {
		return nil
	}
	var secretRef *LocalObjectReference
	if s.SecretRef != nil {
		ref := ConvertLocalObjectReference(*s.SecretRef)
		secretRef = &ref
	}
	return &StorageOSVolumeSource{
		VolumeName:      s.VolumeName,
		VolumeNamespace: s.VolumeNamespace,
		FSType:          s.FSType,
		ReadOnly:        s.ReadOnly,
		SecretRef:       secretRef,
	}
}

type EphemeralVolumeSource struct {
	VolumeClaimTemplate *PersistentVolumeClaimTemplate
}

type PersistentVolumeClaimTemplate struct {
	ObjectMeta ObjectMeta
	Spec       PersistentVolumeClaimSpec
}

// Added new type for VolumeResourceRequirements
type VolumeResourceRequirements struct {
	Limits   map[string]resource.Quantity
	Requests map[string]resource.Quantity
}

// Existing ResourceRequirements for other uses
type ResourceRequirements struct {
	Limits   map[string]resource.Quantity
	Requests map[string]resource.Quantity
}

type TypedLocalObjectReference struct {
	APIGroup *string
	Kind     string
	Name     string
}

type TypedObjectReference struct {
	APIGroup  *string
	Kind      string
	Name      string
	Namespace *string
}

// Added conversion function for VolumeResourceRequirements
func ConvertVolumeResourceRequirements(vrr corev1.VolumeResourceRequirements) VolumeResourceRequirements {
	limits := make(map[string]resource.Quantity)
	for k, v := range vrr.Limits {
		limits[string(k)] = v
	}
	requests := make(map[string]resource.Quantity)
	for k, v := range vrr.Requests {
		requests[string(k)] = v
	}
	return VolumeResourceRequirements{
		Limits:   limits,
		Requests: requests,
	}
}

// Existing conversion function
func ConvertResourceRequirements(rr corev1.ResourceRequirements) ResourceRequirements {
	limits := make(map[string]resource.Quantity)
	for k, v := range rr.Limits {
		limits[string(k)] = v
	}
	requests := make(map[string]resource.Quantity)
	for k, v := range rr.Requests {
		requests[string(k)] = v
	}
	return ResourceRequirements{
		Limits:   limits,
		Requests: requests,
	}
}

func ConvertVolumeMode(vm *corev1.PersistentVolumeMode) *string {
	if vm == nil {
		return nil
	}
	s := string(*vm)
	return &s
}

func ConvertAccessModes(am []corev1.PersistentVolumeAccessMode) []string {
	if am == nil {
		return nil
	}
	res := make([]string, len(am))
	for i, mode := range am {
		res[i] = string(mode)
	}
	return res
}

func ConvertTypedLocalObjectReference(tlor *corev1.TypedLocalObjectReference) *TypedLocalObjectReference {
	if tlor == nil {
		return nil
	}
	return &TypedLocalObjectReference{
		APIGroup: tlor.APIGroup,
		Kind:     tlor.Kind,
		Name:     tlor.Name,
	}
}

func ConvertTypedObjectReference(tor *corev1.TypedObjectReference) *TypedObjectReference {
	if tor == nil {
		return nil
	}
	return &TypedObjectReference{
		APIGroup:  tor.APIGroup,
		Kind:      tor.Kind,
		Name:      tor.Name,
		Namespace: tor.Namespace,
	}
}

// Updated to call ConvertVolumeResourceRequirements
func ConvertPersistentVolumeClaimSpec(spec corev1.PersistentVolumeClaimSpec) PersistentVolumeClaimSpec {
	return PersistentVolumeClaimSpec{
		AccessModes:      ConvertAccessModes(spec.AccessModes),
		Selector:         ConvertLabelSelector(spec.Selector),
		Resources:        ConvertVolumeResourceRequirements(spec.Resources),
		VolumeName:       spec.VolumeName,
		StorageClassName: spec.StorageClassName,
		VolumeMode:       ConvertVolumeMode(spec.VolumeMode),
		DataSource:       ConvertTypedLocalObjectReference(spec.DataSource),
		DataSourceRef:    ConvertTypedObjectReference(spec.DataSourceRef),
	}
}

func ConvertPersistentVolumeClaimTemplate(pvct *corev1.PersistentVolumeClaimTemplate) *PersistentVolumeClaimTemplate {
	if pvct == nil {
		return nil
	}
	return &PersistentVolumeClaimTemplate{
		ObjectMeta: ConvertObjectMeta(&pvct.ObjectMeta),
		Spec:       ConvertPersistentVolumeClaimSpec(pvct.Spec),
	}
}

func ConvertEphemeralVolumeSource(s *corev1.EphemeralVolumeSource) *EphemeralVolumeSource {
	if s == nil {
		return nil
	}
	return &EphemeralVolumeSource{
		VolumeClaimTemplate: ConvertPersistentVolumeClaimTemplate(s.VolumeClaimTemplate),
	}
}

// Corrected ImageVolumeSource Definition
type ImageVolumeSource struct {
	Reference  string
	PullPolicy string
}

func ConvertPullPolicy(p *corev1.PullPolicy) string {
	if p == nil {
		return ""
	}
	return string(*p)
}

// Corrected ImageVolumeSource Conversion
func ConvertImageVolumeSource(s *corev1.ImageVolumeSource) *ImageVolumeSource {
	if s == nil {
		return nil
	}
	return &ImageVolumeSource{
		Reference:  s.Reference,
		PullPolicy: ConvertPullPolicy(&s.PullPolicy),
	}
}

// --- Updated VolumeSource including all types ---
type VolumeSource struct {
	HostPath              *HostPathVolumeSource
	EmptyDir              *EmptyDirVolumeSource
	GCEPersistentDisk     *GCEPersistentDiskVolumeSource
	AWSElasticBlockStore  *AWSElasticBlockStoreVolumeSource
	GitRepo               *GitRepoVolumeSource
	Secret                *SecretVolumeSource
	NFS                   *NFSVolumeSource
	ISCSI                 *ISCSIVolumeSource
	Glusterfs             *GlusterfsVolumeSource
	PersistentVolumeClaim *PersistentVolumeClaimVolumeSource
	RBD                   *RBDVolumeSource
	FlexVolume            *FlexVolumeSource
	Flocker               *FlockerVolumeSource
	DownwardAPI           *DownwardAPIVolumeSource
	FC                    *FCVolumeSource
	AzureFile             *AzureFileVolumeSource
	ConfigMap             *ConfigMapVolumeSource
	VsphereVolume         *VsphereVirtualDiskVolumeSource
	Quobyte               *QuobyteVolumeSource
	AzureDisk             *AzureDiskVolumeSource
	PhotonPersistentDisk  *PhotonPersistentDiskVolumeSource
	Projected             *ProjectedVolumeSource
	PortworxVolume        *PortworxVolumeSource
	ScaleIO               *ScaleIOVolumeSource
	StorageOS             *StorageOSVolumeSource
	CSI                   *CSIVolumeSource
	Ephemeral             *EphemeralVolumeSource
	Image                 *ImageVolumeSource
}

// --- EmptyDir ---
type EmptyDirVolumeSource struct {
	Medium    string // corev1.StorageMedium
	SizeLimit *resource.Quantity
}

func ConvertStorageMedium(sm corev1.StorageMedium) string {
	return string(sm)
}

func ConvertEmptyDirVolumeSource(s *corev1.EmptyDirVolumeSource) *EmptyDirVolumeSource {
	if s == nil {
		return nil
	}
	return &EmptyDirVolumeSource{
		Medium:    ConvertStorageMedium(s.Medium),
		SizeLimit: s.SizeLimit,
	}
}

// --- Secret ---
type SecretVolumeSource struct {
	SecretName  string
	Items       []KeyToPath
	DefaultMode *int32
	Optional    *bool
}

func ConvertSecretVolumeSource(s *corev1.SecretVolumeSource) *SecretVolumeSource {
	if s == nil {
		return nil
	}
	return &SecretVolumeSource{
		SecretName:  s.SecretName,
		Items:       ConvertKeyToPaths(s.Items),
		DefaultMode: s.DefaultMode,
		Optional:    s.Optional,
	}
}

// --- KeyToPath ---
type KeyToPath struct {
	Key  string
	Path string
	Mode *int32
}

func ConvertKeyToPath(item corev1.KeyToPath) KeyToPath {
	return KeyToPath{
		Key:  item.Key,
		Path: item.Path,
		Mode: item.Mode,
	}
}

func ConvertKeyToPaths(items []corev1.KeyToPath) []KeyToPath {
	if items == nil {
		return nil
	}
	result := make([]KeyToPath, len(items))
	for i, item := range items {
		result[i] = ConvertKeyToPath(item)
	}
	return result
}

// --- NFS ---
type NFSVolumeSource struct {
	Server   string
	Path     string
	ReadOnly bool
}

func ConvertNFSVolumeSource(s *corev1.NFSVolumeSource) *NFSVolumeSource {
	if s == nil {
		return nil
	}
	return &NFSVolumeSource{
		Server:   s.Server,
		Path:     s.Path,
		ReadOnly: s.ReadOnly,
	}
}

// --- PersistentVolumeClaim ---
type PersistentVolumeClaimVolumeSource struct {
	ClaimName string
	ReadOnly  bool
}

func ConvertPersistentVolumeClaimVolumeSource(s *corev1.PersistentVolumeClaimVolumeSource) *PersistentVolumeClaimVolumeSource {
	if s == nil {
		return nil
	}
	return &PersistentVolumeClaimVolumeSource{
		ClaimName: s.ClaimName,
		ReadOnly:  s.ReadOnly,
	}
}

// --- DownwardAPI ---
type DownwardAPIVolumeSource struct {
	Items       []DownwardAPIVolumeFile
	DefaultMode *int32
}

func ConvertDownwardAPIVolumeSource(s *corev1.DownwardAPIVolumeSource) *DownwardAPIVolumeSource {
	if s == nil {
		return nil
	}
	return &DownwardAPIVolumeSource{
		Items:       ConvertDownwardAPIVolumeFiles(s.Items),
		DefaultMode: s.DefaultMode,
	}
}

type DownwardAPIVolumeFile struct {
	Path             string
	FieldRef         *ObjectFieldSelector
	ResourceFieldRef *ResourceFieldSelector
	Mode             *int32
}

func ConvertDownwardAPIVolumeFile(item corev1.DownwardAPIVolumeFile) DownwardAPIVolumeFile {
	return DownwardAPIVolumeFile{
		Path:             item.Path,
		FieldRef:         ConvertObjectFieldSelector(item.FieldRef),
		ResourceFieldRef: ConvertResourceFieldSelector(item.ResourceFieldRef),
		Mode:             item.Mode,
	}
}

func ConvertDownwardAPIVolumeFiles(items []corev1.DownwardAPIVolumeFile) []DownwardAPIVolumeFile {
	if items == nil {
		return nil
	}
	result := make([]DownwardAPIVolumeFile, len(items))
	for i, item := range items {
		result[i] = ConvertDownwardAPIVolumeFile(item)
	}
	return result
}

// --- ConfigMap ---
type ConfigMapVolumeSource struct {
	LocalObjectReference // Name is embedded
	Items                []KeyToPath
	DefaultMode          *int32
	Optional             *bool
}

// Assuming LocalObjectReference is defined elsewhere or needs to be added
type LocalObjectReference struct {
	Name string
}

func ConvertLocalObjectReference(lor corev1.LocalObjectReference) LocalObjectReference {
	return LocalObjectReference{Name: lor.Name}
}

func ConvertConfigMapVolumeSource(s *corev1.ConfigMapVolumeSource) *ConfigMapVolumeSource {
	if s == nil {
		return nil
	}
	return &ConfigMapVolumeSource{
		LocalObjectReference: LocalObjectReference{Name: s.Name}, // Explicitly convert embedded Name
		Items:                ConvertKeyToPaths(s.Items),
		DefaultMode:          s.DefaultMode,
		Optional:             s.Optional,
	}
}

// --- Projected ---
type ProjectedVolumeSource struct {
	Sources     []VolumeProjection
	DefaultMode *int32
}

func ConvertProjectedVolumeSource(s *corev1.ProjectedVolumeSource) *ProjectedVolumeSource {
	if s == nil {
		return nil
	}
	return &ProjectedVolumeSource{
		Sources:     ConvertVolumeProjections(s.Sources),
		DefaultMode: s.DefaultMode,
	}
}

type VolumeProjection struct {
	Secret              *SecretProjection
	DownwardAPI         *DownwardAPIProjection
	ConfigMap           *ConfigMapProjection
	ServiceAccountToken *ServiceAccountTokenProjection
	ClusterTrustBundle  *ClusterTrustBundleProjection
}

// Corrected ClusterTrustBundleProjection Definition
type ClusterTrustBundleProjection struct {
	Name       *string // Corrected field
	SignerName *string
	Path       string
	Optional   *bool
}

// Removed ConvertLabelSelector (was incorrect context)

// Corrected ClusterTrustBundleProjection Conversion
func ConvertClusterTrustBundleProjection(p *corev1.ClusterTrustBundleProjection) *ClusterTrustBundleProjection {
	if p == nil {
		return nil
	}
	return &ClusterTrustBundleProjection{
		Name:       p.Name,
		SignerName: p.SignerName,
		Path:       p.Path,
		Optional:   p.Optional,
	}
}

func ConvertVolumeProjection(p corev1.VolumeProjection) VolumeProjection {
	return VolumeProjection{
		Secret:              ConvertSecretProjection(p.Secret),
		DownwardAPI:         ConvertDownwardAPIProjection(p.DownwardAPI),
		ConfigMap:           ConvertConfigMapProjection(p.ConfigMap),
		ServiceAccountToken: ConvertServiceAccountTokenProjection(p.ServiceAccountToken),
		ClusterTrustBundle:  ConvertClusterTrustBundleProjection(p.ClusterTrustBundle),
	}
}

func ConvertVolumeProjections(sources []corev1.VolumeProjection) []VolumeProjection {
	if sources == nil {
		return nil
	}
	result := make([]VolumeProjection, len(sources))
	for i, p := range sources {
		result[i] = ConvertVolumeProjection(p)
	}
	return result
}

// --- SecretProjection ---
type SecretProjection struct {
	LocalObjectReference // Name embedded
	Items                []KeyToPath
	Optional             *bool
}

func ConvertSecretProjection(p *corev1.SecretProjection) *SecretProjection {
	if p == nil {
		return nil
	}
	return &SecretProjection{
		LocalObjectReference: LocalObjectReference{Name: p.Name}, // Explicitly convert embedded Name
		Items:                ConvertKeyToPaths(p.Items),
		Optional:             p.Optional,
	}
}

// --- DownwardAPIProjection ---
type DownwardAPIProjection struct {
	Items []DownwardAPIVolumeFile
}

func ConvertDownwardAPIProjection(p *corev1.DownwardAPIProjection) *DownwardAPIProjection {
	if p == nil {
		return nil
	}
	return &DownwardAPIProjection{
		Items: ConvertDownwardAPIVolumeFiles(p.Items),
	}
}

// --- ConfigMapProjection ---
type ConfigMapProjection struct {
	LocalObjectReference // Name embedded
	Items                []KeyToPath
	Optional             *bool
}

func ConvertConfigMapProjection(p *corev1.ConfigMapProjection) *ConfigMapProjection {
	if p == nil {
		return nil
	}
	return &ConfigMapProjection{
		LocalObjectReference: LocalObjectReference{Name: p.Name}, // Explicitly convert embedded Name
		Items:                ConvertKeyToPaths(p.Items),
		Optional:             p.Optional,
	}
}

// --- ServiceAccountTokenProjection ---
type ServiceAccountTokenProjection struct {
	Audience          string
	ExpirationSeconds *int64
	Path              string
}

func ConvertServiceAccountTokenProjection(p *corev1.ServiceAccountTokenProjection) *ServiceAccountTokenProjection {
	if p == nil {
		return nil
	}
	return &ServiceAccountTokenProjection{
		Audience:          p.Audience,
		ExpirationSeconds: p.ExpirationSeconds,
		Path:              p.Path,
	}
}

// --- CSI ---
type CSIVolumeSource struct {
	Driver               string
	ReadOnly             *bool
	FSType               *string
	VolumeAttributes     map[string]string
	NodePublishSecretRef *LocalObjectReference
}

func ConvertCSIVolumeSource(s *corev1.CSIVolumeSource) *CSIVolumeSource {
	if s == nil {
		return nil
	}
	var secretRef *LocalObjectReference
	if s.NodePublishSecretRef != nil {
		ref := ConvertLocalObjectReference(*s.NodePublishSecretRef)
		secretRef = &ref
	}
	return &CSIVolumeSource{
		Driver:               s.Driver,
		ReadOnly:             s.ReadOnly,
		FSType:               s.FSType,
		VolumeAttributes:     s.VolumeAttributes,
		NodePublishSecretRef: secretRef,
	}
}

// --- ObjectFieldSelector ---
type ObjectFieldSelector struct {
	APIVersion string
	FieldPath  string
}

func ConvertObjectFieldSelector(s *corev1.ObjectFieldSelector) *ObjectFieldSelector {
	if s == nil {
		return nil
	}
	return &ObjectFieldSelector{
		APIVersion: s.APIVersion,
		FieldPath:  s.FieldPath,
	}
}

// --- ResourceFieldSelector ---
type ResourceFieldSelector struct {
	ContainerName string
	Resource      string
	Divisor       resource.Quantity
}

func ConvertResourceFieldSelector(s *corev1.ResourceFieldSelector) *ResourceFieldSelector {
	if s == nil {
		return nil
	}
	return &ResourceFieldSelector{
		ContainerName: s.ContainerName,
		Resource:      s.Resource,
		Divisor:       s.Divisor,
	}
}

// --- Main Volume Conversion (Updated) ---
func ConvertVolume(vol corev1.Volume) Volume {
	return Volume{
		Name: vol.Name,
		VolumeSource: VolumeSource{
			HostPath:              ConvertHostPathVolumeSource(vol.HostPath),
			EmptyDir:              ConvertEmptyDirVolumeSource(vol.EmptyDir),
			GCEPersistentDisk:     ConvertGCEPersistentDiskVolumeSource(vol.GCEPersistentDisk),
			AWSElasticBlockStore:  ConvertAWSElasticBlockStoreVolumeSource(vol.AWSElasticBlockStore),
			GitRepo:               ConvertGitRepoVolumeSource(vol.GitRepo),
			Secret:                ConvertSecretVolumeSource(vol.Secret),
			NFS:                   ConvertNFSVolumeSource(vol.NFS),
			ISCSI:                 ConvertISCSIVolumeSource(vol.ISCSI),
			Glusterfs:             ConvertGlusterfsVolumeSource(vol.Glusterfs),
			PersistentVolumeClaim: ConvertPersistentVolumeClaimVolumeSource(vol.PersistentVolumeClaim),
			RBD:                   ConvertRBDVolumeSource(vol.RBD),
			FlexVolume:            ConvertFlexVolumeSource(vol.FlexVolume),
			Flocker:               ConvertFlockerVolumeSource(vol.Flocker),
			DownwardAPI:           ConvertDownwardAPIVolumeSource(vol.DownwardAPI),
			FC:                    ConvertFCVolumeSource(vol.FC),
			AzureFile:             ConvertAzureFileVolumeSource(vol.AzureFile),
			ConfigMap:             ConvertConfigMapVolumeSource(vol.ConfigMap),
			VsphereVolume:         ConvertVsphereVirtualDiskVolumeSource(vol.VsphereVolume),
			Quobyte:               ConvertQuobyteVolumeSource(vol.Quobyte),
			AzureDisk:             ConvertAzureDiskVolumeSource(vol.AzureDisk),
			PhotonPersistentDisk:  ConvertPhotonPersistentDiskVolumeSource(vol.PhotonPersistentDisk),
			Projected:             ConvertProjectedVolumeSource(vol.Projected),
			PortworxVolume:        ConvertPortworxVolumeSource(vol.PortworxVolume),
			ScaleIO:               ConvertScaleIOVolumeSource(vol.ScaleIO),
			StorageOS:             ConvertStorageOSVolumeSource(vol.StorageOS),
			CSI:                   ConvertCSIVolumeSource(vol.CSI),
			Ephemeral:             ConvertEphemeralVolumeSource(vol.Ephemeral),
			Image:                 ConvertImageVolumeSource(vol.Image),
		},
	}
}

func ConvertVolumes(volumes []corev1.Volume) []Volume {
	if volumes == nil {
		return nil
	}
	result := make([]Volume, len(volumes))
	for i, vol := range volumes {
		result[i] = ConvertVolume(vol)
	}
	return result
}

// -- Removed placeholder pod types --
