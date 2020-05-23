package service

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"os"
)

type image struct {
	Name	string		`json:"name"`
	Size	string		`json:"size"`
	ModTime	string		`json:"modTime"`
}

var Images = make(map[string]image)

func (i *image) Remove(order Order) error {
	content := order.Content
	image := image{}
	err := json.Unmarshal([]byte(content), &image)
	if err != nil {
		log.Errorf("Json unmarshal error %v", err)
		ErrorPublic(err)
	}

	imageName := image.Name + ".tar"
	imagePath := DefaultRootPath + "/" +imageName
	err = os.Remove(imagePath)
	if err != nil {
		log.Errorf("Remove Image %s error %v", imagePath, err)
		ErrorPublic(err)
	}
	return nil
}