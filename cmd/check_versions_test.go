package cmd

import "testing"

func TestGetRequiredVersions(t *testing.T) {
	cases := Services{
		{
			Name:    "cert-manager",
			Version: "1.4.0",
		},
		{
			Name:    "coredns",
			Version: "1.8.4",
		},
		{
			Name:    "kube-proxy",
			Version: "1.19.6-eksbuild.2",
		},
		{
			Name:    "kube-state-metrics",
			Version: "2.1.0",
		},
		{
			Name:    "metrics-server",
			Version: "0.5.0",
		},
		{
			Name:    "node-problem-detector",
			Version: "0.8.9",
		},
		{
			Name:    "nvidia-device-plugin",
			Version: "0.9.0",
		},
		{
			Name:    "cluster-autoscaler",
			Version: "1.19.0",
		},
	}

	versions, err := loadKubernetesVersions()
	if err != nil {
		t.Fatal(err)
	}

	for _, test := range cases {
		actual, err := getRequiredVersion(versions, "v1.19", test.Name)
		if err != nil {
			t.Fatal(err)
		}
		assertVersion(t, test.Name, test.Version, actual)
	}
}

func assertVersion(t *testing.T, component, expected, actual string) {
	t.Helper()
	if expected != actual {
		t.Fatalf("expected %s version %s, got %s", component, expected, actual)
	}
}
