package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
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

	// compute InfluxDB URL
	sessionurl, err := url.Parse(Config.InfluxServer)
	sessionurl.Path = "write"
	q := sessionurl.Query()
	q.Set("db", Config.InfluxDB)
	sessionurl.RawQuery = q.Encode()

	req, err := http.NewRequest(http.MethodPost, sessionurl.String(), bytes.NewBufferString(lines))
	if err != nil {
		panic(err)
	}
	if len(Config.InfluxUser) > 0 {
		req.SetBasicAuth(Config.InfluxUser, Config.InfluxPasswd)
	}

	if !dryrun {
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			panic(err)
		}
		if resp.StatusCode != http.StatusOK {
			panic(errors.New(fmt.Sprintf("InfluxDB return %d", resp.StatusCode)))
		}
	} else {
		fmt.Println("DRYRUN! No actual call to InfluxDB has been made. The following call would have been made without --dryrun")
		requestDump, err := httputil.DumpRequest(req, true)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(string(requestDump))
	}
}
