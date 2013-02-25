PACKAGE

package prometheus
    import "prometheus"

    Package prometheus implements the prometheus client library.

    To update the list of endpoints for a particular job:

		client := prometheus.Client{"http://url/to/prometheus"}
	  response, err := client.UpdateEndpoints("job-name", []prometheus.Endpoint{{
	    BaseLabels: map[string]string{"label1": "value1", "label2": "value2"},
	    Endpoints:  []string{"http://example.com:8080/metrics.json", "http://example.com:8081/metrics.json"},
	  }, {
	    BaseLabels: map[string]string{"label3": "value3"},
	    Endpoints:  []string{"http://example.com:8082/metrics.json", "http://example.com:8083/metrics.json"},


CONSTANTS

const (
    ENDPOINTS_URL = "/api/jobs/%s/endpoints"
)


FUNCTIONS

func EndpointsToJson(endpoints []Endpoint) ([]byte, error)
    marshal a list of endpoints


func Put(url string, data []byte) (*http.Response, error)
    http PUT to given url



TYPES

type Client struct {
    Url string
}


func (client *Client) UpdateEndpoints(job string, endpoints []Endpoint) (*http.Response, error)
    replace the current list of endpoints with the given new list


type Endpoint struct {
    // a set of labels
    BaseLabels map[string]string `json:"baseLabels"`

    // a group of endpoints
    Endpoints []string `json:"endpoints"`
}



