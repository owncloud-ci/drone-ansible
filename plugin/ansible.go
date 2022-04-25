package plugin

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

var (
	ansibleFolder = "/etc/ansible"
	ansibleConfig = "/etc/ansible/ansible.cfg"
)

var ansibleContent = `
[defaults]
host_key_checking = False
`

func (p *Plugin) ansibleConfig() error {
	if err := os.MkdirAll(ansibleFolder, os.ModePerm); err != nil {
		return errors.Wrap(err, "failed to create ansible directory")
	}

	if err := ioutil.WriteFile(ansibleConfig, []byte(ansibleContent), 0o600); err != nil {
		return errors.Wrap(err, "failed to create ansible config")
	}

	return nil
}

func (p *Plugin) privateKey() error {
	tmpfile, err := ioutil.TempFile("", "privateKey")
	if err != nil {
		return errors.Wrap(err, "failed to create private key file")
	}

	if _, err := tmpfile.Write([]byte(p.settings.PrivateKey)); err != nil {
		return errors.Wrap(err, "failed to write private key file")
	}

	if err := tmpfile.Close(); err != nil {
		return errors.Wrap(err, "failed to close private key file")
	}

	p.settings.PrivateKeyFile = tmpfile.Name()
	return nil
}

func (p *Plugin) vaultPass() error {
	tmpfile, err := ioutil.TempFile("", "vaultPass")
	if err != nil {
		return errors.Wrap(err, "failed to create vault password file")
	}

	if _, err := tmpfile.Write([]byte(p.settings.VaultPassword)); err != nil {
		return errors.Wrap(err, "failed to write vault password file")
	}

	if err := tmpfile.Close(); err != nil {
		return errors.Wrap(err, "failed to close vault password file")
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
		return errors.New("failed to find playbook files")
	}

	p.settings.Playbooks = *cli.NewStringSlice(playbooks...)
	return nil
}

func (p *Plugin) versionCommand() *exec.Cmd {
	args := []string{
		"--version",
	}

	return exec.Command(
		"ansible",
		args...,
	)
}

func (p *Plugin) requirementsCommand() *exec.Cmd {
	args := []string{
		"install",
		"--upgrade",
		"--requirement",
		p.settings.Requirements,
	}

	return exec.Command(
		"pip",
		args...,
	)
}

func (p *Plugin) galaxyCommand() *exec.Cmd {
	args := []string{
		"install",
		"--force",
		"--role-file",
		p.settings.Galaxy,
	}

	if p.settings.Verbose > 0 {
		args = append(args, fmt.Sprintf("-%s", strings.Repeat("v", p.settings.Verbose)))
	}

	return exec.Command(
		"ansible-galaxy",
		args...,
	)
}

func (p *Plugin) ansibleCommand(inventory string) *exec.Cmd {
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

		return exec.Command(
			"ansible-playbook",
			args...,
		)
	}

	if p.settings.SyntaxCheck {
		args = append(args, "--syntax-check")
		args = append(args, p.settings.Playbooks.Value()...)

		return exec.Command(
			"ansible-playbook",
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

	if p.settings.Forks != 5 {
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

	return exec.Command(
		"ansible-playbook",
		args...,
	)
}

func trace(cmd *exec.Cmd) {
	fmt.Println("$", strings.Join(cmd.Args, " "))
}
