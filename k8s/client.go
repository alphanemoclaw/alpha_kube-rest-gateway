package k8s

import (
	"fmt"
	"log/slog"

	"github.com/alpha-kube-rest-gateway/config"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var Clientset *kubernetes.Clientset

func InitK8sClient(cfg *config.Config) error {
	slog.Info("Initialising Kubernetes client", "path", cfg.KubeconfigPath)

	k8sCfg, err := clientcmd.BuildConfigFromFlags("", cfg.KubeconfigPath)
	if err != nil {
		// Try in-cluster config as fallback if path is empty or fails
		slog.Warn("Failed to build config from flags, trying in-cluster config", "error", err)
		k8sCfg, err = rest.InClusterConfig()
		if err != nil {
			return fmt.Errorf("failed to load kubernetes config: %w", err)
		}
	}

	Clientset, err = kubernetes.NewForConfig(k8sCfg)
	if err != nil {
		return fmt.Errorf("failed to create kubernetes clientset: %w", err)
	}

	return nil
}
