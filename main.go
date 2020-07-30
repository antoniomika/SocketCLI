// Package main represents the main entrypoint of the socketcli application.
package main

import (
	"github.com/antoniomika/socketcli/cmd"
	"github.com/sirupsen/logrus"
)

// main will start the socketcli command lifecycle and spawn the socketcli services.
func main() {
	err := cmd.Execute()
	if err != nil {
		logrus.Println("Unable to execute root command:", err)
	}
}
