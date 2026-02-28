package main

import (
	"os"

	"github.com/mgcis-cn/ibookfs/cmd/apiserver/app"

	"k8s.io/component-base/cli"
)

func main() {
	cmd := app.NewAPIServerCommand()
	code := cli.Run(cmd)
	os.Exit(code)
}
