package main

import (
	"log"
	"wenwenxiong/go-ipam-client/cmd/app"
)

func main() {
	cmd := app.NewIpamClientCommand()
	if err := cmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}
