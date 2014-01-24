package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/coreos/go-etcd/etcd"
	"os"
	"strings"
)

const (
	keyRoot = "/mciaas"
)

type handlerFunc func(client *etcd.Client, store *EtcdStore, data *handlerData) (*etcd.Response, error)

type handlerData struct {
	docSetId string
	docId    string
	doc      map[string]interface{}
}

// A basic storage for general use -- a map of document sets
type EtcdStore struct {
	Peers   string // of the form "ipaddr:port[,ipaddr:port[...]]"
	Debug   bool   // if true, dump debug, otherwise silent
	KeyBase string // the base key hierarchy
}

// dumpCURL blindly dumps all curl output to os.Stderr
func dumpCURL(client *etcd.Client) {
	client.OpenCURL()
	for {
		fmt.Fprintf(os.Stderr, "Curl-Example: %s\n", client.RecvCURL())
	}
}

func trimsplit(s, sep string) []string {
	raw := strings.Split(s, ",")
	trimmed := make([]string, 0)
	for _, r := range raw {
		trimmed = append(trimmed, strings.TrimSpace(r))
	}
	return trimmed
}

// wraps the function handlers to for etcd requests.
func (s *EtcdStore) etcdHandler(fn handlerFunc, data *handlerData) (*etcd.Response, error) {
	client := etcd.NewClient(trimsplit(s.Peers, ","))

	if s.Debug {
		go dumpCURL(client)
	}

	// Sync cluster.
	ok := client.SyncCluster()
	if s.Debug {
		fmt.Fprintf(os.Stderr, "Cluster-Peers: %s\n",
			strings.Join(client.GetCluster(), " "))
	}

	if !ok {
		errMsg := fmt.Sprintf("Cannot sync with etcd cluster.")
		return nil, errors.New(errMsg)
	} else {
		// Execute handler function.
		return fn(client, s, data)
	}
}

// gets all keys and subdirectories from a directory
func enumDir(client *etcd.Client, store *EtcdStore, data *handlerData) (*etcd.Response, error) {
	key := store.KeyBase + "/" + data.docSetId + "/"
	recursive := true
	sort := false
	return client.Get(key, sort, recursive)
}

func mkDir(client *etcd.Client, store *EtcdStore, data *handlerData) (*etcd.Response, error) {
	key := store.KeyBase + "/" + data.docSetId
	ttl := 0
	return client.CreateDir(key, uint64(ttl))
}

// sets a document into specified subdirectory
func setDoc(client *etcd.Client, store *EtcdStore, data *handlerData) (*etcd.Response, error) {
	key := store.KeyBase + "/" + data.docSetId + "/" + data.docId

	bytes, err := json.Marshal(data.doc)
	if err != nil {
		fmt.Println("error:", err)
		return nil, err
	}
	value := string(bytes)

	ttl := 0
	return client.Set(key, value, uint64(ttl))
}

// sets a document into specified subdirectory
func deleteDoc(client *etcd.Client, store *EtcdStore, data *handlerData) (*etcd.Response, error) {
	key := store.KeyBase + "/" + data.docSetId + "/" + data.docId
	return client.Delete(key, true)
}

// sets a document into specified subdirectory
func getDoc(client *etcd.Client, store *EtcdStore, data *handlerData) (*etcd.Response, error) {
	key := store.KeyBase + "/" + data.docSetId + "/" + data.docId
	return client.Get(key, false, false)
}

// Store methods
func (this *EtcdStore) Open(id string) error {
	data := &handlerData{
		docSetId: id,
	}
	resp, err := this.etcdHandler(enumDir, data)
	if err == nil && !resp.Node.Dir {
		return errors.New("The key exists but is not a directory.")
	}

	// don't care if the key doesn't exist as we can create it later on a put
	return nil
}

func (this *EtcdStore) Close() error {
	return nil
}

func (this *EtcdStore) GetDocument(docSetId string, docId string) (map[string]interface{}, error) {
	data := &handlerData{
		docSetId: docSetId,
		docId:    docId,
	}
	resp, err := this.etcdHandler(getDoc, data)
	if err != nil {
		return nil, errors.New("Not implemented yet!")
	}

	value := resp.Node.Value
	var doc map[string]interface{}
	if err := json.Unmarshal([]byte(value), &doc); err == nil {
		return doc, err
	} else {
		return nil, err
	}
}

// Deletes a document from the store if it exists
func (this *EtcdStore) DeleteDocument(docSetId string, docId string) error {
	data := &handlerData{
		docSetId: docSetId,
		docId:    docId,
	}
	_, err := this.etcdHandler(deleteDoc, data)
	return err
}

// Puts a new document into the store
func (this *EtcdStore) PutDocument(docSetId string, doc map[string]interface{}) (string, error) {
	key, err := GenerateUUID()
	if err != nil {
		fmt.Println("error:", err)
		return "", err
	}

	data := &handlerData{
		docSetId: docSetId,
		docId:    key,
		doc:      doc,
	}
	if _, err = this.etcdHandler(setDoc, data); err == nil {
		return key, nil
	} else {
		return "", err
	}
}

// replaces an existing Document with the new interface (Document)
func (this *EtcdStore) UpdateDocument(docSetId string, docId string, doc map[string]interface{}) error {
	return errors.New("Not implemented yet!")
}
