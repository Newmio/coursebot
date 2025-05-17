package core

import "cbot/pkg"

type FileVaultImpl struct {
}

func CreateFileVault() pkg.FileVault {
	return &FileVaultImpl{}
}
