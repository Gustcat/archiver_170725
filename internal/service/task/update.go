package task

import (
	"context"
	"errors"
	"github.com/Gustcat/archiver_170725/internal/model"
)

var ErrOverSourcesLimit = errors.New("over task limit")

const SourcesLimit = 3

func (s *serv) Update(ctx context.Context, id int64, source model.FileUrl) error {

	return nil
}
