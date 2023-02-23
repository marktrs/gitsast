package git

import (
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestIsFileTypeIgnored(t *testing.T) {
	for _, ft := range ignoredFileTypes {
		assert.True(t, isFileTypeIgnored(filepath.Join(uuid.New().Domain().String(), ft)))
	}
}
