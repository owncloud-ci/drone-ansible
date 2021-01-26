def main(ctx):
    before = testing(ctx)

    stages = [
        linux(ctx, "amd64"),
        linux(ctx, "arm64"),
        linux(ctx, "arm"),
    ]

    after = manifest(ctx) + notification(ctx)

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
            },
            {
                "name": "lint",
                "image": "golang:1.15",
                "commands": [
                    "go run golang.org/x/lint/golint -set_exit_status ./...",
                ],

            },
            {
                "name": "vet",
                "image": "golang:1.15",
                "commands": [
                    "go vet ./...",
                ],
            },
            {
                "name": "test",
                "image": "golang:1.15",
                "commands": [
                    "go test -cover ./...",
                ],
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
                "dockerfile": "docker/Dockerfile.linux.%s" % (arch),
                "repo": "owncloudci/%s" % (ctx.repo.name),
                "username": {
                    "from_secret": "docker_username",
                },
                "password": {
                    "from_secret": "docker_password",
                },
                "auto_tag": True,
                "auto_tag_suffix": "linux-%s" % (arch),
            },
        })

    return {
        "kind": "pipeline",
        "type": "docker",
        "name": "linux-%s" % (arch),
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
        'name': 'notify',
        'image': 'plugins/slack',
        'settings': {
          'webhook': {
            'from_secret': 'private_rocketchat',
          },
          'channel': 'builds',
        },
      }
    ],
    "depends_on": [
        "manifest",
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
