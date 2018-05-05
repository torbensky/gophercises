package main

import (
	"encoding/binary"
	"errors"

	"github.com/boltdb/bolt"
)

// TODOBucket is the name of the BoltDB bucket storing our todos
var TODOBucket = []byte("todos")

// Contains information about a TODO list task
type task struct {
	Description string
}

// The Generic API for TODO's
type TodoService interface {
	addTask(*task) error
	removeTaskNum(int) (*task, error)
	listTasks() []*task
	close() error
}

// boltStore is a boltdb based persistence store for todos
type boltStore struct {
	db *bolt.DB
}

// newBolt constructs a new instance of the boltdb todo store implementation
// dbFile is a path to the TODO task store (a boltdb file)
func newBolt(dbFile string) (*boltStore, error) {
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

func (s *boltStore) addTask(t *task) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(TODOBucket)
		id, err := b.NextSequence()
		if err != nil {
			return err
		}
		return b.Put(itob(id), []byte(t.Description))
	})
}

func (s *boltStore) removeTaskNum(i int) (*task, error) {
	var t *task
	err := s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(TODOBucket)

		// Find the ith task
		c := b.Cursor()
		var taskId []byte
		for k, v := c.First(); k != nil; k, v = c.Next() {
			i--
			if i == 0 {
				// we are at the ith item
				taskId = k
				t = &task{
					Description: string(v), // copies bytes
				}
				break
			}
		}

		if taskId == nil {
			return errors.New("invalid task number")
		}

		// Delete it from store
		return b.Delete(taskId)
	})

	return t, err
}

// listTasks lists all the tasks in the database
func (s *boltStore) listTasks() []*task {
	var tasks []*task
	s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(TODOBucket)
		if b != nil {
			c := b.Cursor()
			for k, v := c.First(); k != nil; k, v = c.Next() {
				tasks = append(tasks, &task{
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
