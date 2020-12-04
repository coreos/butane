---
layout: default
nav_order: 9
---

# Developing FCCT
{: .no_toc }

1. TOC
{:toc}

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
