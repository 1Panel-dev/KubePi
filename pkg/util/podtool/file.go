package podtool

type File struct {
	Name      string `json:"name"`
	Size      string `json:"size"`
	Mode      string `json:"mode"`
	IsDir     bool   `json:"isDir"`
	ModTime   string `json:"modTime"`
	User      string `json:"user"`
	UserGroup string `json:"group"`
	Link      string `json:"link"`
	Path      string `json:"path"`
}
