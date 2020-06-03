package asn

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type Prefixes struct {
	Prefix string `json:"prefix"`
}

type Data struct {
	Prefixes []Prefixes `json:"ipv4_prefixes"`
}

type Response struct {
	Status string `json:"status"`
	Data   Data   `json:"data"`
}

func ListPrefixes(ASN int) []string {
	// Initialize empty list.
	var list []string

	// Build URL.
	url := "https://api.bgpview.io/asn/AS" + strconv.Itoa(ASN) + "/prefixes"

	// Setup HTTP GET request.
	client := &http.Client{Timeout: time.Second * 5}
	req, _ := http.NewRequest("GET", url, nil)

	// Perform GET request.
	resp, err := client.Do(req)

	// Check for errors.
	if err != nil {
		fmt.Println(err)
	}

	// Ensure to close body at end of execution.
	defer resp.Body.Close()

	// Read output.
	body, err := ioutil.ReadAll(resp.Body)

	// Check for errors.
	if err != nil {
		fmt.Println(err)
	}

	// Create JSON response.
	var jsonResp Response

	// Prase JSON response.
	err = json.Unmarshal([]byte(string(body)), &jsonResp)

	if err != nil {
		fmt.Println(err)
	}

	// Loop through each prefix and add to list.
	for _, v := range jsonResp.Data.Prefixes {
		list = append(list, v.Prefix)
	}

	return list
}
