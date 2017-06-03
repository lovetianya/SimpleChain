package peer

import (
	"os"
	"testing"
)

func TestDBAccess(t *testing.T) {
	var s State = State{
		path:   "",
		dbName: "test.db",
		tbName: "tbTest",
		db:     nil}

	err := s.Open()
	if err != nil {
		t.Errorf("Db[%s/%s] open failed!", s.dbName, s.tbName)
	}
	defer os.Remove(s.dbName)
	defer s.Close()

	key := "no answer"
	value, err := s.Get([]byte(key))
	if err != nil {
		t.Errorf("Db get error KEY: [%s], ERROR: [%s]", key)
	}
	if value != nil && len(value) != 0 {
		t.Errorf("Db get with non-exist key[%s] should return nil", key)
	}

	err = s.Close()
	if err != nil {
		t.Errorf("Close db [%v] failed!", s.dbName)
	}

	err = s.Open()
	if err != nil {
		t.Errorf("Open db [%v] failed after close!", s.dbName)
	}

	testData := map[string]string{
		"key1": "data1",
		"key2": "data2",
		"key3": "data3",
		"key4": "data4",
		"key5": "data5",
		"key6": "data6",
	}

	for k, v := range testData {
		err = s.Put([]byte(k), []byte(v))
		if err != nil {
			t.Errorf("put value [%s, %s] failed", k, v)
		}
	}

	for k, v := range testData {
		value, err := s.Get([]byte(k))
		if err != nil || string(value) != v {
			t.Errorf("get value failed! err:[%v], key[%v], value[%v]",
				err, k, value)
		}
	}
}
