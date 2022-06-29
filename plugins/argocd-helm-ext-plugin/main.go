package main

import (
	"os"

	"github.com/Xfers/nom-nom/argocd-helm-ext-plugin/cmd"
)

func main() {
	// os.Setenv("ARGOCD_APP_NAME", "demo-app")
	// os.Setenv("HELM_REPO_URL", "http://helm-repo.local/charts/packages")
	// os.Setenv("HELM_CHART", "nginx")
	// os.Setenv("HELM_CHART_VERSION", "12.0.5")
	// os.Setenv("HELM_VALUE_FILES", "values.yaml secrets.yaml")
	// os.Setenv("HELM_VALUES", "image.tag=123 slack.ts=abdcd")

	cmd.ArgoCDHelmExtPlugin()

	os.Exit(0)
}
