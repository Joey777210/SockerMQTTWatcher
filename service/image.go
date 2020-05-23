package service

import (
	"errors"
	log "github.com/Sirupsen/logrus"
	"os"
)

type image struct {
	Name	string		`json:"name"`
	Size	string		`json:"size"`
	ModTime	string		`json:"modTime"`
}

var Images = make(map[string]image)

func (i *image) Remove(iName string) error {
	if iName == "" {
		return errors.New("empty name error")
	}
	i.Name = iName

	imageName := i.Name + ".tar"
	imagePath := DefaultRootPath + "/" +imageName
	err := os.Remove(imagePath)
	if err != nil {
		log.Errorf("Remove Image %s error %v", imagePath, err)
		ErrorPublic(err)
	}
	delete(Images, i.Name)
	return nil
}