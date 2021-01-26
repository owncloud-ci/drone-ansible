# drone-ansible

[![Build Status](https://drone.owncloud.com/api/badges/owncloud-ci/hugo/status.svg)](https://drone.owncloud.com/owncloud-ci/hugo/)
[![Docker Hub](https://img.shields.io/badge/docker-latest-blue.svg?logo=docker&logoColor=white)](https://hub.docker.com/r/owncloudci/hugo)

Drone plugin to provision infrastructure with [Ansible](https://www.ansible.com/).

## Build

Build the binary with the following command:

```console
export GOOS=linux
export GOARCH=amd64
export CGO_ENABLED=0
export GO111MODULE=on

go build -v -a -tags netgo -o release/linux/amd64/drone-ansible
```

## Docker

Build the Docker image with the following command:

```console
docker build \
  --label org.label-schema.build-date=$(date -u +"%Y-%m-%dT%H:%M:%SZ") \
  --label org.label-schema.vcs-ref=$(git rev-parse --short HEAD) \
  --file docker/Dockerfile.linux.amd64 --tag plugins/ansible .
```

## Usage

```console
docker run --rm \
  -e PLUGIN_PRIVATE_KEY="$(cat ~/.ssh/id_rsa)" \
  -e PLUGIN_PLAYBOOK="deployment/playbook.yml" \
  -e PLUGIN_INVENTORY="deployment/hosts.yml" \
  -v $(pwd):$(pwd) \
  -w $(pwd) \
  plugins/ansible
```

## License

This project is licensed under the Apache 2.0 License - see the [LICENSE](https://github.com/owncloud-ci/drone-ansible/blob/master/LICENSE) file for details.
