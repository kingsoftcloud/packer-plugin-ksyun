package epc

import (
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
	ksyun "github.com/kingsoftcloud/packer-plugin-ksyun/builder"
)

type Artifact struct {
	*ksyun.Artifact
	// ksyun connection for performing API stuff.
	Client *ClientEpcWrapper
}

func (k *Artifact) Destroy() error {
	errors := make([]error, 0)
	//delete
	//removeImages:=make(map[string]interface{})
	//removeImages["ImageId"]=k.KsyunImageId
	//_, err := k.Client.RemoveImages(&removeImages)
	//if err != nil {
	//	errors = append(errors, err)
	//}
	if len(errors) > 0 {
		if len(errors) == 1 {
			return errors[0]
		} else {
			return &packersdk.MultiError{Errors: errors}
		}
	}
	return nil
}
