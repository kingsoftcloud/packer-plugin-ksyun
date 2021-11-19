package main

import (
	"fmt"
	"github.com/hashicorp/packer-plugin-sdk/plugin"
	kec "github.com/kingsoftcloud/packer-plugin-ksyun/builder/kec"
	"os"
)

func main() {
	pps := plugin.NewSet()
	pps.RegisterBuilder("kec", new(kec.Builder))
	err := pps.Run()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
