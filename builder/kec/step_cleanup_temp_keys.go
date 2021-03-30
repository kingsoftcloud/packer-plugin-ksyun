package kec

import (
	"context"
	"fmt"
	"github.com/hashicorp/packer-plugin-sdk/communicator"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
	"log"
	"strings"
)

type StepCleanupKsyunTempKeys struct {
	Comm                  *communicator.Config
	SSHTemporaryPublicKey *string
}

func (s *StepCleanupKsyunTempKeys) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	if !s.Comm.SSHClearAuthorizedKeys {
		return multistep.ActionContinue
	}

	if s.Comm.Type != "ssh" {
		return multistep.ActionContinue
	}

	if s.Comm.SSHTemporaryKeyPairName == "" {
		return multistep.ActionContinue
	}

	comm := state.Get("communicator").(packersdk.Communicator)
	ui := state.Get("ui").(packersdk.Ui)

	cmd := new(packersdk.RemoteCmd)
	temporaryPublicKey := strings.ReplaceAll(*s.SSHTemporaryPublicKey, "/", "\\/")

	ui.Say(fmt.Sprintf("Trying to remove ephemeral keys %s from authorized_keys files", s.Comm.SSHTemporaryKeyPairName))

	cmd.Command = fmt.Sprintf("sed -i.bak '/%s/d' ~/.ssh/authorized_keys; rm ~/.ssh/authorized_keys.bak", temporaryPublicKey)
	if err := cmd.RunWithUi(ctx, comm, ui); err != nil {
		log.Printf("Error cleaning up ~/.ssh/authorized_keys; please clean up keys manually: %s", err)
	}
	cmd = new(packersdk.RemoteCmd)
	cmd.Command = fmt.Sprintf("sudo sed -i.bak '/%s/d' /root/.ssh/authorized_keys; sudo rm /root/.ssh/authorized_keys.bak", temporaryPublicKey)
	if err := cmd.RunWithUi(ctx, comm, ui); err != nil {
		log.Printf("Error cleaning up /root/.ssh/authorized_keys; please clean up keys manually: %s", err)
	}

	return multistep.ActionContinue
}

func (s *StepCleanupKsyunTempKeys) Cleanup(state multistep.StateBag) {
}
