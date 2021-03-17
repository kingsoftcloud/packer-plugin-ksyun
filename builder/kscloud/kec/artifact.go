package kec

import (
	"fmt"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
)

type Artifact struct {
	// A map of regions to kingcloud image id.
	KingcloudImageId string
	// BuilderId is the unique ID for the builder that created this kingcloud image
	BuilderIdValue string
	// kingcloud connection for performing API stuff.
	Client *ClientWrapper
}

func (k*Artifact) BuilderId() string {
	return k.BuilderIdValue
}

func (k*Artifact) Files() []string {
	return nil
}

func (k*Artifact) Id() string{
	return  k.KingcloudImageId
}

func (k*Artifact) String() string {
	return fmt.Sprintf("Kingcloud images were created:%s", k.KingcloudImageId)
}

func (k*Artifact) State(name string) interface{}{
	return nil
}

func (k*Artifact) Destroy() error{
	errors := make([]error, 0)
	//delete
	//removeImages:=make(map[string]interface{})
	//removeImages["ImageId"]=k.KingcloudImageId
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