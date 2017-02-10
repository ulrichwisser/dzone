package main

import (
	"bytes"
	"fmt"
	"net/http"
)

func influx() {

	// Total stats for IIS Anycast
	lines := "Anycast "
	comma := ""
	for _, status := range statsType {
		lines = lines + fmt.Sprintf(comma+"%s=%di", status, stats[status])
		comma = ","
	}
	lines = lines + "\n"

	// Stats per IIS Anycast customer
	for cust := range custstats {
		lines = lines + fmt.Sprintf("AnycastCustomers,cust=%s ", cust)
		comma = ""
		for _, status := range statsType {
			lines = lines + fmt.Sprintf(comma+"%s=%di", status, custstats[cust][status])
			comma = ","
		}
		lines = lines + "\n"
	}

	// Stats per TLD
	lines = lines + "AnycastTld "
	comma = ""
	for tld := range statsTld {
		lines = lines + comma + fmt.Sprintf("%s=%di", tld, statsTld[tld])
		comma = ","
	}
	lines = lines + "\n"

	// stats per TLD and customer
	for cust := range statsTldCustomer {
		lines = lines + fmt.Sprintf("AncastTldCustomers,cust=%s ", cust)
		comma = ""
		for tld := range statsTldCustomer[cust] {
			lines = lines + comma + fmt.Sprintf("%s=%di", tld, statsTldCustomer[cust][tld])
			comma = ","
		}
		lines = lines + "\n"
	}

	// save to InfluxDB
	if !dryrun {
		_, err := http.Post(fmt.Sprintf("http://%s:%s/write?db=%s", Config.InfluxServer, Config.InfluxPort, Config.InfluxDB), "application/x-www-form-urlencoded", bytes.NewBufferString(lines))
		if err != nil {
			panic(err)
		}
	} else {
		fmt.Printf("http://%s:%s/write?db=%s\n", Config.InfluxServer, Config.InfluxPort, Config.InfluxDB)
		fmt.Println(lines)
	}
}
