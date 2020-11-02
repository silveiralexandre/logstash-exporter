package tower

import (
	"fmt"

	jsoniter "github.com/json-iterator/go"
	"github.com/silveiralexandre/logstash-exporter/access"
	"github.com/silveiralexandre/logstash-exporter/requests"
)

const (
	baseURI               = "api/v2"
	organizationsEndpoint = "organizations"
)

var (
	jsoniterJSON = jsoniter.ConfigCompatibleWithStandardLibrary
)

// Organizations represents Ansible Tower Organizations data provided by its API
type Organizations struct {
	Count    int64       `json:"count"`
	Next     interface{} `json:"next"`
	Previous interface{} `json:"previous"`
	Results  []struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"results"`
}

// Pull will extract list of organizations from Tower API
func (o *Organizations) Pull(a access.Credential) (Organizations, error) {
	uri := fmt.Sprintf("https://%v/%v/%v", a.Host, baseURI, organizationsEndpoint)

	r := requests.Client{}
	r, err := r.Setup(a, uri)
	if err != nil {
		return *o, err
	}

	b, err := r.Get(a)
	if err != nil {
		return *o, err
	}

	if len(b) == 0 || string(b) == "{}" {
		return *o, fmt.Errorf("Received empty response from URI: '%v'", uri)
	}

	err = jsoniterJSON.Unmarshal(b, &o)
	if err != nil {
		return *o, fmt.Errorf("Failed to unmarshal credentials file as a JSON structure: %v", err)
	}

	if len(o.Results) == 0 {
		return *o, fmt.Errorf("Empty list of organizations")
	}

	return *o, nil
}
