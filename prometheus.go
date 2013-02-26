//	Package prometheus implements a client library to manage Prometheus server job state.
//
//	To update the list of endpoints for a particular job:
//
//		client := prometheus.Client{Url: "http://localhost:8080"}
//		err := client.UpdateEndpoints("job-name", []prometheus.TargetGroup{{
//			BaseLabels: map[string]string{"label1": "value1", "label2": "value2"},
//			Endpoints: []string{
//				"http://example.com:8080/metrics.json",
//				"http://example.com:8081/metrics.json",
//			},
//		}, {
//			BaseLabels: map[string]string{"label3": "value3"},
//			Endpoints: []string{
//				"http://example.com:8082/metrics.json",
//				"http://example.com:8083/metrics.json",
//			},
//		}})
package prometheus

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

const (
	EndpointsUrl   = "/api/jobs/%s/endpoints"
	DefaultTimeout = 30
)

type Client struct {
	Url     string
	Timeout time.Duration
}

type TargetGroup struct {
	// a set of labels
	BaseLabels map[string]string `json:"baseLabels"`

	// a group of endpoints
	Endpoints []string `json:"endpoints"`
}

// http PUT to given url
func put(url string, data []byte, timeout time.Duration) (response *http.Response, err error) {
	client := &http.Client{}
	request, err := http.NewRequest("PUT", url, bytes.NewReader(data))
	if err != nil {
		return
	}

	ready := make(chan bool)
	go func() {
		response, err = client.Do(request)
		ready <- true
	}()

	select {
	case <-ready:
		break
	case <-time.After(time.Second * timeout):
		err = fmt.Errorf("timeout reached while trying to update %s", url)
		return
	}
	return
}

// marshal a list of endpoints
func targetGroupsToJson(targetGroups []TargetGroup) (targetGroupsJson []byte, err error) {
	targetGroupsJson, err = json.Marshal(targetGroups)
	return
}

// replace the current list of endpoints with the given new list
func (client *Client) UpdateEndpoints(job string, targetGroups []TargetGroup) (err error) {
	url := fmt.Sprintf(client.Url+EndpointsUrl, url.QueryEscape(job))
	timeout := client.Timeout
	if timeout == 0 {
		timeout = DefaultTimeout
	}

	targetGroupsJson, err := targetGroupsToJson(targetGroups)
	_, err = put(url, targetGroupsJson, timeout)
	return
}
