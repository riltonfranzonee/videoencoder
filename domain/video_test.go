package domain_test

import (
	"encoder/domain"
	"testing"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

func TestEnsureVideoIsNotEmpty(t *testing.T) {
	video := domain.NewVideo()
	err := video.Validate()
	require.Error(t, err)
}

func TestEnsureVideoIdIsValid(t *testing.T) {
	video := domain.NewVideo()
	video.ID = "abc"
	video.ResourceID = "a"
	video.FilePath = "/test"
	video.CreatedAt = time.Now()
	err := video.Validate()
	require.Error(t, err)
}

func TestEnsureVideoIsCreated(t *testing.T) {
	video := domain.NewVideo()
	video.ID = uuid.NewV4().String()
	video.ResourceID = "a"
	video.FilePath = "/test"
	video.CreatedAt = time.Now()
	err := video.Validate()
	require.Nil(t, err)
}
