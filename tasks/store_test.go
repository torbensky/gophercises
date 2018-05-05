package main

import (
	"crypto/rand"
	"encoding/hex"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBoltStore(t *testing.T) {
	tmpDbFile := tempFileName("taskstest", ".db")
	b, err := newBolt(tmpDbFile)
	defer b.Close()
	assert.NoError(t, err)

	err = b.addTask(&task{
		Description: "test 123",
	})
	assert.NoError(t, err)

	tasks := b.listTasks()
	assert.Equal(t, 1, len(tasks))
	assert.Equal(t, "test 123", tasks[0].Description)

	// Remove the task we just added
	tsk, err := b.removeTaskNum(1)
	assert.NoError(t, err)
	assert.Equal(t, "test 123", tsk.Description)

	// Try removing non-existent task - should error
	tsk, err = b.removeTaskNum(1)
	assert.Error(t, err)
	assert.Nil(t, tsk)
}

// tempFileName generates a temporary filename for use in testing or whatever
// source: https://stackoverflow.com/questions/28005865/golang-generate-unique-filename-with-extension
func tempFileName(prefix, suffix string) string {
	randBytes := make([]byte, 16)
	rand.Read(randBytes)
	return filepath.Join(os.TempDir(), prefix+hex.EncodeToString(randBytes)+suffix)
}
