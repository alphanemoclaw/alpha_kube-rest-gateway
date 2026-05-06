package models

import "time"

type PodResponse struct {
	Namespace string `json:"namespace"`
	Count     int    `json:"count"`
	Pods      []Pod  `json:"pods"`
}

type Pod struct {
	Name              string            `json:"name"`
	Namespace         string            `json:"namespace"`
	Phase             string            `json:"phase"`
	PodIP             string            `json:"pod_ip"`
	NodeName          string            `json:"node_name"`
	CreationTimestamp *time.Time        `json:"creation_timestamp"`
	Labels            map[string]string `json:"labels"`
	Containers        []Container       `json:"containers"`
	Conditions        []Condition       `json:"conditions"`
}

type Container struct {
	Name         string          `json:"name"`
	Image        string          `json:"image"`
	Ports        []ContainerPort `json:"ports"`
	Ready        bool            `json:"ready"`
	RestartCount int32           `json:"restart_count"`
	State        ContainerState  `json:"state,omitempty"`
	LastState    ContainerState  `json:"last_state,omitempty"`
}

type ContainerPort struct {
	ContainerPort int32  `json:"container_port"`
	Protocol      string `json:"protocol"`
}

type ContainerState struct {
	State      string     `json:"state,omitempty"`
	Reason     string     `json:"reason,omitempty"`
	Message    string     `json:"message,omitempty"`
	ExitCode   int32      `json:"exit_code,omitempty"`
	Signal     int32      `json:"signal,omitempty"`
	StartedAt  *time.Time `json:"started_at,omitempty"`
	FinishedAt *time.Time `json:"finished_at,omitempty"`
}

type Condition struct {
	Type   string `json:"type"`
	Status string `json:"status"`
}

type PodStatusResponse struct {
	Name               string      `json:"name"`
	Namespace          string      `json:"namespace"`
	Phase              string      `json:"phase"`
	Reason             string      `json:"reason,omitempty"`
	Message            string      `json:"message,omitempty"`
	HostIP             string      `json:"host_ip,omitempty"`
	PodIP              string      `json:"pod_ip,omitempty"`
	NodeName           string      `json:"node_name,omitempty"`
	StartTime          *time.Time  `json:"start_time,omitempty"`
	QOSClass           string      `json:"qos_class,omitempty"`
	RestartPolicy      string      `json:"restart_policy,omitempty"`
	ServiceAccountName string      `json:"service_account_name,omitempty"`
	InitContainers     []Container `json:"init_containers"`
	Containers         []Container `json:"containers"`
	Conditions         []Condition `json:"conditions"`
}

type ServiceResponse struct {
	Namespace string    `json:"namespace"`
	Count     int       `json:"count"`
	Services  []Service `json:"services"`
}

type Service struct {
	Name              string            `json:"name"`
	Namespace         string            `json:"namespace"`
	Type              string            `json:"type"`
	ClusterIP         string            `json:"cluster_ip"`
	ExternalIPs       []string          `json:"external_ips"`
	LoadBalancerIP    string            `json:"load_balancer_ip"`
	Ports             []ServicePort     `json:"ports"`
	Selector          map[string]string `json:"selector"`
	CreationTimestamp *time.Time        `json:"creation_timestamp"`
}

type ServicePort struct {
	Name       string `json:"name"`
	Protocol   string `json:"protocol"`
	Port       int32  `json:"port"`
	TargetPort string `json:"target_port"`
	NodePort   int32  `json:"node_port"`
}

type EndpointResponse struct {
	Namespace string     `json:"namespace"`
	Count     int        `json:"count"`
	Endpoints []Endpoint `json:"endpoints"`
}

type Endpoint struct {
	Name              string            `json:"name"`
	Namespace         string            `json:"namespace"`
	CreationTimestamp *time.Time        `json:"creation_timestamp"`
	Labels            map[string]string `json:"labels"`
	Subsets           []EndpointSubset  `json:"subsets"`
}

type EndpointSubset struct {
	Addresses         []EndpointAddress `json:"addresses"`
	NotReadyAddresses []EndpointAddress `json:"not_ready_addresses"`
	Ports             []EndpointPort    `json:"ports"`
}

type EndpointAddress struct {
	IP         string `json:"ip"`
	Hostname   string `json:"hostname,omitempty"`
	NodeName   string `json:"node_name,omitempty"`
	TargetKind string `json:"target_kind,omitempty"`
	TargetName string `json:"target_name,omitempty"`
}

type EndpointPort struct {
	Name     string `json:"name"`
	Port     int32  `json:"port"`
	Protocol string `json:"protocol"`
}

type EndpointSliceResponse struct {
	Namespace      string          `json:"namespace"`
	Count          int             `json:"count"`
	EndpointSlices []EndpointSlice `json:"endpointslices"`
}

type EndpointSlice struct {
	Name              string                 `json:"name"`
	Namespace         string                 `json:"namespace"`
	CreationTimestamp *time.Time             `json:"creation_timestamp"`
	Labels            map[string]string      `json:"labels"`
	AddressType       string                 `json:"address_type"`
	Ports             []EndpointSlicePort    `json:"ports"`
	Endpoints         []EndpointSliceBackend `json:"endpoints"`
}

type EndpointSlicePort struct {
	Name        string `json:"name,omitempty"`
	Port        *int32 `json:"port,omitempty"`
	Protocol    string `json:"protocol,omitempty"`
	AppProtocol string `json:"app_protocol,omitempty"`
}

type EndpointSliceBackend struct {
	Addresses  []string         `json:"addresses"`
	Conditions map[string]*bool `json:"conditions"`
	Hostname   string           `json:"hostname,omitempty"`
	NodeName   string           `json:"node_name,omitempty"`
	TargetKind string           `json:"target_kind,omitempty"`
	TargetName string           `json:"target_name,omitempty"`
	Zone       string           `json:"zone,omitempty"`
}

type IngressResponse struct {
	Namespace string    `json:"namespace"`
	Count     int       `json:"count"`
	Ingresses []Ingress `json:"ingresses"`
}

type Ingress struct {
	Name              string            `json:"name"`
	Namespace         string            `json:"namespace"`
	CreationTimestamp *time.Time        `json:"creation_timestamp"`
	Labels            map[string]string `json:"labels"`
	ClassName         string            `json:"class_name,omitempty"`
	TLSHosts          []string          `json:"tls_hosts"`
	Rules             []IngressRule     `json:"rules"`
	LoadBalancerIPs   []string          `json:"load_balancer_ips"`
}

type IngressRule struct {
	Host        string `json:"host"`
	Path        string `json:"path"`
	PathType    string `json:"path_type"`
	ServiceName string `json:"service_name"`
	ServicePort string `json:"service_port"`
}

type DeploymentResponse struct {
	Namespace   string       `json:"namespace"`
	Count       int          `json:"count"`
	Deployments []Deployment `json:"deployments"`
}

type Deployment struct {
	Name              string            `json:"name"`
	Namespace         string            `json:"namespace"`
	CreationTimestamp *time.Time        `json:"creation_timestamp"`
	Labels            map[string]string `json:"labels"`
	Replicas          Replicas          `json:"replicas"`
	Strategy          string            `json:"strategy"`
	Conditions        []DepCondition    `json:"conditions"`
}

type Replicas struct {
	Desired   int32 `json:"desired"`
	Ready     int32 `json:"ready"`
	Available int32 `json:"available"`
	Updated   int32 `json:"updated"`
}

type DepCondition struct {
	Type    string `json:"type"`
	Status  string `json:"status"`
	Reason  string `json:"reason"`
	Message string `json:"message"`
}

type ReplicaSetResponse struct {
	Namespace   string       `json:"namespace"`
	Count       int          `json:"count"`
	ReplicaSets []ReplicaSet `json:"replicasets"`
}

type ReplicaSet struct {
	Name              string            `json:"name"`
	Namespace         string            `json:"namespace"`
	CreationTimestamp *time.Time        `json:"creation_timestamp"`
	Labels            map[string]string `json:"labels"`
	Replicas          Replicas          `json:"replicas"`
	OwnerReferences   []OwnerReference  `json:"owner_references"`
}

type StatefulSetResponse struct {
	Namespace    string        `json:"namespace"`
	Count        int           `json:"count"`
	StatefulSets []StatefulSet `json:"statefulsets"`
}

type StatefulSet struct {
	Name              string            `json:"name"`
	Namespace         string            `json:"namespace"`
	CreationTimestamp *time.Time        `json:"creation_timestamp"`
	Labels            map[string]string `json:"labels"`
	Replicas          Replicas          `json:"replicas"`
	ServiceName       string            `json:"service_name"`
}

type DaemonSetResponse struct {
	Namespace  string      `json:"namespace"`
	Count      int         `json:"count"`
	DaemonSets []DaemonSet `json:"daemonsets"`
}

type DaemonSet struct {
	Name                   string            `json:"name"`
	Namespace              string            `json:"namespace"`
	CreationTimestamp      *time.Time        `json:"creation_timestamp"`
	Labels                 map[string]string `json:"labels"`
	DesiredNumberScheduled int32             `json:"desired_number_scheduled"`
	CurrentNumberScheduled int32             `json:"current_number_scheduled"`
	NumberReady            int32             `json:"number_ready"`
	UpdatedNumberScheduled int32             `json:"updated_number_scheduled"`
	NumberAvailable        int32             `json:"number_available"`
	NumberUnavailable      int32             `json:"number_unavailable"`
}

type OwnerReference struct {
	Kind string `json:"kind"`
	Name string `json:"name"`
}

type JobResponse struct {
	Namespace string `json:"namespace"`
	Count     int    `json:"count"`
	Jobs      []Job  `json:"jobs"`
}

type Job struct {
	Name              string            `json:"name"`
	Namespace         string            `json:"namespace"`
	CreationTimestamp *time.Time        `json:"creation_timestamp"`
	Labels            map[string]string `json:"labels"`
	Completions       *int32            `json:"completions,omitempty"`
	Parallelism       *int32            `json:"parallelism,omitempty"`
	Active            int32             `json:"active"`
	Succeeded         int32             `json:"succeeded"`
	Failed            int32             `json:"failed"`
	Conditions        []DepCondition    `json:"conditions"`
}

type CronJobResponse struct {
	Namespace string    `json:"namespace"`
	Count     int       `json:"count"`
	CronJobs  []CronJob `json:"cronjobs"`
}

type CronJob struct {
	Name               string            `json:"name"`
	Namespace          string            `json:"namespace"`
	CreationTimestamp  *time.Time        `json:"creation_timestamp"`
	Labels             map[string]string `json:"labels"`
	Schedule           string            `json:"schedule"`
	Suspend            *bool             `json:"suspend,omitempty"`
	Active             []string          `json:"active"`
	LastScheduleTime   *time.Time        `json:"last_schedule_time,omitempty"`
	LastSuccessfulTime *time.Time        `json:"last_successful_time,omitempty"`
}

type NodeResponse struct {
	Count int    `json:"count"`
	Nodes []Node `json:"nodes"`
}

type Node struct {
	Name              string            `json:"name"`
	CreationTimestamp *time.Time        `json:"creation_timestamp"`
	Labels            map[string]string `json:"labels"`
	Addresses         map[string]string `json:"addresses"`
	Capacity          map[string]string `json:"capacity"`
	OSImage           string            `json:"os_image"`
	KernelVersion     string            `json:"kernel_version"`
	KubeletVersion    string            `json:"kubelet_version"`
	ContainerRuntime  string            `json:"container_runtime"`
	Conditions        []Condition       `json:"conditions"`
}

type NamespaceResponse struct {
	Count      int         `json:"count"`
	Namespaces []Namespace `json:"namespaces"`
}

type Namespace struct {
	Name              string            `json:"name"`
	Phase             string            `json:"phase"`
	CreationTimestamp *time.Time        `json:"creation_timestamp"`
	Labels            map[string]string `json:"labels"`
}

type LogResponse struct {
	PodName            string   `json:"pod_name"`
	Namespace          string   `json:"namespace"`
	Container          string   `json:"container"`
	TailLinesRequested int64    `json:"tail_lines_requested"`
	LinesReturned      int      `json:"lines_returned"`
	Logs               []string `json:"logs"`
}

type EventResponse struct {
	Namespace string  `json:"namespace"`
	Count     int     `json:"count"`
	Events    []Event `json:"events"`
}

type Event struct {
	Name           string         `json:"name"`
	Namespace      string         `json:"namespace"`
	Type           string         `json:"type"`
	Reason         string         `json:"reason"`
	Message        string         `json:"message"`
	EventCount     int32          `json:"count"`
	FirstTimestamp *time.Time     `json:"first_timestamp"`
	LastTimestamp  *time.Time     `json:"last_timestamp"`
	InvolvedObject InvolvedObject `json:"involved_object"`
	Source         Source         `json:"source"`
}

type InvolvedObject struct {
	Kind      string `json:"kind"`
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

type Source struct {
	Component string `json:"component"`
	Host      string `json:"host"`
}

type PersistentVolumeClaimResponse struct {
	Namespace string                  `json:"namespace"`
	Count     int                     `json:"count"`
	PVCs      []PersistentVolumeClaim `json:"pvcs"`
}

type PersistentVolumeClaim struct {
	Name              string            `json:"name"`
	Namespace         string            `json:"namespace"`
	CreationTimestamp *time.Time        `json:"creation_timestamp"`
	Labels            map[string]string `json:"labels"`
	Phase             string            `json:"phase"`
	StorageClassName  string            `json:"storage_class_name,omitempty"`
	VolumeName        string            `json:"volume_name,omitempty"`
	AccessModes       []string          `json:"access_modes"`
	Capacity          map[string]string `json:"capacity"`
	Requested         map[string]string `json:"requested"`
}

type PersistentVolumeResponse struct {
	Count int                `json:"count"`
	PVs   []PersistentVolume `json:"pvs"`
}

type PersistentVolume struct {
	Name              string            `json:"name"`
	CreationTimestamp *time.Time        `json:"creation_timestamp"`
	Labels            map[string]string `json:"labels"`
	Phase             string            `json:"phase"`
	StorageClassName  string            `json:"storage_class_name,omitempty"`
	Capacity          map[string]string `json:"capacity"`
	AccessModes       []string          `json:"access_modes"`
	ReclaimPolicy     string            `json:"reclaim_policy"`
	ClaimNamespace    string            `json:"claim_namespace,omitempty"`
	ClaimName         string            `json:"claim_name,omitempty"`
	VolumeMode        string            `json:"volume_mode,omitempty"`
	Reason            string            `json:"reason,omitempty"`
}

type StorageClassResponse struct {
	Count          int            `json:"count"`
	StorageClasses []StorageClass `json:"storageclasses"`
}

type StorageClass struct {
	Name                 string            `json:"name"`
	CreationTimestamp    *time.Time        `json:"creation_timestamp"`
	Labels               map[string]string `json:"labels"`
	Provisioner          string            `json:"provisioner"`
	ReclaimPolicy        string            `json:"reclaim_policy,omitempty"`
	VolumeBindingMode    string            `json:"volume_binding_mode,omitempty"`
	AllowVolumeExpansion *bool             `json:"allow_volume_expansion,omitempty"`
	Parameters           map[string]string `json:"parameters"`
}

type NetworkPolicyResponse struct {
	Namespace       string          `json:"namespace"`
	Count           int             `json:"count"`
	NetworkPolicies []NetworkPolicy `json:"networkpolicies"`
}

type NetworkPolicy struct {
	Name              string            `json:"name"`
	Namespace         string            `json:"namespace"`
	CreationTimestamp *time.Time        `json:"creation_timestamp"`
	Labels            map[string]string `json:"labels"`
	PodSelector       map[string]string `json:"pod_selector"`
	PolicyTypes       []string          `json:"policy_types"`
	IngressRuleCount  int               `json:"ingress_rule_count"`
	EgressRuleCount   int               `json:"egress_rule_count"`
}

type ResourceQuotaResponse struct {
	Namespace      string          `json:"namespace"`
	Count          int             `json:"count"`
	ResourceQuotas []ResourceQuota `json:"resourcequotas"`
}

type ResourceQuota struct {
	Name              string            `json:"name"`
	Namespace         string            `json:"namespace"`
	CreationTimestamp *time.Time        `json:"creation_timestamp"`
	Labels            map[string]string `json:"labels"`
	Hard              map[string]string `json:"hard"`
	Used              map[string]string `json:"used"`
}

type LimitRangeResponse struct {
	Namespace   string       `json:"namespace"`
	Count       int          `json:"count"`
	LimitRanges []LimitRange `json:"limitranges"`
}

type LimitRange struct {
	Name              string            `json:"name"`
	Namespace         string            `json:"namespace"`
	CreationTimestamp *time.Time        `json:"creation_timestamp"`
	Labels            map[string]string `json:"labels"`
	Items             []LimitRangeItem  `json:"items"`
}

type LimitRangeItem struct {
	Type                 string            `json:"type"`
	Max                  map[string]string `json:"max"`
	Min                  map[string]string `json:"min"`
	Default              map[string]string `json:"default"`
	DefaultRequest       map[string]string `json:"default_request"`
	MaxLimitRequestRatio map[string]string `json:"max_limit_request_ratio"`
}

type ConfigMapResponse struct {
	Namespace  string              `json:"namespace"`
	Count      int                 `json:"count"`
	ConfigMaps []ConfigMapMetadata `json:"configmaps"`
}

type ConfigMapMetadata struct {
	Name              string            `json:"name"`
	Namespace         string            `json:"namespace"`
	CreationTimestamp *time.Time        `json:"creation_timestamp"`
	Labels            map[string]string `json:"labels"`
	DataKeys          []string          `json:"data_keys"`
	BinaryDataKeys    []string          `json:"binary_data_keys"`
	Immutable         *bool             `json:"immutable,omitempty"`
}

type PodMetricsResponse struct {
	Namespace string       `json:"namespace"`
	Count     int          `json:"count"`
	Pods      []PodMetrics `json:"pods"`
}

type PodMetrics struct {
	Name       string                     `json:"name"`
	Namespace  string                     `json:"namespace"`
	Timestamp  *time.Time                 `json:"timestamp"`
	Window     string                     `json:"window"`
	Containers []ContainerResourceMetrics `json:"containers"`
}

type ContainerResourceMetrics struct {
	Name  string            `json:"name"`
	Usage map[string]string `json:"usage"`
}

type NodeMetricsResponse struct {
	Count int           `json:"count"`
	Nodes []NodeMetrics `json:"nodes"`
}

type NodeMetrics struct {
	Name      string            `json:"name"`
	Timestamp *time.Time        `json:"timestamp"`
	Window    string            `json:"window"`
	Usage     map[string]string `json:"usage"`
}
