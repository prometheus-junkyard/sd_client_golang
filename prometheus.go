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
	"net"
	"net/http"
	"net/url"
	"time"
)

const (
	targetsUrl     = "/api/jobs/%s/targets"
	defaultTimeout = 30
)

type Client struct {
	url     string
	timeout time.Duration
}

type TargetGroup struct {
	// a set of labels
	BaseLabels map[string]string `json:"baseLabels"`

	// a group of endpoints
	Endpoints []string `json:"endpoints"`
}

func transport(netw, addr string, timeout time.Duration) (connection net.Conn, err error) {
	if timeout == 0 {
		timeout = defaultTimeout
	}
	deadline := time.Now().Add(timeout)
	connection, err = net.DialTimeout(netw, addr, timeout)
	if err == nil {
		connection.SetDeadline(deadline)
	}
	return
}

func New(url string) Client {
	return Client{url: url}
}

// marshal a list of target groups
func targetGroupsToJson(targetGroups []TargetGroup) (targetGroupsJson []byte, err error) {
	targetGroupsJson, err = json.Marshal(targetGroups)
	return
}

// http PUT to given url
func (client *Client) put(path string, data []byte) (response *http.Response, err error) {
	httpClient := http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) { return transport(netw, addr, client.timeout) },
		},
	}
	url := client.url + path
	request, err := http.NewRequest("PUT", url, bytes.NewReader(data))
	if err != nil {
		return
	}

	response, err = httpClient.Do(request)
	return
}

// replace the current list of target groups with the given new list
func (client *Client) UpdateEndpoints(job string, targetGroups []TargetGroup) (err error) {
	path := fmt.Sprintf(targetsUrl, url.QueryEscape(job))
	targetGroupsJson, err := targetGroupsToJson(targetGroups)
	_, err = client.put(path, targetGroupsJson)
	return
}

func (client *Client) SetTimeout(timeout time.Duration) {
	client.timeout = timeout
}
