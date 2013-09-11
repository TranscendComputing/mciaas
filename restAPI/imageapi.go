package restAPI

import (
	"github.com/TranscendComputing/mciaas/store"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ant0ine/go-json-rest"
	"github.com/peterbourgon/mergemap"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

type ImageAPI struct {
	rootPath string
	storage  store.Store
}

type FileInfo struct {
	name string
	size int64
}

type FileList []FileInfo

func (this *ImageAPI) runPacker(path string) error {
	var stdout, stderr bytes.Buffer

	cmd := exec.Command("packer", "build", path)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Dir = this.rootPath
	err := cmd.Start()

	if err == nil {
		// make the channel to watch the process
		waitDone := make(chan error)

		// start packer in the background
		go func() {
			waitDone <- cmd.Wait()
		}()
	}

	return err
}

func (this *ImageAPI) Delete(w *rest.ResponseWriter, r *rest.Request) {
	userId := r.PathParam("user")
	docId := r.PathParam("docId")
	dir := filepath.Join(this.rootPath, userId, docId)
	err := os.RemoveAll(dir)
	if err == nil {
		w.WriteJson(IdResponse{docId})
	} else {
		rest.Error(w, err.Error(), http.StatusNotFound)
	}
}

// Put on the imageAPI object requests that Packer create the
// physical image file and any ancillary files Packer produces.
func (this *ImageAPI) Put(w *rest.ResponseWriter, r *rest.Request) {
	// get the requested run file
	userId := r.PathParam("user")
	docId := r.PathParam("docId")

	this.storage.Open(userId)
	doc, err := this.storage.GetDocument(userId, docId)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// force non-overridable parameters in the packer template
	merge := fmt.Sprintf(
		"{{\"builders\":{\"qemu\":{\"output_directory\":%s}}}}",
		filepath.Join(this.rootPath, userId, docId))
	bytes := []byte(merge)

	var m map[string]interface{}
	if err = json.Unmarshal(bytes, &m); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// merge the documents
	merged := mergemap.Merge(doc, m)

	// store the doc into a temporary file
	path := filepath.Join("/tmp", docId)
	store.WriteJSONFile(path, merged)

	// run Packer
	err = this.runPacker(path)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		w.WriteJson(IdResponse{docId})
	}
}

func (this *ImageAPI) Get(w *rest.ResponseWriter, r *rest.Request) {
	userId := r.PathParam("user")
	docId := r.PathParam("docId")
	dir := filepath.Join(this.rootPath, userId, docId)
	if dirList, err := ioutil.ReadDir(dir); err == nil {
		fileList := make(map[string]interface{}, 0)
		for _, fileInfo := range dirList {
			fileList[fileInfo.Name()] = fileInfo.Size()
		}
	} else {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (this *ImageAPI) Post(w *rest.ResponseWriter, r *rest.Request) {
	this.Put(w, r)
}

func (this *ImageAPI) GetImageFile(w *rest.ResponseWriter, r *rest.Request) {
}

func (this *ImageAPI) SetStorage(storage store.Store) {
	this.storage = storage
}

func (this *ImageAPI) SetRootPath(rootPath string) {
	this.rootPath = rootPath
}
