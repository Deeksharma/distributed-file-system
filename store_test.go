package main

import (
	"bytes"
	"fmt"
	"io"
	"testing"
)

func TestPathTransformFunc(t *testing.T) {
	key := "momsbestpicture"
	pathKey := CASPathTransformFunc(key)
	//fmt.Println(pathName)
	expectedOriginal := "6804429f74181a63c50c3d81d733a12f14a353ff"
	expectedPathName := "68044/29f74/181a6/3c50c/3d81d/733a1/2f14a/353ff"
	if pathKey.PathName != expectedPathName {
		t.Errorf("path name should be %s, got %s", expectedPathName, pathKey.PathName)
	}
	if pathKey.FileName != expectedOriginal {
		t.Errorf("original should be %s, got %s", expectedOriginal, pathKey.FileName)
	}
}

func TestNewStore(t *testing.T) {
	opts := StoreOpts{PathTransformFunc: CASPathTransformFunc}
	s := NewStore(opts)

	key := "momsbestpicture"

	data := []byte("some jpeg value")
	if err := s.writeStream(key, bytes.NewReader(data)); err != nil {
		t.Error(err)
	}

	r, err := s.Read(key)
	if err != nil {
		t.Error(err)
	}
	b, _ := io.ReadAll(r)
	if string(b) != string(data) {
		t.Errorf("data should be %s, got %s", data, b)
	}
	fmt.Println(string(b))

}
