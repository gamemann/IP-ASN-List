package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

type JSONList struct {
	ASNs     []int    `json:"ASN"`
	Prefixes []string `json:"Prefix"`
}

func Handler(w http.ResponseWriter, r *http.Request, api int) {
	// Retrieve list value.
	list := filepath.Base(r.PostFormValue("list"))

	// Retrieve prefix or ASN value.
	var prefix string
	var asn int

	if api == 1 {
		prefix = r.PostFormValue("prefix")

		// Check prefix.
		_, _, err := net.ParseCIDR(prefix)

		// Check for errors.
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)

			fmt.Fprintf(w, "Error parsing prefix. Failed ParseCIDR inspection.")

			return
		}
	} else {
		var err error

		asn, err = strconv.Atoi(r.PostFormValue("asn"))

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)

			fmt.Fprintf(w, "Error parsing ASN as integer.")

			return
		}
	}

	// Try reading list file.
	file, err := os.Open("lists/" + list + ".json")

	// Check for errors.
	if err != nil {
		w.WriteHeader(http.StatusNotFound)

		fmt.Fprintf(w, "List not found - "+list+".")

		return
	}

	// Create stat.
	stat, _ := file.Stat()

	// Make data.
	data := make([]byte, stat.Size())

	// Read data.
	_, err = file.Read(data)

	// Check for errors.
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		fmt.Fprintf(w, "Error reading list file.")

		file.Close()

		return
	}

	// Close file for now.
	file.Close()

	// Create JSONList object.
	var jsonobj JSONList

	// Parse JSON from JSON file.
	err = json.Unmarshal([]byte(data), &jsonobj)

	// Check for errors when parsing.
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		fmt.Fprintf(w, "Error parsing JSON data within list.")

		return
	}

	// Check what API type (prefix or ASN).
	if api == 1 {
		// Check whether we're adding or removing prefix.
		if r.Method == "PUT" && !prefixExist(jsonobj, prefix) {
			// Append prefix to JSON object.
			jsonobj.Prefixes = append(jsonobj.Prefixes, prefix)
		} else if r.Method == "DELETE" {
			// Remove prefix from JSON object.
			jsonobj = remPrefix(jsonobj, prefix)
		}
	} else if api == 2 {
		// Check whether we're adding or removing prefix.
		if r.Method == "PUT" && !asnExist(jsonobj, asn) {
			// Append prefix to JSON object.
			jsonobj.ASNs = append(jsonobj.ASNs, asn)
		} else if r.Method == "DELETE" {
			// Remove ASN from JSON object.
			jsonobj = remASN(jsonobj, asn)
		}
	}

	// Convert object to JSON string (uses pretty-print).
	newdata, err := json.MarshalIndent(jsonobj, "", "    ")

	// Check for errors.
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		fmt.Fprintf(w, "Error converting new JSON object to JSON string.")

		return
	}

	// Write new JSON string to file.
	err = ioutil.WriteFile("lists/"+list+".json", newdata, 0644)

	// Check for errors.
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		fmt.Fprintf(w, "Error writing new JSON data to list file.")

		return
	}

	// Write Status OK header since we're good.
	w.WriteHeader(http.StatusOK)
}

func prefixExist(data JSONList, prefix string) bool {
	// Loop through all prefixes and if it matches, return true. Otherwise, return false.
	for _, v := range data.Prefixes {
		if v == prefix {
			return true
		}
	}

	return false
}

func asnExist(data JSONList, asn int) bool {
	// Loop through all ASNs and if it matches, return true. Otherwise, return false.
	for _, v := range data.ASNs {
		if v == asn {
			return true
		}
	}

	return false
}

func remPrefix(data JSONList, prefix string) JSONList {
	// Loop through all prefixes, match value, and remove.
	for i, v := range data.Prefixes {
		if v == prefix {
			data.Prefixes = append(data.Prefixes[:i], data.Prefixes[i+1:]...)

			break
		}
	}

	return data
}

func remASN(data JSONList, asn int) JSONList {
	// Loop through all prefixes, match value, and remove.
	for i, v := range data.ASNs {
		if v == asn {
			data.ASNs = append(data.ASNs[:i], data.ASNs[i+1:]...)

			break
		}
	}

	return data
}
