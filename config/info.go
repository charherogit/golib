package config

import (
	"encoding/base64"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"runtime/debug"
	"time"
)

var (
	metaB64   string
	marshaled string

	buildInfo = struct {
		err       error
		BuildInfo *debug.BuildInfo `yaml:"build_info"`
		StartTime time.Time        `yaml:"start_time"`
		Commit    string           `yaml:"commit"`
		BuildTime string           `yaml:"build_time"`
		Submodule []string         `yaml:"submodule"`
	}{
		StartTime: time.Now(),
	}
)

func initBuildInfo() {
	i, ok := debug.ReadBuildInfo()
	if ok {
		i.Deps = nil
		buildInfo.BuildInfo = i
	}

	res, err := base64.StdEncoding.DecodeString(metaB64)
	if err != nil {
		buildInfo.err = err
		return
	}
	if err := yaml.Unmarshal(res, &buildInfo); err != nil {
		buildInfo.err = err
		return
	}
	if m, err := yaml.Marshal(&buildInfo); err != nil {
		buildInfo.err = err
		return
	} else {
		marshaled = string(m)
	}
	if len(os.Args) > 1 && os.Args[1] == `ABCDEFG` {
		_, _ = fmt.Fprint(os.Stderr, marshaled)
		os.Exit(5)
	}
}

func GetBuildInfo() string {
	if buildInfo.err != nil {
		return buildInfo.err.Error()
	}
	return marshaled
}
