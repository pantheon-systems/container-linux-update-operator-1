package k8sutil

import (
	"fmt"
	"os"

	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// isInCluster is the same heuristic used by rest.InClusterConfig() to determine whether it's an in-cluster execution environment.
func isInCluster() bool {
	return os.Getenv("KUBERNETES_SERVICE_HOST") != "" && os.Getenv("KUBERNETES_SERVICE_PORT") != ""
}

// kubernetesConfig loads a Kubernetes config using in-cluster configuration if it detects it's running inside the cluster.
// Otherwise it uses the default loading rules, such as the well-known path and the environment variable.
func kubernetesConfig() (*rest.Config, error) {
	var config *rest.Config
	var err error

	if isInCluster() {
		config, err = rest.InClusterConfig()
		if err != nil {
			return nil, err
		}
		return config, nil
	}

	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		loadingRules,
		&clientcmd.ConfigOverrides{})
	config, err = kubeConfig.ClientConfig()
	if err != nil {
		return nil, err
	}

	return config, nil
}

// GetClient returns a Kubernetes client (clientset) from the kubeconfig path
// or from the in-cluster service account environment.
func GetClient() (*kubernetes.Clientset, error) {
	conf, err := kubernetesConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get Kubernetes client config: %v", err)
	}
	return kubernetes.NewForConfig(conf)
}
