# Developing FCCT

## Creating a release

Create a [release checklist](https://github.com/coreos/fcct/issues/new?template=release-checklist.md) and follow those steps.

## Bumping spec versions

Up to this point in FCCT development, FCC versions and `base` versions have been 1:1 mapped onto Ignition specs, and we have not had any distro-specific sugar. This checklist therefore describes bumping the Ignition spec version, `base` version, and `config` version, while leaving the distro sugar version unchanged. If your scenario is different, modify to taste.

### Stabilize Ignition spec version

- Bump `go.mod` for new Ignition release and update vendor.
- Update imports. Drop `-experimental` from Ignition spec versions in `base/vB_exp/translate_test.go`.

### Bump base version

- Rename `base/vB_exp` to `base/vB` and update `package` statements. Update imports.
- Copy `base/vB` to `base/vB+1_exp`.
- Update `package` statements in `base/vB+1_exp`.

### Bump config version

- Rename `config/vC_exp` to `config/vC` and update `package` statements. Update imports.
- Drop `-experimental` from `registry` in `config/config.go`.
- Drop `-experimental` from examples in `docs/`.
- Copy `config/vC` to `config/vC+1_exp`.
- Update `package` statements in `config/vC+1_exp`. Bump its base dependency to `base/vB+1_exp`.
- Import `config/vC+1_exp` in `config/config.go` and add `fcos+C+1-experimental` to `registry`.

### Bump Ignition spec version

- Add translation function for experimental spec `I+1` to `distro/fcos/vF`. Revendor Ignition.
- Bump Ignition types imports and rename `ToIgnI` and `TestToIgnI` functions in `base/vB+1_exp`. Bump Ignition spec versions in `base/vB+1_exp/translate_test.go`.
- Bump Ignition types imports in `config/vC+1_exp`. Update `Translate` to call `ToIgnI+1` functions. Update versions in `TranslateBytes` comment.

### Update docs

- Copy the `C-exp` spec doc to `C+1-exp`. Update the header and the version numbers in the description of the `version` field.
- Rename the `C-exp` spec doc to `C`. Update the header, delete the experimental config warning, and update the version numbers in the description of the `version` field.
- Update `docs/migrating-configs.md` for the new spec version. Copy the relevant section from Ignition's `doc/migrating-configs.md`, convert the configs to FCCs, convert field names to snake case, and update wording as needed. Add subsections for any new FCC-specific features.
