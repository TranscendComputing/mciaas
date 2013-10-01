package main

import (
	"github.com/TranscendComputing/mciaas/store"
	"github.com/ant0ine/go-json-rest"
	"path/filepath"
)

type JSONTemplates struct {
	BasePath string
}

func (this *JSONTemplates) ListTemplateTypes(w *rest.ResponseWriter, r *rest.Request) {
	store.SendJSONFile(w, r, filepath.Join(this.BasePath, "list.json"))
}

func (this *JSONTemplates) ListBuilderTemplates(w *rest.ResponseWriter, r *rest.Request) {
	store.SendJSONFile(w, r, filepath.Join(this.BasePath, "builders", "list.json"))
}

func (this *JSONTemplates) GetBuilderTemplate(w *rest.ResponseWriter, r *rest.Request) {
	builder := r.PathParam("type")
	jsonPath := filepath.Join(this.BasePath, "builders", builder+".json")
	store.SendJSONFile(w, r, jsonPath)
}

func (this *JSONTemplates) ListProvisionerTemplates(w *rest.ResponseWriter, r *rest.Request) {
	store.SendJSONFile(w, r, filepath.Join(this.BasePath, "provisioners", "list.json"))
}

func (this *JSONTemplates) GetProvisionerTemplate(w *rest.ResponseWriter, r *rest.Request) {
	provisioner := r.PathParam("type")
	jsonPath := filepath.Join(this.BasePath, "provisioners", provisioner+".json")
	store.SendJSONFile(w, r, jsonPath)
}

func (this *JSONTemplates) ListPostprocessorTemplates(w *rest.ResponseWriter, r *rest.Request) {
	store.SendJSONFile(w, r, filepath.Join(this.BasePath, "postprocessors", "list.json"))
}

func (this *JSONTemplates) ListPostprocessorTemplate(w *rest.ResponseWriter, r *rest.Request) {
	postprocessor := r.PathParam("type")
	jsonPath := filepath.Join(this.BasePath, "postprocessors", postprocessor+".json")
	store.SendJSONFile(w, r, jsonPath)
}
