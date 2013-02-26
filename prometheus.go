//	Package prometheus implements the prometheus client library.
//
//	To update the list of endpoints for a particular job:
//
//		client := prometheus.Client{"http://url/to/prometheus"}
//	  response, err := client.UpdateEndpoints("job-name", []prometheus.Endpoint{{
//	    BaseLabels: map[string]string{"label1": "value1", "label2": "value2"},
//	    Endpoints:  []string{"http://example.com:8080/metrics.json", "http://example.com:8081/metrics.json"},
//	  }, {
//	    BaseLabels: map[string]string{"label3": "value3"},
//	    Endpoints:  []string{"http://example.com:8082/metrics.json", "http://example.com:8083/metrics.json"},
package prometheus

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const (
	EndpointsUrl = "/api/jobs/%s/endpoints"
)

type Client struct {
	Url string
}

type TargetGroup struct {
	// a set of labels
	BaseLabels map[string]string `json:"baseLabels"`

	// a group of endpoints
	Endpoints []string `json:"endpoints"`
}

// http PUT to given url
func Put(url string, data []byte) (*http.Response, error) {
	client := &http.Client{}
	request, err := http.NewRequest("PUT", url, bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("Could not create request to %s: %s", url, err)
	}

	response, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("Could not send request to %s: %s", url, err)
	}
	return response, nil
}

// marshal a list of endpoints
func targetGroupsToJson(targetGroups []TargetGroup) ([]byte, error) {
	targetGroupsJson, err := json.Marshal(targetGroups)

	if err != nil {
		return nil, fmt.Errorf("Could not marshal data: %s", err)
	}
	return targetGroupsJson, nil
}

// replace the current list of endpoints with the given new list
func (client *Client) UpdateEndpoints(job string, targetGroups []TargetGroup) (*http.Response, error) {
	url := fmt.Sprintf(client.Url+EndpointsUrl, url.QueryEscape(job))

	targetGroupsJson, err := targetGroupsToJson(targetGroups)
	if err != nil {
		return nil, err
	}
	return Put(url, targetGroupsJson)
}
