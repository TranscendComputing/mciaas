package store

import (
	"fmt"
	"github.com/nu7hatch/gouuid"
)

// A DocumentSet is a map of Documents keyed by string (uuid)
type DocumentSet map[string]interface{}

// A Store interface must support the following methods
type Store interface {
	Open(id string) (*DocumentSet, error)
	Close() error
	GetDocument(docSetId string, docId string) (map[string]interface{}, error)
	DeleteDocument(docSetId string, docId string) error
	PutDocument(docSetId string, doc map[string]interface{}) (string, error)
	UpdateDocument(docSetId string, docId string, doc map[string]interface{}) error
}

func GenerateUUID() (string, error) {
	u4, err := uuid.NewV4()

	if err != nil {
		fmt.Println("error:", err)
	}

	return u4.String(), err
}
