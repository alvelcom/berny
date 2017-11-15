package main

import (
	"log"
	"os"

	"github.com/alvelcom/redoubt/cmd/harvest"
	"github.com/alvelcom/redoubt/cmd/serve"
)

func init() {
	log.SetFlags(0)
}

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Command name is expected")
	}

	switch os.Args[1] {
	case "serve":
		serve.Main(os.Args[2:])
	case "harvest":
		harvest.Main(os.Args[2:])
	}
}
