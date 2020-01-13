// Package xml specifies functions and structs
// for writing dictionaries to xml files and
// android project res folder
package xml

import (
	"encoding/xml"
	"github.com/Semior001/androidstringstocsv/converter/general"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const (
	// ValuesPrefix defines the default prefix for values folder
	ValuesPrefix = "values-"
	// StringsFilename defines the default filename for android string constants file
	StringsFilename = "strings.xml"
	// ExportFileMode defines the default permissions for created file
	ExportFileMode = 0750
)

// StringEntry struct defines a node of <string></string> tag in xml file
type StringEntry struct {
	XMLName xml.Name `xml:"string"`    // name of xml tag
	Name    string   `xml:"name,attr"` // name attribute of xml tag
	Value   string   `xml:",innerxml"` // value of xml string tag
}

// ResourcesEntry struct defines a node of <resources></resources> tag in xml file
type ResourcesEntry struct {
	XMLName xml.Name      `xml:"resources"` // name of xml tag
	Strings []StringEntry `xml:"string"`    // strings itself
}

// ReadXMLFile unmarshals structure of strings.xml file and returns its content
func ReadXMLFile(path string) (r *ResourcesEntry, err error) {
	var reader *os.File
	var byteArray []byte
	var res ResourcesEntry

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

// WriteToXMLFile marshals and writes xml structure of ResourcesEntry to the specified file
func (r *ResourcesEntry) WriteToXMLFile(path string) (file *os.File, err error) {
	file, err = os.Create(path)
	if err != nil {
		return
	}
	byteArray, err := xml.MarshalIndent(*r, "", "	")
	if err != nil {
		return nil, err
	}
	byteArray = []byte(xml.Header + string(byteArray))
	err = ioutil.WriteFile(path, byteArray, ExportFileMode)
	return file, err
}

// ConvertToDictionary converts the given ResourcesEntry to the dictionary map[code]translation
func (r *ResourcesEntry) ConvertToDictionary() (d general.Dictionary) {
	d = make(general.Dictionary)

	for _, entry := range (*r).Strings {
		d[entry.Name] = entry.Value
	}

	return
}

// convertDictionaryToResources converts the given dictionary map[code]translation to the ResourcesEntry
func convertDictionaryToResources(d general.Dictionary) (r ResourcesEntry) {
	r = ResourcesEntry{
		Strings: []StringEntry{},
	}
	for name, value := range d {
		r.Strings = append(r.Strings, StringEntry{
			Name:  name,
			Value: value,
		})
	}
	return
}

// exportDictionaryToXML writes the given dictionary to the xml file at the given path
func exportDictionaryToXML(path string, d general.Dictionary) (files *os.File, err error) {
	r := convertDictionaryToResources(d)
	files, err = r.WriteToXMLFile(path)
	return
}

// WriteResFolder writes the given set of dictionaries to the res folder at the given path
func WriteResFolder(path string, dicts general.Dictionaries, override bool) (files []*os.File, err error) {
	err = os.Mkdir(path, ExportFileMode)
	if (err != nil && err != os.ErrExist) || (err == os.ErrExist && !override) { // todo
		return nil, err
	}

	files = []*os.File{}

	for langCode, d := range dicts {
		valPath := filepath.Join(path, ValuesPrefix+langCode)

		err = os.Mkdir(valPath, ExportFileMode)
		if (err != nil && err != os.ErrExist) || (err == os.ErrExist && !override) { // todo
			return
		}

		var file *os.File

		file, err = exportDictionaryToXML(filepath.Join(valPath, StringsFilename), d)
		files = append(files, file)
		if (err != nil && err != os.ErrExist) || (err == os.ErrExist && !override) { // todo
			return
		}
	}

	return
}

// ReadResFolder reads and unmarshals all strings.xml files in the "res" folder
func ReadResFolder(path string) (dicts general.Dictionaries, err error) {
	contents, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	dicts = make(general.Dictionaries)

	for _, entry := range contents {
		// skip if it is not a directory, that starts with "values-"
		if !entry.IsDir() || !strings.HasPrefix(entry.Name(), ValuesPrefix) {
			continue
		}

		var res *ResourcesEntry
		// reading xml structure
		res, err = ReadXMLFile(filepath.Join(path, entry.Name(), StringsFilename))

		if err != nil {
			return
		}

		langCode := entry.Name()[len(ValuesPrefix):]

		dicts[langCode] = (*res).ConvertToDictionary()
	}

	return
}
