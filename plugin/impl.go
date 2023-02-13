package plugin

import (
	"os"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	"golang.org/x/sys/execabs"
)

// Settings for the Plugin.
type Settings struct {
	Requirements      string
	Galaxy            string
	Inventories       cli.StringSlice
	Playbooks         cli.StringSlice
	Limit             string
	SkipTags          string
	StartAtTask       string
	Tags              string
	ExtraVars         cli.StringSlice
	ModulePath        cli.StringSlice
	Check             bool
	Diff              bool
	FlushCache        bool
	ForceHandlers     bool
	ListHosts         bool
	ListTags          bool
	ListTasks         bool
	SyntaxCheck       bool
	Forks             int
	VaultID           string
	VaultPassword     string
	VaultPasswordFile string
	Verbose           int
	PrivateKey        string
	PrivateKeyFile    string
	User              string
	Connection        string
	Timeout           int
	SSHCommonArgs     string
	SFTPExtraArgs     string
	SCPExtraArgs      string
	SSHExtraArgs      string
	Become            bool
	BecomeMethod      string
	BecomeUser        string
}

// Validate handles the settings validation of the plugin.
func (p *Plugin) Validate() error {
	if len(p.settings.Playbooks.Value()) == 0 {
		return errors.New("you must provide a playbook")
	}

	if len(p.settings.Inventories.Value()) == 0 {
		return errors.New("you must provide an inventory")
	}

	return nil
}

// Execute provides the implementation of the plugin.
func (p *Plugin) Execute() error {
	if err := p.playbooks(); err != nil {
		return err
	}

	if err := p.ansibleConfig(); err != nil {
		return err
	}

	if p.settings.PrivateKey != "" {
		if err := p.privateKey(); err != nil {
			return err
		}

		defer os.Remove(p.settings.PrivateKeyFile)
	}

	if p.settings.VaultPassword != "" {
		if err := p.vaultPass(); err != nil {
			return err
		}

		defer os.Remove(p.settings.VaultPasswordFile)
	}

	commands := []*execabs.Cmd{
		p.versionCommand(),
	}

	if p.settings.Requirements != "" {
		commands = append(commands, p.requirementsCommand())
	}

	if p.settings.Galaxy != "" {
		commands = append(commands, p.galaxyCommand())
	}

	for _, inventory := range p.settings.Inventories.Value() {
		commands = append(commands, p.ansibleCommand(inventory))
	}

	for _, cmd := range commands {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		cmd.Env = os.Environ()
		cmd.Env = append(cmd.Env, "ANSIBLE_FORCE_COLOR=1")

		trace(cmd)

		if err := cmd.Run(); err != nil {
			return err
		}
	}

	return nil
}
