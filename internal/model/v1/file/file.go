package file

import "io"

type Request struct {
	Cluster       string    `json:"cluster" validate:"required"`
	Folder        string    `json:"folder"`
	PodName       string    `json:"podName" validate:"required"`
	ContainerName string    `json:"containerName"`
	Namespace     string    `json:"namespace" validate:"required"`
	Path          string    `json:"path"`
	OldPath       string    `json:"oldPath"`
	Commands      []string  `json:"-"`
	Stdin         io.Reader `json:"-"`
	Content       string    `json:"content"`
	FilePath      string    `json:"filePath"`
}
