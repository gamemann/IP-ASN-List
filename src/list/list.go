package list

import (
	"../asn"
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type JSONList struct {
	ASNs     []int    `json:"ASN"`
	Prefixes []string `json:"Prefix"`
}

type List struct {
	Name     string
	Prefixes []string
}

type Lists struct {
	Lists []List
}

func UpdateLists(lists *Lists) {
	// Loop through each list and update.
	for _, v := range lists.Lists {
		UpdateList(&v)
	}
}

func UpdateList(list *List) bool {
	// Open file.
	file, err := os.OpenFile("public/"+list.Name+".txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	// Check for errors.
	if err != nil {
		fmt.Println(err)

		return false
	}

	// Remove contents of file.
	file.Truncate(0)
	file.Seek(0, 0)

	// Loop through each prefix and write them to the file as a new line.
	w := bufio.NewWriter(file)

	for _, v := range list.Prefixes {
		_, _ = w.WriteString(v + "\n")
	}

	// Flush buffer.
	w.Flush()

	// Close file.
	file.Close()

	return true
}

func ParseList(name string) JSONList {
	// Initiate empty List struct.
	var list JSONList

	// Open JSON list file.
	file, err := os.Open("lists/" + name + ".json")

	// Check for errors.
	if err != nil {
		fmt.Println("Error opening list \"" + name + "\".")

		return list
	}

	// Defer file close.
	defer file.Close()

	// Create stat.
	stat, _ := file.Stat()

	// Make data.
	data := make([]byte, stat.Size())

	// Read data.
	_, err = file.Read(data)

	// Check for errors.
	if err != nil {
		fmt.Println("Error reading config file.")

		return list
	}

	// Parse JSON data.
	err = json.Unmarshal([]byte(string(data)), &list)

	// Check for errors.
	if err != nil {
		fmt.Println("Error parsing JSON Data.")

		return list
	}

	return list
}

func ExtractList(json JSONList) []string {
	// Create empty prefixes list.
	var prefixes []string

	// Loop through each AS and do appropriate lookups. After lookups, append to list.
	for _, v := range json.ASNs {
		prefixes = append(prefixes, asn.ListPrefixes(v)...)
	}

	// Loop through each additional prefix and add.
	for _, v := range json.Prefixes {
		prefixes = append(prefixes, v)
	}

	return prefixes
}

func GetLists() Lists {
	// Initialize empty lists slice.
	var lists Lists

	// Read a directory.
	files, err := ioutil.ReadDir("lists/")

	// Check for errors.
	if err != nil {
		fmt.Println(err)

		return lists
	}

	for _, f := range files {
		// Check if this is a directory.
		if f.IsDir() {
			continue
		}

		// Get file name without extension.
		name := strings.TrimSuffix(f.Name(), filepath.Ext(f.Name()))

		// Parse specific list.
		jsonlists := ParseList(name)

		// Extract all IPs including ASN prefixes.
		prefixes := ExtractList(jsonlists)

		// Create list.
		var list List

		// Assign name.
		list.Name = name

		// Append prefixes to list.
		list.Prefixes = append(list.Prefixes, prefixes...)

		// Add list to lists variable.
		lists.Lists = append(lists.Lists, list)
	}

	return lists
}
