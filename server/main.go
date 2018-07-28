package main

import (
	"log"

	client "github.com/omgitsotis/pocket-stats/server/client"
)

func main() {
	log.Fatal(client.ServeAPI())
}
