package utils_test

import (
	"encoder/framework/utils"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsJon(t *testing.T) {
	json := `{
		"id" : "123456",
		"file_path": "test.mp4",
		"status": "pending"  
	}`

	err := utils.IsJson(json)
	require.Nil(t, err)

	json = `invalid text`
	err = utils.IsJson(json)
	require.NotNil(t, err)
}
