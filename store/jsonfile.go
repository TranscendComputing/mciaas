package store

import (
	"encoding/json"
	"github.com/ant0ine/go-json-rest"
	"io/ioutil"
	"os"
)

func LoadJSONFile(jsonFile string) (map[string]interface{}, error) {
	f, err := os.Open(jsonFile)
	if err != nil {
		return nil, err
	}

	dec := json.NewDecoder(f)
	var m map[string]interface{}
	if err = dec.Decode(&m); err != nil {
		return nil, err
	}

	return m, nil
}

func WriteJSONFile(path string, jsonOb interface{}) error {
	if bytes, err := json.Marshal(jsonOb); err == nil {
		return ioutil.WriteFile(path, bytes, 0660)
	} else {
		return err
	}
}

func SendJSONFile(w *rest.ResponseWriter, r *rest.Request, filePath string) {
	m, err := LoadJSONFile(filePath)
	if err == nil {
		w.WriteJson(&m)
	} else {
		m["Error"] = err
	}
}
