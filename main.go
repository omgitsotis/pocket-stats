package main

import (
	"log"

	client "github.com/omgitsotis/pocket-stats/client"
)

func main() {
	log.Fatal(client.ServeAPI())
}
