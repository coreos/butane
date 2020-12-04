---
layout: default
has_children: true
nav_order: 4
has_toc: false
---

# Configuration specifications

CoreOS Configs must conform to a specific variant and version of the `fcct` schema, specified with the `variant` and `version` fields in the configuration.

See the [Upgrading Configs](migrating-configs.md) page for instructions to update a configuration to the latest specification.

## Stable specification versions

We recommend that you always use the latest **stable** specification for your operating system to benefit from new features and bug fixes. The following **stable** specification versions are currently supported in `fcct`:

- Fedora CoreOS (`fcos`)
  - [v1.3.0](config-fcos-v1_3.md)
  - [v1.2.0](config-fcos-v1_2.md)
  - [v1.1.0](config-fcos-v1_1.md)
  - [v1.0.0](config-fcos-v1_0.md)
- RHEL CoreOS (`rhcos`)
  - [v0.1.0](config-rhcos-v0_1.md)

## Experimental specification versions

Do not use **experimental** specifications for anything beyond **development and testing** as they are subject to change **without warning or announcement**. The following **experimental** specification versions are currently available in `fcct`:

- Fedora CoreOS (`fcos`)
  - [v1.4.0-experimental](config-fcos-v1_4-exp.md)
- RHEL CoreOS (`rhcos`)
  - [v0.2.0-experimental](config-rhcos-v0_2-exp.md)

## FCC specifications and Ignition specifications

Each version of the FCC specification corresponds to a version of the Ignition specification:

| FCC variant | FCC version        | Ignition spec      |
|-------------|--------------------|--------------------|
| `fcos`      | 1.0.0              | 3.0.0              |
| `fcos`      | 1.1.0              | 3.1.0              |
| `fcos`      | 1.2.0              | 3.2.0              |
| `fcos`      | 1.3.0              | 3.2.0              |
| `fcos`      | 1.4.0-experimental | 3.3.0-experimental |
| `rhcos`     | 0.1.0              | 3.2.0              |
| `rhcos`     | 0.2.0-experimental | 3.3.0-experimental |
