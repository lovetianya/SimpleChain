package peer

import (
	"bytes"
	"fmt"
	"github.com/boltdb/bolt"
)

type State struct {
	path   string
	dbName string
	tbName string
	db     *bolt.DB
}

func (s *State) Open() error {
	db, err := bolt.Open(s.path+s.dbName, 0600, nil)
	if err != nil {
		return fmt.Errorf("Db open[%v] failed! Reason: %v",
			s.dbName, err.Error())
	}
	s.db = db

	err = s.db.Update(func(tx *bolt.Tx) error {
		_, err = tx.CreateBucketIfNotExists([]byte(s.tbName))
		if err != nil {
			return fmt.Errorf("create bucket [%s] failed", s.dbName)
		}
		return nil
	})
	return err
}

func (s *State) Close() error {
	if s.db != nil {
		s.db.Close()
	}
	return nil
}

// TODO: get and put might be deadlock.

func (s *State) Get(key []byte) ([]byte, error) {
	if s.db == nil {
		err := s.Open()
		if err != nil {
			s.Close()
			return nil, err
		}
	}

	value := bytes.NewBuffer(nil)

	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(s.tbName))
		if b == nil {
			return fmt.Errorf("tx bucket[%s] got failed!", s.tbName)
		}
		v := b.Get(key)
		if v != nil {
			value.Write(v)
		}
		return nil
	})

	if err == nil {
		return value.Bytes(), nil
	}
	return nil, fmt.Errorf("Got value failed (key = [%v]) in state.Get", key)
}

func (s *State) Put(key []byte, value []byte) error {
	if s.db == nil {
		err := s.Open()
		if err == nil {
			s.Close()
			return err
		}
	}

	err := s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(s.tbName))
		err := b.Put(key, value)
		return err
	})

	if err == nil {
		return nil
	}

	return fmt.Errorf("Put value (%v, %v) failed in state.Put", key, value)
}
