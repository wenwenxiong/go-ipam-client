package main

import (
	"github.com/wenwenxiong/go-ipam-client/cmd/app"
	"log"
)

func main() {
	cmd := app.NewIpamClientCommand()
	if err := cmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}
