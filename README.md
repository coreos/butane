# Fedora CoreOS Config Transpiler

The Fedora CoreOS Config Transpiler (FCCT) translates human readable Fedora CoreOS Configs (FCCs)
into machine readable [Ignition](https://github.com/coreos/ignition) Configs. See the [getting
started](docs/getting-started.md) guide for how to use FCCT and the [configuration
specifications](docs/specs.md) for everything FCCs support.

### Project Layout

Each config spec is composed of a base and distro component. The base components
roughly mirror the Ignition spec and are distro agnostic. The distro components
contain sugar for common configuration of the host (e.g. etcd) and are not
distro-independent.

Each distro and base component are versioned independently with each new
version getting it's own package. These versions are not exposed to the user.

Each fcos config version has it's own version which is independent of the
versions of the base and distro components that compose it. However a major
or minor bump of either component necessitates a corresponding bump in the fcos
config version.

`internal/`
  main, non-exported code

`base/`
  Contains distro-agnostic code. Each package here targets only one Ignition
  spec versions.

`distro/`
  Contains distro-specific code. Each package here can target multiple Ignition
  versions if it makes sense.

`config/`
  Contains the top level Translate() function that determines which version to
  parse and emit.

`config/common/`
  Contains the common bits and functions for all spec versions. This means the
  (un)marshaling helpers and the version+variant struct to be included in every
  user facing spec

`config/vX_Y/`
  Contains user facing definitions of the spec. Each is composed by combining a
  base and distro package with the common version+variant. Each of the defines
  their own translate function to be registered in the config/ package.
