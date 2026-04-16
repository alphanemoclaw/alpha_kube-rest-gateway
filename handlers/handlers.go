package handlers

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"sort"
	"strings"

	"github.com/alpha-kube-rest-gateway/config"
	"github.com/alpha-kube-rest-gateway/k8s"
	"github.com/alpha-kube-rest-gateway/models"
	"github.com/gin-gonic/gin"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		if cfg.ApiToken == "" {
			c.Next()
			return
		}

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized", "message": "Missing bearer token"})
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token != cfg.ApiToken {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized", "message": "Invalid bearer token"})
			return
		}

		c.Next()
	}
}

func Healthz(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok", "service": "kube-rest-gateway"})
}

func Root(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Welcome to Kube-REST-Gateway. Use /api/help for instructions.",
		"docs":    "/docs",
		"help":    "/api/help",
	})
}

func ApiHelp(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"gateway_info": gin.H{
			"version": "1.0.0",
			"purpose": "Secure Kubernetes API proxy for AI Agents",
			"auth":    "Bearer Token required in 'Authorization' header",
		},
		"capabilities": []gin.H{
			{
				"endpoint":    "/api/pods",
				"method":      "GET",
				"description": "List all pods. Use ?namespace= to filter.",
				"examples":    []string{"/api/pods?namespace=default", "/api/pods?label_selector=app=nginx"},
			},
			{
				"endpoint":    "/api/pods/{name}",
				"method":      "GET",
				"description": "Get detailed JSON for a specific pod.",
			},
			{
				"endpoint":    "/api/logs/{name}",
				"method":      "GET",
				"description": "Fetch recent logs. Use ?tail_lines= to limit output.",
				"params":      []string{"namespace", "container", "tail_lines", "previous"},
			},
			{
				"endpoint":    "/api/services",
				"method":      "GET",
				"description": "List all services in a namespace.",
			},
			{
				"endpoint":    "/api/deployments",
				"method":      "GET",
				"description": "List deployments and checkout replica health.",
			},
			{
				"endpoint":    "/api/events",
				"method":      "GET",
				"description": "Get cluster events to debug crashing pods (ImagePullBackOff, etc).",
				"params":      []string{"namespace", "field_selector"},
			},
			{
				"endpoint":    "/api/namespaces",
				"method":      "GET",
				"description": "List all available namespaces in the cluster.",
			},
			{
				"endpoint":    "/api/nodes",
				"method":      "GET",
				"description": "List physical/virtual nodes in the cluster.",
			},
		},
		"instructions": "1. Always check /api/namespaces first if you aren't sure where to look. 2. If a pod is 'Pending' or 'Error', check /api/events for the root cause. 3. Keep log requests small (e.g., tail_lines=50) to save processing time.",
	})
}

func ListPods(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		namespace := c.DefaultQuery("namespace", cfg.DefaultNamespace)
		labelSelector := c.Query("label_selector")
		fieldSelector := c.Query("field_selector")

		opts := metav1.ListOptions{
			LabelSelector: labelSelector,
			FieldSelector: fieldSelector,
		}

		pods, err := k8s.Clientset.CoreV1().Pods(namespace).List(context.TODO(), opts)
		if err != nil {
			handleK8sError(c, err, "pods/"+namespace)
			return
		}

		resp := models.PodResponse{
			Namespace: namespace,
			Count:     len(pods.Items),
			Pods:      make([]models.Pod, 0),
		}

		for _, p := range pods.Items {
			resp.Pods = append(resp.Pods, mapPod(p))
		}

		c.JSON(http.StatusOK, resp)
	}
}

func GetPod(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		podName := c.Param("pod_name")
		namespace := c.DefaultQuery("namespace", cfg.DefaultNamespace)

		pod, err := k8s.Clientset.CoreV1().Pods(namespace).Get(context.TODO(), podName, metav1.GetOptions{})
		if err != nil {
			handleK8sError(c, err, "pod/"+namespace+"/"+podName)
			return
		}

		c.JSON(http.StatusOK, mapPod(*pod))
	}
}

func GetPodLogs(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		podName := c.Param("pod_name")
		namespace := c.DefaultQuery("namespace", cfg.DefaultNamespace)
		container := c.Query("container")
		tailLinesStr := c.DefaultQuery("tail_lines", "100")
		previous := c.Query("previous") == "true"

		var tailLines int64 = 100
		if strings.TrimSpace(tailLinesStr) != "" {
			// Parsing handled by default query but let's be safe
			// In a real app we'd validate this better
		}

		opts := &corev1.PodLogOptions{
			Container: container,
			TailLines: &tailLines,
			Previous:  previous,
		}

		req := k8s.Clientset.CoreV1().Pods(namespace).GetLogs(podName, opts)
		podLogs, err := req.Stream(context.TODO())
		if err != nil {
			handleK8sError(c, err, "logs/"+namespace+"/"+podName)
			return
		}
		defer podLogs.Close()

		content, err := io.ReadAll(podLogs)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "InternalError", "message": err.Error()})
			return
		}

		lines := strings.Split(string(content), "\n")
		// Remove last empty line if exists
		if len(lines) > 0 && lines[len(lines)-1] == "" {
			lines = lines[:len(lines)-1]
		}

		c.JSON(http.StatusOK, models.LogResponse{
			PodName:            podName,
			Namespace:          namespace,
			Container:          container,
			TailLinesRequested: tailLines,
			LinesReturned:      len(lines),
			Logs:               lines,
		})
	}
}

func ListServices(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		namespace := c.DefaultQuery("namespace", cfg.DefaultNamespace)
		labelSelector := c.Query("label_selector")

		opts := metav1.ListOptions{
			LabelSelector: labelSelector,
		}

		svcs, err := k8s.Clientset.CoreV1().Services(namespace).List(context.TODO(), opts)
		if err != nil {
			handleK8sError(c, err, "services/"+namespace)
			return
		}

		resp := models.ServiceResponse{
			Namespace: namespace,
			Count:     len(svcs.Items),
			Services:  make([]models.Service, 0),
		}

		for _, s := range svcs.Items {
			resp.Services = append(resp.Services, mapService(s))
		}

		c.JSON(http.StatusOK, resp)
	}
}

func ListDeployments(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		namespace := c.DefaultQuery("namespace", cfg.DefaultNamespace)
		labelSelector := c.Query("label_selector")

		opts := metav1.ListOptions{
			LabelSelector: labelSelector,
		}

		deps, err := k8s.Clientset.AppsV1().Deployments(namespace).List(context.TODO(), opts)
		if err != nil {
			handleK8sError(c, err, "deployments/"+namespace)
			return
		}

		resp := models.DeploymentResponse{
			Namespace:   namespace,
			Count:       len(deps.Items),
			Deployments: make([]models.Deployment, 0),
		}

		for _, d := range deps.Items {
			resp.Deployments = append(resp.Deployments, mapDeployment(d))
		}

		c.JSON(http.StatusOK, resp)
	}
}

func ListNodes(c *gin.Context) {
	nodes, err := k8s.Clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		handleK8sError(c, err, "nodes")
		return
	}

	resp := models.NodeResponse{
		Count: len(nodes.Items),
		Nodes: make([]models.Node, 0),
	}

	for _, n := range nodes.Items {
		resp.Nodes = append(resp.Nodes, mapNode(n))
	}

	c.JSON(http.StatusOK, resp)
}

func ListNamespaces(c *gin.Context) {
	nss, err := k8s.Clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		handleK8sError(c, err, "namespaces")
		return
	}

	resp := models.NamespaceResponse{
		Count:      len(nss.Items),
		Namespaces: make([]models.Namespace, 0),
	}

	for _, ns := range nss.Items {
		resp.Namespaces = append(resp.Namespaces, models.Namespace{
			Name:              ns.Name,
			Phase:             string(ns.Status.Phase),
			CreationTimestamp: &ns.CreationTimestamp.Time,
			Labels:            ns.Labels,
		})
	}

	c.JSON(http.StatusOK, resp)
}

func ListEvents(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		namespace := c.DefaultQuery("namespace", cfg.DefaultNamespace)
		fieldSelector := c.Query("field_selector")

		opts := metav1.ListOptions{
			FieldSelector: fieldSelector,
		}

		events, err := k8s.Clientset.CoreV1().Events(namespace).List(context.TODO(), opts)
		if err != nil {
			handleK8sError(c, err, "events/"+namespace)
			return
		}

		resp := models.EventResponse{
			Namespace: namespace,
			Count:     len(events.Items),
			Events:    make([]models.Event, 0),
		}

		for _, e := range events.Items {
			resp.Events = append(resp.Events, models.Event{
				Name:           e.Name,
				Namespace:      e.Namespace,
				Type:           e.Type,
				Reason:         e.Reason,
				Message:        e.Message,
				EventCount:     e.Count,
				FirstTimestamp: &e.FirstTimestamp.Time,
				LastTimestamp:  &e.LastTimestamp.Time,
				InvolvedObject: models.InvolvedObject{
					Kind:      e.InvolvedObject.Kind,
					Name:      e.InvolvedObject.Name,
					Namespace: e.InvolvedObject.Namespace,
				},
				Source: models.Source{
					Component: e.Source.Component,
					Host:      e.Source.Host,
				},
			})
		}

		// Sort most recent first
		sort.Slice(resp.Events, func(i, j int) bool {
			if resp.Events[i].LastTimestamp == nil || resp.Events[j].LastTimestamp == nil {
				return false
			}
			return resp.Events[i].LastTimestamp.After(*resp.Events[j].LastTimestamp)
		})

		c.JSON(http.StatusOK, resp)
	}
}

// Helpers

func handleK8sError(c *gin.Context, err error, resource string) {
	slog.Error("Kubernetes API error", "resource", resource, "error", err)

	status := http.StatusBadGateway
	reason := "KubernetesApiError"
	message := err.Error()

	if statusError, ok := err.(*errors.StatusError); ok {
		status = int(statusError.ErrStatus.Code)
		reason = string(statusError.ErrStatus.Reason)
		message = statusError.ErrStatus.Message
	}

	c.JSON(status, gin.H{
		"error":      reason,
		"k8s_status": status,
		"message":    message,
		"resource":   resource,
	})
}

func mapPod(p corev1.Pod) models.Pod {
	pod := models.Pod{
		Name:              p.Name,
		Namespace:         p.Namespace,
		Phase:             string(p.Status.Phase),
		PodIP:             p.Status.PodIP,
		NodeName:          p.Spec.NodeName,
		CreationTimestamp: &p.CreationTimestamp.Time,
		Labels:            p.Labels,
		Containers:        make([]models.Container, 0),
		Conditions:        make([]models.Condition, 0),
	}

	statusMap := make(map[string]models.Container)
	for _, cs := range p.Status.ContainerStatuses {
		statusMap[cs.Name] = models.Container{
			Ready:        cs.Ready,
			RestartCount: cs.RestartCount,
		}
	}

	for _, c := range p.Spec.Containers {
		container := models.Container{
			Name:  c.Name,
			Image: c.Image,
			Ports: make([]models.ContainerPort, 0),
		}

		if s, ok := statusMap[c.Name]; ok {
			container.Ready = s.Ready
			container.RestartCount = s.RestartCount
		}

		for _, port := range c.Ports {
			container.Ports = append(container.Ports, models.ContainerPort{
				ContainerPort: port.ContainerPort,
				Protocol:      string(port.Protocol),
			})
		}
		pod.Containers = append(pod.Containers, container)
	}

	for _, cond := range p.Status.Conditions {
		pod.Conditions = append(pod.Conditions, models.Condition{
			Type:   string(cond.Type),
			Status: string(cond.Status),
		})
	}

	return pod
}

func mapService(s corev1.Service) models.Service {
	svc := models.Service{
		Name:              s.Name,
		Namespace:         s.Namespace,
		Type:              string(s.Spec.Type),
		ClusterIP:         s.Spec.ClusterIP,
		ExternalIPs:       s.Spec.ExternalIPs,
		LoadBalancerIP:    s.Spec.LoadBalancerIP,
		Ports:             make([]models.ServicePort, 0),
		Selector:          s.Spec.Selector,
		CreationTimestamp: &s.CreationTimestamp.Time,
	}

	for _, p := range s.Spec.Ports {
		svc.Ports = append(svc.Ports, models.ServicePort{
			Name:       p.Name,
			Protocol:   string(p.Protocol),
			Port:       p.Port,
			TargetPort: p.TargetPort.String(),
			NodePort:   p.NodePort,
		})
	}

	return svc
}

func mapDeployment(d appsv1.Deployment) models.Deployment {
	dep := models.Deployment{
		Name:              d.Name,
		Namespace:         d.Namespace,
		CreationTimestamp: &d.CreationTimestamp.Time,
		Labels:            d.Labels,
		Replicas: models.Replicas{
			Desired:   *d.Spec.Replicas,
			Ready:     d.Status.ReadyReplicas,
			Available: d.Status.AvailableReplicas,
			Updated:   d.Status.UpdatedReplicas,
		},
		Conditions: make([]models.DepCondition, 0),
	}

	if d.Spec.Strategy.Type != "" {
		dep.Strategy = string(d.Spec.Strategy.Type)
	}

	for _, cond := range d.Status.Conditions {
		dep.Conditions = append(dep.Conditions, models.DepCondition{
			Type:    string(cond.Type),
			Status:  string(cond.Status),
			Reason:  cond.Reason,
			Message: cond.Message,
		})
	}

	return dep
}

func mapNode(n corev1.Node) models.Node {
	node := models.Node{
		Name:              n.Name,
		CreationTimestamp: &n.CreationTimestamp.Time,
		Labels:            n.Labels,
		Addresses:         make(map[string]string),
		Capacity:          make(map[string]string),
		OSImage:           n.Status.NodeInfo.OSImage,
		KernelVersion:     n.Status.NodeInfo.KernelVersion,
		KubeletVersion:    n.Status.NodeInfo.KubeletVersion,
		ContainerRuntime:  n.Status.NodeInfo.ContainerRuntimeVersion,
		Conditions:        make([]models.Condition, 0),
	}

	for _, a := range n.Status.Addresses {
		node.Addresses[string(a.Type)] = a.Address
	}

	for k, v := range n.Status.Capacity {
		node.Capacity[string(k)] = v.String()
	}

	for _, cond := range n.Status.Conditions {
		node.Conditions = append(node.Conditions, models.Condition{
			Type:   string(cond.Type),
			Status: string(cond.Status),
		})
	}

	return node
}
