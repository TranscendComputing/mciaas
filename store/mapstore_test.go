package store

import (
	"testing"
)

func TestMapStore(t *testing.T) {
	testMap := map[string]interface{}{
		"test1key": "test1value",
	}

	storage := MapStore{}

	// creates the document set in the storage
	_, err := storage.Open("test")

	id, err := storage.PutDocument("test", testMap)
	if err != nil {
		t.Errorf("Expected to store test document without error: %s.", id)
	}

	testMap, err = storage.GetDocument("test", id)
	if err != nil || testMap["test1key"] != "test1value" {
		t.Errorf("Expected to retreieve test document without error.")
	}

	err = storage.DeleteDocument("test", id)
	if err != nil {
		t.Errorf("Expected to retreieve test document without error.")
	}
}
