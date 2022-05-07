package util

import "time"

type File struct {
	Name    string      `json:"Name"`
	Path    string      `json:"Path"`
	Size    int64       `json:"Size"`
	Mode    string      `json:"Mode"`
	ModTime time.Time   `json:"ModTime"`
	IsDir   bool        `json:"IsDir"`
	Sys     interface{} `json:"SysInfo"`
}
