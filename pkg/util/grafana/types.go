package grafana

// 获取文件夹列表返回数据
type GetFoldersResp struct {
	Id    int    `json:"id"`
	Uid   string `json:"uid"`
	Title string `json:"title"`
}

// 获取仪表盘列表返回数据
type GetDashboardsResp struct {
	Id          int           `json:"id"`
	Uid         string        `json:"uid"`
	Title       string        `json:"title"`
	Uri         string        `json:"uri"`
	Url         string        `json:"url"`
	Slug        string        `json:"slug"`
	Type        string        `json:"type"`
	Tags        []interface{} `json:"tags"`
	IsStarred   bool          `json:"isStarred"`
	FolderId    int           `json:"folderId"`
	FolderUid   string        `json:"folderUid"`
	FolderTitle string        `json:"folderTitle"`
	FolderUrl   string        `json:"folderUrl"`
	SortMeta    int           `json:"sortMeta"`
	IsDeleted   bool          `json:"isDeleted"`
}
