---
layout: default
nav_order: 3
---

# Examples
{: .no_toc }

1. TOC
{:toc}

Here you can find a bunch of simple examples for using `fcct`, with some explanations about what they do. The examples here are in no way comprehensive, for a full list of all the options present in `fcct` check out the [configuration specification][spec].

## Users and groups

This example modifies the existing `core` user and sets its ssh key.

<!-- fedora-coreos-config -->
```yaml
variant: fcos
version: 1.1.0
passwd:
  users:
    - name: core
      ssh_authorized_keys:
        - key1
```

This example creates one user, `user1` and sets up one ssh public key for the user. The user is also given the home directory `/home/user1`, but it's not created, the user is added to the `wheel` and `plugdev` groups, and the user's shell is set to `/bin/bash`.

<!-- fedora-coreos-config -->
```yaml
variant: fcos
version: 1.1.0
passwd:
  users:
    - name: user1
      ssh_authorized_keys:
        - key1
      home_dir: /home/user1
      no_create_home: true
      groups:
        - wheel
        - plugdev
      shell: /bin/bash
```

## Storage and files

### Files

This example creates a file at `/opt/file` with the contents `Hello, world!`, permissions 0644 (so readable and writable by the owner, and only readable by everyone else), and the file is owned by user uid 500 and gid 501.

<!-- fedora-coreos-config -->
```yaml
variant: fcos
version: 1.1.0
storage:
  files:
    - path: /opt/file
      contents:
        inline: Hello, world!
      mode: 0644
      user:
        id: 500
      group:
        id: 501
```

This example fetches a gzip-compressed file from `http://example.com/file2`, makes sure that the _uncompressed_ contents match the provided sha512 hash, and writes it to `/opt/file2`.

<!-- fedora-coreos-config -->
```yaml
variant: fcos
version: 1.1.0
storage:
  files:
    - path: /opt/file2
      contents:
        source: http://example.com/file2
        compression: gzip
        verification:
          hash: sha512-4ee6a9d20cc0e6c7ee187daffa6822bdef7f4cebe109eff44b235f97e45dc3d7a5bb932efc841192e46618f48a6f4f5bc0d15fd74b1038abf46bf4b4fd409f2e
      mode: 0644
```

This example creates a file at `/opt/file3` whose contents are read from a local file `local-file3` on the system running FCCT. The path of the local file is relative to a _files-dir_ which must be specified via the `-d`/`--files-dir` option to FCCT.

<!-- fedora-coreos-config -->
```yaml
variant: fcos
version: 1.1.0
storage:
  files:
    - path: /opt/file3
      contents:
        local: local-file3
      mode: 0644
```

### Directory trees

Consider a directory tree at `~/fcc/tree` on the system running FCCT:

```
file
overridden-file
directory/file
directory/symlink -> ../file
```

This example copies that directory tree to `/etc/files` on the target system. The ownership and mode for `overridden-file` are explicitly set by the config. All other filesystem objects are owned by `root:root`, directory modes are set to 0755, and file modes are set to 0755 if the source file is executable or 0644 otherwise. The example must be transpiled with `--files-dir ~/fcc`.

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

### Filesystems and Partitions

This example creates a single partition spanning all of the sdb device then creates a btrfs filesystem on it to use as /var. Finally it creates the mount unit for systemd so it gets mounted on boot.

<!-- fedora-coreos-config -->
```yaml
variant: fcos
version: 1.1.0
storage:
  disks:
    - device: /dev/sdb
      wipe_table: true
      partitions:
        - number: 1
          label: var
  filesystems:
    - path: /var
      device: /dev/disk/by-partlabel/var
      format: btrfs
      wipe_filesystem: true
      label: var
      with_mount_unit: true
```

## systemd units

This example adds a drop-in for the `serial-getty@ttyS0` unit, turning on autologin on `ttyS0` by overriding the `ExecStart=` defined in the default unit. More information on systemd dropins can be found in [the systemd docs][dropins].

<!-- fedora-coreos-config -->
```yaml
variant: fcos
version: 1.1.0
systemd:
  units:
    - name: serial-getty@ttyS0.service
      dropins:
        - name: autologin.conf
          contents: |
            [Service]
            TTYVTDisallocate=no
            ExecStart=
            ExecStart=-/usr/sbin/agetty --autologin core --noclear %I $TERM
```

This example creates a new systemd unit called hello.service, enables it so it will run on boot, and defines the contents to simply echo `"Hello, World!"`.

<!-- fedora-coreos-config -->
```yaml
variant: fcos
version: 1.1.0
systemd:
  units:
    - name: hello.service
      enabled: true
      contents: |
        [Unit]
        Description=A hello world unit!
        [Service]
        Type=oneshot
        ExecStart=/usr/bin/echo "Hello, World!"
        [Install]
        WantedBy=multi-user.target
```

[spec]: specs.md
[dropins]: https://www.freedesktop.org/software/systemd/man/systemd.unit.html#Description
