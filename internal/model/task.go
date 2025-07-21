package model

import "archive/zip"

type Status string

const (
	StatusNew        Status = "new"
	StatusInProgress Status = "in_progress"
	StatusDone       Status = "done"
)

type TaskResult struct {
	Status      Status   `json:"status"`
	ArchiveLink *FileUrl `json:"archive"`
}

type Task struct {
	Sources      []FileUrl
	Status       Status
	ArchiveLink  FileUrl
	DownloadDown chan struct{}
	ZipWriter    *zip.Writer
}

type TaskId struct {
	ID string `json:"id" uri:"id" binding:"required"`
}

type SourceRequest struct {
	Source FileUrl
}
