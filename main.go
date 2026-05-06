package main

import (
	"fmt"
	"log"

	"github.com/alpha-kube-rest-gateway/config"
	"github.com/alpha-kube-rest-gateway/handlers"
	"github.com/alpha-kube-rest-gateway/k8s"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()

	if err := k8s.InitK8sClient(cfg); err != nil {
		log.Fatalf("Failed to initialize Kubernetes client: %v", err)
	}

	r := gin.Default()

	// Public routes
	r.GET("/healthz", handlers.Healthz)
	r.GET("/", handlers.Root)

	// API endpoints (Authenticated)
	api := r.Group("/api")
	api.Use(handlers.AuthMiddleware(cfg))
	{
		api.GET("/help", handlers.ApiHelp)
		api.GET("/pods", handlers.ListPods(cfg))
		api.GET("/pods/:pod_name/status", handlers.GetPodStatus(cfg))
		api.GET("/pods/:pod_name", handlers.GetPod(cfg))
		api.GET("/logs/:pod_name", handlers.GetPodLogs(cfg))
		api.GET("/services", handlers.ListServices(cfg))
		api.GET("/deployments", handlers.ListDeployments(cfg))
		api.GET("/replicasets", handlers.ListReplicaSets(cfg))
		api.GET("/statefulsets", handlers.ListStatefulSets(cfg))
		api.GET("/daemonsets", handlers.ListDaemonSets(cfg))
		api.GET("/jobs", handlers.ListJobs(cfg))
		api.GET("/cronjobs", handlers.ListCronJobs(cfg))
		api.GET("/nodes", handlers.ListNodes)
		api.GET("/namespaces", handlers.ListNamespaces)
		api.GET("/events", handlers.ListEvents(cfg))
		api.GET("/endpoints", handlers.ListEndpoints(cfg))
		api.GET("/endpointslices", handlers.ListEndpointSlices(cfg))
		api.GET("/ingresses", handlers.ListIngresses(cfg))
		api.GET("/pvcs", handlers.ListPVCs(cfg))
		api.GET("/pvs", handlers.ListPVs)
		api.GET("/storageclasses", handlers.ListStorageClasses)
		api.GET("/networkpolicies", handlers.ListNetworkPolicies(cfg))
		api.GET("/resourcequotas", handlers.ListResourceQuotas(cfg))
		api.GET("/limitranges", handlers.ListLimitRanges(cfg))
		api.GET("/configmaps", handlers.ListConfigMaps(cfg))
		api.GET("/metrics/pods", handlers.ListPodMetrics(cfg))
		api.GET("/metrics/nodes", handlers.ListNodeMetrics)
	}

	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	fmt.Printf("Starting server on %s\n", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
