package pkg_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/d2iq-labs/helm-list-images/pkg"
	"github.com/d2iq-labs/helm-list-images/pkg/k8s"
)

func Test_getImages(t *testing.T) {
	imageClient := pkg.Images{
		ImageRegex: pkg.ImageRegex,
	}
	imageClient.SetLogger("info")

	helmTemplate := `
---
# Source: prometheus/charts/prometheus/templates/alertmanager/clusterrole.yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
 labels:
   component: "alertmanager"
   app: prometheus
   release: prometheus-standalone
   chart: prometheus-14.4.1
   heritage: Helm
 name: prometheus-standalone-alertmanager
rules:
 []
---
# Source: prometheus/charts/prometheus/templates/pushgateway/clusterrole.yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
 labels:
   component: "pushgateway"
   app: prometheus
   release: prometheus-standalone
   chart: prometheus-14.4.1
   heritage: Helm
 name: prometheus-standalone-pushgateway
rules:
 []
---
# Source: prometheus/charts/prometheus/charts/kube-state-metrics/templates/clusterrolebinding.yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
 labels:
   app.kubernetes.io/name: kube-state-metrics
   helm.sh/chart: kube-state-metrics-3.1.1
   app.kubernetes.io/managed-by: Helm
   app.kubernetes.io/instance: prometheus-standalone
 name: prometheus-standalone-kube-state-metrics
roleRef:
 apiGroup: rbac.authorization.k8s.io
 kind: ClusterRole
 name: prometheus-standalone-kube-state-metrics
subjects:
- kind: ServiceAccount
 name: prometheus-standalone-kube-state-metrics
 namespace: test
---
# Source: tracing/templates/jaeger/configmap.yaml
apiVersion: v1
kind: ConfigMap
metadata:
 name: jaeger-ca-cert
data:
   CA_CERTIFICATE: |
       -----BEGIN CERTIFICATE-----
		OCOIRRGVEGHEIGHEnwoircne20394809234nfh834retitneh83t5ljfKHD&$&$
		OCOIRRGVEGHEIGHEnwoircne20394809234nfh834retitneh83t5ljfKHD&$&$
		OCOIRRGVEGHEIGHEnwoircne20394809234nfh834retitneh83t5ljfKHD&$&$
		OCOIRRGVEGHEIGHEnwoircne20394809234nfh834retitneh83t5ljfKHD&$&$
		OCOIRRGVEGHEIGHEnwoircne20394809234nfh834retitneh83t5ljfKHD&$&$
		OCOIRRGVEGHEIGHEnwoircne20394809234nfh834retitneh83t5ljfKHD&$&$
		OCOIRRGVEGHEIGHEnwoircne20394809234nfh834retitneh83t5ljfKHD&$&$
		OCOIRRGVEGHEIGHEnwoircne20394809234nfh834retitneh83t5ljfKHD&$&$
		OCOIRRGVEGHEIGHEnwoircne20394809234nfh834retitneh83t5ljfKHD&$&$
		OCOIRRGVEGHEIGHEnwoircne20394809234nfh834retitneh83t5ljfKHD&$&$
       -----END CERTIFICATE-----
`

	t.Run("should be able to split rendered templates to individual templates", func(t *testing.T) {
		//nolint:lll // Long YAML lines
		expected := []string{
			"\napiVersion: rbac.authorization.k8s.io/v1\nkind: ClusterRole\nmetadata:\n labels:\n   component: \"alertmanager\"\n   app: prometheus\n   release: prometheus-standalone\n   chart: prometheus-14.4.1\n   heritage: Helm\n name: prometheus-standalone-alertmanager\nrules:\n []\n",
			"\napiVersion: rbac.authorization.k8s.io/v1\nkind: ClusterRole\nmetadata:\n labels:\n   component: \"pushgateway\"\n   app: prometheus\n   release: prometheus-standalone\n   chart: prometheus-14.4.1\n   heritage: Helm\n name: prometheus-standalone-pushgateway\nrules:\n []\n",
			"\napiVersion: rbac.authorization.k8s.io/v1\nkind: ClusterRoleBinding\nmetadata:\n labels:\n   app.kubernetes.io/name: kube-state-metrics\n   helm.sh/chart: kube-state-metrics-3.1.1\n   app.kubernetes.io/managed-by: Helm\n   app.kubernetes.io/instance: prometheus-standalone\n name: prometheus-standalone-kube-state-metrics\nroleRef:\n apiGroup: rbac.authorization.k8s.io\n kind: ClusterRole\n name: prometheus-standalone-kube-state-metrics\nsubjects:\n- kind: ServiceAccount\n name: prometheus-standalone-kube-state-metrics\n namespace: test\n",
			"\napiVersion: v1\nkind: ConfigMap\nmetadata:\n name: jaeger-ca-cert\ndata:\n   CA_CERTIFICATE: |\n       -----BEGIN CERTIFICATE-----\n\t\tOCOIRRGVEGHEIGHEnwoircne20394809234nfh834retitneh83t5ljfKHD&$&$\n\t\tOCOIRRGVEGHEIGHEnwoircne20394809234nfh834retitneh83t5ljfKHD&$&$\n\t\tOCOIRRGVEGHEIGHEnwoircne20394809234nfh834retitneh83t5ljfKHD&$&$\n\t\tOCOIRRGVEGHEIGHEnwoircne20394809234nfh834retitneh83t5ljfKHD&$&$\n\t\tOCOIRRGVEGHEIGHEnwoircne20394809234nfh834retitneh83t5ljfKHD&$&$\n\t\tOCOIRRGVEGHEIGHEnwoircne20394809234nfh834retitneh83t5ljfKHD&$&$\n\t\tOCOIRRGVEGHEIGHEnwoircne20394809234nfh834retitneh83t5ljfKHD&$&$\n\t\tOCOIRRGVEGHEIGHEnwoircne20394809234nfh834retitneh83t5ljfKHD&$&$\n\t\tOCOIRRGVEGHEIGHEnwoircne20394809234nfh834retitneh83t5ljfKHD&$&$\n\t\tOCOIRRGVEGHEIGHEnwoircne20394809234nfh834retitneh83t5ljfKHD&$&$\n       -----END CERTIFICATE-----\n",
		}
		actual := imageClient.GetTemplates([]byte(helmTemplate))
		assert.Equal(t, len(expected), len(actual))
		assert.ElementsMatch(t, expected, actual)
	})
}

func Test_getImagesFromKind(t *testing.T) {
	t.Run("should be able to get all images from struct kind", func(t *testing.T) {
		kindObj := []*k8s.Image{
			{
				Kind:  "DaemonSet",
				Name:  "prometheus-standalone-node-exporter",
				Image: []string{"quay.io/prometheus/node-exporter:v1.1.2"},
			},
			{
				Kind:  "Deployment",
				Name:  "prometheus-standalone-server",
				Image: []string{"jimmidyson/configmap-reload:v0.5.0"},
			},
			{
				Kind:  "StatefulSet",
				Name:  "prometheus-standalone-kube-state-metrics",
				Image: []string{"k8s.gcr.io/kube-state-metrics/kube-state-metrics:v2.0.0"},
			},
		}

		expected := []string{
			"quay.io/prometheus/node-exporter:v1.1.2",
			"jimmidyson/configmap-reload:v0.5.0",
			"k8s.gcr.io/kube-state-metrics/kube-state-metrics:v2.0.0",
		}
		images := pkg.GetImagesFromKind(kindObj)
		assert.ElementsMatch(t, expected, images)
	})
}

func TestImages_SetRelease(t *testing.T) {
	t.Run("Should be able to set the release", func(t *testing.T) {
		imageClient := pkg.Images{}
		imageClient.SetRelease("testRelease")

		assert.Equal(t, imageClient.GetRelease(), "testRelease")
	})
}

func TestImages_SetChart(t *testing.T) {
	t.Run("Should be able to set the chart", func(t *testing.T) {
		imageClient := pkg.Images{}
		imageClient.SetChart("testChart")

		assert.Equal(t, imageClient.GetChart(), "testChart")
	})
}

func TestImages_SetWriter(t *testing.T) {
	t.Run("Should be able to set the writer", func(t *testing.T) {
		imageClient := pkg.Images{}
		imageClient.SetWriter(os.Stdout)

		assert.NotNil(t, imageClient.GetWriter())
	})
}
