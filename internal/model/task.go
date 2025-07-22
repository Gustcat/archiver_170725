package model

import (
	"archive/zip"
	"os"
	"sync"
)

type Status string

const (
	StatusNew        Status = "new"
	StatusInProgress Status = "in_progress"
	StatusDone       Status = "done"
	StatusFailed     Status = "failed"
)

type TaskResult struct {
	Status              Status  `json:"status"`
	ArchiveLink         *string `json:"archive,omitempty"`
	ErrorSourcesMessage *string `json:"error_sources_message,omitempty"`
}

type Task struct {
	Mu              sync.RWMutex
	Sources         []string
	Status          Status
	ArchiveLink     string
	DownloadDown    []error
	ZipWriter       *zip.Writer
	Archive         *os.File
	ZipWriterClosed bool
	ArchiveClosed   bool
}

type TaskId struct {
	ID string `json:"id" uri:"id" binding:"required"`
}

type SourceRequest struct {
	Source string
}
