package main

import (
	"encoding/csv"
	"encoding/xml"
	"io/ioutil"
	"os"
	"strings"
)

const (
	valuesPrefix    = "values-"
	stringsFilename = "strings.xml"
	csvExportHeader = "code \\ language"
	exportFileMode  = 0750
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

// ValuesFile struct defines a values type in android framework - map[langCode]resource
type ValuesFile map[string]Resources

// unmarshals structure of strings.xml file and returns its content
func readXMLFile(path string) (r *Resources, err error) {
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
func readResFolder(path string) (ValuesFile, error) {
	contents, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	vals := make(ValuesFile)

	for _, entry := range contents {
		// skip if it is not a directory, that starts with "values-"
		if !entry.IsDir() || !strings.HasPrefix(entry.Name(), valuesPrefix) {
			continue
		}

		// reading xml structure
		res, err := readXMLFile(path + "/" + entry.Name() + "/" + stringsFilename)

		if err != nil {
			return nil, err
		}

		langCode := entry.Name()[len(valuesPrefix):]

		vals[langCode] = *res
	}

	return vals, err
}

// converts the slice of ValuesFile to the map[name]map[langCode]value
func convertValuesToMap(vals ValuesFile) (m map[string]map[string]string) {

	m = make(map[string]map[string]string)
	// filling val names
	for langCode, res := range vals {
		for _, str := range res.Strings {
			if m[str.Name] == nil {
				m[str.Name] = make(map[string]string)
			}
			m[str.Name][langCode] = str.Value
		}
	}
	return
}

// converts the map[name]map[langCode]value to the matrix of strings, like that:
//      , lang1, lang2, lang3 ...;
// name1,  val1,  val2,  val3 ...;
// name2,  val1,  val2,  val3 ...;
//   ...,   ...,   ...,   ... ...
func convertMapToStringsMatrix(m map[string]map[string]string) (s [][]string) {
	// if we get empty map - just do nothing
	if m == nil {
		return nil
	}

	// filling first line - headers
	for _, langVal := range m {
		row := []string{csvExportHeader}
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

// reads CSV file
func readCSVFile(path string) (vals [][]string, err error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	vals, err = csv.NewReader(file).ReadAll()
	return vals, err
}

// converts matrix of strings (example below) to the map[name]map[langCode]value
//      , lang1, lang2, lang3 ...;
// name1,  val1,  val2,  val3 ...;
// name2,  val1,  val2,  val3 ...;
//   ...,   ...,   ...,   ... ...
func convertStringsMatrixToMap(vals [][]string) (m map[string]map[string]string) {
	if vals == nil {
		return nil
	}
	m = make(map[string]map[string]string)

	// first read names
	for _, row := range vals[1:] {
		m[row[0]] = make(map[string]string)
		for idx, val := range row[1:] {
			// wtf? idx can be zero despite I start from row + 1 idx?
			m[row[0]][vals[0][idx+1]] = val
		}
	}
	return
}

// converts map[name]map[langCode]value to the slice of ValuesFile
func convertMapToValues(m map[string]map[string]string) (vals ValuesFile) {

	vals = make(ValuesFile)

	for name, langVal := range m {
		for langCode, val := range langVal {
			res, ok := vals[langCode]

			if !ok {
				res = Resources{
					Strings: []StringEntry{},
				}
			}

			res.Strings = append(
				vals[langCode].Strings,
				StringEntry{
					Name:  name,
					Value: val,
				},
			)

			vals[langCode] = res

		}
	}

	return
}

// marshals and writes xml structure of Resources to the specified file
func writeToXMLFile(path string, r Resources) (file *os.File, err error) {
	file, err = os.Create(path)
	if err != nil {
		return
	}
	byteArray, err := xml.MarshalIndent(r, "", "	")
	if err != nil {
		return nil, err
	}
	byteArray = []byte(xml.Header + string(byteArray))
	err = ioutil.WriteFile(path, byteArray, exportFileMode)
	return file, err
}

// marshals and writes whole xml structure of given ValuesFile to the specified path
func writeResFolder(path string, vals ValuesFile) (files []*os.File, err error) {
	err = os.Mkdir(path, exportFileMode)
	if err != nil {
		return nil, err
	}

	files = []*os.File{}

	for langCode, res := range vals {
		valPath := path + "/" + valuesPrefix + langCode

		err = os.Mkdir(valPath, exportFileMode)
		if err != nil {
			return nil, err
		}

		var file *os.File

		file, err = writeToXMLFile(valPath+"/"+stringsFilename, res)
		files = append(files, file)
		if err != nil {
			return files, err
		}

	}
	return
}
