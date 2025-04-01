package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func CASPathTransformFunc(root string, key string) PathKey {
	hash := sha1.Sum([]byte(key)) // [20]byte => convert this to slice then just use hash[:]
	hashStr := hex.EncodeToString(hash[:])

	blockSize := 5
	sliceLen := len(hashStr) / blockSize

	paths := make([]string, sliceLen)

	for i := 0; i < sliceLen; i++ {
		from, to := i*blockSize, (i+1)*blockSize
		paths[i] = hashStr[from:to]
	}

	return PathKey{
		PathName: filepath.Join(root, strings.Join(paths, "/")),
		FileName: hashStr,
	}
}

type ParthTransformFunc func(string, string) PathKey
type StoreOpts struct {
	// Root is the folder name of the root, containing all the folders/files of the system
	Root              string
	PathTransformFunc ParthTransformFunc
}

var DefaultPathTransformFunction = func(root string, key string) PathKey {
	return PathKey{
		PathName: filepath.Join(root, key),
		FileName: filepath.Join(root, key),
	}
}

const defaultRootFolderName = "deeksha-distributed-store"

type Store struct {
	StoreOpts
}

func NewStore(opts StoreOpts) *Store {
	if opts.PathTransformFunc == nil {
		opts.PathTransformFunc = DefaultPathTransformFunction
	}
	if len(opts.Root) == 0 {
		opts.Root = defaultRootFolderName
	}
	return &Store{
		StoreOpts: opts,
	}
}

type PathKey struct {
	PathName string
	FileName string
}

func (p PathKey) FullPath() string {
	return fmt.Sprintf("%s/%s", p.PathName, p.FileName)
}

func (s *Store) Has(key string) bool {
	pathKey := s.PathTransformFunc(s.Root, key)

	_, err := os.Stat(pathKey.FullPath())
	return !errors.Is(err, os.ErrNotExist)
}

func (s *Store) Clear() error {
	return os.RemoveAll(s.Root)
}

func (s *Store) Delete(key string) error {
	pathKey := s.PathTransformFunc(s.Root, key)

	defer func() {
		log.Printf("deleted [%s] from disk", pathKey.FileName)
	}()
	absPath, _ := filepath.Abs(pathKey.PathName)
	// remove all is not deleting all the directories
	if err := os.RemoveAll(pathKey.FullPath()); err != nil {
		return err
	}

	err := os.RemoveAll(absPath)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) Write(key string, r io.Reader) error {
	return s.writeStream(key, r)
}

func (s *Store) Read(key string) (io.Reader, error) {
	f, err := s.readStream(key)
	defer f.Close()
	if err != nil {
		return nil, err
	}
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, f)
	return buf, err
}

func (s *Store) readStream(key string) (io.ReadCloser, error) {
	pathKey := s.PathTransformFunc(s.Root, key)
	//fullFilePathName := pathKey.FileName

	f, err := os.Open(pathKey.FullPath())
	if err != nil {
		return nil, err
	}

	return f, nil
}

// content addressable storage where we can store anything and transformation on keys
func (s *Store) writeStream(key string, r io.Reader) error {
	pathKey := s.PathTransformFunc(s.Root, key)

	if err := os.MkdirAll(pathKey.PathName, os.ModePerm); err != nil {
		return err
	}

	// the contents that we are putting cannot be the filename
	//buf := new(bytes.Buffer)
	//io.Copy(buf, r)
	//
	//filenameBytes := md5.Sum(buf.Bytes())
	fullFilePathName := pathKey.FullPath()

	f, err := os.Create(fullFilePathName)
	if err != nil {
		return err
	}
	defer f.Close()

	n, err := io.Copy(f, r)
	if err != nil {
		return err
	}

	log.Printf("written (%d) bytes to the disk: %s", n, fullFilePathName)
	return nil
}
