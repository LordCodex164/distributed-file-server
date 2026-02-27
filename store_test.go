package main

import (
	"bytes"
	"fmt"
	"io"
	"testing"
	"time"
)

func TestStore(t *testing.T) {
	// key := "testing_key"
	// originalPathName := "17de3/f4ae7/0d310/2ad86/340f8/5f132/a395c/68348"
	// originalFileName := "17de3f4ae70d3102ad86340f85f132a395c68348"
	// pathKey := CASPathTransformFunc(key)

	// if pathKey.Filename != originalFileName {
	// 	t.Errorf("have %s want %s", pathKey.Filename, originalFileName)
	// }

	// if pathKey.Pathname != originalPathName {
	// 	t.Errorf("have %s want %s", pathKey.Filename, originalPathName)
	// }

	s := newStore()

	key := "test_key"

	data := ([]byte("New Bytes"))

	defer func ()  {
		time.Sleep(3 * time.Second)
		tearDown(s, t)
	}()

	for i := 0; i < 10; i++ {
		//generate an id for every iteration
		id := generateID()
		if _, err := s.Write(id, key, bytes.NewReader(data)); err != nil {
			fmt.Printf("error: %s writing to the file", err)
			return
		}

		//check if the file exist
		if ok := s.Has(id, key); !ok {
			t.Errorf("file of key %s does not exist", key)
			return
		} else {
			t.Logf("file of key %s exist", key)
		}

		//then we read the file
		if size, file_data, err := s.Read(id, key); err != nil {
			fmt.Printf("error reading the file")
		} else {
			r, err := io.ReadAll(file_data)
			if err != nil {
				fmt.Printf("error reading the file data")
			}
			if string(r) != string(data) {
				fmt.Println("size", size)
				t.Errorf("want %s have %s", data, r)
			}
		}

		if err := s.Delete(id, key); err != nil {
			t.Error("error deleting file")
			return
		}
	}

}

func tearDown(s *Store, t *testing.T) {
	if err := s.Clear(); err != nil {
		t.Error("error clearing the root")
	}
	t.Log("The root cleared successfully")
}

func newStore() *Store {
	opts := StoreOpts{
		PathTransformFunc: DefaultPathTransformFunc,
		Root:              "src",
	}
	return &Store{
		StoreOpts: opts,
	}
}
