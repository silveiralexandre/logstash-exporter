package tower

import (
	"fmt"
	"log"

	jsoniter "github.com/json-iterator/go"
	"github.com/silveiralexandre/logstash-exporter/access"
)

// ReportInventoryMetadata presents to stdout a predefined set of inventory data variables as a JSON array
// for Logstash exec plugin consumption with JSON codec
func ReportInventoryMetadata(limit, timeout, retries *int) {
	a := access.Credential{}
	a, err := a.PullEnvTower(limit, timeout, retries)
	if err != nil {
		log.Fatal(err)
	}

	o := Organizations{}
	o, err = o.Pull(a)
	if err != nil {
		log.Fatal(err)
	}

	i := Inventories{}
	vd, err := i.PullVariableData(a, o)
	if err != nil {
		log.Fatal(err)
	}

	b, err := jsoniter.Marshal(vd)
	if err != nil {
		log.Fatal(err)
	}

	if len(b) != 0 {
		fmt.Print(string(b))
	}

}
