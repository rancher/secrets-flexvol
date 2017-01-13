package main

import (
	"os"

	flexvol "github.com/rancher/rancher-flexvol"
	"github.com/rancher/secrets-flexvol/secrets"
)

var VERSION = "v0.0.0-dev"

func main() {
	backend := &secrets.FlexVolume{}

	app := flexvol.NewApp(backend)
	app.Version = VERSION

	app.Run(os.Args)
}
