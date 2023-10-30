# drone-ansible

[![Build Status](https://drone.owncloud.com/api/badges/owncloud-ci/drone-ansible/status.svg)](https://drone.owncloud.com/owncloud-ci/drone-ansible)
[![Docker Hub](https://img.shields.io/docker/v/owncloudci/drone-ansible?logo=docker&label=dockerhub&sort=semver&logoColor=white)](https://hub.docker.com/r/owncloudci/drone-ansible)
[![GitHub contributors](https://img.shields.io/github/contributors/owncloud-ci/drone-ansible)](https://github.com/owncloud-ci/drone-ansible/graphs/contributors)
[![Source: GitHub](https://img.shields.io/badge/source-github-blue.svg?logo=github&logoColor=white)](https://github.com/owncloud-ci/drone-ansible)
[![License: Apache-2.0](https://img.shields.io/github/license/owncloud-ci/drone-ansible)](https://github.com/owncloud-ci/drone-ansible/blob/main/LICENSE)

Drone plugin to provision infrastructure with [Ansible](https://www.ansible.com/).

## Versioning

The tags follow the major version of Docker, e.g. `8`, and the minor and patch parts reflect the `version` of the plugin. A full example would be `8.5.2`. Minor versions can introduce breaking changes, while patch versions can be considered non-breaking.

## Usage

```yaml
kind: pipeline
type: docker
name: default

steps:
  - name: ansible
    image: owncloudci/drone-ansible
    settings:
      playbook: deployment/playbook.yml
      private_key:
        from_secret: ansible_private_key
      inventory: deployment/hosts.yml
```

## Build

Build the binary with the following command:

```console
make build
```

Build the Docker image with the following command:

```console
docker build --file Dockerfile.multiarch --tag owncloudci/drone-ansible .
```

## Test

```console
docker run --rm \
  -e PLUGIN_PRIVATE_KEY="$(cat ~/.ssh/id_rsa)" \
  -e PLUGIN_PLAYBOOK="deployment/playbook.yml" \
  -e PLUGIN_INVENTORY="deployment/hosts.yml" \
  -v $(pwd):$(pwd) \
  -w $(pwd) \
  owncloudci/drone-ansible --dry-run
```

## Releases

Create and push the new tag to trigger the CI release process:

```console
git tag v2.10.3
git push origin v2.10.3
```

## License

This project is licensed under the Apache 2.0 License - see the [LICENSE](https://github.com/owncloud-ci/drone-ansible/blob/main/LICENSE) file for details.

## Copyright

```text
Copyright (c) 2022 ownCloud GmbH
```
