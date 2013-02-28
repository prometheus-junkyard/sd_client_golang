package prometheus_test

import (
	"github.com/prometheus/sd_client_golang"
	"log"
	"time"
)

func ExampleClient_UpdateEndpoints() {
	timeout, err := time.ParseDuration("2s")
	if err != nil {
		log.Fatalf("Can't parse timeout")
	}
	client := prometheus.New("http://localhost:8080")
	client.SetTimeout(timeout)
	err = client.UpdateEndpoints("job-name", []prometheus.TargetGroup{{
		BaseLabels: map[string]string{"label1": "value1", "label2": "value2"},
		Endpoints: []string{
			"http://example.com:8080/metrics.json",
			"http://example.com:8081/metrics.json",
		},
	}, {
		BaseLabels: map[string]string{"label3": "value3"},
		Endpoints: []string{
			"http://example.com:8082/metrics.json",
			"http://example.com:8083/metrics.json",
		},
	}})

	if err != nil {
		log.Fatalf("Error updating endpoints: %s", err)
	}
	// Output:
}
