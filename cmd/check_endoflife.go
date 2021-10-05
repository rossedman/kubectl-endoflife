package cmd

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"code.hq.twilio.com/platform-base/kubectl-check/pkg/client"
	"code.hq.twilio.com/platform-base/kubectl-check/pkg/endoflife"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/printers"
)

var (
	product     string
	silent      bool
	expiryRange int
)

func init() {
	rootCmd.AddCommand(endOfLifeCmd)
	endOfLifeCmd.PersistentFlags().StringVarP(&product, "product", "p", "kubernetes", "the product to lookup, supported values: kubernetes, amazon-eks")
	endOfLifeCmd.PersistentFlags().BoolVarP(&silent, "silent", "s", false, "silence the output, only provide exit codes")
	endOfLifeCmd.PersistentFlags().IntVarP(&expiryRange, "expiry-range", "e", 0, "set a range which the command should exit 1, this is days within the expiration date")
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

		// get cluster version from current context
		current, _, err := GetClusterVersion()
		if err != nil {
			return err
		}

		// get current version data for EKS
		client := endoflife.NewClient(endoflife.BaseURL, &http.Client{
			Timeout: time.Second * 2,
		})

		// determine the product type
		prod, err := endoflife.GetProduct(product)
		if err != nil {
			return err
		}

		// retrieve the endoflife data
		resp, err := client.Get(prod, current)
		if err != nil {
			return err
		}

		// get date from response
		date, err := resp.ConvertToTime()
		if err != nil {
			return err
		}

		// calculate how many days left
		days, err := endoflife.GetDaysUntilEnd(date)
		if err != nil {
			return err
		}

		// check if within threshold
		threshold, err := endoflife.InExpiryRange(date, expiryRange)
		if err != nil {
			return err
		}

		// check if expired
		expired, err := endoflife.IsExpired(date)
		if err != nil {
			return err
		}

		// print if silent not set
		if !silent {
			// append to table output
			t.Rows = append(t.Rows, metav1.TableRow{
				Cells: []interface{}{prod.String(), current, resp.EOL, days},
			})
			p := printers.NewTablePrinter(printers.PrintOptions{})
			p.PrintObj(t, os.Stdout)
		}

		if expired || threshold {
			os.Exit(1)
		}

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
