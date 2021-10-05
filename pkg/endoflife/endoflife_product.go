package endoflife

const (
	KubernetesProduct Product = iota
	AmazonEKSProduct
)

type Product int64

func (p Product) String() string {
	switch p {
	case KubernetesProduct:
		return "kubernetes"
	case AmazonEKSProduct:
		return "amazon-eks"
	}
	return "unknown"
}
