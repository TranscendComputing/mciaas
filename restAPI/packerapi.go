package restAPI

import (
	"github.com/TranscendComputing/mciaas/store"
	"github.com/ant0ine/go-json-rest"
	"github.com/peterbourgon/mergemap"
	"net/http"
)

type PackerRestAPI struct {
	storage *store.MapStore
}

func (this *PackerRestAPI) Delete(w *rest.ResponseWriter, r *rest.Request) {
	userId := r.PathParam("user")
	docId := r.PathParam("docId")
	this.storage.Open(userId)
	this.storage.DeleteDocument(userId, docId)
	w.WriteJson(&IdResponse{docId})
}

// The document of interest is the JSON payload
func (this *PackerRestAPI) Get(w *rest.ResponseWriter, r *rest.Request) {
	userId := r.PathParam("user")
	docId := r.PathParam("docId")
	(*this.storage).Open(userId)
	doc, err := (*this.storage).GetDocument(userId, docId)
	if err == nil {
		w.WriteJson(&doc)
	} else {
		ProcessError(w, err, http.StatusInternalServerError)
	}
}

// The document of interest is the JSON payload
func (this *PackerRestAPI) Put(w *rest.ResponseWriter, r *rest.Request) {
	var empty interface{}
	userId := r.PathParam("user")
	err := r.DecodeJsonPayload(&empty)
	if err == nil {
		doc := empty
		(*this.storage).Open(userId)
		docId, err := this.storage.PutDocument(userId, doc.(map[string]interface{}))
		if err == nil {
			w.WriteJson(&IdResponse{docId})
		}
	}
	ProcessError(w, err, http.StatusInternalServerError)
}

func (this *PackerRestAPI) Post(w *rest.ResponseWriter, r *rest.Request) {
	userId := r.PathParam("user")
	docId := r.PathParam("docId")

	// first, locate the existing document
	(*this.storage).Open(userId)
	doc, err := this.storage.GetDocument(userId, docId)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// second, decode the post payload to merge 'into' to prior document
	var snippet interface{}
	err = r.DecodeJsonPayload(&snippet)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// merge the documents
	merged := mergemap.Merge(doc, snippet.(map[string]interface{}))

	//replace the document in the map
	err = this.storage.UpdateDocument(userId, docId, merged)
	if err == nil {
		w.WriteJson(&IdResponse{docId})
	} else {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (this *PackerRestAPI) SetStorage(s *store.MapStore) {
	this.storage = s
}
