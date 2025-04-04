package main

import (
	"bytes"
	"fmt"
	"io"
	"testing"
)

func TestPathTransformFunc(t *testing.T) {
	key := "momsbestpicture"
	pathKey := CASPathTransformFunc(defaultRootFolderName, key)
	//log.Println(pathName)
	expectedOriginal := "6804429f74181a63c50c3d81d733a12f14a353ff"
	expectedPathName := "deeksha-distributed-store/68044/29f74/181a6/3c50c/3d81d/733a1/2f14a/353ff"
	if pathKey.PathName != expectedPathName {
		t.Errorf("path name should be %s, got %s", expectedPathName, pathKey.PathName)
	}
	if pathKey.FileName != expectedOriginal {
		t.Errorf("original should be %s, got %s", expectedOriginal, pathKey.FileName)
	}
}

func TestStore_Delete(t *testing.T) {
	s := newStore()
	defer teardown(t, s)
	key := "momsbestpicture"

	data := []byte("some jpeg value")
	if err := s.writeStream(key, bytes.NewReader(data)); err != nil {
		t.Error(err)
	}

	if err := s.Delete(key); err != nil {
		t.Error(err)
	}
}

func TestStore(t *testing.T) {
	//opts := StoreOpts{PathTransformFunc: CASPathTransformFunc}
	//s := NewStore(opts)
	s := newStore()
	defer teardown(t, s)

	for i := 0; i < 10; i++ {
		key := fmt.Sprintf("momsbestpicture_%d", i)

		data := []byte("some jpeg value")
		if err := s.writeStream(key, bytes.NewReader(data)); err != nil {
			t.Error(err)
		}

		if ok := s.Has(key); !ok {
			t.Errorf("key should exist %s", key)
		}

		r, err := s.Read(key)
		if err != nil {
			t.Error(err)
		}
		b, _ := io.ReadAll(r)
		if string(b) != string(data) {
			t.Errorf("data should be %s, got %s", data, b)
		}
		log.Println(string(b))

		if err := s.Delete(key); err != nil {
			t.Error(err)
		}
	}
}

func newStore() *Store {
	return &Store{StoreOpts{
		Root:              defaultRootFolderName,
		PathTransformFunc: CASPathTransformFunc,
	}}
}

func teardown(t *testing.T, s *Store) {
	if err := s.Clear(); err != nil {
		t.Error(err)
	}
}
