package domain_test

import (
	"encoder/domain"
	"testing"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

func TestNewJob(t *testing.T) {
	video := domain.NewVideo()
	video.ID = uuid.NewV4().String()
	video.FilePath = "/test"
	video.CreatedAt = time.Now()

	job, err := domain.NewJob(video.FilePath, "Converted", video)
	require.NotNil(t, job)
	require.Nil(t, err)
}
