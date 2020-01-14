package csv

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConvertingDictToSlices(t *testing.T) {
	vals := convertDictionariesToSlices(map[string]map[string]string{
		"tl": map[string]string{
			"test_str": "Test translation",
		},
	})
	assert.Equal(t, [][]string{
		{SlicesHeader, "tl"},
		{"test_str", "Test translation"},
	}, vals)
}

func TestVConvertingSlicesToDict(t *testing.T) {
	dicts := convertSlicesToDictionaries([][]string{
		{SlicesHeader, "tl"},
		{"test_str", "Test translation"},
	})
	assert.Equal(t, map[string]map[string]string{
		"tl": map[string]string{
			"test_str": "Test translation",
		},
	}, dicts)
}

func TestCSVReadWrite(t *testing.T) {
	defer os.RemoveAll("/tmp/androidstringscsv.test")
	_, err := writeSlicesToCSVFile("/tmp/androidstringscsv.test", [][]string{
		{SlicesHeader, "tl"},
		{"test_str", "Test translation"},
	})
	require.NoError(t, err)
	assert.FileExists(t, "/tmp/androidstringscsv.test", "function didn't create file")

	vals, err := readSlicesFromCSVFile("/tmp/androidstringscsv.test")
	require.NoError(t, err)
	assert.Equal(t, [][]string{
		{SlicesHeader, "tl"},
		{"test_str", "Test translation"},
	}, vals)
}

func TestDictReadWrite(t *testing.T) {
	defer os.RemoveAll("/tmp/androidstringscsv.test")
	_, err := WriteCSVFile("/tmp/androidstringscsv.test", map[string]map[string]string{
		"tl": map[string]string{
			"test_str": "Test translation",
		},
	}, true)
	require.NoError(t, err)
	assert.FileExists(t, "/tmp/androidstringscsv.test", "function didn't create file")

	dicts, err := ReadCSVFile("/tmp/androidstringscsv.test")
	assert.Equal(t, map[string]map[string]string{
		"tl": map[string]string{
			"test_str": "Test translation",
		},
	}, dicts)
}
