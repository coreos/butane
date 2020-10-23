---
layout: default
nav_order: 2
---

# Getting started

`fcct`, the Fedora CoreOS Config Transpiler, is a tool that consumes a Fedora CoreOS Config and produces an Ignition config, which is a JSON document that can be given to a Fedora CoreOS machine when it first boots. Using this config, a machine can be told to create users, create filesystems, set up the network, install systemd units, and more.

Fedora CoreOS Configs are YAML files conforming to `fcct`'s schema. For more information on the schema, take a look at the [configuration specifications][spec].

### Getting FCCT

`fcct` can be downloaded as a standalone binary or run as a container with docker or podman.

#### Standalone binary

Download the latest version of `fcct` and the detached signature from the [releases page](https://github.com/coreos/fcct/releases). Verify it with gpg:

```
gpg --verify <detached sig> <fcct binary>
```
You may need to download the [Fedora signing keys](https://getfedora.org/static/fedora.gpg) and import them with `gpg --import <key>` if you have not already done so.

New releases of `fcct` are backwards compatible with old releases unless otherwise noted.

#### Container

This example uses podman, but docker can also be used.

```bash
# Pull the latest release
podman pull quay.io/coreos/fcct:release

# Run fcct using standard in and standard out
podman run -i --rm quay.io/coreos/fcct:release --pretty --strict < your_config.fcc > transpiled_config.ign

# Run fcct using files.
podman run --rm -v /path/to/your_config.fcc:/config.fcc:z quay.io/coreos/fcct:release --pretty --strict /config.fcc > transpiled_config.ign
```

### Writing and using Fedora CoreOS Configs

As a simple example, let's use `fcct` to set the authorized ssh key for the `core` user on a Fedora CoreOS machine.

<!-- fedora-coreos-config -->
```yaml
variant: fcos
version: 1.1.0
passwd:
  users:
    - name: core
      ssh_authorized_keys:
        - ssh-rsa AAAAB3NzaC1yc...
```

In this above file, you'll want to set the `ssh-rsa AAAAB3NzaC1yc...` line to be your ssh public key (which is probably the contents of `~/.ssh/id_rsa.pub`, if you're on Linux).

If we take this file and give it to `fcct`:

```
$ ./bin/amd64/fcct example.yaml

{"ignition":{"config":{"replace":{"source":null,"verification":{}}},"security":{"tls":{}},"timeouts":{},"version":"3.0.0"},"passwd":{"users":[{"name":"core","sshAuthorizedKeys":["ssh-rsa ssh-rsa AAAAB3NzaC1yc..."]}]},"storage":{},"systemd":{}}
```

We can see that it produces a JSON file. This file isn't intended to be human-friendly, and will definitely be a pain to read/edit (especially if you have multi-line things like systemd units). Luckily, you shouldn't have to care about this file! Just provide it to a booting Fedora CoreOS machine and [Ignition][ignition], the utility inside of Fedora CoreOS that receives this file, will know what to do with it.

The method by which this file is provided to a Fedora CoreOS machine depends on the environment in which the machine is running. For instructions on a given provider, head over to the [list of supported platforms for Ignition][supported-platforms].

To see some examples for what else `fcct` can do, head over to the [examples][examples].

[spec]: specs.md
[ignition]: https://coreos.github.io/ignition/
[supported-platforms]: https://coreos.github.io/ignition/supported-platforms/
[examples]: examples.md
