package task

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/Gustcat/archiver_170725/internal/logger"
	"github.com/Gustcat/archiver_170725/internal/model"
	"io"
	"log/slog"
	"strconv"
)

var UnsupportedFileType = errors.New("unsupported file type")
var NoLinkConnection = errors.New("no connection on link")

func (s *serv) Update(ctx context.Context, id string, sourceLink string) error {
	const op = "service.task.Update"
	log := logger.LogFromContextAddOP(ctx, op)

	task, sourcesCount, err := s.taskRepo.Update(ctx, id, sourceLink)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Debug("Start check and download resource")
	go func() {
		serr := s.checkAndDownloadSource(ctx, sourceLink, task, sourcesCount)
		s.taskRepo.ProcessDownloadSignal(ctx, task, serr)
	}()

	return nil
}

func (s *serv) checkAndDownloadSource(ctx context.Context, sourceLink string, task *model.Task, sourcesCount int) error {
	const op = "service.task.checkAndDownloadSource"
	log := logger.LogFromContextAddOP(ctx, op)

	log.Debug("Try to get %s link", slog.String("link", sourceLink))
	resp, err := s.client.Get(sourceLink)
	if err != nil {
		return fmt.Errorf("failed to get file from %s: %w", sourceLink, err)
	}
	defer func(Body io.ReadCloser) {
		if closeErr := Body.Close(); closeErr != nil {
			log.Error("failed to close response body", slog.String("error", closeErr.Error())) //
			if err == nil {
				err = closeErr
			} else {
				err = errors.Join(err, closeErr)
			}
		}
	}(resp.Body)

	header := make([]byte, 4)
	n, err := io.ReadFull(resp.Body, header)
	if err != nil {
		return UnsupportedFileType
	}

	log.Debug("Header of resource %s", slog.String("header", fmt.Sprintf("%x", header)))
	extension := ""
	switch {
	case isJPEG(header[:n]):
		extension = ".jpg"
	case isPDF(header[:n]):
		extension = ".pdf"
	default:
		return UnsupportedFileType
	}

	sourceNumberStr := strconv.Itoa(sourcesCount + 1)
	fileName := sourceNumberStr + extension

	buf := &bytes.Buffer{}

	if _, err := buf.Write(header[:n]); err != nil {
		return fmt.Errorf("failed write header to buffer: %w", err)
	}

	if _, err := io.Copy(buf, resp.Body); err != nil {
		return fmt.Errorf("failed copy body to buffer: %w", err)
	}

	err = s.taskRepo.DownloadSource(ctx, task, fileName, buf.Bytes())
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func isJPEG(header []byte) bool {
	if len(header) >= 3 && header[0] == 0xFF && header[1] == 0xD8 && header[2] == 0xFF {
		return true
	}

	return false
}

func isPDF(header []byte) bool {
	if len(header) >= 4 && bytes.Equal(header[:4], []byte("%PDF")) {
		return true
	}

	return false
}
