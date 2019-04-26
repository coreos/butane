module github.com/ajeddeloh/fcct

go 1.12

require (
	github.com/ajeddeloh/go-json v0.0.0-20170920214419-6a2fe990e083
	github.com/ajeddeloh/yaml v0.0.0-20141224210557-6b16a5714269
	github.com/coreos/go-semver v0.3.0
	github.com/coreos/ignition/v2 v2.0.0-alpha
	github.com/go-yaml/yaml v2.1.0+incompatible
)

replace github.com/coreos/ignition/v2 => ../ignition
