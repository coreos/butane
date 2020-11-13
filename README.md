# Fedora CoreOS Config Transpiler

The Fedora CoreOS Config Transpiler (FCCT) translates human readable Fedora CoreOS Configs (FCCs)
into machine readable [Ignition](https://github.com/coreos/ignition) Configs. See the [getting
started](docs/getting-started.md) guide for how to use FCCT and the [configuration
specifications](docs/specs.md) for everything FCCs support.

### Project Layout

Internally, FCCT has a versioned `base` component which contains support for
a particular Ignition spec version, plus distro-independent sugar. New base
functionality is added only to the experimental base package. Eventually the
experimental base package is stabilized and a new experimental package
created. The base component is versioned independently of any particular
distro, and its versions are not exposed to the user. Client code should
not need to import anything from `base`.

Each FCC variant/version pair corresponds to a `config` package, which
derives either from a `base` package or from another `config` package. New
functionality is similarly added only to an experimental config version,
which is eventually stabilized and a new experimental version created.
(This will often happen when the underlying package is stabilized.) A
`config` package can contain sugar or validation logic specific to a distro
(for example, additional properties for configuring etcd).

Packages outside the FCCT repository can implement additional FCC versions
by deriving from a `base` or `config` package and registering their
variant/version pair with `config`.

`config/`
  Top-level `TranslateBytes()` function that determines which config version
  to parse and emit. Clients should typically use this to translate FCCs.

`config/common/`
  Common definitions for all spec versions, including translate options
  structs and error definitions.

`config/*/vX_Y/`
  User facing definitions of the spec. Each is derived from another config
  package or from a base package. Each one defines its own translate
  functions to be registered in the `config` package. Clients can use
  these directly if they want to translate a specific spec version.

`config/util/`
  Utility code for implementing config packages, including the
  (un)marshaling helpers. Clients don't need to import this unless they're
  implementing an out-of-tree config version.

`base/`
  Distro-agnostic code targeting individual Ignition spec versions. Clients
  don't need to import this unless they're implementing an out-of-tree
  config version.

`internal/`
  `main`, non-exported code.
