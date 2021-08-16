package api

import (
	"github.com/graymeta/stow"
	_ "github.com/graymeta/stow/azure"
	_ "github.com/graymeta/stow/google"
	_ "github.com/graymeta/stow/s3"
	log "github.com/sirupsen/logrus"
)

type Storage interface {
}

type StorageImpl struct {
}

func (s *StorageImpl) StoreFile(path string, date []byte) {
	var v stow.ConfigMap
	log.Info(v)
}

// Dial dials stow storage.
// See stow.Dial for more information.
func Dial(kind string, config stow.Config) (stow.Location, error) {
	v := stow.Location{}
	log.Info(v)
	// v.

	return stow.Dial(kind, config)
}
