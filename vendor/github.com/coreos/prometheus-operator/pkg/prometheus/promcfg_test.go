// Copyright 2017 The prometheus-operator Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package prometheus

import (
	"bytes"
	"testing"

	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/pkg/api/v1"

	"github.com/coreos/prometheus-operator/pkg/client/monitoring/v1alpha1"
)

func TestConfigGeneration(t *testing.T) {
	cfg, err := generateTestConfig()
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 1000; i++ {
		testcfg, err := generateTestConfig()
		if err != nil {
			t.Fatal(err)
		}

		if !bytes.Equal(cfg, testcfg) {
			t.Fatalf("Config generation is not deterministic.\n\n\nFirst generation: \n\n%s\n\nDifferent generation: \n\n%s\n\n", string(cfg), string(testcfg))
		}
	}
}

func generateTestConfig() ([]byte, error) {
	replicas := int32(1)
	return generateConfig(
		&v1alpha1.Prometheus{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test",
				Namespace: "default",
			},
			Spec: v1alpha1.PrometheusSpec{
				Alerting: v1alpha1.AlertingSpec{
					Alertmanagers: []v1alpha1.AlertmanagerEndpoints{
						{
							Name:      "alertmanager-main",
							Namespace: "default",
							Port:      intstr.FromString("web"),
						},
					},
				},
				Replicas: &replicas,
				ServiceMonitorSelector: &metav1.LabelSelector{
					MatchLabels: map[string]string{
						"group": "group1",
					},
				},
				RuleSelector: &metav1.LabelSelector{
					MatchLabels: map[string]string{
						"role": "rulefile",
					},
				},
				Resources: v1.ResourceRequirements{
					Requests: v1.ResourceList{
						v1.ResourceMemory: resource.MustParse("400Mi"),
					},
				},
			},
		},
		makeServiceMonitors(),
		1,
		map[string]BasicAuthCredentials{},
	)
}

func makeServiceMonitors() map[string]*v1alpha1.ServiceMonitor {
	res := map[string]*v1alpha1.ServiceMonitor{}

	res["servicemonitor1"] = &v1alpha1.ServiceMonitor{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "testservicemonitor1",
			Namespace: "default",
			Labels: map[string]string{
				"group": "group1",
			},
		},
		Spec: v1alpha1.ServiceMonitorSpec{
			Selector: metav1.LabelSelector{
				MatchLabels: map[string]string{
					"group": "group1",
				},
			},
			Endpoints: []v1alpha1.Endpoint{
				v1alpha1.Endpoint{
					Port:     "web",
					Interval: "30s",
				},
			},
		},
	}

	res["servicemonitor2"] = &v1alpha1.ServiceMonitor{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "testservicemonitor2",
			Namespace: "default",
			Labels: map[string]string{
				"group": "group2",
			},
		},
		Spec: v1alpha1.ServiceMonitorSpec{
			Selector: metav1.LabelSelector{
				MatchLabels: map[string]string{
					"group":  "group2",
					"group3": "group3",
				},
			},
			Endpoints: []v1alpha1.Endpoint{
				v1alpha1.Endpoint{
					Port:     "web",
					Interval: "30s",
				},
			},
		},
	}

	return res
}
