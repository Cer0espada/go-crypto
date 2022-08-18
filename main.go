package main

import (
	"os"

	"go-crypto/cli"
)

func main() {
	defer os.Exit(0)

	cmd := cli.CommandLine{}
	cmd.Run()
}


//concesus algorithms / proof of work algorithms -- forcing the network to do work to add work to the chain 
// bitcoin -- get fees by powering the network, makes the blocks more secure,proof of work validation, when a user does work to sign of block they show proof, the work is hard but the proof that the work was done is easy