package service

type image struct {
	Name	string		`json:"name"`
	Size	string		`json:"size"`
	ModTime	string		`json:"modTime"`
}

var Images = make(map[string]image)