package store

import (
	"errors"
	"fmt"
)

// A basic storage for general use -- a map of document sets
type MapStore map[string]*DocumentSet

// Store methods
func (this *MapStore) Open(id string) (*DocumentSet, error) {
	if ds, ok := (*this)[id]; ok {
		return ds, nil
	} else {
		ds = &DocumentSet{}
		(*this)[id] = ds
		return ds, nil
	}
}

func (this *MapStore) Close() error {
	return nil
}

func (this *MapStore) GetDocument(docSetId string, docId string) (map[string]interface{}, error) {
	var docRet map[string]interface{}
	var errRet error

	if dsTest, ok := (*this)[docSetId]; ok {
		if docTest, ok := (*dsTest)[docId]; ok {
			docRet = docTest.(map[string]interface{})
		} else {
			errMsg := fmt.Sprintf("Document does not exist for id: %s", docId)
			errRet = errors.New(errMsg)
		}
	}
	return docRet, errRet
}

// Deletes a document from the store if it exists
func (this *MapStore) DeleteDocument(docSetId string, docId string) {
	if dsTest, ok := (*this)[docSetId]; ok {
		delete(*dsTest, docId)
	}
}

// Puts a new document into the store
func (this *MapStore) PutDocument(docSetId string, doc map[string]interface{}) (string, error) {
	var idRet string
	var errRet error

	if dsTest, ok := (*this)[docSetId]; ok {
		if u, err := GenerateUUID(); err != nil {
			fmt.Println("error:", err)
			errRet = err
		} else {
			(*dsTest)[u] = doc
			(*this)[u] = dsTest
			idRet = u
		}
	} else {
		errMsg := fmt.Sprintf("Document set does not exist for id: %s", docSetId)
		errRet = errors.New(errMsg)
	}

	return idRet, errRet
}

// replaces an existing Document with the new interface (Document)
func (this *MapStore) UpdateDocument(docSetId string, docId string, doc map[string]interface{}) error {
	var errRet error

	if _, err := this.GetDocument(docSetId, docId); err == nil {
		ds := (*this)[docSetId]
		(*ds)[docId] = doc
	} else {
		errMsg := fmt.Sprintf("Document does not exist for id: %s in set %s", docId, docSetId)
		errRet = errors.New(errMsg)
	}

	return errRet
}
