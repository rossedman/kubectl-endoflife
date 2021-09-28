package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"code.hq.twilio.com/redman/kubectl-tks/pkg/client"
	"code.hq.twilio.com/redman/kubectl-tks/pkg/endoflife"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/printers"
)

var (
	getEKS        bool
	getKubernetes bool
)

func init() {
	rootCmd.AddCommand(endOfLifeCmd)
	endOfLifeCmd.PersistentFlags().BoolVarP(&getEKS, "eks", "e", false, "retrieve EKS endoflife data")
	endOfLifeCmd.PersistentFlags().BoolVarP(&getKubernetes, "kubernetes", "k", true, "retrieve Kubernetes endoflife data")
}

// AmazonEKSRelease represents the data found at
// the endoflife.date endpoint: https://endoflife.date/api/amazon-eks.json
type AmazonEKSRelease struct {
	EOL     string `json:"eol"`
	Release string `json:"release"`
	Latest  string `json:"latest"`
}

// KubernetesRelease ...
type KubernetesRelease struct {
	EOL     string `json:"eol"`
	Release string `json:"release"`
	Latest  string `json:"latest"`
}

var endOfLifeCmd = &cobra.Command{
	Use:   "endoflife",
	Short: "checks end of life date for cluster",
	RunE: func(c *cobra.Command, args []string) error {

		// setup table output
		t := &metav1.Table{
			ColumnDefinitions: []metav1.TableColumnDefinition{
				{Name: "Type", Type: "string"},
				{Name: "Version", Type: "string"},
				{Name: "EOL Date", Type: "string"},
				{Name: "Days Left", Type: "string"},
			},
		}

		current, _, err := GetClusterVersion()
		if err != nil {
			return err
		}

		// get current version data for EKS
		client := endoflife.NewClient()

		// get EKS versions
		if getEKS {
			eks, err := client.GetAmazonEKS(current)
			if err != nil {
				return err
			}

			days, err := eks.GetDaysUntilEnd()
			if err != nil {
				return err
			}

			t.Rows = append(t.Rows, metav1.TableRow{
				Cells: []interface{}{"EKS", current, eks.EOL, days},
			})
		}

		// get upstream kubernetes versions
		if getKubernetes {
			k8s, err := client.GetKubernetes(current)
			if err != nil {
				return err
			}

			days, err := k8s.GetDaysUntilEnd()
			if err != nil {
				return err
			}

			t.Rows = append(t.Rows, metav1.TableRow{
				Cells: []interface{}{"Kubernetes", current, k8s.EOL, days},
			})
		}

		// print
		p := printers.NewTablePrinter(printers.PrintOptions{})
		p.PrintObj(t, os.Stdout)
		return nil
	},
}

func GetClusterVersion() (string, string, error) {
	// get cluster version
	clientset := client.InitClient()
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
