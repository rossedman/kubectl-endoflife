package client

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	// utilities for kubernetes integration
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

// ClientConfig
func ClientConfig() clientcmd.ClientConfig {
	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		&clientcmd.ConfigOverrides{})
}

// InitClient - Kubernetes Client
func InitClient() *kubernetes.Clientset {
	clientConfig := ClientConfig()
	config, err := clientConfig.ClientConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}

	return kubernetes.NewForConfigOrDie(config)
}

// GetClusterVersion will return the current cluster version
// as well as the next incremental version
func GetClusterVersion() (string, string, error) {
	// get cluster version
	clientset := InitClient()
	ver, err := clientset.ServerVersion()
	if err != nil {
		return "", "", err
	}

	// determine cluster version for endoflife.date
	minor, err := strconv.Atoi(strings.ReplaceAll(ver.Minor, "+", ""))
	if err != nil {
		return "", "", err
	}

	current := fmt.Sprintf("%s.%s", ver.Major, strconv.Itoa(minor))
	next := fmt.Sprintf("%s.%s", ver.Major, strconv.Itoa(minor+1))

	return current, next, nil
}
