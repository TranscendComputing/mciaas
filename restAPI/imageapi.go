package restAPI

import (
	"bytes"
	"fmt"
	"github.com/TranscendComputing/mciaas/store"
	"github.com/ant0ine/go-json-rest"
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

func (this *ImageAPI) setOverrides(
	doc map[string]interface{},
	userId string,
	docId string) map[string]interface{} {

	// set any output_directory and port info in builders
	// note, this will fail hard (panic) if no builders exist
	builders := doc["builders"]
	for idx := range builders.([]interface{}) {
		b := builders.([]interface{})[idx]
		m := b.(map[string]interface{})
		if m["type"] == "qemu" {
			builderName := m["name"]
			m["output_directory"] = fmt.Sprintf("output_%s", builderName)
			m["http_directory"] = filepath.Join(this.rootPath,
				userId, docId, "httpfiles")
			m["http_port_min"] = 10000
			m["http_port_max"] = 10999
			m["ssh_host_port_min"] = 11000
			m["ssh_host_port_max"] = 11999
		}
	}

	return doc
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
	merged := this.setOverrides(doc, userId, docId)

	// store the doc merged document into a temporary file
	path := filepath.Join(this.rootPath, userId, docId)
	fmt.Printf("imageAPI: creating dir: %s", path)
	err = os.MkdirAll(path, 0750)
	if err == nil {
		path = filepath.Join(path, "build.json")
		fmt.Printf("imageAPI: writing file: %s", path)
		err = store.WriteJSONFile(path, merged)

		if err == nil {
			// run Packer
			fmt.Printf("imageAPI: running packer with json file: %s", path)
			err = this.runPacker(path)
		}
	}

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
		w.WriteJson(fileList)
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
