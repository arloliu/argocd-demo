package cmd

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/Xfers/nom-nom/argocd-helm-ext-plugin/pkg/config"
	"github.com/Xfers/nom-nom/argocd-helm-ext-plugin/pkg/helm"
)

const usage = `Usage of argocd-helm-ext-plugin:
  -d, --development  development mode
  -c, --include-crds include custom resource
  -h, --help         prints help information
`

func ArgoCDHelmExtPlugin() {
	var dev_mode bool
	var opts config.Options

	flag.BoolVar(&dev_mode, "d", false, "development mode")
	flag.BoolVar(&dev_mode, "development", false, "development mode")
	flag.BoolVar(&opts.IncludeCrds, "c", false, "include custom resource")
	flag.BoolVar(&opts.IncludeCrds, "include-crds", false, "include custom resource")

	flag.Usage = func() { fmt.Print(usage) }

	flag.Parse()

	if !dev_mode {
		logFile, err := os.OpenFile("/tmp/plugin.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			os.Exit(1)
		}
		log.SetOutput(logFile)
	}
	log.SetFlags(log.LstdFlags)

	cli, err := helm.New(&opts)
	if err != nil {
		log.Fatalf("create helm cli error: %s\n", err.Error())
	}

	if err := cli.GenerateTemplate(); err != nil {
		log.Fatalf("Generate %s v%s template fail\n", cli.Chart, cli.ChartVersion)
	}
}
