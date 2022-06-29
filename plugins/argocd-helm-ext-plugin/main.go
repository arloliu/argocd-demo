package main

import (
	"os"

	"github.com/Xfers/nom-nom/argocd-helm-ext-plugin/cmd"
)

func main() {
	cmd.ArgoCDHelmExtPlugin()
	os.Exit(0)
}
