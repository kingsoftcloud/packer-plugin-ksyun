package ksyun

import "fmt"

type Artifact struct {
	// A map of regions to ksyun image id.
	KsyunImageId string
	// BuilderId is the unique ID for the builder that created this ksyun image
	BuilderIdValue string
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
	return nil
}
