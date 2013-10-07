package restAPI

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
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

func (this *ImageAPI) runPacker(
	jsonPath string,
	cwdPath string) error {

	var stdout, stderr bytes.Buffer

	// setup logging and the (iso) cache directory
	log := "PACKER_LOG=1"
	logDir := filepath.Dir(jsonPath)
	logPath := fmt.Sprintf("PACKER_LOG_PATH=%s",
		filepath.Join(logDir, "packer.log"))
	cacheDir := fmt.Sprintf("PACKER_CACHE_DIR=%s",
		filepath.Join(this.rootPath, "packer_cache"))
	env := append(os.Environ(), log, logPath, cacheDir)

	cmd := exec.Command("packer", "build", jsonPath)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Dir = cwdPath
	cmd.Env = env
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

func writeFile(filePath string, bytes []byte) error {
	dir := filepath.Dir(filePath)
	fmt.Printf("INFO: creating %s", dir)
	if err := os.MkdirAll(dir, 0755); err != nil {
		fmt.Printf("ERROR: %s\n", err)
		return err
	}

	fmt.Printf("INFO: opening %s\n", filePath)
	fo, err := os.Create(filePath)
	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
		return err
	}
	defer fo.Close()

	_, err = fo.Write(bytes)
	return err
}

func processFiles(tgtMap map[string]interface{}, tgtDir string) error {
	files := tgtMap["mciaas_files"]
	if files == nil {
		return nil
	}

	for fileName, val := range files.(map[string]interface{}) {
		m := val.(map[string]interface{})
		content := m["content"].(string)
		contentType := m["type"].(string)
		var bytes []byte
		if contentType == "base64" {
			decoded, err := base64.StdEncoding.DecodeString(content)
			if err == nil {
				bytes = decoded
			} else {
				return err
			}
		} else {
			bytes = []byte(content)
		}

		// write the file
		filePath := filepath.Join(tgtDir, fileName)
		if err := writeFile(filePath, bytes); err != nil {
			return err
		}
	}

	return nil
}

// Return a (deep) copy of a document.
// For now, this is an unoptimized copy until time permits
// writing a deepcopy package. Most public deep copies depend on
// some functionality from the 'unsafe' package that no longer exist.
func deepcopy(doc map[string]interface{}) (map[string]interface{}, error) {
	// empty interface into which to place the copy
	var docCopy interface{}

	if bytes, err := json.Marshal(doc); err != nil {
		return nil, err
	} else {
		if err := json.Unmarshal(bytes, &docCopy); err != nil {
			return nil, err
		}
	}

	return docCopy.(map[string]interface{}), nil
}

func (this *ImageAPI) setOverrides(
	doc map[string]interface{},
	userId string,
	docId string) (map[string]interface{}, error) {

	if merged, err := deepcopy(doc); err == nil {
		// set any output_directory and port info in builders
		// note, this will fail hard (panic) if no builders exist
		subMap := merged["builders"]
		for idx := range subMap.([]interface{}) {
			b := subMap.([]interface{})[idx]
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
			err = processFiles(m, m["http_directory"].(string))
			delete(m, "mciaas_files")
		}
		return merged, err
	} else {
		return nil, err
	}
}

func (this *ImageAPI) Delete(w *rest.ResponseWriter, r *rest.Request) {
	userId := r.PathParam("user")
	docId := r.PathParam("docId")
	dir := filepath.Join(this.rootPath, userId, docId)
	err := os.RemoveAll(dir)
	if err == nil {
		err = w.WriteJson(IdResponse{docId})
	}
	ProcessError(w, err, http.StatusNotFound)
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
		ProcessError(w, err, http.StatusInternalServerError)
		return
	}

	// force non-overridable parameters in the packer template
	merged, err := this.setOverrides(doc, userId, docId)
	if err != nil {
		ProcessError(w, err, http.StatusInternalServerError)
		return
	}

	// store the doc merged document into a temporary file
	buildPath := filepath.Join(this.rootPath, userId, docId)
	fmt.Printf("imageAPI: creating dir: %s\n", buildPath)
	err = os.MkdirAll(buildPath, 0750)
	if err == nil {
		jsonPath := filepath.Join(buildPath, "build.json")
		fmt.Printf("imageAPI: writing file: %s\n", jsonPath)
		err = store.WriteJSONFile(jsonPath, merged)

		if err == nil {
			// run Packer
			fmt.Printf("imageAPI: running packer with json file: %s\n",
				jsonPath)
			err = this.runPacker(jsonPath, buildPath)
		}
	}

	if err == nil {
		w.WriteJson(IdResponse{docId})
	} else {
		ProcessError(w, err, http.StatusInternalServerError)
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
		ProcessError(w, err, http.StatusInternalServerError)
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
