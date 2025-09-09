[![Build Status][travis-image]][travis-url]
[![Appveyor Status][appveyor-image]][appveyor-url]
[![GoDoc][godoc-image]][godoc-url]
[![GitHub release][release-image]][release-url]

# xignore

A golang package for pattern matching of file paths. Like gitignore, dockerignore chefignore.


## Requirements

* Golang ≥ 1.11


## Use

```golang
result, _ := xignore.DirMatches("/workspace/my_project", &xignore.MatchesOptions{
	Ignorefile: ".gitignore",
	Nested: true, // Handle nested ignorefile
})

// ignorefile rules matched files
fmt.Printf("%#v\n", result.MatchedFiles)
// ignorefile rules unmatched files
fmt.Printf("%#v\n", result.UnmatchedFiles)
// ignorefile rules matched dirs
fmt.Printf("%#v\n", result.MatchedDirs)
// ignorefile rules unmatched dirs
fmt.Printf("%#v\n", result.UnmatchedDirs)
```


## LICENSE
[MIT](https://github.com/zealic/xignore/blob/master/LICENSE.txt)


## Reference

* https://git-scm.com/docs/gitignore
* https://github.com/moby/moby/blob/master/pkg/fileutils/fileutils.go

[travis-image]:   https://travis-ci.org/zealic/xignore.svg
[travis-url]:     https://travis-ci.org/zealic/xignore
[appveyor-image]: https://ci.appveyor.com/api/projects/status/avl6x1qease9bjic?svg=true
[appveyor-url]:   https://ci.appveyor.com/project/zealic/xignore
[godoc-image]:    https://godoc.org/github.com/zealic/xignore?status.svg
[godoc-url]:      https://godoc.org/github.com/zealic/xignore
[release-image]:  https://img.shields.io/github/release/zealic/xignore.svg
[release-url]:    https://github.com/zealic/xignore/releases
