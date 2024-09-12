package assets

import (
	"embed"
	"errors"
	"io"
	"io/fs"
	"time"
)

//go:embed css img js webfonts
var AssetsFS embed.FS

// ================================================================ //
// ================================================================ //

func NewAssetsFS() *WrapperFS {
	time, _ := time.Parse(time.RFC3339, "2024-01-01T00:00:00Z")
	return &WrapperFS{FS: AssetsFS, FixedModTime: time}
}

type WrapperFS struct {
	embed.FS
	FixedModTime time.Time
}

type FileWrapper struct {
	fs.File
	fixedModTime time.Time
}

type FileInfoWrapper struct {
	fs.FileInfo
	fixedModTime time.Time
}

func (f *WrapperFS) Open(name string) (fs.File, error) {
	file, err := f.FS.Open(name)
	if err != nil {
		return nil, err
	}
	return &FileWrapper{File: file, fixedModTime: f.FixedModTime}, nil
}

func (f *FileWrapper) Stat() (fs.FileInfo, error) {
	fileInfo, err := f.File.Stat()
	return &FileInfoWrapper{FileInfo: fileInfo, fixedModTime: f.fixedModTime}, err
}

func (f *FileWrapper) Read(p []byte) (n int, err error) {
	return f.File.Read(p)
}

func (f *FileWrapper) Seek(offset int64, whence int) (int64, error) {
	i, ok := f.File.(io.Seeker)
	if !ok {
		return 0, errors.New("Seek not implemented")
	}
	return i.Seek(offset, whence)
}

func (f *FileInfoWrapper) ModTime() time.Time {
	return f.fixedModTime
}
