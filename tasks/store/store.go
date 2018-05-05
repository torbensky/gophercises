package store

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"time"

	"github.com/coreos/bolt"
)

// TODOBucket is the name of the BoltDB bucket storing our todos
var TODOBucket = []byte("todos")

// Contains information about a TODO list Task
type Task struct {
	CompletedAt time.Time
	Description string
}

// The Generic API for TODO's
type TodoService interface {
	AddTask(*Task) error
	RemoveTaskNum(int) (*Task, error)
	CompleteTaskNum(int) (*Task, error)
	ListUncompleteTasks() []*Task
	ListTasksCompletedAfter(t time.Time) []*Task
	Close()
}

// boltStore is a boltdb based persistence store for todos
type boltStore struct {
	db *bolt.DB
}

// newBolt constructs a new instance of the boltdb todo store implementation
// dbFile is a path to the TODO Task store (a boltdb file)
func NewBolt(dbFile string) (TodoService, error) {
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		return nil, err
	}

	// Initialize bucket right away
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(TODOBucket)
		return err
	})

	return &boltStore{
		db: db,
	}, nil
}

// AddTask adds a Task to the end of the current todo list
func (s *boltStore) AddTask(t *Task) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(TODOBucket)
		id, err := b.NextSequence()
		if err != nil {
			return err
		}
		taskBytes, err := serializeTask(t)
		if err != nil {
			return err
		}

		return b.Put(itob(id), taskBytes)
	})
}

// CompleteTaskNum removes the i'th Task from the current todo list
func (s *boltStore) CompleteTaskNum(i int) (*Task, error) {
	var task *Task
	err := s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(TODOBucket)

		// Find the ith Task
		c := b.Cursor()
		var taskId []byte
		for k, v := c.First(); k != nil; k, v = c.Next() {
			i--
			if i == 0 {
				// we are at the ith item
				taskId = k
				t, err := deserializeTask(v)
				if err != nil {
					return err
				}
				task = t
				break
			}
		}

		if taskId == nil {
			return errors.New("invalid task number")
		}

		// Set completed
		task.CompletedAt = time.Now()

		// Save update to store
		taskBytes, err := serializeTask(task)
		if err != nil {
			return err
		}
		return b.Put(taskId, taskBytes)
	})

	return task, err
}

// RemoveTaskNum removes the i'th Task from the current todo list
func (s *boltStore) RemoveTaskNum(i int) (*Task, error) {
	var task *Task
	err := s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(TODOBucket)

		// Find the ith Task
		c := b.Cursor()
		var taskId []byte
		for k, v := c.First(); k != nil; k, v = c.Next() {
			i--
			if i == 0 {
				// we are at the ith item
				taskId = k
				t, err := deserializeTask(v)
				if err != nil {
					return err
				}
				task = t
				break
			}
		}

		if taskId == nil {
			return errors.New("invalid task number")
		}

		// Delete it from store
		return b.Delete(taskId)
	})

	return task, err
}

// ListTasksCompletedAfter lists tasks that were completed after the specified time
func (s *boltStore) ListTasksCompletedAfter(t time.Time) []*Task {
	var tasks []*Task
	// TODO: (ironic) maybe dont silently suppress this deserialization error. Out of date db probably...
	s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(TODOBucket)
		if b != nil {
			c := b.Cursor()
			for k, v := c.First(); k != nil; k, v = c.Next() {
				task, err := deserializeTask(v)
				if err != nil {
					return err
				}
				if !task.CompletedAt.IsZero() {
					tasks = append(tasks, task)
				}
			}
		}

		return nil
	})
	return tasks
}

// ListUncompleteTasks lists all the Tasks in the database
func (s *boltStore) ListUncompleteTasks() []*Task {
	var tasks []*Task
	s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(TODOBucket)
		if b != nil {
			c := b.Cursor()
			for k, v := c.First(); k != nil; k, v = c.Next() {
				task, err := deserializeTask(v)
				if err != nil {
					return err
				}
				// Skip completed tasks
				if task.CompletedAt.IsZero() {
					tasks = append(tasks, task)
				}
			}
		}

		return nil
	})

	return tasks
}

func deserializeTask(b []byte) (*Task, error) {
	var t Task
	err := json.Unmarshal(b, &t)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func serializeTask(t *Task) ([]byte, error) {
	return json.Marshal(t)
}

// Close closes the underlying boltdb instance
func (s *boltStore) Close() {
	s.db.Close()
}

// itob returns an 8-byte big endian representation of v.
// source: https://github.com/boltdb/bolt#autoincrementing-integer-for-the-bucket
func itob(v uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, v)
	return b
}
