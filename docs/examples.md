# Examples

Here you can find a bunch of simple examples for using `fcct`, with some explanations about what they do. The examples here are in no way comprehensive, for a full list of all the options present in `fcct` check out the [configuration specification][spec].

## Users and groups

This example modifies the existing `core` user and sets its ssh key.

```yaml fedora-coreos-config
variant: fcos
version: 1.0.0
passwd:
  users:
    - name: core
      ssh_authorized_keys:
        - key1
```

This example creates one user, `user1` and sets up one ssh public key for the user. The user is also given the home directory `/home/user1`, but it's not created, the user is added to the `wheel` and `plugdev` groups, and the user's shell is set to `/bin/bash`.

```yaml fedora-coreos-config
variant: fcos
version: 1.0.0
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

```yaml fedora-coreos-config
variant: fcos
version: 1.0.0
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

This example fetches a gzip-compressed file from `http://example.com/file2`, makes sure that it matches the provided sha512 hash, and writes it to `/opt/file2`.

```yaml fedora-coreos-config
variant: fcos
version: 1.0.0
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

## systemd units

This example adds a drop-in for the `serial-getty@ttyS0` unit, turning on autologin on `ttyS0` by overriding the `ExecStart=` defined in the default unit. More information on systemd dropins can be found in [the systemd docs][dropins].

```yaml fedora-coreos-config
variant: fcos
version: 1.0.0
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

```yaml fedora-coreos-config
variant: fcos
version: 1.0.0
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

### Filesystems and Partitions

This example creates a single partition spanning all of the sdb device then creates a btrfs filesystem on it to use as /var. Finally it creates the mount unit for systemd so it gets mounted on boot.

```yaml fedora-coreos-config
variant: fcos
version: 1.0.0
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
systemd:
  units:
    - name: var.mount
      enabled: true
      contents: |
        [Unit]
        Before=local-fs.target
        [Mount]
        Where=/var
        What=/dev/disk/by-partlabel/var
        [Install]
        WantedBy=local-fs.target
```

[spec]: configuration.md
[dropins]: https://www.freedesktop.org/software/systemd/man/systemd.unit.html#Description
