// Package general declares some general
// structs for working with dictionaries
package general

// Dictionary defines a single dictionary in
// format map[code]translation
type Dictionary map[string]string

// Dictionaries defines a set of dictionaries
// in format map[languageCode]Dictionary
type Dictionaries map[string]Dictionary
