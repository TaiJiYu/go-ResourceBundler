package api

import "go-resource-bundler/component"

type IPacker interface {
	Add(key string, data []byte)
	AddResourceDir(dir string)
	AddFile(key string, path string)
	PackResource()
	Save() error
}

// packer
func NewPacker(o ...component.PackOption) IPacker { return component.NewPacker(o...) }

// unpacker
func UnpackerInit(o ...component.UnpackOption) error { return component.UnpackerInit(o...) }
func Key(key string) []byte                          { return component.Key(key) }
func Close()                                         { component.Close() }
func Show()                                          { component.Show() }
