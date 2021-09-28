package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"code.hq.twilio.com/redman/kubectl-tks/pkg/client"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/printers"
)

func init() {
	checkCmd.AddCommand(endOfLifeCmd)
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

		// get cluster version
		clientset := client.InitClient()
		ver, err := clientset.ServerVersion()
		if err != nil {
			return err
		}

		// determine cluster version for endoflife.date
		minor, err := strconv.Atoi(strings.ReplaceAll(ver.Minor, "+", ""))
		if err != nil {
			return err
		}

		version := fmt.Sprintf("%s.%s", ver.Major, strconv.Itoa(minor))
		versionNext := fmt.Sprintf("%s.%s", ver.Major, strconv.Itoa(minor+1))

		log.Println(version, versionNext)

		// get current version data for EKS
		eksDaysLeft, eksEOL, err := GetEKSFormatted(version)
		if err != nil {
			return err
		}
		t.Rows = append(t.Rows, metav1.TableRow{
			Cells: []interface{}{"EKS", version, eksEOL, eksDaysLeft},
		})

		// get next version data for EKS
		eksNextDaysLeft, eksNextEOL, err := GetEKSFormatted(versionNext)
		if err != nil {
			return err
		}
		t.Rows = append(t.Rows, metav1.TableRow{
			Cells: []interface{}{"EKS", versionNext, eksNextEOL, eksNextDaysLeft},
		})

		// get current version data for kubernetes
		k8sDaysLeft, k8sEOL, err := GetKubernetesFormatted(version)
		if err != nil {
			return err
		}
		t.Rows = append(t.Rows, metav1.TableRow{
			Cells: []interface{}{"Kubernetes", version, k8sEOL, k8sDaysLeft},
		})

		// get next version data for Kubernetes
		k8sNextDaysLeft, k8sNextEOL, err := GetKubernetesFormatted(versionNext)
		if err != nil {
			return err
		}
		t.Rows = append(t.Rows, metav1.TableRow{
			Cells: []interface{}{"Kubernetes", versionNext, k8sNextEOL, k8sNextDaysLeft},
		})

		// print
		p := printers.NewTablePrinter(printers.PrintOptions{})
		p.PrintObj(t, os.Stdout)
		return nil
	},
}

// GetEKSFormatted returns the days left and the eol date formatted
// for print output
func GetEKSFormatted(version string) (float64, string, error) {
	eks, err := GetEKSVersion(version)
	if err != nil {
		return 0, "", err
	}

	// calculate date difference
	today := time.Now()
	endDate, err := time.Parse("2006-01-02", eks.EOL)
	if err != nil {
		return 0, "", err
	}

	return math.Round(endDate.Sub(today).Hours() / 24), eks.EOL, nil
}

// GetEKSVersion retrieves the EOL data for a single release
func GetEKSVersion(version string) (AmazonEKSRelease, error) {
	eks := AmazonEKSRelease{}

	client := http.Client{
		Timeout: time.Second * 2,
	}

	url := fmt.Sprintf("https://endoflife.date/api/amazon-eks/%s.json", version)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return eks, err
	}

	res, err := client.Do(req)
	if err != nil {
		return eks, err
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return eks, err
	}

	err = json.Unmarshal(body, &eks)
	if err != nil {
		return eks, err
	}

	return eks, nil
}

// GetKubernetesFormatted returns the days left and the eol date formatted
// for print output
func GetKubernetesFormatted(version string) (float64, string, error) {
	eks, err := GetKubernetesVersion(version)
	if err != nil {
		return 0, "", err
	}

	// calculate date difference
	today := time.Now()
	endDate, err := time.Parse("2006-01-02", eks.EOL)
	if err != nil {
		return 0, "", err
	}

	return math.Round(endDate.Sub(today).Hours() / 24), eks.EOL, nil
}

// GetKubernetesVersion retrieves the EOL data for a single release
func GetKubernetesVersion(version string) (AmazonEKSRelease, error) {
	eks := AmazonEKSRelease{}

	client := http.Client{
		Timeout: time.Second * 2,
	}

	url := fmt.Sprintf("https://endoflife.date/api/kubernetes/%s.json", version)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return eks, err
	}

	res, err := client.Do(req)
	if err != nil {
		return eks, err
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return eks, err
	}

	err = json.Unmarshal(body, &eks)
	if err != nil {
		return eks, err
	}

	return eks, nil
}
