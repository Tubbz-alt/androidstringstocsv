package main

import (
	"encoding/csv"
	"encoding/xml"
	"io/ioutil"
	"os"
	"strings"
)

const (
	valuesPrefix = "values-"
)

// StringEntry struct defines a node of <string></string> tag in xml file
type StringEntry struct {
	XMLName xml.Name `xml:"string"`    // name of xml tag
	Name    string   `xml:"name,attr"` // name attribute of xml tag
	Value   string   `xml:",innerxml"` // value of xml string tag
}

// Resources struct defines a node of <resources></resources> tag in xml file
type Resources struct {
	XMLName xml.Name      `xml:"resources"` // name of xml tag
	Strings []StringEntry `xml:"string"`    // strings itself
}

// ValuesFile struct defines a values type in android framework
type ValuesFile struct {
	LanguageCode string    // language code of the strings.xml file
	Content      Resources // content (strings) itself
}

// unmarshals structure of strings.xml file and returns its content
func readFromXMLFile(path string) (r *Resources, err error) {
	var reader *os.File
	var byteArray []byte
	var res Resources

	if reader, err = os.Open(path); err != nil {
		return nil, err
	}

	if byteArray, err = ioutil.ReadAll(reader); err != nil {
		return nil, err
	}

	err = xml.Unmarshal(byteArray, &res)
	r = &res
	return r, err
}

// reads and unmarshals all strings.xml files in the "res" folder
func readFromResFolder(path string) ([]ValuesFile, error) {
	contents, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	vals := []ValuesFile{}

	for _, entry := range contents {
		// skip if it is not a directory, that starts with "values-"
		if !entry.IsDir() || !strings.HasPrefix(entry.Name(), valuesPrefix) {
			continue
		}

		// reading xml structure
		res, err := readFromXMLFile(path + "/" + entry.Name() + "/strings.xml")

		if err != nil {
			return nil, err
		}

		vals = append(vals, ValuesFile{
			LanguageCode: entry.Name()[len(valuesPrefix):],
			Content:      *res,
		})
	}

	return vals, err
}

// converts the slice of ValuesFile to the map[langCode]map[name]value
func convertValuesToMap(vals []ValuesFile) (m map[string]map[string]string) {
	m = make(map[string]map[string]string)
	// filling val names
	for _, v := range vals {
		for _, s := range v.Content.Strings {
			if m[s.Name] == nil {
				m[s.Name] = make(map[string]string)
			}
			m[s.Name][v.LanguageCode] = s.Value
		}
	}
	return
}

// converts the map[langCode]map[name]value to the matrix of strings, like that:
//      , lang1, lang2, lang3 ...;
// name1,  val1,  val2,  val3 ...;
// name2,  val1,  val2,  val3 ...;
//   ...,   ...,   ...,   ... ...
func convertResourcesToStrings(m map[string]map[string]string) (s [][]string) {
	// if we get empty map - just do nothing
	if m == nil {
		return nil
	}

	// filling first line - headers
	for _, langVal := range m {
		row := []string{"code \\ language"}
		for lang := range langVal {
			row = append(row, lang)
		}
		s = append(s, row)
		break
	}

	// filling values itself
	for name := range m {
		row := []string{name}
		for _, langCode := range s[0][1:] {
			row = append(row, m[name][langCode])
		}
		s = append(s, row)
	}
	return
}

// writes the specified structure to the csv file
func writeToCSVFile(path string, vals [][]string) (file *os.File, err error) {
	// creating the csv file itself
	file, err = os.Create(path)
	if err != nil {
		return nil, err
	}

	csvWriter := csv.NewWriter(file)
	err = csvWriter.WriteAll(vals)
	if err != nil {
		return nil, err
	}
	return file, nil
}
