package endoflife

import (
	"fmt"
	"time"
)

const (
	KubernetesProduct Product = iota + 1
	AmazonEKSProduct
)

// Product represents a product that is
// available in endoflife for querying
type Product int64

// String produces the name of the product
// note: this currently matches the slug
// used in endoflife and is used to construct
// requests
func (p Product) String() string {
	switch p {
	case KubernetesProduct:
		return "kubernetes"
	case AmazonEKSProduct:
		return "amazon-eks"
	default:
		return "unknown"
	}
}

// GetProduct retrieves a product using the string
// based representation of it
func GetProduct(product string) (Product, error) {
	switch product {
	case "kubernetes":
		return KubernetesProduct, nil
	case "amazon-eks":
		return AmazonEKSProduct, nil
	default:
		return 0, fmt.Errorf("no product with name %s", product)
	}
}

// ProductResponse represents what is returns from the
// the endoflife API once unmarshalled. This only captures
// the endoflife data
type ProductResponse struct {
	EOL string `json:"eol"`
}

// ConvertToTime provides the time correctly formatted
// and parsed to use for comparisons
func (p *ProductResponse) ConvertToTime() (time.Time, error) {
	return time.Parse("2006-01-02", p.EOL)
}
