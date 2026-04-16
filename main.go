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

	// Middleware
	r.Use(handlers.AuthMiddleware(cfg))

	// Health check
	r.GET("/healthz", handlers.Healthz)
	r.GET("/", handlers.Root)

	// API endpoints
	api := r.Group("/api")
	{
		api.GET("/help", handlers.ApiHelp)
		api.GET("/pods", handlers.ListPods(cfg))
		api.GET("/pods/:pod_name", handlers.GetPod(cfg))
		api.GET("/logs/:pod_name", handlers.GetPodLogs(cfg))
		api.GET("/services", handlers.ListServices(cfg))
		api.GET("/deployments", handlers.ListDeployments(cfg))
		api.GET("/nodes", handlers.ListNodes)
		api.GET("/namespaces", handlers.ListNamespaces)
		api.GET("/events", handlers.ListEvents(cfg))
	}

	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	fmt.Printf("Starting server on %s\n", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
