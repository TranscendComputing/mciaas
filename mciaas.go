// This is the main package for the mciaas application.
package main

import (
	"flag"
	"fmt"
	"github.com/TranscendComputing/mciaas/restAPI"
	"github.com/TranscendComputing/mciaas/store"
	"github.com/ant0ine/go-json-rest"
	"net/http"
	"os"
)

var listenPort int
var rootPath string

func realMain() int {
	genericTemplates := JSONTemplates{"templates"}
	packerTemplates := restAPI.PackerRestAPI{}
	storage := &store.MapStore{}
	packerTemplates.SetStorage(storage)
	images := restAPI.ImageAPI{}
	images.SetStorage(storage)
	images.SetRootPath(rootPath)

	handler := rest.ResourceHandler{
		EnableRelaxedContentType: true,
	}
	handler.SetRoutes(
		rest.RouteObjectMethod("GET", "/templates", &genericTemplates, "ListTemplateTypes"),
		rest.RouteObjectMethod("GET", "/templates/builders", &genericTemplates, "ListBuilderTemplates"),
		rest.RouteObjectMethod("GET", "/templates/builders/:type", &genericTemplates, "GetBuilderTemplate"),
		rest.RouteObjectMethod("GET", "/templates/provisioners", &genericTemplates, "ListProvisionerTemplates"),
		rest.RouteObjectMethod("GET", "/templates/provisioners/:type", &genericTemplates, "GetProvisionerTemplate"),
		rest.RouteObjectMethod("GET", "/templates/postprocessors", &genericTemplates, "ListPostprocessorTemplates"),
		rest.RouteObjectMethod("GET", "/templates/postprocessors/:type", &genericTemplates, "ListPostprocessorTemplates"),
		rest.RouteObjectMethod("DELETE", "/packer/:user/:docId", &packerTemplates, "Delete"),
		rest.RouteObjectMethod("GET", "/packer/:user/:docId", &packerTemplates, "Get"),
		rest.RouteObjectMethod("POST", "/packer/:user/:docId", &packerTemplates, "Post"),
		rest.RouteObjectMethod("PUT", "/packer/:user", &packerTemplates, "Put"),
		rest.RouteObjectMethod("DELETE", "/image/:user/:docId", &images, "Delete"),
		rest.RouteObjectMethod("GET", "/image/:user/:docId", &images, "Get"),
		rest.RouteObjectMethod("POST", "/image/:user/:docId", &images, "Post"),
		rest.RouteObjectMethod("PUT", "/image/:user/:docId", &images, "Put"),
	)
	http.ListenAndServe(fmt.Sprintf(":%v", listenPort), &handler)

	return 0
}

func initFlags() {
	const (
		defaultPort = 8080
		portUsage   = "port on which to listen for HTTP(S) requests"
		defaultRoot = "/tmp"
		rootUsage   = "root directory in which to run all Packer operations."
	)
	flag.IntVar(&listenPort, "port", defaultPort, portUsage)
	flag.IntVar(&listenPort, "p", defaultPort, portUsage+" (shorthand)")
	flag.StringVar(&rootPath, "rootpath", defaultRoot, portUsage)
	flag.StringVar(&rootPath, "r", defaultRoot, rootUsage+" (shorthand)")
}

func main() {
	// Delegate to realMain so defer operations can happen (os.Exit exits
	// the program without servicing defer statements)
	initFlags()
	flag.Parse()
	os.Exit(realMain())
}
