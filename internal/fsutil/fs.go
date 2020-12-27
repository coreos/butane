package fsutil

import (
	"errors"
	"io/fs"
	"os"
	"path"
	"reflect"
)

// ErrNotImplemented is the sentinel error returned when an optional
// feature is not implemented.
var ErrNotImplemented = errors.New("not implemented")

// An FS provides access to a hierarchical file system.
type FS interface {
	fs.FS
	fs.StatFS
	ReadLinkFS
}

// ReadLinkFS is the interface implemented by a file system
// that provides an implementation of ReadLink.
type ReadLinkFS interface {
	fs.FS

	// ReadLink returns the destination of the named symbolic link.
	// If there is an error, it will be of type *fs.PathError.
	ReadLink(name string) (string, error)
}

// DirFS returns a file system for the tree of files rooted at the directory dir.
//
// Note that DirFS("/prefix") only guarantees that the Open calls it makes to the
// operating system will begin with "/prefix": DirFS("/prefix").Open("file") is the
// same as os.Open("/prefix/file"). So if /prefix/file is a symbolic link pointing outside
// the /prefix tree, then using DirFS does not stop the access any more than using
// os.Open does. DirFS is therefore not a general substitute for a chroot-style
// security mechanism when the directory tree contains arbitrary content.
func DirFS(dir string) FS {
	return dirFS(path.Clean(dir))
}

type dirFS string

func (dir dirFS) Open(name string) (fs.File, error) {
	if name = path.Clean(name); !fs.ValidPath(name) {
		return nil, &fs.PathError{Op: "open", Path: name, Err: fs.ErrInvalid}
	}
	f, err := os.Open(string(dir) + "/" + name)
	if err != nil {
		return nil, err // nil fs.File
	}
	return f, nil
}

func (dir dirFS) ReadLink(name string) (string, error) {
	if name = path.Clean(name); !fs.ValidPath(name) {
		return "", &fs.PathError{Op: "readlink", Path: name, Err: fs.ErrInvalid}
	}
	return os.Readlink(string(dir) + "/" + name)
}

func (dir dirFS) Stat(name string) (fs.FileInfo, error) {
	if name = path.Clean(name); !fs.ValidPath(name) {
		return nil, &fs.PathError{Op: "stat", Path: name, Err: fs.ErrInvalid}
	}
	return os.Stat(string(dir) + "/" + name)
}

// ReadLink returns the destination of the named symbolic link.
// If there is an error, it will be of type *fs.PathError.
//
// If fsys implements ReadLinkFS, ReadLink calls fs.ReadLink.
// If fsys is the type returned by os.DirFS, ReadLink calls os.Readlink.
// Otherwise an error is returned.
func ReadLink(fsys fs.FS, name string) (string, error) {
	name = path.Clean(name)
	if fsys, ok := fsys.(ReadLinkFS); ok {
		return fsys.ReadLink(name)
	}
	if fsys, ok := toDirFS(fsys); ok {
		return fsys.ReadLink(name)
	}
	return "", &fs.PathError{Op: "readlink", Path: name, Err: ErrNotImplemented}
}

// Stat returns a FileInfo describing the named file from the file system.
//
// If fsys implements StatFS, Stat calls fsys.Stat.
// If fsys is the type returned by os.DirFS, ReadLink calls os.Stat.
// Otherwise, Stat opens the file to stat it.
func Stat(fsys fs.FS, name string) (fs.FileInfo, error) {
	name = path.Clean(name)
	if fsys, ok := fsys.(fs.StatFS); ok {
		return fsys.Stat(name)
	}
	if fsys, ok := toDirFS(fsys); ok {
		return fsys.Stat(name)
	}
	return fs.Stat(fsys, name)
}

var osDirFSType = reflect.TypeOf(os.DirFS(""))

func toDirFS(fsys fs.FS) (FS, bool) {
	if reflect.TypeOf(fsys) == osDirFSType {
		dir := reflect.ValueOf(fsys).String()
		return dirFS(dir), true
	}
	return nil, false
}
