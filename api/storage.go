package api

import (
	"errors"
	"github.com/graymeta/stow"
	_ "github.com/graymeta/stow/google"
	"github.com/graymeta/stow/local"
	_ "github.com/graymeta/stow/s3"
	"github.com/netsage-project/grafana-dashboard-manager/config"
	log "github.com/sirupsen/logrus"
	"github.com/thoas/go-funk"
	"io/ioutil"
	"net/url"
	"os"
	"strings"
)

type Storage interface {
	ReadFile(string) ([]byte, error)
	WriteFile(string, []byte) error
	ReadDir(string) ([]string, error)
}

type LocalStorageImpl struct {
	locations map[string]stow.Location
	config    *config.GrafanaConfig
}

func (s *LocalStorageImpl) ReadDir(path string) ([]string, error) {
	location := s.getLocation(path)
	log.Info(location)
	containers, _, err := location.Containers("", stow.CursorStart, 10)
	if err != nil {
		return nil, err
	}
	files := []string{}
	for _, container := range containers {
		items, f, _ := container.Items("", "", 1000)
		log.Info(f)
		for _, item := range items {
			itemUrl := item.URL()
			files = append(files, itemUrl.Path)
		}
	}

	return files, nil
}

func (s *LocalStorageImpl) getLocation(path string) stow.Location {
	if val, ok := s.locations[path]; ok {
		return val
	}
	return nil
}
func (s *LocalStorageImpl) WriteFile(path string, data []byte) error {
	l, err := s.getLocationByPath(path)
	if err != nil {
		log.Fatalf("No valid location found")
	}
	itemUrl, err := url.Parse(path)
	if err != nil {
		panic(err)
	}
	item, err := l.ItemByURL(itemUrl)
	r, err := item.Open()
	if err != nil {
		return err
	}
	defer r.Close()
	err = ioutil.WriteFile(path, data, os.FileMode(int(0666)))
	return err
}

func (s *LocalStorageImpl) getLocationByPath(path string) (stow.Location, error) {
	var l stow.Location
	for _, key := range funk.Keys(s.locations).([]string) {
		if strings.Contains(path, key) {
			l = s.getLocation(key)
			return l, nil
		}
	}
	return nil, errors.New("no Valid Location Found")
}

func (s *LocalStorageImpl) ReadFile(path string) ([]byte, error) {
	l, err := s.getLocationByPath(path)
	if err != nil {
		log.Fatalf("No valid location found")
	}

	itemUrl, err := url.Parse(path)
	if err != nil {
		panic(err)
	}
	log.Infof("Reading data to %s using %s", path, l)
	item, err := l.ItemByURL(itemUrl)
	r, err := item.Open()
	if err != nil {
		return nil, err
	}
	defer r.Close()
	data, err := ioutil.ReadAll(r)
	return data, nil
}

func NewLocalStorageImpl(grafanaConf *config.GrafanaConfig) *LocalStorageImpl {
	locations := make(map[string]stow.Location)
	for _, i := range []string{grafanaConf.OutputDataSource, grafanaConf.OutputDashboard} {
		cfg := stow.ConfigMap{"path": i}
		l, err := stow.Dial(local.Kind, cfg)
		if err != nil {
			log.Fatalf("Cannot instantiate a valid Location %s", err)
		}
		locations[i] = l

	}
	return &LocalStorageImpl{
		locations: locations,
		config:    grafanaConf,
	}
}
