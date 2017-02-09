package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"
)

const PAGESIZE = 100
const STATSSIZE = 5

var Config *Configuration

var statsType [STATSSIZE]string = [STATSSIZE]string{"SYNC", "NOTPROV", "TRANSFER_FAILING", "CONFLICT", "TOTAL"}
var stats map[string]uint
var custstats map[string]map[string]uint

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
		wg.Add(1)
		go count(data, wg)
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

func main() {
	stats = make(map[string]uint, STATSSIZE)
	for _, i := range statsType {
		stats[i] = 0
	}
	custstats = make(map[string]map[string]uint)

	Config = readDefaultConfigFiles()
	Config = joinConfig(Config, parseFlags())

	startSession()
	getZoneList()

	influx()
}
