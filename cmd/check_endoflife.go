package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
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
	output      string
)

func init() {
	rootCmd.AddCommand(endOfLifeCmd)
	endOfLifeCmd.PersistentFlags().StringVarP(&product, "product", "p", "kubernetes", "the product to lookup, supported values: kubernetes, amazon-eks")
	endOfLifeCmd.PersistentFlags().BoolVarP(&silent, "silent", "s", false, "silence the output, only provide exit codes")
	endOfLifeCmd.PersistentFlags().IntVarP(&expiryRange, "expiry-range", "e", 0, "set a range which the command should exit 1, this is days within the expiration date")
	endOfLifeCmd.PersistentFlags().StringVarP(&output, "output", "o", "table", "set output, supported values: table, json")
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
		current, _, err := client.GetClusterVersion()
		if err != nil {
			return err
		}

		// create endoflife client
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

		// print table output
		if !silent && output == "table" {
			// append to table output
			t.Rows = append(t.Rows, metav1.TableRow{
				Cells: []interface{}{prod.String(), current, resp.EOL, days},
			})
			p := printers.NewTablePrinter(printers.PrintOptions{})
			p.PrintObj(t, os.Stdout)
		}

		// print json output
		if !silent && output == "json" {
			output := struct {
				Type     string  `json:"type"`
				Version  string  `json:"version"`
				EOL      string  `json:"eol-date"`
				DaysLeft float64 `json:"days-left"`
			}{
				Type:     prod.String(),
				Version:  current,
				EOL:      resp.EOL,
				DaysLeft: days,
			}
			j, err := json.MarshalIndent(&output, "", "  ")
			if err != nil {
				return err
			}
			fmt.Println(string(j))
		}

		if expired || threshold {
			os.Exit(1)
		}

		return nil
	},
}
