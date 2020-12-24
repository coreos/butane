---
layout: default
nav_order: 9
---

# Developing FCCT
{: .no_toc }

1. TOC
{:toc}

## Project layout

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

- `config/` &mdash;
  Top-level `TranslateBytes()` function that determines which config version
  to parse and emit. Clients should typically use this to translate FCCs.

- `config/common/` &mdash;
  Common definitions for all spec versions, including translate options
  structs and error definitions.

- `config/*/vX_Y/` &mdash;
  User facing definitions of the spec. Each is derived from another config
  package or from a base package. Each one defines its own translate
  functions to be registered in the `config` package. Clients can use
  these directly if they want to translate a specific spec version.

- `config/util/` &mdash;
  Utility code for implementing config packages, including the
  (un)marshaling helpers. Clients don't need to import this unless they're
  implementing an out-of-tree config version.

- `base/` &mdash;
  Distro-agnostic code targeting individual Ignition spec versions. Clients
  don't need to import this unless they're implementing an out-of-tree
  config version.

- `internal/` &mdash;
  `main`, non-exported code.

## Creating a release

Create a [release checklist](https://github.com/coreos/fcct/issues/new?template=release-checklist.md) and follow those steps.

## Bumping spec versions

This checklist describes bumping the Ignition spec version, `base` version, and `config` version. If your scenario is different, modify to taste.

### Stabilize Ignition spec version

- Bump `go.mod` for new Ignition release and update vendor.
- Update imports. Drop `-experimental` from Ignition spec versions in `base/vB_exp/translate_test.go`.

### Bump base version

- Rename `base/vB_exp` to `base/vB` and update `package` statements. Update imports.
- Copy `base/vB` to `base/vB+1_exp`.
- Update `package` statements in `base/vB+1_exp`.

### Bump config version

- Rename `config/fcos/vC_exp` to `config/fcos/vC` and update `package` statements. Update imports.
- Drop `-experimental` from `init()` in `config/config.go`.
- Drop `-experimental` from examples in `docs/`.
- Copy `config/fcos/vC` to `config/fcos/vC+1_exp`.
- Update `package` statements in `config/fcos/vC+1_exp`. Bump its base dependency to `base/vB+1_exp`.
- Import `config/vC+1_exp` in `config/config.go` and add `fcos` `C+1-experimental` to `init()`.

### Bump Ignition spec version

- Bump Ignition types imports and rename `ToIgnI` and `TestToIgnI` functions in `base/vB+1_exp`. Bump Ignition spec versions in `base/vB+1_exp/translate_test.go`.
- Bump Ignition types imports in `config/fcos/vC+1_exp`. Update `ToIgnI` function names, `util` calls, and header comments to `ToIgnI+1`.

### Update docs

- Copy the `C-exp` spec doc to `C+1-exp`. Update the header and the version numbers in the description of the `version` field.
- Rename the `C-exp` spec doc to `C`. Update the header, delete the experimental config warning, and update the version numbers in the description of the `version` field. Update the `nav_order` to one less than the previous stable release.
- Update `docs/specs.md`.
- Update `docs/migrating-configs.md` for the new spec version. Copy the relevant section from Ignition's `doc/migrating-configs.md`, convert the configs to FCCs, convert field names to snake case, and update wording as needed. Add subsections for any new FCC-specific features.
