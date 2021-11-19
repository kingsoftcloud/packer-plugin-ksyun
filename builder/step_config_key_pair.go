package ksyun

import (
	"context"
	"fmt"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
	"strings"
)

type StepConfigKsyunKeyPair struct {
	CommonConfig          *CommonConfig
	keyId                 string
	SSHTemporaryPublicKey *string
}

func (s *StepConfigKsyunKeyPair) Run(ctx context.Context, stateBag multistep.StateBag) multistep.StepAction {
	ui := stateBag.Get("ui").(packersdk.Ui)
	client := stateBag.Get("ksyun_client").(*ClientWrapper)
	if s.CommonConfig.Comm.SSHPrivateKeyFile != "" {
		if s.CommonConfig.Comm.SSHKeyPairName == "" {
			return Halt(stateBag,
				fmt.Errorf(fmt.Sprintf("ssh_keypair_name is empty")), "")
		}
		ui.Say("Using existing SSH private key")
		privateKeyBytes, err := s.CommonConfig.Comm.ReadSSHPrivateKeyFile()
		if err != nil {
			stateBag.Put("error", err)
			return multistep.ActionHalt
		}
		s.CommonConfig.Comm.SSHPrivateKey = privateKeyBytes
		return multistep.ActionContinue
	}

	if s.CommonConfig.Comm.SSHAgentAuth && s.CommonConfig.Comm.SSHKeyPairName == "" {
		ui.Say("Using SSH Agent with key pair in source image")
		return multistep.ActionContinue
	}

	if s.CommonConfig.Comm.SSHAgentAuth && s.CommonConfig.Comm.SSHKeyPairName != "" {
		ui.Say(fmt.Sprintf("Using SSH Agent for existing key pair %s", s.CommonConfig.Comm.SSHKeyPairName))
		return multistep.ActionContinue
	}

	if s.CommonConfig.Comm.SSHTemporaryKeyPairName == "" {
		ui.Say("Not using temporary keypair")
		s.CommonConfig.Comm.SSHKeyPairName = ""
		return multistep.ActionContinue
	}

	ui.Say(fmt.Sprintf("Using SSH Agent for create new key pair %s", s.CommonConfig.Comm.SSHTemporaryKeyPairName))
	//create ssh Key
	createSSHKey := make(map[string]interface{})
	createSSHKey["KeyName"] = s.CommonConfig.Comm.SSHTemporaryKeyPairName
	resp, err := client.SksClient.CreateKey(&createSSHKey)
	if err != nil {
		return Halt(stateBag, err, "Error creating new keypair")
	}
	if resp != nil {
		s.CommonConfig.Comm.SSHKeyPairName = GetSdkValue(stateBag, "Key.KeyId", *resp).(string)
		privateKey := GetSdkValue(stateBag, "PrivateKey", *resp).(string)
		publicKey := GetSdkValue(stateBag, "Key.PublicKey", *resp).(string)
		s.CommonConfig.Comm.SSHPrivateKey = []byte(privateKey)
		*(s.SSHTemporaryPublicKey) = strings.Split(publicKey, " ")[1]
		s.keyId = s.CommonConfig.Comm.SSHKeyPairName
	}
	return multistep.ActionContinue
}

func (s *StepConfigKsyunKeyPair) Cleanup(stateBag multistep.StateBag) {
	// if key not create by packer plugin no need to clean
	if s.keyId != "" {
		client := stateBag.Get("ksyun_client").(*ClientWrapper)
		ui := stateBag.Get("ui").(packersdk.Ui)
		//delete ssh Key
		ui.Say(fmt.Sprintf("Deleting temporary keypair %s ", s.CommonConfig.Comm.SSHKeyPairName))
		deleteSSHKey := make(map[string]interface{})
		deleteSSHKey["KeyId"] = s.keyId
		_, err := client.SksClient.DeleteKey(&deleteSSHKey)
		if err != nil {
			ui.Error(fmt.Sprintf(
				"Error cleaning up keypair. Please delete the key manually: name = %s, id = %s",
				s.CommonConfig.Comm.SSHTemporaryKeyPairName, s.keyId))
		}
	}
}
