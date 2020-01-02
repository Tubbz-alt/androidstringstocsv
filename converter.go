package main

import (
	"encoding/xml"
	"io/ioutil"
	"os"
)

// StringEntry struct defines a node of <string></string> tag in xml file
type StringEntry struct {
	XMLName xml.Name `xml:"string"`    // name of xml tag
	Name    string   `xml:"name,attr"` // name attribute of xml tag
	Value   string   `xml:",innerxml"` // value of xml string tag
}

// Resources struct defines a node of <resources></resources> tag in xml file
type Resources struct {
	XMLName xml.Name      `xml:"resources"` //
	Strings []StringEntry `xml:"string"`    //
}

func readFromXML(path string) (r *Resources, err error) {
	reader, err := os.Open(path)
	if err != nil {
		return
	}
	byteArray, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	var res Resources
	err = xml.Unmarshal(byteArray, &res)
	r = &res
	return
}
