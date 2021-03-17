package main

import (
	"fmt"
	"github.com/hashicorp/packer-plugin-sdk/plugin"
	ksyun "github.com/kingsoftcloud/packer-plugin-ksyun/builder/kec"
	"os"
)

func main() {
	pps := plugin.NewSet()
	pps.RegisterBuilder(plugin.DEFAULT_NAME, new(ksyun.Builder))
	err := pps.Run()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
