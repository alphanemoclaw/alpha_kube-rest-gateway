package handlers

import (
	"context"
	"net/http"
	"sort"
	"time"

	"github.com/alpha-kube-rest-gateway/config"
	"github.com/alpha-kube-rest-gateway/k8s"
	"github.com/alpha-kube-rest-gateway/models"
	"github.com/gin-gonic/gin"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	discoveryv1 "k8s.io/api/discovery/v1"
	networkingv1 "k8s.io/api/networking/v1"
	storagev1 "k8s.io/api/storage/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	metricsv1beta1 "k8s.io/metrics/pkg/apis/metrics/v1beta1"
)

func GetPodStatus(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		podName := c.Param("pod_name")
		namespace := c.DefaultQuery("namespace", cfg.DefaultNamespace)

		pod, err := k8s.Clientset.CoreV1().Pods(namespace).Get(context.TODO(), podName, metav1.GetOptions{})
		if err != nil {
			handleK8sError(c, err, "pod-status/"+namespace+"/"+podName)
			return
		}

		c.JSON(http.StatusOK, mapPodStatus(*pod))
	}
}

func ListReplicaSets(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		namespace := c.DefaultQuery("namespace", cfg.DefaultNamespace)
		labelSelector := c.Query("label_selector")

		replicaSets, err := k8s.Clientset.AppsV1().ReplicaSets(namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: labelSelector})
		if err != nil {
			handleK8sError(c, err, "replicasets/"+namespace)
			return
		}

		resp := models.ReplicaSetResponse{Namespace: namespace, Count: len(replicaSets.Items), ReplicaSets: make([]models.ReplicaSet, 0, len(replicaSets.Items))}
		for _, replicaSet := range replicaSets.Items {
			resp.ReplicaSets = append(resp.ReplicaSets, mapReplicaSet(replicaSet))
		}

		c.JSON(http.StatusOK, resp)
	}
}

func ListStatefulSets(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		namespace := c.DefaultQuery("namespace", cfg.DefaultNamespace)
		labelSelector := c.Query("label_selector")

		statefulSets, err := k8s.Clientset.AppsV1().StatefulSets(namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: labelSelector})
		if err != nil {
			handleK8sError(c, err, "statefulsets/"+namespace)
			return
		}

		resp := models.StatefulSetResponse{Namespace: namespace, Count: len(statefulSets.Items), StatefulSets: make([]models.StatefulSet, 0, len(statefulSets.Items))}
		for _, statefulSet := range statefulSets.Items {
			resp.StatefulSets = append(resp.StatefulSets, mapStatefulSet(statefulSet))
		}

		c.JSON(http.StatusOK, resp)
	}
}

func ListDaemonSets(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		namespace := c.DefaultQuery("namespace", cfg.DefaultNamespace)
		labelSelector := c.Query("label_selector")

		daemonSets, err := k8s.Clientset.AppsV1().DaemonSets(namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: labelSelector})
		if err != nil {
			handleK8sError(c, err, "daemonsets/"+namespace)
			return
		}

		resp := models.DaemonSetResponse{Namespace: namespace, Count: len(daemonSets.Items), DaemonSets: make([]models.DaemonSet, 0, len(daemonSets.Items))}
		for _, daemonSet := range daemonSets.Items {
			resp.DaemonSets = append(resp.DaemonSets, mapDaemonSet(daemonSet))
		}

		c.JSON(http.StatusOK, resp)
	}
}

func ListJobs(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		namespace := c.DefaultQuery("namespace", cfg.DefaultNamespace)
		labelSelector := c.Query("label_selector")

		jobs, err := k8s.Clientset.BatchV1().Jobs(namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: labelSelector})
		if err != nil {
			handleK8sError(c, err, "jobs/"+namespace)
			return
		}

		resp := models.JobResponse{Namespace: namespace, Count: len(jobs.Items), Jobs: make([]models.Job, 0, len(jobs.Items))}
		for _, job := range jobs.Items {
			resp.Jobs = append(resp.Jobs, mapJob(job))
		}

		c.JSON(http.StatusOK, resp)
	}
}

func ListCronJobs(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		namespace := c.DefaultQuery("namespace", cfg.DefaultNamespace)
		labelSelector := c.Query("label_selector")

		cronJobs, err := k8s.Clientset.BatchV1().CronJobs(namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: labelSelector})
		if err != nil {
			handleK8sError(c, err, "cronjobs/"+namespace)
			return
		}

		resp := models.CronJobResponse{Namespace: namespace, Count: len(cronJobs.Items), CronJobs: make([]models.CronJob, 0, len(cronJobs.Items))}
		for _, cronJob := range cronJobs.Items {
			resp.CronJobs = append(resp.CronJobs, mapCronJob(cronJob))
		}

		c.JSON(http.StatusOK, resp)
	}
}

func ListEndpoints(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		namespace := c.DefaultQuery("namespace", cfg.DefaultNamespace)
		labelSelector := c.Query("label_selector")

		endpoints, err := k8s.Clientset.CoreV1().Endpoints(namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: labelSelector})
		if err != nil {
			handleK8sError(c, err, "endpoints/"+namespace)
			return
		}

		resp := models.EndpointResponse{Namespace: namespace, Count: len(endpoints.Items), Endpoints: make([]models.Endpoint, 0, len(endpoints.Items))}
		for _, endpoint := range endpoints.Items {
			resp.Endpoints = append(resp.Endpoints, mapEndpoint(endpoint))
		}

		c.JSON(http.StatusOK, resp)
	}
}

func ListEndpointSlices(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		namespace := c.DefaultQuery("namespace", cfg.DefaultNamespace)
		labelSelector := c.Query("label_selector")

		endpointSlices, err := k8s.Clientset.DiscoveryV1().EndpointSlices(namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: labelSelector})
		if err != nil {
			handleK8sError(c, err, "endpointslices/"+namespace)
			return
		}

		resp := models.EndpointSliceResponse{Namespace: namespace, Count: len(endpointSlices.Items), EndpointSlices: make([]models.EndpointSlice, 0, len(endpointSlices.Items))}
		for _, endpointSlice := range endpointSlices.Items {
			resp.EndpointSlices = append(resp.EndpointSlices, mapEndpointSlice(endpointSlice))
		}

		c.JSON(http.StatusOK, resp)
	}
}

func ListIngresses(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		namespace := c.DefaultQuery("namespace", cfg.DefaultNamespace)
		labelSelector := c.Query("label_selector")

		ingresses, err := k8s.Clientset.NetworkingV1().Ingresses(namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: labelSelector})
		if err != nil {
			handleK8sError(c, err, "ingresses/"+namespace)
			return
		}

		resp := models.IngressResponse{Namespace: namespace, Count: len(ingresses.Items), Ingresses: make([]models.Ingress, 0, len(ingresses.Items))}
		for _, ingress := range ingresses.Items {
			resp.Ingresses = append(resp.Ingresses, mapIngress(ingress))
		}

		c.JSON(http.StatusOK, resp)
	}
}

func ListPVCs(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		namespace := c.DefaultQuery("namespace", cfg.DefaultNamespace)
		labelSelector := c.Query("label_selector")

		pvcs, err := k8s.Clientset.CoreV1().PersistentVolumeClaims(namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: labelSelector})
		if err != nil {
			handleK8sError(c, err, "pvcs/"+namespace)
			return
		}

		resp := models.PersistentVolumeClaimResponse{Namespace: namespace, Count: len(pvcs.Items), PVCs: make([]models.PersistentVolumeClaim, 0, len(pvcs.Items))}
		for _, pvc := range pvcs.Items {
			resp.PVCs = append(resp.PVCs, mapPVC(pvc))
		}

		c.JSON(http.StatusOK, resp)
	}
}

func ListPVs(c *gin.Context) {
	labelSelector := c.Query("label_selector")

	pvs, err := k8s.Clientset.CoreV1().PersistentVolumes().List(context.TODO(), metav1.ListOptions{LabelSelector: labelSelector})
	if err != nil {
		handleK8sError(c, err, "pvs")
		return
	}

	resp := models.PersistentVolumeResponse{Count: len(pvs.Items), PVs: make([]models.PersistentVolume, 0, len(pvs.Items))}
	for _, pv := range pvs.Items {
		resp.PVs = append(resp.PVs, mapPV(pv))
	}

	c.JSON(http.StatusOK, resp)
}

func ListStorageClasses(c *gin.Context) {
	labelSelector := c.Query("label_selector")

	storageClasses, err := k8s.Clientset.StorageV1().StorageClasses().List(context.TODO(), metav1.ListOptions{LabelSelector: labelSelector})
	if err != nil {
		handleK8sError(c, err, "storageclasses")
		return
	}

	resp := models.StorageClassResponse{Count: len(storageClasses.Items), StorageClasses: make([]models.StorageClass, 0, len(storageClasses.Items))}
	for _, storageClass := range storageClasses.Items {
		resp.StorageClasses = append(resp.StorageClasses, mapStorageClass(storageClass))
	}

	c.JSON(http.StatusOK, resp)
}

func ListNetworkPolicies(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		namespace := c.DefaultQuery("namespace", cfg.DefaultNamespace)
		labelSelector := c.Query("label_selector")

		networkPolicies, err := k8s.Clientset.NetworkingV1().NetworkPolicies(namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: labelSelector})
		if err != nil {
			handleK8sError(c, err, "networkpolicies/"+namespace)
			return
		}

		resp := models.NetworkPolicyResponse{Namespace: namespace, Count: len(networkPolicies.Items), NetworkPolicies: make([]models.NetworkPolicy, 0, len(networkPolicies.Items))}
		for _, networkPolicy := range networkPolicies.Items {
			resp.NetworkPolicies = append(resp.NetworkPolicies, mapNetworkPolicy(networkPolicy))
		}

		c.JSON(http.StatusOK, resp)
	}
}

func ListResourceQuotas(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		namespace := c.DefaultQuery("namespace", cfg.DefaultNamespace)
		labelSelector := c.Query("label_selector")

		resourceQuotas, err := k8s.Clientset.CoreV1().ResourceQuotas(namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: labelSelector})
		if err != nil {
			handleK8sError(c, err, "resourcequotas/"+namespace)
			return
		}

		resp := models.ResourceQuotaResponse{Namespace: namespace, Count: len(resourceQuotas.Items), ResourceQuotas: make([]models.ResourceQuota, 0, len(resourceQuotas.Items))}
		for _, resourceQuota := range resourceQuotas.Items {
			resp.ResourceQuotas = append(resp.ResourceQuotas, mapResourceQuota(resourceQuota))
		}

		c.JSON(http.StatusOK, resp)
	}
}

func ListLimitRanges(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		namespace := c.DefaultQuery("namespace", cfg.DefaultNamespace)
		labelSelector := c.Query("label_selector")

		limitRanges, err := k8s.Clientset.CoreV1().LimitRanges(namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: labelSelector})
		if err != nil {
			handleK8sError(c, err, "limitranges/"+namespace)
			return
		}

		resp := models.LimitRangeResponse{Namespace: namespace, Count: len(limitRanges.Items), LimitRanges: make([]models.LimitRange, 0, len(limitRanges.Items))}
		for _, limitRange := range limitRanges.Items {
			resp.LimitRanges = append(resp.LimitRanges, mapLimitRange(limitRange))
		}

		c.JSON(http.StatusOK, resp)
	}
}

func ListConfigMaps(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		namespace := c.DefaultQuery("namespace", cfg.DefaultNamespace)
		labelSelector := c.Query("label_selector")

		configMaps, err := k8s.Clientset.CoreV1().ConfigMaps(namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: labelSelector})
		if err != nil {
			handleK8sError(c, err, "configmaps/"+namespace)
			return
		}

		resp := models.ConfigMapResponse{Namespace: namespace, Count: len(configMaps.Items), ConfigMaps: make([]models.ConfigMapMetadata, 0, len(configMaps.Items))}
		for _, configMap := range configMaps.Items {
			resp.ConfigMaps = append(resp.ConfigMaps, mapConfigMapMetadata(configMap))
		}

		c.JSON(http.StatusOK, resp)
	}
}

func ListPodMetrics(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		if k8s.MetricsClientset == nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "MetricsUnavailable", "message": "metrics clientset is not initialized"})
			return
		}

		namespace := c.DefaultQuery("namespace", cfg.DefaultNamespace)
		labelSelector := c.Query("label_selector")

		podMetrics, err := k8s.MetricsClientset.MetricsV1beta1().PodMetricses(namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: labelSelector})
		if err != nil {
			handleK8sError(c, err, "metrics/pods/"+namespace)
			return
		}

		resp := models.PodMetricsResponse{Namespace: namespace, Count: len(podMetrics.Items), Pods: make([]models.PodMetrics, 0, len(podMetrics.Items))}
		for _, metric := range podMetrics.Items {
			resp.Pods = append(resp.Pods, mapPodMetrics(metric))
		}

		c.JSON(http.StatusOK, resp)
	}
}

func ListNodeMetrics(c *gin.Context) {
	if k8s.MetricsClientset == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "MetricsUnavailable", "message": "metrics clientset is not initialized"})
		return
	}

	labelSelector := c.Query("label_selector")
	nodeMetrics, err := k8s.MetricsClientset.MetricsV1beta1().NodeMetricses().List(context.TODO(), metav1.ListOptions{LabelSelector: labelSelector})
	if err != nil {
		handleK8sError(c, err, "metrics/nodes")
		return
	}

	resp := models.NodeMetricsResponse{Count: len(nodeMetrics.Items), Nodes: make([]models.NodeMetrics, 0, len(nodeMetrics.Items))}
	for _, metric := range nodeMetrics.Items {
		resp.Nodes = append(resp.Nodes, mapNodeMetrics(metric))
	}

	c.JSON(http.StatusOK, resp)
}

func mapPodStatus(p corev1.Pod) models.PodStatusResponse {
	return models.PodStatusResponse{
		Name:               p.Name,
		Namespace:          p.Namespace,
		Phase:              string(p.Status.Phase),
		Reason:             p.Status.Reason,
		Message:            p.Status.Message,
		HostIP:             p.Status.HostIP,
		PodIP:              p.Status.PodIP,
		NodeName:           p.Spec.NodeName,
		StartTime:          timePtr(p.Status.StartTime),
		QOSClass:           string(p.Status.QOSClass),
		RestartPolicy:      string(p.Spec.RestartPolicy),
		ServiceAccountName: p.Spec.ServiceAccountName,
		InitContainers:     mapContainers(p.Spec.InitContainers, p.Status.InitContainerStatuses),
		Containers:         mapContainers(p.Spec.Containers, p.Status.ContainerStatuses),
		Conditions:         mapPodConditions(p.Status.Conditions),
	}
}

func mapContainers(containers []corev1.Container, statuses []corev1.ContainerStatus) []models.Container {
	statusMap := make(map[string]corev1.ContainerStatus, len(statuses))
	for _, status := range statuses {
		statusMap[status.Name] = status
	}

	result := make([]models.Container, 0, len(containers))
	for _, container := range containers {
		mapped := models.Container{
			Name:  container.Name,
			Image: container.Image,
			Ports: make([]models.ContainerPort, 0, len(container.Ports)),
		}

		if status, ok := statusMap[container.Name]; ok {
			mapped.Ready = status.Ready
			mapped.RestartCount = status.RestartCount
			mapped.State = mapContainerState(status.State)
			mapped.LastState = mapContainerState(status.LastTerminationState)
		}

		for _, port := range container.Ports {
			mapped.Ports = append(mapped.Ports, models.ContainerPort{
				ContainerPort: port.ContainerPort,
				Protocol:      string(port.Protocol),
			})
		}

		result = append(result, mapped)
	}

	return result
}

func mapContainerState(state corev1.ContainerState) models.ContainerState {
	switch {
	case state.Running != nil:
		return models.ContainerState{State: "running", StartedAt: &state.Running.StartedAt.Time}
	case state.Waiting != nil:
		return models.ContainerState{State: "waiting", Reason: state.Waiting.Reason, Message: state.Waiting.Message}
	case state.Terminated != nil:
		return models.ContainerState{
			State:      "terminated",
			Reason:     state.Terminated.Reason,
			Message:    state.Terminated.Message,
			ExitCode:   state.Terminated.ExitCode,
			Signal:     state.Terminated.Signal,
			StartedAt:  &state.Terminated.StartedAt.Time,
			FinishedAt: &state.Terminated.FinishedAt.Time,
		}
	default:
		return models.ContainerState{}
	}
}

func mapPodConditions(conditions []corev1.PodCondition) []models.Condition {
	result := make([]models.Condition, 0, len(conditions))
	for _, condition := range conditions {
		result = append(result, models.Condition{Type: string(condition.Type), Status: string(condition.Status)})
	}
	return result
}

func mapReplicaSet(replicaSet appsv1.ReplicaSet) models.ReplicaSet {
	desired := int32(0)
	if replicaSet.Spec.Replicas != nil {
		desired = *replicaSet.Spec.Replicas
	}
	return models.ReplicaSet{
		Name:              replicaSet.Name,
		Namespace:         replicaSet.Namespace,
		CreationTimestamp: &replicaSet.CreationTimestamp.Time,
		Labels:            replicaSet.Labels,
		Replicas: models.Replicas{
			Desired:   desired,
			Ready:     replicaSet.Status.ReadyReplicas,
			Available: replicaSet.Status.AvailableReplicas,
			Updated:   replicaSet.Status.FullyLabeledReplicas,
		},
		OwnerReferences: mapOwnerReferences(replicaSet.OwnerReferences),
	}
}

func mapStatefulSet(statefulSet appsv1.StatefulSet) models.StatefulSet {
	desired := int32(0)
	if statefulSet.Spec.Replicas != nil {
		desired = *statefulSet.Spec.Replicas
	}
	return models.StatefulSet{
		Name:              statefulSet.Name,
		Namespace:         statefulSet.Namespace,
		CreationTimestamp: &statefulSet.CreationTimestamp.Time,
		Labels:            statefulSet.Labels,
		Replicas: models.Replicas{
			Desired:   desired,
			Ready:     statefulSet.Status.ReadyReplicas,
			Available: statefulSet.Status.AvailableReplicas,
			Updated:   statefulSet.Status.UpdatedReplicas,
		},
		ServiceName: statefulSet.Spec.ServiceName,
	}
}

func mapDaemonSet(daemonSet appsv1.DaemonSet) models.DaemonSet {
	return models.DaemonSet{
		Name:                   daemonSet.Name,
		Namespace:              daemonSet.Namespace,
		CreationTimestamp:      &daemonSet.CreationTimestamp.Time,
		Labels:                 daemonSet.Labels,
		DesiredNumberScheduled: daemonSet.Status.DesiredNumberScheduled,
		CurrentNumberScheduled: daemonSet.Status.CurrentNumberScheduled,
		NumberReady:            daemonSet.Status.NumberReady,
		UpdatedNumberScheduled: daemonSet.Status.UpdatedNumberScheduled,
		NumberAvailable:        daemonSet.Status.NumberAvailable,
		NumberUnavailable:      daemonSet.Status.NumberUnavailable,
	}
}

func mapJob(job batchv1.Job) models.Job {
	conditions := make([]models.DepCondition, 0, len(job.Status.Conditions))
	for _, condition := range job.Status.Conditions {
		conditions = append(conditions, models.DepCondition{
			Type:    string(condition.Type),
			Status:  string(condition.Status),
			Reason:  condition.Reason,
			Message: condition.Message,
		})
	}

	return models.Job{
		Name:              job.Name,
		Namespace:         job.Namespace,
		CreationTimestamp: &job.CreationTimestamp.Time,
		Labels:            job.Labels,
		Completions:       job.Spec.Completions,
		Parallelism:       job.Spec.Parallelism,
		Active:            job.Status.Active,
		Succeeded:         job.Status.Succeeded,
		Failed:            job.Status.Failed,
		Conditions:        conditions,
	}
}

func mapCronJob(cronJob batchv1.CronJob) models.CronJob {
	active := make([]string, 0, len(cronJob.Status.Active))
	for _, ref := range cronJob.Status.Active {
		active = append(active, ref.Name)
	}

	return models.CronJob{
		Name:               cronJob.Name,
		Namespace:          cronJob.Namespace,
		CreationTimestamp:  &cronJob.CreationTimestamp.Time,
		Labels:             cronJob.Labels,
		Schedule:           cronJob.Spec.Schedule,
		Suspend:            cronJob.Spec.Suspend,
		Active:             active,
		LastScheduleTime:   timePtr(cronJob.Status.LastScheduleTime),
		LastSuccessfulTime: timePtr(cronJob.Status.LastSuccessfulTime),
	}
}

func mapEndpoint(endpoint corev1.Endpoints) models.Endpoint {
	subsets := make([]models.EndpointSubset, 0, len(endpoint.Subsets))
	for _, subset := range endpoint.Subsets {
		subsets = append(subsets, models.EndpointSubset{
			Addresses:         mapEndpointAddresses(subset.Addresses),
			NotReadyAddresses: mapEndpointAddresses(subset.NotReadyAddresses),
			Ports:             mapEndpointPorts(subset.Ports),
		})
	}

	return models.Endpoint{
		Name:              endpoint.Name,
		Namespace:         endpoint.Namespace,
		CreationTimestamp: &endpoint.CreationTimestamp.Time,
		Labels:            endpoint.Labels,
		Subsets:           subsets,
	}
}

func mapEndpointAddresses(addresses []corev1.EndpointAddress) []models.EndpointAddress {
	result := make([]models.EndpointAddress, 0, len(addresses))
	for _, address := range addresses {
		mapped := models.EndpointAddress{IP: address.IP, Hostname: address.Hostname}
		if address.NodeName != nil {
			mapped.NodeName = *address.NodeName
		}
		if address.TargetRef != nil {
			mapped.TargetKind = address.TargetRef.Kind
			mapped.TargetName = address.TargetRef.Name
		}
		result = append(result, mapped)
	}
	return result
}

func mapEndpointPorts(ports []corev1.EndpointPort) []models.EndpointPort {
	result := make([]models.EndpointPort, 0, len(ports))
	for _, port := range ports {
		result = append(result, models.EndpointPort{Name: port.Name, Port: port.Port, Protocol: string(port.Protocol)})
	}
	return result
}

func mapEndpointSlice(endpointSlice discoveryv1.EndpointSlice) models.EndpointSlice {
	ports := make([]models.EndpointSlicePort, 0, len(endpointSlice.Ports))
	for _, port := range endpointSlice.Ports {
		mapped := models.EndpointSlicePort{Port: port.Port}
		if port.Name != nil {
			mapped.Name = *port.Name
		}
		if port.Protocol != nil {
			mapped.Protocol = string(*port.Protocol)
		}
		if port.AppProtocol != nil {
			mapped.AppProtocol = *port.AppProtocol
		}
		ports = append(ports, mapped)
	}

	backends := make([]models.EndpointSliceBackend, 0, len(endpointSlice.Endpoints))
	for _, endpoint := range endpointSlice.Endpoints {
		backend := models.EndpointSliceBackend{
			Addresses: endpoint.Addresses,
			Conditions: map[string]*bool{
				"ready":       endpoint.Conditions.Ready,
				"serving":     endpoint.Conditions.Serving,
				"terminating": endpoint.Conditions.Terminating,
			},
		}
		if endpoint.Hostname != nil {
			backend.Hostname = *endpoint.Hostname
		}
		if endpoint.NodeName != nil {
			backend.NodeName = *endpoint.NodeName
		}
		if endpoint.TargetRef != nil {
			backend.TargetKind = endpoint.TargetRef.Kind
			backend.TargetName = endpoint.TargetRef.Name
		}
		if endpoint.Zone != nil {
			backend.Zone = *endpoint.Zone
		}
		backends = append(backends, backend)
	}

	return models.EndpointSlice{
		Name:              endpointSlice.Name,
		Namespace:         endpointSlice.Namespace,
		CreationTimestamp: &endpointSlice.CreationTimestamp.Time,
		Labels:            endpointSlice.Labels,
		AddressType:       string(endpointSlice.AddressType),
		Ports:             ports,
		Endpoints:         backends,
	}
}

func mapIngress(ingress networkingv1.Ingress) models.Ingress {
	tlsHosts := make([]string, 0)
	for _, tls := range ingress.Spec.TLS {
		tlsHosts = append(tlsHosts, tls.Hosts...)
	}
	sort.Strings(tlsHosts)

	rules := make([]models.IngressRule, 0)
	for _, rule := range ingress.Spec.Rules {
		if rule.HTTP == nil {
			continue
		}
		for _, path := range rule.HTTP.Paths {
			mapped := models.IngressRule{
				Host:        rule.Host,
				Path:        path.Path,
				ServiceName: path.Backend.Service.Name,
				ServicePort: path.Backend.Service.Port.String(),
			}
			if path.PathType != nil {
				mapped.PathType = string(*path.PathType)
			}
			rules = append(rules, mapped)
		}
	}

	loadBalancerIPs := make([]string, 0)
	for _, item := range ingress.Status.LoadBalancer.Ingress {
		if item.IP != "" {
			loadBalancerIPs = append(loadBalancerIPs, item.IP)
		}
		if item.Hostname != "" {
			loadBalancerIPs = append(loadBalancerIPs, item.Hostname)
		}
	}

	mapped := models.Ingress{
		Name:              ingress.Name,
		Namespace:         ingress.Namespace,
		CreationTimestamp: &ingress.CreationTimestamp.Time,
		Labels:            ingress.Labels,
		TLSHosts:          tlsHosts,
		Rules:             rules,
		LoadBalancerIPs:   loadBalancerIPs,
	}
	if ingress.Spec.IngressClassName != nil {
		mapped.ClassName = *ingress.Spec.IngressClassName
	}
	return mapped
}

func mapPVC(pvc corev1.PersistentVolumeClaim) models.PersistentVolumeClaim {
	storageClassName := ""
	if pvc.Spec.StorageClassName != nil {
		storageClassName = *pvc.Spec.StorageClassName
	}

	return models.PersistentVolumeClaim{
		Name:              pvc.Name,
		Namespace:         pvc.Namespace,
		CreationTimestamp: &pvc.CreationTimestamp.Time,
		Labels:            pvc.Labels,
		Phase:             string(pvc.Status.Phase),
		StorageClassName:  storageClassName,
		VolumeName:        pvc.Spec.VolumeName,
		AccessModes:       mapAccessModes(pvc.Spec.AccessModes),
		Capacity:          mapResourceList(pvc.Status.Capacity),
		Requested:         mapResourceList(pvc.Spec.Resources.Requests),
	}
}

func mapPV(pv corev1.PersistentVolume) models.PersistentVolume {
	mapped := models.PersistentVolume{
		Name:              pv.Name,
		CreationTimestamp: &pv.CreationTimestamp.Time,
		Labels:            pv.Labels,
		Phase:             string(pv.Status.Phase),
		StorageClassName:  pv.Spec.StorageClassName,
		Capacity:          mapResourceList(pv.Spec.Capacity),
		AccessModes:       mapAccessModes(pv.Spec.AccessModes),
		ReclaimPolicy:     string(pv.Spec.PersistentVolumeReclaimPolicy),
		Reason:            pv.Status.Reason,
	}
	if pv.Spec.ClaimRef != nil {
		mapped.ClaimNamespace = pv.Spec.ClaimRef.Namespace
		mapped.ClaimName = pv.Spec.ClaimRef.Name
	}
	if pv.Spec.VolumeMode != nil {
		mapped.VolumeMode = string(*pv.Spec.VolumeMode)
	}
	return mapped
}

func mapStorageClass(storageClass storagev1.StorageClass) models.StorageClass {
	mapped := models.StorageClass{
		Name:                 storageClass.Name,
		CreationTimestamp:    &storageClass.CreationTimestamp.Time,
		Labels:               storageClass.Labels,
		Provisioner:          storageClass.Provisioner,
		AllowVolumeExpansion: storageClass.AllowVolumeExpansion,
		Parameters:           storageClass.Parameters,
	}
	if storageClass.ReclaimPolicy != nil {
		mapped.ReclaimPolicy = string(*storageClass.ReclaimPolicy)
	}
	if storageClass.VolumeBindingMode != nil {
		mapped.VolumeBindingMode = string(*storageClass.VolumeBindingMode)
	}
	return mapped
}

func mapNetworkPolicy(networkPolicy networkingv1.NetworkPolicy) models.NetworkPolicy {
	policyTypes := make([]string, 0, len(networkPolicy.Spec.PolicyTypes))
	for _, policyType := range networkPolicy.Spec.PolicyTypes {
		policyTypes = append(policyTypes, string(policyType))
	}
	return models.NetworkPolicy{
		Name:              networkPolicy.Name,
		Namespace:         networkPolicy.Namespace,
		CreationTimestamp: &networkPolicy.CreationTimestamp.Time,
		Labels:            networkPolicy.Labels,
		PodSelector:       networkPolicy.Spec.PodSelector.MatchLabels,
		PolicyTypes:       policyTypes,
		IngressRuleCount:  len(networkPolicy.Spec.Ingress),
		EgressRuleCount:   len(networkPolicy.Spec.Egress),
	}
}

func mapResourceQuota(resourceQuota corev1.ResourceQuota) models.ResourceQuota {
	return models.ResourceQuota{
		Name:              resourceQuota.Name,
		Namespace:         resourceQuota.Namespace,
		CreationTimestamp: &resourceQuota.CreationTimestamp.Time,
		Labels:            resourceQuota.Labels,
		Hard:              mapResourceList(resourceQuota.Status.Hard),
		Used:              mapResourceList(resourceQuota.Status.Used),
	}
}

func mapLimitRange(limitRange corev1.LimitRange) models.LimitRange {
	items := make([]models.LimitRangeItem, 0, len(limitRange.Spec.Limits))
	for _, item := range limitRange.Spec.Limits {
		items = append(items, models.LimitRangeItem{
			Type:                 string(item.Type),
			Max:                  mapResourceList(item.Max),
			Min:                  mapResourceList(item.Min),
			Default:              mapResourceList(item.Default),
			DefaultRequest:       mapResourceList(item.DefaultRequest),
			MaxLimitRequestRatio: mapResourceList(item.MaxLimitRequestRatio),
		})
	}
	return models.LimitRange{
		Name:              limitRange.Name,
		Namespace:         limitRange.Namespace,
		CreationTimestamp: &limitRange.CreationTimestamp.Time,
		Labels:            limitRange.Labels,
		Items:             items,
	}
}

func mapConfigMapMetadata(configMap corev1.ConfigMap) models.ConfigMapMetadata {
	dataKeys := make([]string, 0, len(configMap.Data))
	for key := range configMap.Data {
		dataKeys = append(dataKeys, key)
	}
	sort.Strings(dataKeys)

	binaryDataKeys := make([]string, 0, len(configMap.BinaryData))
	for key := range configMap.BinaryData {
		binaryDataKeys = append(binaryDataKeys, key)
	}
	sort.Strings(binaryDataKeys)

	return models.ConfigMapMetadata{
		Name:              configMap.Name,
		Namespace:         configMap.Namespace,
		CreationTimestamp: &configMap.CreationTimestamp.Time,
		Labels:            configMap.Labels,
		DataKeys:          dataKeys,
		BinaryDataKeys:    binaryDataKeys,
		Immutable:         configMap.Immutable,
	}
}

func mapPodMetrics(metric metricsv1beta1.PodMetrics) models.PodMetrics {
	containers := make([]models.ContainerResourceMetrics, 0, len(metric.Containers))
	for _, container := range metric.Containers {
		containers = append(containers, models.ContainerResourceMetrics{
			Name:  container.Name,
			Usage: mapResourceList(container.Usage),
		})
	}
	return models.PodMetrics{
		Name:       metric.Name,
		Namespace:  metric.Namespace,
		Timestamp:  &metric.Timestamp.Time,
		Window:     metric.Window.String(),
		Containers: containers,
	}
}

func mapNodeMetrics(metric metricsv1beta1.NodeMetrics) models.NodeMetrics {
	return models.NodeMetrics{
		Name:      metric.Name,
		Timestamp: &metric.Timestamp.Time,
		Window:    metric.Window.String(),
		Usage:     mapResourceList(metric.Usage),
	}
}

func mapOwnerReferences(ownerReferences []metav1.OwnerReference) []models.OwnerReference {
	result := make([]models.OwnerReference, 0, len(ownerReferences))
	for _, ownerReference := range ownerReferences {
		result = append(result, models.OwnerReference{Kind: ownerReference.Kind, Name: ownerReference.Name})
	}
	return result
}

func mapAccessModes(accessModes []corev1.PersistentVolumeAccessMode) []string {
	result := make([]string, 0, len(accessModes))
	for _, accessMode := range accessModes {
		result = append(result, string(accessMode))
	}
	return result
}

func mapResourceList(resources corev1.ResourceList) map[string]string {
	result := make(map[string]string, len(resources))
	for key, value := range resources {
		result[string(key)] = value.String()
	}
	return result
}

func timePtr(t *metav1.Time) *time.Time {
	if t == nil {
		return nil
	}
	return &t.Time
}
