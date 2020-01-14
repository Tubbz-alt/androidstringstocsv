package main

import (
	"encoding/json"
	"fmt"
	"github.com/Semior001/androidstringstocsv/converter/xml"
	"os"
)

var xmlHeader string = "<?xml version=\"1.0\" encoding=\"utf-8\"?>"

const (
	helpString = `asc - android strings converter by semior001

Usage: asc [COMMAND] [FROM] [TO]

Commands:
	xml2csv  - convert android xml strings folders to csv file
	csv2xml  - convert csv file to android xml "values" folders

From - path to the "res" folder in your android project

To - where to put the output (csv file in case of "xml2csv", 
	folders with "values-xx" in case of "csv2xml")
`
)

// just print help
func help() {
	fmt.Println(helpString)
}

func main() {
	if len(os.Args) < 4 || os.Args[1] == "help" {
		help()
		return
	}

	var (
		command string = os.Args[1]
		from    string = os.Args[2]
		to      string = os.Args[3]
	)

	switch command {
	case "xml2csv":
		res, err := xml.ReadXMLFile(from)
		if err != nil {
			panic(err)
		}
		r, _ := json.Marshal(res)
		fmt.Println(string(r))
	case "csv2xml":
		fmt.Println(to)
	default:
		help()
		return
	}
}
