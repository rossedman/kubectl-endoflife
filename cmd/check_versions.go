package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"code.hq.twilio.com/redman/kubectl-tks/pkg/client"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/printers"
	"k8s.io/client-go/kubernetes"

	_ "embed"

	"github.com/blang/semver/v4"
)

var (
	kubeVersion string
)

func init() {
	rootCmd.AddCommand(getVersionsCmd)
	getVersionsCmd.PersistentFlags().StringVarP(&kubeVersion, "kube-version", "k", "v1.19", "the version of dependencies to check against")
}

type Service struct {
	Name    string
	Version string
}

type Services []Service

// KubernetesVersions represents the Kubernetes
// version with a mapping of services and versions
// that can be used for reference
type KubernetesVersions map[string]Services

//go:embed config/components.json
var componentsConfig []byte

var getVersionsCmd = &cobra.Command{
	Use:   "versions",
	Short: "checks compatible versions for upgrade",
	RunE: func(c *cobra.Command, args []string) error {
		clientset := client.InitClient()

		// setup table output
		t := &metav1.Table{
			ColumnDefinitions: []metav1.TableColumnDefinition{
				{Name: "Service", Type: "string"},
				{Name: "Out Of Date", Type: "bool"},
				{Name: "Current Version", Type: "string"},
				{Name: "Required Version", Type: "string"},
			},
		}

		// get all the current versions of deployments running in kube-system
		// 	TODO(redman): make namespaces configurable through flags
		svcs, err := getAllServices(clientset, []string{"kube-system", "platform", "cert-manager"})
		if err != nil {
			return err
		}

		// load the data from the embedded config
		//	TODO(redman): allow override of embedded config for custom configuration
		versions, err := loadKubernetesVersions()
		if err != nil {
			return err
		}

		// loop through services
		for _, s := range svcs {
			// walk through required versions and
			// find the matching service
			req, err := getRequiredVersion(versions, kubeVersion, s.Name)
			if err != nil {
				return fmt.Errorf("failed to get required version for %s: %s", s.Name, err)
			}
			if req == "unknown" {
				continue
			}

			// if is latest, depend and set to false
			if s.Version == "latest" {
				t.Rows = append(t.Rows, metav1.TableRow{
					Cells: []interface{}{s.Name, false, s.Version, req},
				})
				continue
			}

			// determine if component is out of date
			o, err := isOutOfDate(req, s.Version)
			if err != nil {
				return fmt.Errorf("failed to calculate out of date for %s: %s", s.Name, err)
			}

			t.Rows = append(t.Rows, metav1.TableRow{
				Cells: []interface{}{s.Name, o, s.Version, req},
			})
		}

		// print
		p := printers.NewTablePrinter(printers.PrintOptions{})
		p.PrintObj(t, os.Stdout)
		return nil
	},
}

// isOutOfDate takes two semantic version numbers and compares
// them to determine whether the current version is above or
// equal to the required version.
func isOutOfDate(required, current string) (bool, error) {
	r, err := semver.Make(strings.ReplaceAll(required, "v", ""))
	if err != nil {
		return false, err
	}

	c, err := semver.Make(strings.ReplaceAll(current, "v", ""))
	if err != nil {
		return false, err
	}

	return r.Compare(c) > 0, nil
}

func getRequiredVersion(versions KubernetesVersions, kubernetesVersion string, serviceName string) (string, error) {
	for _, e := range versions[kubernetesVersion] {
		if e.Name == serviceName {
			return e.Version, nil
		}
	}
	return "unknown", nil
}

func getDeployments(clientset *kubernetes.Clientset, namespace string) (Services, error) {
	deploys, err := clientset.AppsV1().Deployments(namespace).List(context.TODO(), v1.ListOptions{})
	if err != nil {
		return nil, err
	}
	svcs := Services{}
	for _, d := range deploys.Items {
		version := strings.Split(d.Spec.Template.Spec.Containers[0].Image, ":")[1]
		svcs = append(svcs, Service{Name: d.Name, Version: version})
	}
	return svcs, nil
}

func getDaemonSets(clientset *kubernetes.Clientset, namespace string) (Services, error) {
	deploys, err := clientset.AppsV1().DaemonSets(namespace).List(context.TODO(), v1.ListOptions{})
	if err != nil {
		return nil, err
	}
	svcs := Services{}
	for _, d := range deploys.Items {
		version := strings.Split(d.Spec.Template.Spec.Containers[0].Image, ":")[1]
		svcs = append(svcs, Service{Name: d.Name, Version: version})
	}
	return svcs, nil
}

func getAllServices(clientset *kubernetes.Clientset, namespaces []string) (Services, error) {
	svcs := Services{}
	for _, ns := range namespaces {
		deploys, err := getDeployments(clientset, ns)
		if err != nil {
			return nil, err
		}
		svcs = append(svcs, deploys...)

		daemons, err := getDaemonSets(clientset, ns)
		if err != nil {
			return nil, err
		}
		svcs = append(svcs, daemons...)
	}
	return svcs, nil
}

// loadKubernetesVersions take the config that is embedded and unmarshals
// it into a a struct that can be worked with
func loadKubernetesVersions() (KubernetesVersions, error) {
	var versions KubernetesVersions
	if err := json.Unmarshal(componentsConfig, &versions); err != nil {
		return nil, fmt.Errorf("unable to read components configuration: %w", err)
	}
	return versions, nil
}
