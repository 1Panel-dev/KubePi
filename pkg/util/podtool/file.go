package podtool

type File struct {
	Name string `json:"name"`
	//Path    string      `json:"path"`
	Size string `json:"size"`
	Mode string `json:"mode"`
	//ModTime time.Time   `json:"modTime"`
	IsDir bool `json:"isDir"`
	//Sys     interface{} `json:"sys"`
	ModTime   string `json:"modTime"`
	User      string `json:"user"`
	UserGroup string `json:"group"`
}
