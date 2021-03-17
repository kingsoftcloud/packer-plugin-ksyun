package kec

import (
	"context"
	"fmt"
	"github.com/hashicorp/packer-plugin-sdk/communicator"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
)

type stepConfigKingcloudKeyPair struct {
	KingcloudRunConfig *KingcloudRunConfig
	Comm  *communicator.Config
}

func (s *stepConfigKingcloudKeyPair) Run(ctx context.Context, stateBag multistep.StateBag) multistep.StepAction {
	ui := stateBag.Get("ui").(packersdk.Ui)
	client := stateBag.Get("client").(*ClientWrapper)
	if s.Comm.SSHPrivateKeyFile != "" {
		if s.Comm.SSHKeyPairName == "" {
			return Halt(stateBag,
				fmt.Errorf(fmt.Sprintf("ssh_keypair_name is empty")), "")
		}
		ui.Say("Using existing SSH private key")
		privateKeyBytes, err := s.Comm.ReadSSHPrivateKeyFile()
		if err != nil {
			stateBag.Put("error", err)
			return multistep.ActionHalt
		}
		s.Comm.SSHPrivateKey = privateKeyBytes
	}

	if s.Comm.SSHAgentAuth && s.Comm.SSHKeyPairName == "" {
		ui.Say("Using SSH Agent with key pair in source image")
		return multistep.ActionContinue
	}

	if s.Comm.SSHAgentAuth && s.Comm.SSHKeyPairName != "" {
		ui.Say(fmt.Sprintf("Using SSH Agent for existing key pair %s", s.Comm.SSHKeyPairName))
		return multistep.ActionContinue
	}

	if s.Comm.SSHTemporaryKeyPairName == "" {
		ui.Say("Not using temporary keypair")
		s.Comm.SSHKeyPairName = ""
		return multistep.ActionContinue
	}

	ui.Say(fmt.Sprintf("Using SSH Agent for create new key pair %s", s.Comm.SSHTemporaryKeyPairName))
	//create ssh Key
	createSSHKey := make(map[string]interface{})
	createSSHKey["KeyName"] = s.Comm.SSHTemporaryKeyPairName
	resp, err :=client.SksClient.CreateKey(&createSSHKey)
	if err != nil {
		return Halt(stateBag, err, "Error creating new keypair")
	}
	if resp != nil {
		s.Comm.SSHKeyPairName = getSdkValue(stateBag,"Key.KeyId",*resp).(string)
		privateKey := getSdkValue(stateBag,"PrivateKey",*resp).(string)
		s.Comm.SSHPrivateKey = []byte(privateKey)
	}
	return multistep.ActionContinue
}

func (s *stepConfigKingcloudKeyPair) Cleanup(stateBag multistep.StateBag) {
}

