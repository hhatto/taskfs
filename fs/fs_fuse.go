// +build linux

package fs

import (
	"os"
	"time"

	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"
)

type Dir struct {
	nodefs.Node
}

func NewDir() Dir {
	return Dir{
		Node: nodefs.NewDefaultNode(),
	}
}

type File struct {
	nodefs.File
}

func NewFile(p []byte) File {
	return File{
		File: nodefs.NewDataFile(p),
	}
}

func (f *FileInfo) FillAttr(out *fuse.Attr) {
	out.Mode = uint32(f.Mode & os.ModePerm)
	if f.IsDir() {
		out.Mode |= fuse.S_IFDIR
	}
	out.Size = uint64(f.Size)
	out.Atime = uint64(f.LastMod.Unix())
	out.Mtime = uint64(f.LastMod.Unix())
}

func (f *FileInfo) FillDirEntry(out *fuse.DirEntry) {
	out.Name = f.Name
	out.Mode = uint32(f.Mode & os.ModePerm)
	if f.IsDir() {
		out.Mode |= fuse.S_IFDIR
	}
}

func (root *Root) MountAndServe(mtpt string) error {
	opts := nodefs.Options{
		AttrTimeout:  time.Second,
		EntryTimeout: time.Second,
		Debug:        false,
	}
	s, _, err := nodefs.MountRoot(mtpt, root, &opts)
	if err != nil {
		return err
	}
	s.Serve()
	return nil
}

func (root *Root) GetAttr(out *fuse.Attr, file nodefs.File, ctx *fuse.Context) fuse.Status {
	root.FileInfo.FillAttr(out)
	return fuse.OK
}

func (root *Root) OpenDir(ctx *fuse.Context) ([]fuse.DirEntry, fuse.Status) {
	kids := root.ReadDir()
	a := make([]fuse.DirEntry, len(kids))
	for i, kid := range kids {
		kid.FillDirEntry(&a[i])
	}
	return a, fuse.OK
}
