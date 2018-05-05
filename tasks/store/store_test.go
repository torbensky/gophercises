package store

import (
	"crypto/rand"
	"encoding/hex"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBoltStore(t *testing.T) {
	tmpDbFile := tempFileName("taskstest", ".db")
	b, err := NewBolt(tmpDbFile)
	defer b.Close()
	assert.NoError(t, err)

	// Adda  task
	err = b.AddTask(&Task{
		Description: "test 123",
	})
	assert.NoError(t, err)

	// List available tasks
	tasks := b.ListUncompleteTasks()
	assert.Equal(t, 1, len(tasks))
	assert.Equal(t, "test 123", tasks[0].Description)
	tasks = b.ListTasksCompletedAfter(time.Now().AddDate(0, 0, -1))
	assert.Equal(t, 0, len(tasks))

	// Remove the task we just added
	tsk, err := b.RemoveTaskNum(1)
	assert.NoError(t, err)
	assert.Equal(t, "test 123", tsk.Description)

	// Try removing non-existent task - should error
	tsk, err = b.RemoveTaskNum(1)
	assert.Error(t, err)
	assert.Nil(t, tsk)

	// Add and complete task
	err = b.AddTask(&Task{
		Description: "test 456",
	})
	assert.NoError(t, err)
	_, err = b.CompleteTaskNum(1)
	assert.NoError(t, err)

	// Ensure completed lists this task
	tasks = b.ListTasksCompletedAfter(time.Now().AddDate(0, 0, -1))
	assert.Equal(t, 1, len(tasks))
	assert.Equal(t, "test 456", tasks[0].Description)
	// Ensure uncompleted doesn't list this task
	tasks = b.ListUncompleteTasks()
	assert.Equal(t, 0, len(tasks))

}

// tempFileName generates a temporary filename for use in testing or whatever
// source: https://stackoverflow.com/questions/28005865/golang-generate-unique-filename-with-extension
func tempFileName(prefix, suffix string) string {
	randBytes := make([]byte, 16)
	rand.Read(randBytes)
	return filepath.Join(os.TempDir(), prefix+hex.EncodeToString(randBytes)+suffix)
}
