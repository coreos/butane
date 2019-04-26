package v0_1

type CaReference struct {
	Source       string       `yaml:"source"`
	Verification Verification `yaml:"verification"`
}

type Config struct {
	Ignition Ignition `yaml:"ignition"`
	Passwd   Passwd   `yaml:"passwd"`
	Storage  Storage  `yaml:"storage"`
	Systemd  Systemd  `yaml:"systemd"`
}

type ConfigReference struct {
	Source       *string      `yaml:"source"`
	Verification Verification `yaml:"verification"`
}

type Device string

type Directory struct {
	Group     NodeGroup `yaml:"group"`
	Overwrite *bool     `yaml:"overwrite"`
	Path      string    `yaml:"path"`
	User      NodeUser  `yaml:"user"`
	Mode      *int      `yaml:"mode"`
}

type Disk struct {
	Device     string      `yaml:"device"`
	Partitions []Partition `yaml:"partitions"`
	WipeTable  *bool       `yaml:"wipeTable"`
}

type Dropin struct {
	Contents *string `yaml:"contents"`
	Name     string  `yaml:"name"`
}

type File struct {
	Group     NodeGroup      `yaml:"group"`
	Overwrite *bool          `yaml:"overwrite"`
	Path      string         `yaml:"path"`
	User      NodeUser       `yaml:"user"`
	Append    []FileContents `yaml:"append"`
	Contents  FileContents   `yaml:"contents"`
	Mode      *int           `yaml:"mode"`
}

type FileContents struct {
	Compression  *string      `yaml:"compression"`
	Source       *string      `yaml:"source"`
	Verification Verification `yaml:"verification"`
}

type Filesystem struct {
	Device         string             `yaml:"device"`
	Format         *string            `yaml:"format"`
	Label          *string            `yaml:"label"`
	Options        []FilesystemOption `yaml:"options"`
	Path           *string            `yaml:"path"`
	UUID           *string            `yaml:"uuid"`
	WipeFilesystem *bool              `yaml:"wipeFilesystem"`
}

type FilesystemOption string

type Group string

type Ignition struct {
	Config   IgnitionConfig `yaml:"config"`
	Security Security       `yaml:"security"`
	Timeouts Timeouts       `yaml:"timeouts"`
}

type IgnitionConfig struct {
	Merge   []ConfigReference `yaml:"merge"`
	Replace ConfigReference   `yaml:"replace"`
}

type Link struct {
	Group     NodeGroup `yaml:"group"`
	Overwrite *bool     `yaml:"overwrite"`
	Path      string    `yaml:"path"`
	User      NodeUser  `yaml:"user"`
	Hard      *bool     `yaml:"hard"`
	Target    string    `yaml:"target"`
}

type NodeGroup struct {
	ID   *int    `yaml:"id"`
	Name *string `yaml:"name"`
}

type NodeUser struct {
	ID   *int    `yaml:"id"`
	Name *string `yaml:"name"`
}

type Partition struct {
	GUID               *string `yaml:"guid"`
	Label              *string `yaml:"label"`
	Number             int     `yaml:"number"`
	ShouldExist        *bool   `yaml:"shouldExist"`
	SizeMiB            *int    `yaml:"sizeMiB"`
	StartMiB           *int    `yaml:"startMiB"`
	TypeGUID           *string `yaml:"typeGuid"`
	WipePartitionEntry *bool   `yaml:"wipePartitionEntry"`
}

type Passwd struct {
	Groups []PasswdGroup `yaml:"groups"`
	Users  []PasswdUser  `yaml:"users"`
}

type PasswdGroup struct {
	Gid          *int    `yaml:"gid"`
	Name         string  `yaml:"name"`
	PasswordHash *string `yaml:"passwordHash"`
	System       *bool   `yaml:"system"`
}

type PasswdUser struct {
	Gecos             *string            `yaml:"gecos"`
	Groups            []Group            `yaml:"groups"`
	HomeDir           *string            `yaml:"homeDir"`
	Name              string             `yaml:"name"`
	NoCreateHome      *bool              `yaml:"noCreateHome"`
	NoLogInit         *bool              `yaml:"noLogInit"`
	NoUserGroup       *bool              `yaml:"noUserGroup"`
	PasswordHash      *string            `yaml:"passwordHash"`
	PrimaryGroup      *string            `yaml:"primaryGroup"`
	SSHAuthorizedKeys []SSHAuthorizedKey `yaml:"sshAuthorizedKeys"`
	Shell             *string            `yaml:"shell"`
	System            *bool              `yaml:"system"`
	UID               *int               `yaml:"uid"`
}

type Raid struct {
	Devices []Device     `yaml:"devices"`
	Level   string       `yaml:"level"`
	Name    string       `yaml:"name"`
	Options []RaidOption `yaml:"options"`
	Spares  *int         `yaml:"spares"`
}

type RaidOption string

type SSHAuthorizedKey string

type Security struct {
	TLS TLS `yaml:"tls"`
}

type Storage struct {
	Directories []Directory  `yaml:"directories"`
	Disks       []Disk       `yaml:"disks"`
	Files       []File       `yaml:"files"`
	Filesystems []Filesystem `yaml:"filesystems"`
	Links       []Link       `yaml:"links"`
	Raid        []Raid       `yaml:"raid"`
}

type Systemd struct {
	Units []Unit `yaml:"units"`
}

type TLS struct {
	CertificateAuthorities []CaReference `yaml:"certificateAuthorities"`
}

type Timeouts struct {
	HTTPResponseHeaders *int `yaml:"httpResponseHeaders"`
	HTTPTotal           *int `yaml:"httpTotal"`
}

type Unit struct {
	Contents *string  `yaml:"contents"`
	Dropins  []Dropin `yaml:"dropins"`
	Enabled  *bool    `yaml:"enabled"`
	Mask     *bool    `yaml:"mask"`
	Name     string   `yaml:"name"`
}

type Verification struct {
	Hash *string `yaml:"hash"`
}
