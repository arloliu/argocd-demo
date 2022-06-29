package helm

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/Xfers/nom-nom/argocd-helm-ext-plugin/pkg/config"
	"github.com/Xfers/nom-nom/argocd-helm-ext-plugin/pkg/utils"
)

const (
	SIMPLE_VALUE int = iota
	STRING_VALUE
)

var special_helm_variables = map[string]int{
	"image.tag": STRING_VALUE,
}

func getHelmValueType(name string) int {
	helm_val, ok := special_helm_variables[name]
	if ok {
		return helm_val
	}
	return SIMPLE_VALUE
}

type HelmCli struct {
	AppName      string
	RepoUrl      string
	Chart        string
	ChartVersion string
	ValueFiles   []string
	Values       map[string]string
	opts         *config.Options
}

func New(opts *config.Options) (*HelmCli, error) {
	values_files := strings.Split(utils.GetEnv("HELM_VALUE_FILES"), " ")
	for i, v := range values_files {
		values_files[i] = strings.TrimSpace(v)
	}

	values := make(map[string]string)
	values_env_var := strings.Split(utils.GetEnv("HELM_VALUES"), " ")
	for _, v := range values_env_var {
		token := strings.Split(v, "=")
		if len(token) != 2 {
			continue
		}
		values[strings.TrimSpace(token[0])] = strings.TrimSpace(token[1])
	}

	cli := &HelmCli{
		AppName:      os.Getenv("ARGOCD_APP_NAME"),
		RepoUrl:      utils.GetEnv("HELM_REPO_URL"),
		Chart:        utils.GetEnv("HELM_CHART"),
		ChartVersion: utils.GetEnv("HELM_CHART_VERSION"),
		ValueFiles:   values_files,
		Values:       values,
		opts:         opts,
	}

	if len(cli.AppName) == 0 {
		return nil, errors.New("ARGOCD_APP_NAME is empty")
	}
	if len(cli.RepoUrl) == 0 {
		return nil, errors.New("HELM_REPO_URL is empty")
	}
	if len(cli.Chart) == 0 {
		return nil, errors.New("HELM_CHART is empty")
	}
	if len(cli.ChartVersion) == 0 {
		return nil, errors.New("HELM_CHART_VERSION is empty")
	}

	return cli, nil
}

func (h *HelmCli) PullChart() error {
	cmd := exec.Command("helm", "pull", h.Chart, "--untar", "--repo", h.RepoUrl, "--version", h.ChartVersion)
	log.Printf("Execute: %s %s\n", cmd.Path, strings.Join(cmd.Args, " "))

	if err := cmd.Run(); err != nil {
		return errors.New(fmt.Sprintf("pull chart %s error: %s\n", h.Chart, err.Error()))
	}
	return nil
}

func (h *HelmCli) GenerateTemplate() error {
	args := []string{
		"template",
		h.AppName,
		h.Chart,
		"--repo",
		h.RepoUrl,
		"--version",
		h.ChartVersion,
	}

	if h.opts.IncludeCrds {
		args = append(args, "--include-crds")
	}

	for _, value_file := range h.ValueFiles {
		if _, err := os.Stat(value_file); err == nil {
			args = append(args, "-f")
			args = append(args, value_file)
		}
	}

	for key, value := range h.Values {
		value_type := getHelmValueType(key)
		if value_type == STRING_VALUE {
			args = append(args, "--set-string")
		} else {
			args = append(args, "--set")
		}
		args = append(args, fmt.Sprintf("%s=%s", key, value))
	}

	var out bytes.Buffer
	cmd := exec.Command("helm", args...)
	cmd.Stdout = &out

	log.Printf("Execute: %s %s\n", cmd.Path, strings.Join(cmd.Args, " "))

	if err := cmd.Run(); err != nil {
		return errors.New(fmt.Sprintf("pull chart %s error: %s\n", h.Chart, err.Error()))
	}

	fmt.Print(out.String())
	return nil
}
