package tower

import (
	"fmt"
	"sync"

	"github.com/silveiralexandre/logstash-exporter/access"
	"github.com/silveiralexandre/logstash-exporter/requests"
)

const (
	inventoriesEndpoint = "inventories"
)

// Inventories represents Ansible Tower Organization Inventory data
type Inventories struct {
	Count    int64       `json:"count"`
	Next     interface{} `json:"next"`
	Previous interface{} `json:"previous"`
	Results  []struct {
		ID           int64  `json:"id"`
		Name         string `json:"name"`
		Organization int64  `json:"organization"`
		Related      struct {
			Organization string `json:"organization"`
			VariableData string `json:"variable_data"`
		} `json:"related"`
	} `json:"results"`
}

// VariableData represents inventory metadata as predefined by users
type VariableData struct {
	CustomerPrefix     string `json:"customer_prefix,omitempty"`
	CustomerName       string `json:"customer_name,omitempty"`
	AccountName        string `json:"account_name,omitempty"`
	AccountRestriction string `json:"account_restriction,omitempty"`
	Country            string `json:"country,omitempty"`
	Geo                string `json:"geo,omitempty"`
	Market             string `json:"market,omitempty"`
	Sector             string `json:"sector,omitempty"`
	Type               string `json:"type,omitempty"`
	OrgGeo             string `json:"org_geo,omitempty"`
	OrgCountry         string `json:"org_country,omitempty"`
	OrgMarket          string `json:"org_market,omitempty"`
	OrgAccountName     string `json:"org_account_name,omitempty"`
	OrgIndustry        string `json:"org_industry,omitempty"`
	OrgSector          string `json:"org_sector,omitempty"`
	OrgType            string `json:"org_type,omitempty"`
	OrgRestriction     string `json:"org_restriction,omitempty"`
}

// Pull retrieves a list of all accessible inventories from Ansible Tower API
func (inv *Inventories) Pull(a access.Credential, uri string) (*Inventories, error) {
	r := requests.Client{}
	r, err := r.Setup(a, uri)
	if err != nil {
		return nil, err
	}

	b, err := r.Get(a)
	if err != nil {
		return nil, err
	}

	err = jsoniterJSON.Unmarshal(b, &inv)
	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshal credentials file as a JSON structure: %v", err)
	}

	return inv, nil
}

// PullVariableData extracts Inventory Variable Data from Ansible Tower API
func (inv *Inventories) PullVariableData(a access.Credential, o Organizations) ([]VariableData, error) {
	variableDataHrefs, err := inv.pullVariableDataHRefs(a, o)
	if err != nil {
		return nil, err
	}

	hrefs := make(chan string)
	limitOfWorkers := a.Limit

	wg := sync.WaitGroup{}
	wg.Add(limitOfWorkers)

	data := []VariableData{}
	for i := 0; i < limitOfWorkers; i++ {
		go func(workerNum int) {
			defer wg.Done()

			for {
				href, ok := <-hrefs // pull tasks from queue until done
				if !ok {
					return
				}

				v, err := doGetVariableData(a, href)
				if err != nil {
					break
				}

				if v.AccountName != "" {
					data = append(data, v)
				}
			}
		}(i)
	}

	// Initiate queue of tasks
	for i := range variableDataHrefs {
		hrefs <- variableDataHrefs[i]
	}

	close(hrefs)
	wg.Wait()

	return data, nil
}

// Pull together all Ansible Tower Inventory Variable Data hypertext references
func (inv *Inventories) pullVariableDataHRefs(a access.Credential, o Organizations) ([]string, error) {
	inventoryIndices := make(chan int)
	limitOfWorkers := a.Limit

	wg := sync.WaitGroup{}
	wg.Add(limitOfWorkers)

	hrefs := []string{}
	for i := 0; i < limitOfWorkers; i++ {
		counter := 0
		go func(workerNum int) {

			defer wg.Done()

			for {
				counter++
				inventoryIndex, ok := <-inventoryIndices // pull tasks from queue until done

				if !ok {
					return
				}

				uri := fmt.Sprintf("https://%v/%v/%v/%v/%v",
					a.Host,
					baseURI,
					organizationsEndpoint,
					o.Results[inventoryIndex].ID,
					inventoriesEndpoint,
				)

				var err error
				inv, err = inv.Pull(a, uri)
				if err != nil {
					break
				}

				for j := range inv.Results {
					if inv.Results[j].Related.VariableData != "" {
						hrefs = append(hrefs, inv.Results[j].Related.VariableData)
					}
				}
			}
		}(i)
	}

	// Initiate queue of tasks
	for i := range o.Results {
		inventoryIndices <- i
	}

	close(inventoryIndices)
	wg.Wait()

	return hrefs, nil
}

func doGetVariableData(a access.Credential, selfHRef string) (VariableData, error) {
	v := VariableData{}
	uri := fmt.Sprintf("https://%v%v", a.Host, selfHRef)

	r := requests.Client{}
	r, err := r.Setup(a, uri)
	if err != nil {
		return v, err
	}

	b, err := r.Get(a)
	if err != nil {
		return v, err
	}

	err = jsoniterJSON.Unmarshal(b, &v)
	if err != nil {
		return v, err
	}

	if (VariableData{} == v) {
		return v, nil
	}

	v = updateMetadata(v)
	return v, nil
}

func updateMetadata(v VariableData) VariableData {
	if v.Market != "" {
		v.Market = reclassifyMarket(v.Market)
	}

	if v.OrgMarket != "" {
		v.OrgMarket = reclassifyMarket(v.OrgMarket)
	}

	v.Market = mapMarketFromCountry(v.Country)

	return v
}

func reclassifyMarket(market string) string {
	markets := map[string]string{
		"ASEAN IMT":                      "ASEAN",
		"Australia/New Zealand IMT":      "ANZ",
		"Belgium":                        "Benelux",
		"Netherlands":                    "Benelux",
		"Luxembourg IMT":                 "Benelux",
		"Canada IMT":                     "Canada",
		"Central and Eastern Europe IMT": "CEE",
		"DACH IMT":                       "DACH",
		"France IMT":                     "France",
		"Greater China":                  "GCG",
		"India-South Asia IMT":           "ISA",
		"Italy IMT":                      "Italy",
		"Japan IMT":                      "Japan",
		"Japan":                          "Japan",
		"Korea IMT":                      "Korea",
		"Latin America":                  "Latin America",
		"Middle East & Africa":           "MEA",
		"Nordic IMT":                     "Nordic",
		"Spain":                          "SPGI",
		"Portugal":                       "SPGI",
		"Greece":                         "SPGI",
		"Israel IMT":                     "SPGI",
		"United Kingdom and Ireland IMT": "UKI",
		"United Kingdom":                 "UKI",
		"United Kingdom IMT":             "UKI",
		"Ireland IMT":                    "UKI",
		"Ireland":                        "UKI",
		"US IMT":                         "US",
	}

	for k, v := range markets {
		if market == k {
			return v
		}
	}
	return market
}

func mapMarketFromCountry(country string) (market string) {
	countries := map[string]string{
		"Denmark":     "Nordic",
		"Finland":     "Nordic",
		"Sweden":      "Nordic",
		"France":      "France",
		"Portugal":    "SPGI",
		"Spain":       "SPGI",
		"Italy":       "SPGI",
		"Greece":      "SPGI",
		"UK":          "UKI",
		"UKI":         "UKI",
		"Benelux":     "Benelux",
		"Germany":     "DACH",
		"Austria":     "DACH",
		"Switzerland": "DACH",
		"Turkiye":     "MEA",
		"Japan":       "Japan",
	}

	for k, v := range countries {
		if country == k {
			return v
		}
	}
	return "To be defined"
}
