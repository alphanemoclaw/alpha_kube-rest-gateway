package config

import (
	"os"
	"path/filepath"
	"strconv"
)

type Config struct {
	KubeconfigPath       string
	ApiToken             string
	Host                 string
	Port                 int
	LogLevel             string
	DefaultNamespace     string
	DefaultLogTailLines int64
}

func LoadConfig() *Config {
	home, _ := os.UserHomeDir()
	kubeconfigDefault := filepath.Join(home, ".kube", "config")

	port, _ := strconv.Atoi(getEnv("GATEWAY_PORT", "8081"))
	tailLines, _ := strconv.ParseInt(getEnv("DEFAULT_LOG_TAIL_LINES", "100"), 10, 64)

	return &Config{
		KubeconfigPath:       getEnv("KUBECONFIG_PATH", kubeconfigDefault),
		ApiToken:             getEnv("GATEWAY_API_TOKEN", ""),
		Host:                 getEnv("GATEWAY_HOST", "0.0.0.0"),
		Port:                 port,
		LogLevel:             getEnv("LOG_LEVEL", "INFO"),
		DefaultNamespace:     getEnv("DEFAULT_NAMESPACE", "default"),
		DefaultLogTailLines: tailLines,
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
