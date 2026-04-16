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
}

type ContainerPort struct {
	ContainerPort int32  `json:"container_port"`
	Protocol      string `json:"protocol"`
}

type Condition struct {
	Type   string `json:"type"`
	Status string `json:"status"`
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
	PodName             string   `json:"pod_name"`
	Namespace           string   `json:"namespace"`
	Container           string   `json:"container"`
	TailLinesRequested int64    `json:"tail_lines_requested"`
	LinesReturned       int      `json:"lines_returned"`
	Logs                []string `json:"logs"`
}

type EventResponse struct {
	Namespace string  `json:"namespace"`
	Count     int     `json:"count"`
	Events    []Event `json:"events"`
}

type Event struct {
	Name             string         `json:"name"`
	Namespace        string         `json:"namespace"`
	Type             string         `json:"type"`
	Reason           string         `json:"reason"`
	Message          string         `json:"message"`
	EventCount       int32          `json:"count"`
	FirstTimestamp   *time.Time     `json:"first_timestamp"`
	LastTimestamp    *time.Time     `json:"last_timestamp"`
	InvolvedObject   InvolvedObject `json:"involved_object"`
	Source           Source         `json:"source"`
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
