package main

import (
	"os"

	"github.com/goccy/go-yaml"
	"github.com/pulsejet/go-ndn-dv/cmd"
	"github.com/zjkmxy/go-ndn/pkg/log"
)

func main() {
	var cfgFile string = "/etc/ndn/dv.yml"
	if len(os.Args) >= 2 {
		cfgFile = os.Args[1]
	}

	cfgBytes, err := os.ReadFile(cfgFile)
	if err != nil {
		panic(err)
	}

	dc := cmd.DefaultConfig()
	if err = yaml.Unmarshal(cfgBytes, &dc); err != nil {
		panic(err)
	}

	log.SetLevel(log.InfoLevel)

	dve, err := cmd.NewDvExecutor(dc)
	if err != nil {
		panic(err)
	}
	if err = dve.Start(); err != nil {
		panic(err)
	}
}
