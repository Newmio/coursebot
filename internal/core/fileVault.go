package core

import (
	"cbot/pkg"
	"io"
	"os"
	"path"
)

type FileVaultImpl struct {
	mDir string
}

func CreateFileVault() pkg.FileVault {
	return &FileVaultImpl{mDir: "docs/"}
}

func (obj *FileVaultImpl) Create(data io.Reader, name string) error {
	if err := os.MkdirAll(path.Dir(obj.mDir+name), 0755); err != nil {
		return pkg.Trace(err)
	}

	file, err := os.Create(obj.mDir + name)
	if err != nil {
		return pkg.Trace(err)
	}
	defer file.Close()

	if data != nil {
		_, err = io.Copy(file, data)
		if err != nil {
			return pkg.Trace(err)
		}
	}

	return nil
}

func (obj *FileVaultImpl) Get(name string) (*os.File, error) {
	return os.Open(obj.mDir + name)
}

func (obj *FileVaultImpl) GetDir(name string) string {
	return obj.mDir + name
}

func (obj *FileVaultImpl) Delete(name string) error {
	return os.Remove(obj.mDir + name)
}