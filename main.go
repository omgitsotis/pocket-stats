package main

import (
    "log"
    client "github.com/omgitsotis/pocket-stats/client"
)

var code string
var accessToken string

func main() {
    log.Fatal(client.ServeAPI())
}
