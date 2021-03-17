package main

import (
	"fmt"
	"github.com/hashicorp/packer-plugin-sdk/plugin"
	kscloud "github.com/kingsoftcloud/packer-plugin-kscloud/builder/kscloud/kec"
	"os"
)

func main() {
	pps := plugin.NewSet()
	pps.RegisterBuilder(plugin.DEFAULT_NAME, new(kscloud.Builder))
	err := pps.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
