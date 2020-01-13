// Package csv specifies functions and structs
// for writing dictionaries to single csv file
package csv

import (
	"encoding/csv"
	"github.com/Semior001/androidstringstocsv/converter/general"
	"os"
)

const (
	// SlicesHeader defines the default header for exported CSV file
	SlicesHeader = "code \\ language"
)

// writeSlicesToCSVFile writes the specified structure to the csv file
func writeSlicesToCSVFile(path string, vals [][]string) (file *os.File, err error) {
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

// readSlicesFromCSVFile reads CSV file
func readSlicesFromCSVFile(path string) (vals [][]string, err error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	vals, err = csv.NewReader(file).ReadAll()
	return vals, err
}

// convertDictionariesToSlices converts the given set of dictionaries to the matrix of strings, like that:
//          , lang1, lang2, lang3 ...;
//     name1,  val1,  val2,  val3 ...;
//     name2,  val1,  val2,  val3 ...;
//       ...,   ...,   ...,   ... ...
func convertDictionariesToSlices(dicts general.Dictionaries) (vals [][]string) {
	vals = [][]string{{SlicesHeader}}

	// first filling out language codes and names
	for langCode, dict := range dicts {
		vals[0] = append(vals[0], langCode)
		for name := range dict {
			vals = append(vals, []string{name})
		}
	}

	for i := range vals[1:] {
		name := vals[i][0]
		for _, langCode := range vals[i][1:] {
			vals[i] = append(vals[i], dicts[langCode][name])
		}
	}

	return
}

// convertSlicesToDictionaries converts the given matrix of strings (example below) to the set of dictionaries:
//          , lang1, lang2, lang3 ...;
//     name1,  val1,  val2,  val3 ...;
//     name2,  val1,  val2,  val3 ...;
//       ...,   ...,   ...,   ... ...
func convertSlicesToDictionaries(vals [][]string) (dicts general.Dictionaries) {
	dicts = make(general.Dictionaries)

	// first filling out language codes
	for _, langCode := range vals[0][1:] {
		dicts[langCode] = make(general.Dictionary)
	}

	for i, row := range vals[1:] {
		for j, val := range row[1:] {
			langCode := vals[0][j]
			name := vals[i][0]

			dicts[langCode][name] = val
		}
	}

	return
}

// WriteCSVFile writes the given set of dictionaries to the csv file
func WriteCSVFile(path string, dicts general.Dictionaries, override bool) (file *os.File, err error) {
	// creating the csv file itself
	file, err = os.Create(path)
	if (err != nil && err != os.ErrExist) || (err == os.ErrExist && !override) { // todo
		return
	}

	csvWriter := csv.NewWriter(file)
	vals := convertDictionariesToSlices(dicts)
	err = csvWriter.WriteAll(vals)
	if err != nil {
		return
	}
	return
}

// ReadCSVFile reads and unmarshals all words from the given csv file and converts to the
// set of dictionaries
func ReadCSVFile(path string) (dicts general.Dictionaries, err error) {
	vals, err := readSlicesFromCSVFile(path)
	if err != nil {
		return
	}

	dicts = convertSlicesToDictionaries(vals)
	return
}
