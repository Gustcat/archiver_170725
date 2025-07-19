package model

type Status string

const (
	StatusNew        Status = "new"
	StatusInProgress Status = "in_progress"
	StatusDone       Status = "done"
)

type TaskResult struct {
	Status  Status   `json:"status"`
	Archive *FileUrl `json:"archive"`
}

type Task struct {
	Sources []FileUrl
	Status  Status
	Archive *FileUrl
}

type TaskId struct {
	ID int64 `json:"id"`
}

type SourceRequest struct {
	Source FileUrl
}
