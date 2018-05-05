package store

import (
	"encoding/binary"
	"errors"

	"github.com/boltdb/bolt"
)

// TODOBucket is the name of the BoltDB bucket storing our todos
var TODOBucket = []byte("todos")

// Contains information about a TODO list Task
type Task struct {
	Description string
}

// The Generic API for TODO's
type TodoService interface {
	AddTask(*Task) error
	RemoveTaskNum(int) (*Task, error)
	ListTasks() []*Task
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
		return b.Put(itob(id), []byte(t.Description))
	})
}

// RemoveTaskNum removes the i'th Task from the current todo list
func (s *boltStore) RemoveTaskNum(i int) (*Task, error) {
	var t *Task
	err := s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(TODOBucket)

		// Find the ith Task
		c := b.Cursor()
		var TaskId []byte
		for k, v := c.First(); k != nil; k, v = c.Next() {
			i--
			if i == 0 {
				// we are at the ith item
				TaskId = k
				t = &Task{
					Description: string(v), // copies bytes
				}
				break
			}
		}

		if TaskId == nil {
			return errors.New("invalid task number")
		}

		// Delete it from store
		return b.Delete(TaskId)
	})

	return t, err
}

// ListTasks lists all the Tasks in the database
func (s *boltStore) ListTasks() []*Task {
	var tasks []*Task
	s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(TODOBucket)
		if b != nil {
			c := b.Cursor()
			for k, v := c.First(); k != nil; k, v = c.Next() {
				tasks = append(tasks, &Task{
					Description: string(v), // string() copies bytes
				})
			}
		}

		return nil
	})

	return tasks
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
