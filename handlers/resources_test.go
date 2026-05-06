package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/alpha-kube-rest-gateway/config"
	"github.com/alpha-kube-rest-gateway/k8s"
	"github.com/alpha-kube-rest-gateway/models"
	"github.com/gin-gonic/gin"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
	k8stesting "k8s.io/client-go/testing"
	metricsv1beta1 "k8s.io/metrics/pkg/apis/metrics/v1beta1"
	metricsfake "k8s.io/metrics/pkg/client/clientset/versioned/fake"
)

func TestGetPodStatusIncludesContainerWaitingReason(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfg := &config.Config{DefaultNamespace: "default"}
	k8s.Clientset = fake.NewSimpleClientset(&corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{Name: "api-123", Namespace: "default"},
		Spec: corev1.PodSpec{
			NodeName: "node-a",
			Containers: []corev1.Container{
				{Name: "api", Image: "example/api:latest"},
			},
		},
		Status: corev1.PodStatus{
			Phase: corev1.PodPending,
			ContainerStatuses: []corev1.ContainerStatus{
				{
					Name:         "api",
					Image:        "example/api:latest",
					RestartCount: 3,
					State: corev1.ContainerState{Waiting: &corev1.ContainerStateWaiting{
						Reason:  "ImagePullBackOff",
						Message: "failed to pull image",
					}},
				},
			},
		},
	})

	router := gin.New()
	router.GET("/api/pods/:pod_name/status", GetPodStatus(cfg))

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/pods/api-123/status", nil)
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	var resp models.PodStatusResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	if resp.Phase != "Pending" {
		t.Fatalf("expected Pending phase, got %q", resp.Phase)
	}
	if len(resp.Containers) != 1 {
		t.Fatalf("expected 1 container, got %d", len(resp.Containers))
	}
	if resp.Containers[0].State.Reason != "ImagePullBackOff" {
		t.Fatalf("expected ImagePullBackOff reason, got %q", resp.Containers[0].State.Reason)
	}
}

func TestRoutesRegisterWithoutConflicts(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfg := &config.Config{DefaultNamespace: "default"}

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("route registration panicked: %v", r)
		}
	}()

	router := gin.New()
	api := router.Group("/api")
	api.GET("/help", ApiHelp)
	api.GET("/pods", ListPods(cfg))
	api.GET("/pods/:pod_name/status", GetPodStatus(cfg))
	api.GET("/pods/:pod_name", GetPod(cfg))
	api.GET("/logs/:pod_name", GetPodLogs(cfg))
	api.GET("/services", ListServices(cfg))
	api.GET("/deployments", ListDeployments(cfg))
	api.GET("/replicasets", ListReplicaSets(cfg))
	api.GET("/statefulsets", ListStatefulSets(cfg))
	api.GET("/daemonsets", ListDaemonSets(cfg))
	api.GET("/jobs", ListJobs(cfg))
	api.GET("/cronjobs", ListCronJobs(cfg))
	api.GET("/nodes", ListNodes)
	api.GET("/namespaces", ListNamespaces)
	api.GET("/events", ListEvents(cfg))
	api.GET("/endpoints", ListEndpoints(cfg))
	api.GET("/endpointslices", ListEndpointSlices(cfg))
	api.GET("/ingresses", ListIngresses(cfg))
	api.GET("/pvcs", ListPVCs(cfg))
	api.GET("/pvs", ListPVs)
	api.GET("/storageclasses", ListStorageClasses)
	api.GET("/networkpolicies", ListNetworkPolicies(cfg))
	api.GET("/resourcequotas", ListResourceQuotas(cfg))
	api.GET("/limitranges", ListLimitRanges(cfg))
	api.GET("/configmaps", ListConfigMaps(cfg))
	api.GET("/metrics/pods", ListPodMetrics(cfg))
	api.GET("/metrics/nodes", ListNodeMetrics)
}

func TestListConfigMapsReturnsKeysOnly(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfg := &config.Config{DefaultNamespace: "default"}
	k8s.Clientset = fake.NewSimpleClientset(&corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{Name: "app-config", Namespace: "default"},
		Data: map[string]string{
			"APP_MODE": "production",
			"TOKEN":    "should-not-leak",
		},
	})

	router := gin.New()
	router.GET("/api/configmaps", ListConfigMaps(cfg))

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/configmaps", nil)
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	if strings.Contains(rec.Body.String(), "should-not-leak") {
		t.Fatalf("configmap response leaked data value: %s", rec.Body.String())
	}

	var resp models.ConfigMapResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	if len(resp.ConfigMaps) != 1 {
		t.Fatalf("expected 1 configmap, got %d", len(resp.ConfigMaps))
	}
	if len(resp.ConfigMaps[0].DataKeys) != 2 {
		t.Fatalf("expected 2 data keys, got %d", len(resp.ConfigMaps[0].DataKeys))
	}
}

func TestListPodMetricsUsesMetricsClient(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfg := &config.Config{DefaultNamespace: "default"}
	podMetrics := &metricsv1beta1.PodMetricsList{
		Items: []metricsv1beta1.PodMetrics{
			{
				ObjectMeta: metav1.ObjectMeta{Name: "api-123", Namespace: "default"},
				Timestamp:  metav1.Time{Time: time.Unix(100, 0)},
				Window:     metav1.Duration{Duration: time.Minute},
				Containers: []metricsv1beta1.ContainerMetrics{
					{
						Name: "api",
						Usage: corev1.ResourceList{
							corev1.ResourceCPU:    resource.MustParse("50m"),
							corev1.ResourceMemory: resource.MustParse("128Mi"),
						},
					},
				},
			},
		},
	}
	metricsClient := metricsfake.NewSimpleClientset()
	metricsClient.Fake.PrependReactor("list", "pods", func(action k8stesting.Action) (bool, runtime.Object, error) {
		return true, podMetrics, nil
	})
	k8s.MetricsClientset = metricsClient

	router := gin.New()
	router.GET("/api/metrics/pods", ListPodMetrics(cfg))

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/metrics/pods", nil)
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	var resp models.PodMetricsResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	if resp.Count != 1 {
		t.Fatalf("expected 1 pod metric, got %d", resp.Count)
	}
	if resp.Pods[0].Containers[0].Usage["cpu"] != "50m" {
		t.Fatalf("expected cpu 50m, got %q", resp.Pods[0].Containers[0].Usage["cpu"])
	}
}
