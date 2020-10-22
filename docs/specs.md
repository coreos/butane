---
layout: default
has_children: true
nav_order: 4
has_toc: false
---

# Configuration specifications

Fedora CoreOS Configs must conform to a specific version of the `fcct` schema,
specified with the `version: X.Y.Z` field in the configuration.

See the [Upgrading Configs](migrating-configs.md) page for instructions to
update a configuration to the latest specification.

## Stable specification versions

We recommend that you always use the latest **stable** specification to benefit
from new features and bug fixes. The following **stable** specification
versions are currently supported in `fcct`:

- [v1.2.0](configuration-v1_2.md)
- [v1.1.0](configuration-v1_1.md)
- [v1.0.0](configuration-v1_0.md)

## Experimental specification versions

Do not use the **experimental** specification for anything beyond **development
and testing** as it is subject to change **without warning or announcement**.
The following **experimental** specification version is currently available in
`fcct`:

- [v1.3.0-experimental](configuration-v1_3-exp.md)

## FCCT specifications and Ignition specifications

Each version of the FCCT specification corresponds to a version of the Ignition
specification:

| FCCT spec          | Igntion spec       |
|--------------------|--------------------|
| 1.0.0              | 3.0.0              |
| 1.1.0              | 3.1.0              |
| 1.2.0              | 3.2.0              |
| 1.3.0-experimental | 3.3.0-experimental |
