package cmd

import (
	"context"
	"os"
	"sort"
	"strings"

	"code.hq.twilio.com/redman/kubectl-tks/pkg/client"
	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/printers"
)

var (
	kubeconfig string
)

func init() {
	getCmd.AddCommand(getNodesCmd)
	getNodesCmd.PersistentFlags().StringVarP(&kubeconfig, "kubeconfig", "c", "$HOME/.kube/config", "path to kube config file")
}

// ByRole fulfills the package sort interface
// so we can sort by custom values from the v1.Node
// that is returned
type ByRole []v1.Node

func (r ByRole) Len() int      { return len(r) }
func (r ByRole) Swap(i, j int) { r[i], r[j] = r[j], r[i] }

// Less if the strings are already sorted it would mean
// the first value passed in is "lesser" on the alphabet
// making this true
func (r ByRole) Less(i, j int) bool {
	return sort.StringsAreSorted([]string{
		r[i].GetLabels()["twilio.com/role"],
		r[j].GetLabels()["twilio.com/role"],
	})
}

var getNodesCmd = &cobra.Command{
	Use:   "nodes",
	Short: "retrieve cluster nodes",
	RunE: func(c *cobra.Command, args []string) error {
		clientset := client.InitClient()
		// get all nodes
		nodes, err := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			return err
		}

		// sort the nodes by role
		sort.Sort(ByRole(nodes.Items))

		// setup table schema
		t := &metav1.Table{
			ColumnDefinitions: []metav1.TableColumnDefinition{
				{Name: "Role", Type: "string"},
				{Name: "Name", Type: "string"},
				{Name: "SID", Type: "string"},
				{Name: "Status", Type: "string"},
			},
		}
		// add rows
		for _, n := range nodes.Items {
			// deduce status
			status := ""
			for _, c := range n.Status.Conditions {
				if (c.Status == "True") && (c.Type == "Ready") {
					status = string(c.Type)
				}
				if (c.Status == "True") && (c.Type != "Ready") {
					status = strings.Join([]string{status, string(c.Type)}, ",")
				}
			}
			// create row
			t.Rows = append(t.Rows, metav1.TableRow{
				Cells: []interface{}{
					n.GetLabels()["twilio.com/role"],
					n.Name,
					n.GetLabels()["twilio.com/host-sid"],
					status,
				},
			})
		}

		// print
		p := printers.NewTablePrinter(printers.PrintOptions{})
		p.PrintObj(t, os.Stdout)

		return nil
	},
}
