package main

import (
	"flag"

	"github.com/kristofferostlund/no-store/server/server"
)

var address = flag.String(
	"address",
	"0.0.0.0",
	"The address to run the server on",
)

var port = flag.Int(
	"port",
	4000,
	"The port to run the server on",
)

func main() {
	flag.Parse()

	server.Serve(*address, *port)
}
