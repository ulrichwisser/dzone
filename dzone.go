package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

const PAGESIZE = 100
const STATSSIZE = 6

var Config *Configuration

// Number of zones statistics
var statsType [STATSSIZE]string = [STATSSIZE]string{"SYNC", "OUTOFSYNC", "NOTPROV", "TRANSFER_FAILING", "CONFLICT", "TOTAL"}
var stats map[string]uint
var custstats map[string]map[string]uint

// Number of zones by tld statistics
var statsTldType = [...]string{"SE", "NU", "ARPA"}
var TLD_OTHER = "OTHER"
var statsTld map[string]uint
var statsTldCustomer map[string]map[string]uint

var dryrun bool = false

func startSession() {
	sessionurl, err := url.Parse(Config.ServerRoot)
	sessionurl.Path = "sessions"
	buf := bytes.NewBufferString(`{"username":"` + Config.ApiUser + `","password":"` + Config.ApiPasswd + `"}`)
	req, err := http.NewRequest(http.MethodPost, sessionurl.String(), buf)
	if err != nil {
		panic(err)
	}
	req.Header.Add("Authorization", Config.ApiSecret)
	req.Header.Add("Content-type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	var access accessData
	err = json.Unmarshal(body, &access)
	if err != nil {
		panic(err)
	}
	Config.ApiAccessToken = access.Access_token
}

func getZoneListPage(page uint) *zoneList {
	sessionurl, err := url.Parse(Config.ServerRoot)
	sessionurl.Path = "zones"
	q := sessionurl.Query()
	q.Set("name", "")
	q.Set("page", fmt.Sprintf("%d", page))
	q.Set("pageSize", fmt.Sprintf("%d", PAGESIZE))
	sessionurl.RawQuery = q.Encode()

	buf := bytes.NewBufferString("")
	req, err := http.NewRequest(http.MethodGet, sessionurl.String(), buf)
	if err != nil {
		panic(err)
	}
	req.Header.Add("Authorization", Config.ApiUser+":"+Config.ApiAccessToken)
	req.Header.Add("Content-type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var data zoneList
	err = json.Unmarshal(body, &data)
	if err != nil {
		panic(err)
	}

	return &data
}

func getZoneList() {
	var page uint = 0
	var wg *sync.WaitGroup = &sync.WaitGroup{}

	for {
		data := getZoneListPage(page)
		wg.Add(2)
		go count(data, wg)
		go countTld(data, wg)
		if !data.HasNext {
			break
		}
		page++
	}
	wg.Wait()
}

func count(data *zoneList, wg *sync.WaitGroup) {
	for _, item := range data.Items {
		// init customer stats if needed
		if _, ok := custstats[item.CustomerName]; !ok {
			custstats[item.CustomerName] = make(map[string]uint)
			for _, i := range statsType {
				custstats[item.CustomerName][i] = 0
			}
		}

		// stats
		stats[item.Status]++
		stats["TOTAL"]++
		custstats[item.CustomerName][item.Status]++
		custstats[item.CustomerName]["TOTAL"]++
		if item.Conflict {
			stats["CONFLICT"]++
			custstats[item.CustomerName]["CONFLICT"]++
		}
	}
	wg.Done()
}

func countTld(data *zoneList, wg *sync.WaitGroup) {
	for _, item := range data.Items {
		// init customer stats if needed
		if _, ok := statsTldCustomer[item.CustomerName]; !ok {
			statsTldCustomer[item.CustomerName] = make(map[string]uint)
			for _, tld := range statsTldType {
				statsTldCustomer[item.CustomerName][tld] = 0
			}
		}

		foundTld := false
		for _, tld := range statsTldType {
			if strings.HasSuffix(strings.ToUpper(item.Name), "."+tld+".") {
				fmt.Printf("TLD %s %s %s\n", tld, item.CustomerName, item.Name)
				statsTld[tld]++
				statsTldCustomer[item.CustomerName][tld]++
				foundTld = true
				break
			}
		}
		if !foundTld {
			fmt.Printf("TLD %s %s %s\n", TLD_OTHER, item.CustomerName, item.Name)
			statsTld[TLD_OTHER]++
			statsTldCustomer[item.CustomerName][TLD_OTHER]++
		}

	}
	wg.Done()
}

func main() {
	stats = make(map[string]uint, STATSSIZE)
	for _, i := range statsType {
		stats[i] = 0
	}
	custstats = make(map[string]map[string]uint)
	statsTld = make(map[string]uint, len(statsTldType)+1)
	statsTldCustomer = make(map[string]map[string]uint)
	Config = readDefaultConfigFiles()
	Config = joinConfig(Config, parseFlags())

	startSession()
	getZoneList()

	influx()
}
