Release checklist:
 - [ ] Write release notes in NEWS. Get them reviewed and merged
     - [ ] If doing a branched release, also include a PR to merge the NEWS changes into master
 - [ ] Ensure your local copy is up to date with master and your working directory is clean
 - [ ] Ensure you can sign commits and any yubikeys/smartcards are plugged in
 - [ ] Run `./tag_release.sh <vX.Y.z> <git commit hash>`
 - [ ] Push that tag to GitHub
 - [ ] Run `./build_releases`
 - [ ] Sign the release artifacts by running
```
gpg --local-user 0xCDDE268EBB729EC7! --detach-sign --armor <path to artifact>
```
for each release artifact. Do not try to sign all of them at once by globbing. If you do, gpg will sign the combination of all the release artifacts instead of each one individually.

 - [ ] Create a draft release on GitHub and upload all the release artifacts and their signatures. Copy and paste the release notes from NEWS here as well.
 - [ ] Publish the release
 - Update the `release` tag on Quay:
   - [ ] Visit the [Quay tags page](https://quay.io/repository/coreos/fcct?tab=tags) and wait for a versioned tag to appear
   - [ ] Click the gear next to the tag, select "Add New Tag", enter `release`, and confirm
 - Update Fedora RPM:
   - [ ] Create a PR to bump the FCCT spec file in [Fedora](https://src.fedoraproject.org/rpms/fedora-coreos-config-transpiler).
   - [ ] Once that PR merges to `master`, merge `master` into the other relevant branches (e.g. `f31`) then push those.
   - [ ] On each of those branches (including `master`) run `fedpkg build`
   - [ ] Once the builds have finished, submit them to [Bodhi](https://bodhi.fedoraproject.org/updates/new), filling in:
     - `fedora-coreos-config-transpiler` for `Packages`
     - Selecting the build(s) that just completed, except for the Rawhide one (which gets submitted automatically)
     - Copy relevant part of release notes from the GitHub release
     - Leave `Update name` blank
     - `Type`, `Severity` and `Suggestion` can be left as `unspecified` unless it is a security release. In that case select `security` with the appropriate severity.
     - `Stable karma` and `Unstable` karma can be set to `2` and `-1`, respectively.
