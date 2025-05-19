package pkg

import (
	"io"
	"os"
)

type FileVault interface {
	Create(data io.Reader, name string) error
	Get(name string) (*os.File, error)
	GetDir(name string) string
	Delete(name string) error
}
