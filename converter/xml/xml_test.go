package xml

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConvertions(t *testing.T) {
	r := convertDictionaryToResources(map[string]string{
		"test_str": "Test translation",
	})
	assert.Equal(t, ResourcesEntry{
		Strings: []StringEntry{
			StringEntry{
				Name:  "test_str",
				Value: "Test translation",
			},
		},
	}, r)

	d := r.ConvertToDictionary()
	assert.Equal(t, map[string]string{
		"test_str": "Test translation",
	}, d)
}

func TestReadWriteXML(t *testing.T) {
	defer os.RemoveAll("/tmp/androidstringscsv.test")

	_, err := exportDictionaryToXML("/tmp/androidstringscsv.test", map[string]string{
		"test_str": "Test translation",
	})
	require.NoError(t, err)
	assert.FileExists(t, "/tmp/androidstringscsv.test")

	readed, err := ReadXMLFile("/tmp/androidstringscsv.test")
	require.NoError(t, err)
	assert.Equal(t, map[string]string{
		"test_str": "Test translation",
	}, (*readed).ConvertToDictionary())
}

func TestReadWriteRes(t *testing.T) {
	defer os.RemoveAll("/tmp/res")
	_, err := WriteResFolder("/tmp/res", map[string]map[string]string{
		"tl": map[string]string{
			"test_str": "Test translation",
		},
	})
	require.NoError(t, err)
	assert.DirExists(t, "/tmp/res")

	dicts, err := ReadResFolder("/tmp/res")
	require.NoError(t, err)
	assert.Equal(t, map[string]map[string]string{
		"tl": map[string]string{
			"test_str": "Test translation",
		},
	}, dicts)
}
