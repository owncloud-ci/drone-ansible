package main

import (
	"github.com/owncloud-ci/drone-ansible/plugin"
	"github.com/urfave/cli/v2"
)

// settingsFlags has the cli.Flags for the plugin.Settings.
func settingsFlags(settings *plugin.Settings) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "requirements",
			Usage:       "path to python requirements",
			EnvVars:     []string{"PLUGIN_REQUIREMENTS"},
			Destination: &settings.Requirements,
		},
		&cli.StringFlag{
			Name:        "galaxy",
			Usage:       "path to galaxy requirements",
			EnvVars:     []string{"PLUGIN_GALAXY"},
			Destination: &settings.Galaxy,
		},
		&cli.StringSliceFlag{
			Name:        "inventory",
			Usage:       "specify inventory host path",
			EnvVars:     []string{"PLUGIN_INVENTORY", "PLUGIN_INVENTORIES"},
			Destination: &settings.Inventories,
		},
		&cli.StringSliceFlag{
			Name:        "playbook",
			Usage:       "list of playbooks to apply",
			EnvVars:     []string{"PLUGIN_PLAYBOOK", "PLUGIN_PLAYBOOKS"},
			Destination: &settings.Playbooks,
		},
		&cli.StringFlag{
			Name:        "limit",
			Usage:       "further limit selected hosts to an additional pattern",
			EnvVars:     []string{"PLUGIN_LIMIT"},
			Destination: &settings.Limit,
		},
		&cli.StringFlag{
			Name:        "skip-tags",
			Usage:       "only run plays and tasks whose tags do not match",
			EnvVars:     []string{"PLUGIN_SKIP_TAGS"},
			Destination: &settings.SkipTags,
		},
		&cli.StringFlag{
			Name:        "start-at-task",
			Usage:       "start the playbook at the task matching this name",
			EnvVars:     []string{"PLUGIN_START_AT_TASK"},
			Destination: &settings.StartAtTask,
		},
		&cli.StringFlag{
			Name:        "tags",
			Usage:       "only run plays and tasks tagged with these values",
			EnvVars:     []string{"PLUGIN_TAGS"},
			Destination: &settings.Tags,
		},
		&cli.StringSliceFlag{
			Name:        "extra-vars",
			Usage:       "set additional variables as key=value",
			EnvVars:     []string{"PLUGIN_EXTRA_VARS", "ANSIBLE_EXTRA_VARS"},
			Destination: &settings.ExtraVars,
		},
		&cli.StringSliceFlag{
			Name:        "module-path",
			Usage:       "prepend paths to module library",
			EnvVars:     []string{"PLUGIN_MODULE_PATH"},
			Destination: &settings.ModulePath,
		},
		&cli.BoolFlag{
			Name:        "check",
			Usage:       "run a check, do not apply any changes",
			EnvVars:     []string{"PLUGIN_CHECK"},
			Destination: &settings.Check,
		},
		&cli.BoolFlag{
			Name:        "diff",
			Usage:       "show the differences, may print secrets",
			EnvVars:     []string{"PLUGIN_DIFF"},
			Destination: &settings.Diff,
		},
		&cli.BoolFlag{
			Name:        "flush-cache",
			Usage:       "clear the fact cache for every host in inventory",
			EnvVars:     []string{"PLUGIN_FLUSH_CACHE"},
			Destination: &settings.FlushCache,
		},
		&cli.BoolFlag{
			Name:        "force-handlers",
			Usage:       "run handlers even if a task fails",
			EnvVars:     []string{"PLUGIN_FORCE_HANDLERS"},
			Destination: &settings.ForceHandlers,
		},
		&cli.BoolFlag{
			Name:        "list-hosts",
			Usage:       "outputs a list of matching hosts",
			EnvVars:     []string{"PLUGIN_LIST_HOSTS"},
			Destination: &settings.ListHosts,
		},
		&cli.BoolFlag{
			Name:        "list-tags",
			Usage:       "list all available tags",
			EnvVars:     []string{"PLUGIN_LIST_TAGS"},
			Destination: &settings.ListTags,
		},
		&cli.BoolFlag{
			Name:        "list-tasks",
			Usage:       "list all tasks that would be executed",
			EnvVars:     []string{"PLUGIN_LIST_TASKS"},
			Destination: &settings.ListTasks,
		},
		&cli.BoolFlag{
			Name:        "syntax-check",
			Usage:       "perform a syntax check on the playbook",
			EnvVars:     []string{"PLUGIN_SYNTAX_CHECK"},
			Destination: &settings.SyntaxCheck,
		},
		&cli.IntFlag{
			Name:        "forks",
			Usage:       "specify number of parallel processes to use",
			EnvVars:     []string{"PLUGIN_FORKS"},
			Value:       5,
			Destination: &settings.Forks,
		},
		&cli.StringFlag{
			Name:        "vault-id",
			Usage:       "the vault identity to use",
			EnvVars:     []string{"PLUGIN_VAULT_ID", "ANSIBLE_VAULT_ID"},
			Destination: &settings.VaultID,
		},
		&cli.StringFlag{
			Name:        "vault-password",
			Usage:       "the vault password to use",
			EnvVars:     []string{"PLUGIN_VAULT_PASSWORD", "ANSIBLE_VAULT_PASSWORD"},
			Destination: &settings.VaultPassword,
		},
		&cli.IntFlag{
			Name:        "verbose",
			Usage:       "level of verbosity, 0 up to 4",
			EnvVars:     []string{"PLUGIN_VERBOSE"},
			Destination: &settings.Verbose,
		},
		&cli.StringFlag{
			Name:        "private-key",
			Usage:       "use this key to authenticate the connection",
			EnvVars:     []string{"PLUGIN_PRIVATE_KEY", "ANSIBLE_PRIVATE_KEY"},
			Destination: &settings.PrivateKey,
		},
		&cli.StringFlag{
			Name:        "user",
			Usage:       "connect as this user",
			EnvVars:     []string{"PLUGIN_USER", "ANSIBLE_USER"},
			Destination: &settings.User,
		},
		&cli.StringFlag{
			Name:        "connection",
			Usage:       "connection type to use",
			EnvVars:     []string{"PLUGIN_CONNECTION"},
			Destination: &settings.Connection,
		},
		&cli.IntFlag{
			Name:        "timeout",
			Usage:       "override the connection timeout in seconds",
			EnvVars:     []string{"PLUGIN_TIMEOUT"},
			Destination: &settings.Timeout,
		},
		&cli.StringFlag{
			Name:        "ssh-common-args",
			Usage:       "specify common arguments to pass to sftp/scp/ssh",
			EnvVars:     []string{"PLUGIN_SSH_COMMON_ARGS"},
			Destination: &settings.SSHCommonArgs,
		},
		&cli.StringFlag{
			Name:        "sftp-extra-args",
			Usage:       "specify extra arguments to pass to sftp only",
			EnvVars:     []string{"PLUGIN_SFTP_EXTRA_ARGS"},
			Destination: &settings.SFTPExtraArgs,
		},
		&cli.StringFlag{
			Name:        "scp-extra-args",
			Usage:       "specify extra arguments to pass to scp only",
			EnvVars:     []string{"PLUGIN_SCP_EXTRA_ARGS"},
			Destination: &settings.SCPExtraArgs,
		},
		&cli.StringFlag{
			Name:        "ssh-extra-args",
			Usage:       "specify extra arguments to pass to ssh only",
			EnvVars:     []string{"PLUGIN_SSH_EXTRA_ARGS"},
			Destination: &settings.SSHExtraArgs,
		},
		&cli.BoolFlag{
			Name:        "become",
			Usage:       "run operations with become",
			EnvVars:     []string{"PLUGIN_BECOME"},
			Destination: &settings.Become,
		},
		&cli.StringFlag{
			Name:        "become-method",
			Usage:       "privilege escalation method to use",
			EnvVars:     []string{"PLUGIN_BECOME_METHOD", "ANSIBLE_BECOME_METHOD"},
			Destination: &settings.BecomeMethod,
		},
		&cli.StringFlag{
			Name:        "become-user",
			Usage:       "run operations as this user",
			EnvVars:     []string{"PLUGIN_BECOME_USER", "ANSIBLE_BECOME_USER"},
			Destination: &settings.BecomeUser,
		},
	}
}
