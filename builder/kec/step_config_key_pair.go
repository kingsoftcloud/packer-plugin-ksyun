package kec

import (
	"context"
	"fmt"
	"github.com/hashicorp/packer-plugin-sdk/communicator"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/kingsoftcloud/packer-plugin-ksyun/builder"
	"strings"
)

type stepConfigKsyunKeyPair struct {
	KsyunRunConfig        *KsyunKecRunConfig
	Comm                  *communicator.Config
	keyId                 string
	SSHTemporaryPublicKey *string
}

func (s *stepConfigKsyunKeyPair) Run(ctx context.Context, stateBag multistep.StateBag) multistep.StepAction {
	ui := stateBag.Get("ui").(packersdk.Ui)
	client := stateBag.Get("client").(*ClientWrapper)
	if s.Comm.SSHPrivateKeyFile != "" {
		if s.Comm.SSHKeyPairName == "" {
			return ksyun.Halt(stateBag,
				fmt.Errorf(fmt.Sprintf("ssh_keypair_name is empty")), "")
		}
		ui.Say("Using existing SSH private key")
		privateKeyBytes, err := s.Comm.ReadSSHPrivateKeyFile()
		if err != nil {
			stateBag.Put("error", err)
			return multistep.ActionHalt
		}
		s.Comm.SSHPrivateKey = privateKeyBytes
		return multistep.ActionContinue
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
	resp, err := client.SksClient.CreateKey(&createSSHKey)
	if err != nil {
		return ksyun.Halt(stateBag, err, "Error creating new keypair")
	}
	if resp != nil {
		s.Comm.SSHKeyPairName = ksyun.GetSdkValue(stateBag, "Key.KeyId", *resp).(string)
		privateKey := ksyun.GetSdkValue(stateBag, "PrivateKey", *resp).(string)
		publicKey := ksyun.GetSdkValue(stateBag, "Key.PublicKey", *resp).(string)
		s.Comm.SSHPrivateKey = []byte(privateKey)
		*(s.SSHTemporaryPublicKey) = strings.Split(publicKey, " ")[1]
		s.keyId = s.Comm.SSHKeyPairName
	}
	return multistep.ActionContinue
}

func (s *stepConfigKsyunKeyPair) Cleanup(stateBag multistep.StateBag) {
	// if key not create by packer plugin no need to clean
	if s.keyId != "" {
		client := stateBag.Get("client").(*ClientWrapper)
		ui := stateBag.Get("ui").(packersdk.Ui)
		//delete ssh Key
		ui.Say(fmt.Sprintf("Deleting temporary keypair %s ", s.Comm.SSHKeyPairName))
		deleteSSHKey := make(map[string]interface{})
		deleteSSHKey["KeyId"] = s.keyId
		_, err := client.SksClient.DeleteKey(&deleteSSHKey)
		if err != nil {
			ui.Error(fmt.Sprintf(
				"Error cleaning up keypair. Please delete the key manually: name = %s, id = %s",
				s.Comm.SSHTemporaryKeyPairName, s.keyId))
		}
	}
}
