def main(ctx):
    before = testing(ctx)

    stages = [
        linux(ctx, "amd64"),
        linux(ctx, "arm64"),
        linux(ctx, "arm"),
    ]

    after = manifest(ctx) + pushrm(ctx) + release(ctx) + notification(ctx)

    for b in before:
        for s in stages:
            s["depends_on"].append(b["name"])

    for s in stages:
        for a in after:
            a["depends_on"].append(s["name"])

    return before + stages + after

def testing(ctx):
    return [{
        "kind": "pipeline",
        "type": "docker",
        "name": "testing",
        "platform": {
            "os": "linux",
            "arch": "amd64",
        },
        "steps": [
            {
                "name": "staticcheck",
                "image": "golang:1.15",
                "pull": "always",
                "commands": [
                    "go run honnef.co/go/tools/cmd/staticcheck ./...",
                ],
                "volumes": [
                    {
                        "name": "gopath",
                        "path": "/go",
                    },
                ],
            },
            {
                "name": "lint",
                "image": "golang:1.15",
                "commands": [
                    "go run golang.org/x/lint/golint -set_exit_status ./...",
                ],
                "volumes": [
                    {
                        "name": "gopath",
                        "path": "/go",
                    },
                ],
            },
            {
                "name": "vet",
                "image": "golang:1.15",
                "commands": [
                    "go vet ./...",
                ],
                "volumes": [
                    {
                        "name": "gopath",
                        "path": "/go",
                    },
                ],
            },
            {
                "name": "test",
                "image": "golang:1.15",
                "commands": [
                    "go test -cover ./...",
                ],
                "volumes": [
                    {
                        "name": "gopath",
                        "path": "/go",
                    },
                ],
            },
        ],
        "volumes": [
            {
                "name": "gopath",
                "temp": {},
            },
        ],
        "trigger": {
            "ref": [
                "refs/heads/master",
                "refs/tags/**",
                "refs/pull/**",
            ],
        },
    }]

def linux(ctx, arch):
    if ctx.build.event == "tag":
        build = [
            'go build -v -ldflags "-X main.version=%s" -a -tags netgo -o release/linux/%s/drone-ansible ./cmd/drone-ansible' % (ctx.build.ref.replace("refs/tags/v", ""), arch),
        ]
    else:
        build = [
            'go build -v -ldflags "-X main.version=%s" -a -tags netgo -o release/linux/%s/drone-ansible ./cmd/drone-ansible' % (ctx.build.commit[0:8], arch),
        ]

    steps = [
        {
            "name": "environment",
            "image": "golang:1.15",
            "pull": "always",
            "environment": {
                "CGO_ENABLED": "0",
            },
            "commands": [
                "go version",
                "go env",
            ],
        },
        {
            "name": "build",
            "image": "golang:1.15",
            "environment": {
                "CGO_ENABLED": "0",
            },
            "commands": build,
        },
        {
            "name": "executable",
            "image": "golang:1.15",
            "commands": [
                "./release/linux/%s/drone-ansible --help" % (arch),
            ],
        },
    ]

    if ctx.build.event != "pull_request":
        steps.append({
            "name": "docker",
            "image": "plugins/docker",
            "settings": {
                "dockerfile": "docker/Dockerfile.%s" % (arch),
                "repo": "owncloudci/%s" % (ctx.repo.name),
                "username": {
                    "from_secret": "docker_username",
                },
                "password": {
                    "from_secret": "docker_password",
                },
                "auto_tag": True,
                "auto_tag_suffix": "%s" % (arch),
            },
        })
    else:
        steps.append({
            "name": "dryrun",
            "image": "plugins/docker",
            "settings": {
                'dry_run': True,
                "dockerfile": "docker/Dockerfile.%s" % (arch),
                "repo": "owncloudci/%s" % (ctx.repo.name),
                "auto_tag": True,
                "auto_tag_suffix": "%s" % (arch),
            },
        })

    return {
        "kind": "pipeline",
        "type": "docker",
        "name": "build-%s" % (arch),
        "platform": {
            "os": "linux",
            "arch": arch,
        },
        "steps": steps,
        "depends_on": [],
        "trigger": {
            "ref": [
                "refs/heads/master",
                "refs/tags/**",
                "refs/pull/**",
            ],
        },
    }

def manifest(ctx):
    return [{
        "kind": "pipeline",
        "type": "docker",
        "name": "manifest",
        "steps": [
            {
                "name": "manifest",
                "image": "plugins/manifest",
                "settings": {
                    "auto_tag": "true",
                    "username": {
                        "from_secret": "docker_username",
                    },
                    "password": {
                        "from_secret": "docker_password",
                    },
                    "spec": "docker/manifest.tmpl",
                    "ignore_missing": "true",
                },
            },
        ],
        "depends_on": [],
        "trigger": {
            "ref": [
                "refs/heads/master",
                "refs/tags/**",
            ],
        },
    }]

def pushrm(ctx):
  return [{
    "kind": "pipeline",
    "type": "docker",
    "name": "pushrm",
    "steps": [
        {
            "name": "pushrm",
            "image": "chko/docker-pushrm:1",
            "environment": {
                "DOCKER_PASS": {
                    "from_secret": "docker_password",
                },
                "DOCKER_USER": {
                    "from_secret": "docker_username",
                },
                "PUSHRM_FILE": "README.md",
                "PUSHRM_SHORT": "Drone plugin to provision infrastructure with Ansible",
                "PUSHRM_TARGET": "owncloudci/%s" % (ctx.repo.name),
            },
        },
    ],
    "depends_on": [
        "manifest",
    ],
    "trigger": {
        "ref": [
            "refs/heads/master",
            "refs/tags/**",
        ],
        "status": ["success"],
    },
  }]

def release(ctx):
  return [{
    "kind": "pipeline",
    "type": "docker",
    "name": "release",
    "steps": [
        {
            "name": "changelog",
            "image": "thegeeklab/git-chglog",
            "commands": [
                "git fetch -tq",
                "git-chglog --no-color --no-emoji %s" % (ctx.build.ref.replace("refs/tags/v", "") if ctx.build.event == "tag" else "--next-tag unreleased unreleased"),
                "git-chglog --no-color --no-emoji -o CHANGELOG.md %s" % (ctx.build.ref.replace("refs/tags/v", "") if ctx.build.event == "tag" else "--next-tag unreleased unreleased"),
            ]
        },
        {
           "name": "release",
           "image": "plugins/github-release",
           "settings": {
               "api_key": {
                   "from_secret": "github_token",
               },
               "note": "CHANGELOG.md",
               "overwrite": True,
               "title": ctx.build.ref.replace("refs/tags/", ""),
           },
           "when": {
             "ref": [
                "refs/tags/**",
             ],
          },
        }
    ],
    "depends_on": [
        "pushrm",
    ],
    "trigger": {
        "ref": [
            "refs/heads/master",
            "refs/tags/**",
            "refs/pull/**",
        ],
    },
  }]

def notification(ctx):
  return [{
    "kind": "pipeline",
    "type": "docker",
    "name": "notify",
    "clone": {
        "disable": True,
    },
    "steps": [
        {
            "name": "notify",
            "image": "plugins/slack",
            "settings": {
            "webhook": {
                "from_secret": "private_rocketchat",
            },
            "channel": "builds",
            },
        }
    ],
    "depends_on": [
        "release",
    ],
    "trigger": {
        "ref": [
            "refs/heads/master",
            "refs/tags/**",
        ],
        "status": [
            "success",
            "failure",
        ],
    },
  }]
