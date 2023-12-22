package plugin

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/urfave/cli/v2"
	"golang.org/x/sys/execabs"
)

const (
	AnsibleForksDefault = 5

	ansibleFolder = "/etc/ansible"
	ansibleConfig = "/etc/ansible/ansible.cfg"

	pipBin             = "/usr/local/bin/pip"
	ansibleBin         = "/usr/local/bin/ansible"
	ansibleGalaxyBin   = "/usr/local/bin/ansible-galaxy"
	ansiblePlaybookBin = "/usr/local/bin/ansible-playbook"

	strictFilePerm = 0o600
)

const ansibleContent = `
[defaults]
host_key_checking = False
`

var ErrAnsiblePlaybookNotFound = errors.New("playbook not found")

func (p *Plugin) ansibleConfig() error {
	if err := os.MkdirAll(ansibleFolder, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create ansible directory: %w", err)
	}

	if err := os.WriteFile(ansibleConfig, []byte(ansibleContent), strictFilePerm); err != nil {
		return fmt.Errorf("failed to create ansible config: %w", err)
	}

	return nil
}

func (p *Plugin) privateKey() error {
	tmpfile, err := os.CreateTemp("", "privateKey")
	if err != nil {
		return fmt.Errorf("failed to create private key file: %w", err)
	}

	if _, err := tmpfile.Write([]byte(p.settings.PrivateKey)); err != nil {
		return fmt.Errorf("failed to write private key file: %w", err)
	}

	if err := tmpfile.Close(); err != nil {
		return fmt.Errorf("failed to close private key file: %w", err)
	}

	p.settings.PrivateKeyFile = tmpfile.Name()

	return nil
}

func (p *Plugin) vaultPass() error {
	tmpfile, err := os.CreateTemp("", "vaultPass")
	if err != nil {
		return fmt.Errorf("failed to create vault password file: %w", err)
	}

	if _, err := tmpfile.Write([]byte(p.settings.VaultPassword)); err != nil {
		return fmt.Errorf("failed to write vault password file: %w", err)
	}

	if err := tmpfile.Close(); err != nil {
		return fmt.Errorf("failed to close vault password file: %w", err)
	}

	p.settings.VaultPasswordFile = tmpfile.Name()

	return nil
}

func (p *Plugin) playbooks() error {
	var playbooks []string

	for _, p := range p.settings.Playbooks.Value() {
		files, err := filepath.Glob(p)
		if err != nil {
			playbooks = append(playbooks, p)

			continue
		}

		playbooks = append(playbooks, files...)
	}

	if len(playbooks) == 0 {
		return ErrAnsiblePlaybookNotFound
	}

	p.settings.Playbooks = *cli.NewStringSlice(playbooks...)

	return nil
}

func (p *Plugin) versionCommand() *execabs.Cmd {
	args := []string{
		"--version",
	}

	return execabs.Command(
		ansibleBin,
		args...,
	)
}

func (p *Plugin) requirementsCommand() *execabs.Cmd {
	args := []string{
		"install",
		"--upgrade",
		"--requirement",
		p.settings.Requirements,
	}

	return execabs.Command(
		pipBin,
		args...,
	)
}

func (p *Plugin) galaxyCommand() *execabs.Cmd {
	args := []string{
		"install",
		"--force",
		"--role-file",
		p.settings.Galaxy,
	}

	if p.settings.Verbose > 0 {
		args = append(args, fmt.Sprintf("-%s", strings.Repeat("v", p.settings.Verbose)))
	}

	return execabs.Command(
		ansibleGalaxyBin,
		args...,
	)
}

func (p *Plugin) ansibleCommand(inventory string) *execabs.Cmd {
	args := []string{
		"--inventory",
		inventory,
	}

	if len(p.settings.ModulePath.Value()) > 0 {
		args = append(args, "--module-path", strings.Join(p.settings.ModulePath.Value(), ":"))
	}

	if p.settings.VaultID != "" {
		args = append(args, "--vault-id", p.settings.VaultID)
	}

	if p.settings.VaultPasswordFile != "" {
		args = append(args, "--vault-password-file", p.settings.VaultPasswordFile)
	}

	for _, v := range p.settings.ExtraVars.Value() {
		args = append(args, "--extra-vars", v)
	}

	if p.settings.ListHosts {
		args = append(args, "--list-hosts")
		args = append(args, p.settings.Playbooks.Value()...)

		return execabs.Command(
			ansiblePlaybookBin,
			args...,
		)
	}

	if p.settings.SyntaxCheck {
		args = append(args, "--syntax-check")
		args = append(args, p.settings.Playbooks.Value()...)

		return execabs.Command(
			ansiblePlaybookBin,
			args...,
		)
	}

	if p.settings.Check {
		args = append(args, "--check")
	}

	if p.settings.Diff {
		args = append(args, "--diff")
	}

	if p.settings.FlushCache {
		args = append(args, "--flush-cache")
	}

	if p.settings.ForceHandlers {
		args = append(args, "--force-handlers")
	}

	if p.settings.Forks != AnsibleForksDefault {
		args = append(args, "--forks", strconv.Itoa(p.settings.Forks))
	}

	if p.settings.Limit != "" {
		args = append(args, "--limit", p.settings.Limit)
	}

	if p.settings.ListTags {
		args = append(args, "--list-tags")
	}

	if p.settings.ListTasks {
		args = append(args, "--list-tasks")
	}

	if p.settings.SkipTags != "" {
		args = append(args, "--skip-tags", p.settings.SkipTags)
	}

	if p.settings.StartAtTask != "" {
		args = append(args, "--start-at-task", p.settings.StartAtTask)
	}

	if p.settings.Tags != "" {
		args = append(args, "--tags", p.settings.Tags)
	}

	if p.settings.PrivateKeyFile != "" {
		args = append(args, "--private-key", p.settings.PrivateKeyFile)
	}

	if p.settings.User != "" {
		args = append(args, "--user", p.settings.User)
	}

	if p.settings.Connection != "" {
		args = append(args, "--connection", p.settings.Connection)
	}

	if p.settings.Timeout != 0 {
		args = append(args, "--timeout", strconv.Itoa(p.settings.Timeout))
	}

	if p.settings.SSHCommonArgs != "" {
		args = append(args, "--ssh-common-args", p.settings.SSHCommonArgs)
	}

	if p.settings.SFTPExtraArgs != "" {
		args = append(args, "--sftp-extra-args", p.settings.SFTPExtraArgs)
	}

	if p.settings.SCPExtraArgs != "" {
		args = append(args, "--scp-extra-args", p.settings.SCPExtraArgs)
	}

	if p.settings.SSHExtraArgs != "" {
		args = append(args, "--ssh-extra-args", p.settings.SSHExtraArgs)
	}

	if p.settings.Become {
		args = append(args, "--become")
	}

	if p.settings.BecomeMethod != "" {
		args = append(args, "--become-method", p.settings.BecomeMethod)
	}

	if p.settings.BecomeUser != "" {
		args = append(args, "--become-user", p.settings.BecomeUser)
	}

	if p.settings.Verbose > 0 {
		args = append(args, fmt.Sprintf("-%s", strings.Repeat("v", p.settings.Verbose)))
	}

	args = append(args, p.settings.Playbooks.Value()...)

	return execabs.Command(
		ansiblePlaybookBin,
		args...,
	)
}

func trace(cmd *execabs.Cmd) {
	fmt.Println("$", strings.Join(cmd.Args, " "))
}
