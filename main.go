package main

import (
	"github.com/banditml/goat/app"
)

func main() {
	// Run blocks forever doing nothing, waiting for SIGINT/SIGTERM signals.
	app.New().Run()
}
