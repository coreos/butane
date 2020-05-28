# Migrating Between Configuration Versions

Occasionally, there are changes made to Fedora CoreOS configuration that break backward compatibility. While this is not a concern for running machines (since Ignition only runs one time during first boot), it is a concern for those who maintain configuration files. This document serves to detail each of the breaking changes and tries to provide some reasoning for the change. This does not cover all of the changes to the spec - just those that need to be considered when migrating from one version to the next.

## From Version 1.0.0 to 1.1.0

There are no breaking changes between versions 1.0.0 and 1.1.0 of the configuration specification. Any valid 1.0.0 configuration can be updated to a 1.1.0 configuration by changing the version string in the config.

The following is a list of notable new features, deprecations, and changes.

### Compression support for certificate authorities and merged configs

The config `merge` and `replace` sections and the `certificate_authorities` section now support gzip-compressed resources via the `compression` field. `gzip` compression is supported for all URL schemes except `s3`.

```yaml fedora-coreos-config
variant: fcos
version: 1.1.0
ignition:
  config:
    merge:
      - source: https://secure.example.com/example.ign.gz
        compression: gzip
  security:
    tls:
      certificate_authorities:
        - source: https://example.com/ca.pem.gz
          compression: gzip
```

### SHA-256 resource verification

All `verification.hash` fields now support the `sha256` hash type.

```yaml fedora-coreos-config
variant: fcos
version: 1.1.0
storage:
  files:
    - path: /etc/hosts
      mode: 644
      contents:
        source: https://example.com/etc/hosts
        verification:
          hash: sha256-e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
```

### Filesystem mount options

The `filesystems` section gained a new `mount_options` field. It is a list of options Ignition should pass to `mount -o` when mounting the specified filesystem. This is useful for mounting btrfs subvolumes. This field only affects mounting performed by Ignition while it is running; it does not affect mounting of the filesystem by the provisioned system.

```yaml fedora-coreos-config
variant: fcos
version: 1.1.0
storage:
  filesystems:
    - path: /var/data
      device: /dev/vdb1
      wipe_filesystem: false
      format: btrfs
      mount_options:
        - subvolid=5
```

### Custom HTTP headers

The sections which allow fetching a remote URL &mdash; config `merge` and `replace`, `certificate_authorities`, and file `contents` and `append` &mdash; gained a new field called `http_headers`. This field can be set to an array of HTTP headers which will be added to an HTTP or HTTPS request. Custom headers can override Ignition's default headers, and will not be retained across HTTP redirects.

During config merging, if a child config specifies a header `name` but not a corresponding `value`, any header with that `name` in the parent config will be removed.

```yaml fedora-coreos-config
variant: fcos
version: 1.1.0
storage:
  files:
    - path: /etc/hosts
      mode: 0644
      contents:
        source: https://example.com/etc/hosts
        http_headers:
          - name: Authorization
            value: Basic YWxhZGRpbjpvcGVuc2VzYW1l
          - name: User-Agent
            value: Mozilla/5.0 (compatible; MSIE 6.0; Windows NT 5.1)
```

### HTTP proxies

The `ignition` section gained a new field called `proxy`. It allows configuring proxies for HTTP and HTTPS requests, as well as exempting certain hosts from proxying.

The `https_proxy` field specifies the proxy URL for HTTPS requests. The `http_proxy` field specifies the proxy URL for HTTP requests, and also for HTTPS requests if `https_proxy` is not specified. The `no_proxy` field lists specifiers of hosts that should not be proxied, in any of several formats:

- An IP address prefix (`1.2.3.4`)
- An IP address prefix in CIDR notation (`1.2.3.4/8`)
- A domain name, matching the domain and its subdomains (`example.com`)
- A domain name, matching subdomains only (`.example.com`)
- A wildcard matching all hosts (`*`)

IP addresses and domain names can also include a port number (`1.2.3.4:80`).

```yaml fedora-coreos-config
variant: fcos
version: 1.1.0
ignition:
  proxy:
    http_proxy: https://proxy.example.net/
    https_proxy: https://secure.proxy.example.net/
    no_proxy:
     - www.example.net
storage:
  files:
    - path: /etc/hosts
      mode: 0644
      contents:
        source: https://example.com/etc/hosts
```
