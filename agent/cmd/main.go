package main

import (
	"os"

	"agent/internal/cli"

	_ "agent/api"
)

// @title           Monitoring Agent API
// @version         1.0
// @host            localhost:8088
// @BasePath        /
func main() {
	os.Exit(cli.New().Execute())

}
