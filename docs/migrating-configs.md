---
layout: default
nav_order: 5
---

# Upgrading Configs

Occasionally, there are changes made to Fedora CoreOS configuration that break backward compatibility. While this is not a concern for running machines (since Ignition only runs one time during first boot), it is a concern for those who maintain configuration files. This document serves to detail each of the breaking changes and tries to provide some reasoning for the change. This does not cover all of the changes to the spec - just those that need to be considered when migrating from one version to the next.

{: .no_toc }

1. TOC
{:toc}

## From Version 1.2.0 to 1.3.0

There are no breaking changes between versions 1.2.0 and 1.3.0 of the configuration specification. Any valid 1.2.0 configuration can be updated to a 1.3.0 configuration by changing the version string in the config.

The following is a list of notable new features, deprecations, and changes.

### Boot disk mirroring and LUKS

The config gained a new top-level `boot_device` section with `luks` and `mirror` subsections, which provide a simple way to configure encryption and/or mirroring for the boot disk. When `luks` is specified, the root filesystem is encrypted and can be unlocked with any combination of a TPM2 device and network Tang servers. When `mirror` is specified, all default partitions are replicated across multiple disks, allowing the system to survive disk failure. On aarch64 or ppc64le systems, the `layout` field must be set to `aarch64` or `ppc64le` to select the correct partition layout.

<!-- fedora-coreos-config -->
```yaml
variant: fcos
version: 1.3.0
boot_device:
  layout: ppc64le
  mirror:
    devices:
      - /dev/sda
      - /dev/sdb
  luks:
    tang:
      - url: https://tang.example.com
        thumbprint: REPLACE-THIS-WITH-YOUR-TANG-THUMBPRINT
    tpm2: true
    threshold: 2
```

## From Version 1.1.0 to 1.2.0

There are no breaking changes between versions 1.1.0 and 1.2.0 of the configuration specification. Any valid 1.1.0 configuration can be updated to a 1.2.0 configuration by changing the version string in the config.

The following is a list of notable new features, deprecations, and changes.

### Partition resizing

The `partition` section gained a new `resize` field. When true, Ignition will resize an existing partition if it matches the config in all respects except the partition size.

<!-- fedora-coreos-config -->
```yaml
variant: fcos
version: 1.2.0
storage:
  disks:
    - device: /dev/sda
      partitions:
        - label: root
          size_mib: 16384
          resize: true
```

### LUKS encrypted storage

Ignition now supports creating LUKS2 encrypted storage volumes. Volumes can be configured to allow unlocking with any combination of a TPM2 device via Clevis, network Tang servers via Clevis, and static key files. Alternatively, the Clevis configuration can be manually specified with a custom PIN and CFG. If a key file is not specified for a device, an ephemeral one will be created.

<!-- fedora-coreos-config -->
```yaml
variant: fcos
version: 1.2.0
storage:
  luks:
    - name: static-key-example
      device: /dev/sdb
      key_file:
        inline: REPLACE-THIS-WITH-YOUR-KEY-MATERIAL
    - name: tpm-example
      device: /dev/sdc
      clevis:
        tpm2: true
    - name: tang-example
      device: /dev/sdd
      clevis:
        tang:
          - url: https://tang.example.com
            thumbprint: REPLACE-THIS-WITH-YOUR-TANG-THUMBPRINT
  filesystems:
    - path: /var/lib/static_key_example
      device: /dev/disk/by-id/dm-name-static-key-example
      format: ext4
      label: STATIC-EXAMPLE
      with_mount_unit: true
    - path: /var/lib/tpm_example
      device: /dev/disk/by-id/dm-name-tpm-example
      format: ext4
      label: TPM-EXAMPLE
      with_mount_unit: true
    - path: /var/lib/tang_example
      device: /dev/disk/by-id/dm-name-tang-example
      format: ext4
      label: TANG-EXAMPLE
      with_mount_unit: true
```

### User/group deletion

The `passwd` `users` and `groups` sections have a new field `should_exist`. If specified and false, Ignition will delete the specified user or group if it exists.

<!-- fedora-coreos-config -->
```yaml
variant: fcos
version: 1.2.0
passwd:
  users:
    - name: core
      should_exist: false
  groups:
    - name: core
      should_exist: false
```

### Google Cloud Storage URL support

The sections which allow fetching a remote URL now accept Google Cloud Storage (`gs://`) URLs in the `source` field.

<!-- fedora-coreos-config -->
```yaml
variant: fcos
version: 1.2.0
storage:
  files:
    - path: /etc/example
      mode: 0644
      contents:
        source: gs://bucket/object
```

## From Version 1.0.0 to 1.1.0

There are no breaking changes between versions 1.0.0 and 1.1.0 of the configuration specification. Any valid 1.0.0 configuration can be updated to a 1.1.0 configuration by changing the version string in the config.

The following is a list of notable new features, deprecations, and changes.

### Embedding local files in configs

The config `merge` and `replace` sections, the `certificate_authorities` section, and the files `contents` and `append` sections gained a new field called `local`, which is mutually exclusive with the `source` and `inline` fields. It causes the contents of a file from the system running FCCT to be embedded in the config. The specified path is relative to a local _files-dir_, specified with the `-d`/`--files-dir` option to FCCT. If no _files-dir_ is specified, this functionality is unavailable.

<!-- fedora-coreos-config -->
```yaml
variant: fcos
version: 1.1.0
ignition:
  config:
    merge:
      - local: config.ign
  security:
    tls:
      certificate_authorities:
        - local: ca.pem
storage:
  files:
    - path: /opt/file
      contents:
        local: file
      append:
        - local: file-epilogue
      mode: 0644
```

### Embedding directory trees in configs

The `storage` section gained a new subsection called `trees`. It is a list of directory trees on the system running FCCT whose files, directories, and symlinks will be embedded in the config. By default, the resulting filesystem objects are owned by `root:root`, directory modes are set to 0755, and file modes are set to 0755 if the source file is executable or 0644 otherwise. Attributes of files, directories, and symlinks can be overridden by creating an entry in the `files`, `directories`, or `links` section; such `files` entries must omit `contents` and such `links` entries must omit `target`.

Tree paths are relative to a local _files-dir_, specified with the `-d`/`--files-dir` option to FCCT. If no _files-dir_ is specified, this functionality is unavailable.

<!-- fedora-coreos-config -->
```yaml
variant: fcos
version: 1.1.0
storage:
  trees:
    - local: tree
      path: /etc/files
  files:
    - path: /etc/files/overridden-file
      mode: 0600
      user:
        id: 500
      group:
        id: 501
```

### Inline contents on certificate authorities and merged configs

The `certificate_authorities` section now supports inline contents via the `inline` field. The config `merge` and `replace` sections also now support `inline`, but using this functionality is not recommended.

<!-- fedora-coreos-config -->
```yaml
variant: fcos
version: 1.1.0
ignition:
  config:
    merge:
      - inline: |
          {"ignition": {"version": "3.1.0"}}
  security:
    tls:
      certificate_authorities:
        - inline: |
            -----BEGIN CERTIFICATE-----
            MIICzTCCAlKgAwIBAgIJALTP0pfNBMzGMAoGCCqGSM49BAMCMIGZMQswCQYDVQQG
            EwJVUzETMBEGA1UECAwKQ2FsaWZvcm5pYTEWMBQGA1UEBwwNU2FuIEZyYW5jaXNj
            bzETMBEGA1UECgwKQ29yZU9TIEluYzEUMBIGA1UECwwLRW5naW5lZXJpbmcxEzAR
            BgNVBAMMCmNvcmVvcy5jb20xHTAbBgkqhkiG9w0BCQEWDm9lbUBjb3Jlb3MuY29t
            MB4XDTE4MDEyNTAwMDczOVoXDTI4MDEyMzAwMDczOVowgZkxCzAJBgNVBAYTAlVT
            MRMwEQYDVQQIDApDYWxpZm9ybmlhMRYwFAYDVQQHDA1TYW4gRnJhbmNpc2NvMRMw
            EQYDVQQKDApDb3JlT1MgSW5jMRQwEgYDVQQLDAtFbmdpbmVlcmluZzETMBEGA1UE
            AwwKY29yZW9zLmNvbTEdMBsGCSqGSIb3DQEJARYOb2VtQGNvcmVvcy5jb20wdjAQ
            BgcqhkjOPQIBBgUrgQQAIgNiAAQDEhfHEulYKlANw9eR5l455gwzAIQuraa049Rh
            vM7PPywaiD8DobteQmE8wn7cJSzOYw6GLvrL4Q1BO5EFUXknkW50t8lfnUeHveCN
            sqvm82F1NVevVoExAUhDYmMREa6jZDBiMA8GA1UdEQQIMAaHBH8AAAEwHQYDVR0O
            BBYEFEbFy0SPiF1YXt+9T3Jig2rNmBtpMB8GA1UdIwQYMBaAFEbFy0SPiF1YXt+9
            T3Jig2rNmBtpMA8GA1UdEwEB/wQFMAMBAf8wCgYIKoZIzj0EAwIDaQAwZgIxAOul
            t3MhI02IONjTDusl2YuCxMgpy2uy0MPkEGUHnUOsxmPSG0gEBCNHyeKVeTaPUwIx
            AKbyaAqbChEy9CvDgyv6qxTYU+eeBImLKS3PH2uW5etc/69V/sDojqpH3hEffsOt
            9g==
            -----END CERTIFICATE-----
```

### Compression support for certificate authorities and merged configs

The config `merge` and `replace` sections and the `certificate_authorities` section now support gzip-compressed resources via the `compression` field. `gzip` compression is supported for all URL schemes except `s3`.

<!-- fedora-coreos-config -->
```yaml
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

<!-- fedora-coreos-config -->
```yaml
variant: fcos
version: 1.1.0
storage:
  files:
    - path: /etc/hosts
      mode: 0644
      contents:
        source: https://example.com/etc/hosts
        verification:
          hash: sha256-e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
```

### Automatic generation of mount units

The `filesystems` section gained a new `with_mount_unit` field. If `true`, a generic mount unit will be automatically generated for the specified filesystem.

<!-- fedora-coreos-config -->
```yaml
variant: fcos
version: 1.1.0
storage:
  filesystems:
    - path: /var/data
      device: /dev/vdb1
      format: ext4
      with_mount_unit: true
```

### Filesystem mount options

The `filesystems` section gained a new `mount_options` field. It is a list of options Ignition should pass to `mount -o` when mounting the specified filesystem. This is useful for mounting btrfs subvolumes. If the `with_mount_unit` field is `true`, this field also affects mount options used by the provisioned system when mounting the filesystem.

<!-- fedora-coreos-config -->
```yaml
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
      with_mount_unit: true
```

### Custom HTTP headers

The sections which allow fetching a remote URL &mdash; config `merge` and `replace`, `certificate_authorities`, and file `contents` and `append` &mdash; gained a new field called `http_headers`. This field can be set to an array of HTTP headers which will be added to an HTTP or HTTPS request. Custom headers can override Ignition's default headers, and will not be retained across HTTP redirects.

During config merging, if a child config specifies a header `name` but not a corresponding `value`, any header with that `name` in the parent config will be removed.

<!-- fedora-coreos-config -->
```yaml
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

<!-- fedora-coreos-config -->
```yaml
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
