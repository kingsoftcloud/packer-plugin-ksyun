package main

import (
	"fmt"
	"os"

	"github.com/KscSDK/ksc-sdk-go/ksc"
	"github.com/hashicorp/packer-plugin-sdk/plugin"
	"github.com/kingsoftcloud/packer-plugin-ksyun/builder/epc"
	"github.com/kingsoftcloud/packer-plugin-ksyun/builder/kec"
	"github.com/kingsoftcloud/packer-plugin-ksyun/datasource/kmi"
)

var (
	version string
)

func main() {

	ksc.SDKName = "packer-plugin-ksyun"
	ksc.SDKVersion = version

	pps := plugin.NewSet()
	pps.RegisterBuilder("kec", new(kec.Builder))
	pps.RegisterBuilder("epc", new(epc.Builder))
	pps.RegisterDatasource("kmi", new(kmi.Datasource))
	err := pps.Run()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
