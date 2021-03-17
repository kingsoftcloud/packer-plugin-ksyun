package kec

import (
	"fmt"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
)

type Artifact struct {
	// A map of regions to ksyun image id.
	KsyunImageId string
	// BuilderId is the unique ID for the builder that created this ksyun image
	BuilderIdValue string
	// ksyun connection for performing API stuff.
	Client *ClientWrapper
}

func (k *Artifact) BuilderId() string {
	return k.BuilderIdValue
}

func (k *Artifact) Files() []string {
	return nil
}

func (k *Artifact) Id() string {
	return k.KsyunImageId
}

func (k *Artifact) String() string {
	return fmt.Sprintf("Ksyun images were created:%s", k.KsyunImageId)
}

func (k *Artifact) State(name string) interface{} {
	return nil
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
