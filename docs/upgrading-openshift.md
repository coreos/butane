---
title: OpenShift
parent: Upgrading configs
nav_order: 2
---

# Upgrading OpenShift configs

Occasionally, changes are made to OpenShift Butane configs (those that specify `variant: openshift` or, historically, `variant: rhcos`) that break backward compatibility. While this is not a concern for running machines, since Ignition only runs one time during first boot, it is a concern for those who maintain configuration files. This document serves to detail each of the breaking changes and tries to provide some reasoning for the change. This does not cover all of the changes to the spec - just those that need to be considered when migrating from one version to the next.

{: .no_toc }

1. TOC
{:toc}

## From Version 4.11.0 to 4.12.0

There are no breaking changes between versions 4.11.0 and 4.12.0 of the `openshift` configuration specification. Any valid 4.11.0 configuration can be updated to a 4.12.0 configuration by changing the version string in the config.

## From Version 4.10.0 to 4.11.0

There are no breaking changes between versions 4.10.0 and 4.11.0 of the `openshift` configuration specification. Any valid 4.10.0 configuration can be updated to a 4.11.0 configuration by changing the version string in the config.

## From Version 4.9.0 to 4.10.0

There are no breaking changes between versions 4.9.0 and 4.10.0 of the `openshift` configuration specification. Any valid 4.9.0 configuration can be updated to a 4.10.0 configuration by changing the version string in the config. 

### Resource compression

Resource compression, which was disabled in all `openshift` specs in Butane 0.12.1, is re-introduced in this spec version. The `compression` field can be set to `gzip` to decompress gzip-compressed resources. In addition, Butane may automatically compress resources specified with `inline` or `local`.

<!-- butane-config -->
```yaml
variant: openshift
version: 4.10.0
metadata:
  labels:
    machineconfiguration.openshift.io/role: worker
  name: config-openshift
storage:
  files:
    - path: /opt/file2
      contents:
        source: data:;base64,H4sIAAAAAAAC/zSQQY4bMQwE7/OKfsBgXpHccs0DGKntEJBIWSINP3+htfcmQECxq/74ZIeOlR3Vm08sDUhnnChuiyUYOSFVh66idgebxonFiuoHNVf3imAfPqFWtGpNC2SgyT+fBOONJrrcTSBNHykX8DcOmnZIRdf9eNJU+olH6oL5ipkVfHEWDQl1Q7YmvfgbreswXbpPfTN1gC9QULx3r/42eKTEBfzaTMkgdObkx1btmByT/2mVUwNqeHrLERLEc7uCaxFFW/tpRDBxy7tKHLYXYchUiZwX8PtVOIK5S1rASxEWCZQcWiUkYG4Y07XS4jzWjqWGkm3INoffblpUULk492/3tnfITqQVXJ+02a/jKwAA//+jjAk6wQEAAA==
        compression: gzip
      mode: 0644
```

## From Version 4.8.0 to 4.9.0

There are no functionality changes between versions 4.8.0 and 4.9.0 of the `openshift` configuration specification. Any valid 4.8.0 configuration can be updated to a 4.9.0 configuration by changing the version string in the config.

## From `rhcos` Version 0.1.0 to `openshift` Version 4.8.0

The new `openshift` config variant is intended to work both on the OpenShift Container Platform with RHEL CoreOS, and on OKD with Fedora CoreOS. The `rhcos` variant is still accepted by Butane but will not receive new features.

The `openshift` 4.8.0 specification is not backward-compatible with the `rhcos` 0.1.0 specification. It adds new mandatory metadata fields and removes certain Ignition config fields. In addition, `openshift` configs are transpiled to an OpenShift [MachineConfig] rather than an Ignition config by default. A valid `rhcos` 0.1.0 configuration can be updated to an `openshift` 4.8.0 configuration by changing the variant and version strings and then correcting any errors reported during transpilation.

The following is a list of breaking changes and notable new features.

### MachineConfig generation

By default, Butane transpiles an `openshift` Butane config into an OpenShift [MachineConfig]. Butane produces an Ignition config if the `-r` or `--raw` option is specified on the Butane command line.

### Removed config fields

The config no longer allows certain Ignition config fields that are rejected or discouraged by the OpenShift [Machine Config Operator].

In the `storage` section, `directories` and `links` are removed, along with `append` in `files`. Local file trees referenced in `trees` must not contain symlinks.

In the `passwd` section, `groups` is removed. All fields in `users` are removed except for `name` (which must be set to `core`) and `ssh_authorized_keys`.

### MachineConfig metadata fields

The config gained a new top-level `metadata` section containing metadata for the generated [MachineConfig]. The mandatory `name` field specifies a [name for the Kubernetes MachineConfig resource][k8s-names]. The `labels` field specifies a map of key-value pairs to be applied to the MachineConfig resource as [Kubernetes labels][k8s-labels]. The `machineconfiguration.openshift.io/role` label is required.

The `metadata` section is ignored when generating a raw Ignition config using the `-r` or `--raw` option.

<!-- butane-config -->
```yaml
variant: openshift
version: 4.8.0
metadata:
  name: minimal-config
  labels:
    machineconfiguration.openshift.io/role: worker
```

### MCO settings

The config gained a new top-level `openshift` section specifying [configuration][MCO settings] for the [Machine Config Operator]. The `extensions` field lists [RHCOS extension modules] to be installed on the node. The `fips` field enables [FIPS mode] when set to `true`. The `kernel_arguments` field specifies a list of [arguments][kernel arguments] to be added to the kernel command line. The `kernel_type` field can be set to `realtime` to use the [real-time kernel] on the node.

Fields in the `openshift` section are not included in a raw Ignition config generated using the `-r` or `--raw` option.

<!-- butane-config -->
```yaml
variant: openshift
version: 4.8.0
metadata:
  name: config-openshift
  labels:
    machineconfiguration.openshift.io/role: worker
openshift:
  extensions:
    - usbguard
  fips: true
  kernel_arguments:
    - console=ttyS1,115200
  kernel_type: realtime
```

### FIPS configuration for LUKS

When the `fips` field in the `openshift` section is set to `true`, LUKS volumes specified in the config (but not in any referenced configs) are configured to use a cipher compatible with [FIPS 140-2]. This cipher is applied to LUKS volumes specified in the `luks` subsections of the `storage` and `boot_device` sections.

<!-- butane-config -->
```yaml
variant: openshift
version: 4.8.0
metadata:
  name: fips-luks
  labels:
    machineconfiguration.openshift.io/role: worker
openshift:
  fips: true
boot_device:
  luks:
    tpm2: true
```

[FIPS 140-2]: https://csrc.nist.gov/publications/detail/fips/140/2/final
[FIPS mode]: https://docs.openshift.com/container-platform/4.7/installing/installing-fips.html
[k8s-names]: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
[k8s-labels]: https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/
[kernel arguments]: https://docs.openshift.com/container-platform/4.7/post_installation_configuration/machine-configuration-tasks.html#nodes-nodes-kernel-arguments_post-install-machine-configuration-tasks
[Machine Config Operator]: https://docs.openshift.com/container-platform/4.7/post_installation_configuration/machine-configuration-tasks.html#understanding-the-machine-config-operator
[MachineConfig]: https://docs.openshift.com/container-platform/4.7/post_installation_configuration/machine-configuration-tasks.html#machine-config-overviewpost-install-machine-configuration-tasks
[MCO settings]: https://docs.openshift.com/container-platform/4.7/post_installation_configuration/machine-configuration-tasks.html#what-can-you-change-with-machine-configs
[real-time kernel]: https://docs.openshift.com/container-platform/4.7/post_installation_configuration/machine-configuration-tasks.html#nodes-nodes-rtkernel-arguments_post-install-machine-configuration-tasks
[RHCOS extension modules]: https://docs.openshift.com/container-platform/4.7/post_installation_configuration/machine-configuration-tasks.html#rhcos-add-extensions_post-install-machine-configuration-tasks
