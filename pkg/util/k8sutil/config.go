package k8sutil

import (
	"io/ioutil"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"net/url"
	"os"
	"path/filepath"
)

var (
	KubeConfigLocation string
	DefaultConfig      *rest.Config
)

const (
	kubeconfigDir  = ".kube"
	kubeconfigFile = "config"
)

type promptedCredentials struct {
	username string
	password string
}

func GenerateKubernetesConfig() (*rest.Config, error) {
	if KubeConfigLocation != "" {
		return clientcmd.BuildConfigFromFlags("", KubeConfigLocation)
	} else {
		configPath := filepath.Join(homeDir(), kubeconfigDir, kubeconfigFile)
		_, err := os.Stat(configPath)
		if err == nil {
			return clientcmd.BuildConfigFromFlags("", configPath)
		} else {
			return rest.InClusterConfig()
		}
	}
}

func GetClusterEndpoint(kubeconfig string) (string, error) {
	// Let's create a tempfile and line it up for removal
	file, err := ioutil.TempFile(os.TempDir(), "kraken-kubeconfig")
	defer func() {
		_ = os.Remove(file.Name())
	}()

	if err != nil {
		return "", err
	}
	_, err = file.WriteString(kubeconfig)
	if err != nil {
		return "", err
	}

	clusterConfig, err := clientcmd.BuildConfigFromFlags("", file.Name())
	if err != nil {
		return "", err
	}

	hostUrl, err := url.Parse(clusterConfig.Host)
	if err != nil {
		return "", err
	}

	return hostUrl.Hostname(), nil
}
