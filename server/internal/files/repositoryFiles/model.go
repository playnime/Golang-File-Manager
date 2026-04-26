package repositoryfiles

import "time"

type DeleteModel struct {
	id       int
	filename string
	filepath string
	gd       int
}

type DownloadModel struct {
	id       int
	filename string
	filepath string
	gd       int
}

type GiveRightsModel struct {
	user_id     string
	file_name   string
	path        string
	size_bytes  int
	mime_type   string
	gd          int
	uploaded_at time.Time
}
