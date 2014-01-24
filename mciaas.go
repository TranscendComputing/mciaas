// This is the main package for the mciaas application.
package main

import (
	"os"
)

var app ApplicationContext

func main() {
	// Delegate to realMain so defer operations can happen (os.Exit exits
	// the program without servicing defer statements)
	app.initialize()
	os.Exit(app.realMain())
}
