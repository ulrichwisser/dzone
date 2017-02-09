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
		comma := ""
		for _, status := range statsType {
			lines = lines + fmt.Sprintf(comma+"%s=%di", status, custstats[cust][status])
			comma = ","
		}
		lines = lines + "\n"
	}

	// save to InfluxDB
	_, err := http.Post(fmt.Sprintf("http://%s:%s/write?db=%s", Config.InfluxServer, Config.InfluxPort, Config.InfluxDB), "application/x-www-form-urlencoded", bytes.NewBufferString(lines))
	if err != nil {
		panic(err)
	}
}
