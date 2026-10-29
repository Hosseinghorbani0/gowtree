package main

import (
	"os"

	"github.com/hosseinghorbani0/gowtree/internal/app"
)

func main() {
	os.Exit(app.Run(os.Args[1:]))
}
