package peer

import (
	"bytes"
	"fmt"
	"github.com/boltdb/bolt"
)

type Ledger struct {
	path   string
	dbName string
	tbName string
	db     *bolt.DB
}

func (l *Ledger) Open() error {
	db, err := bolt.Open(l.path+l.dbName, 0600, nil)
	if err != nil {
		return fmt.Errorf("Create db [%v] failed! Reason: %v",
			l.dbName, err.Error())
	}
	l.db = db

	err = l.db.Update(func(tx *bolt.Tx) error {
		_, err = tx.CreateBucketIfNotExists([]byte(l.tbName))
		if err != nil {
			return fmt.Errorf("Create bucket [%s] failed! Reason: %v",
				l.tbName, err.Error())
		}
		return nil
	})

	return err
}

func (l *Ledger) Close() error {
	if l.db != nil {
		l.db.Close()
	}
	return nil
}

func (l *Ledger) Get(key []byte) ([]byte, error) {
	if l.db == nil {
		err := l.Open()
		if err != nil {
			l.Close()
			return nil, fmt.Errorf("Db [%s] open failed! Reason: %v",
				l.dbName, err.Error())
		}
	}

	value := bytes.NewBuffer(nil)

	err := l.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(l.tbName))
		if b == nil {
			return fmt.Errorf("Tx bucket[%s] got failed in leger.Got!", l.tbName)
		}
		v := b.Get(key)
		if v != nil {
			value.Write(v)
		}
		return nil
	})
	return value.Bytes(), err
}

func (l *Ledger) Put(key []byte, value []byte) error {
	if l.db == nil {
		err := l.Open()
		if err != nil {
			l.Close()
			return fmt.Errorf("Db [%s] open failed! Reason: %v",
				l.dbName, err.Error())
		}
	}

	err := l.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(l.tbName))
		if b == nil {
			return fmt.Errorf("Tx bucket[%s] got failed in ledger.Put!", l.tbName)
		}
		err := b.Put(key, value)
		return err
	})

	if err == nil {
		return nil
	}
	return fmt.Errorf("Put value (%v, %v) failed in ledger.Put", key, value)
}
