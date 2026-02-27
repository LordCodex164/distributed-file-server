package main

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type PathKey struct {
	Pathname string
	Filename string
}

func CASPathTransformFunc(key string) PathKey {
	hash := sha1.Sum([]byte(key))
	fmt.Println(">>hash", hash)
	hashStr := hex.EncodeToString(hash[:])
	fmt.Println(">>hashStr", hashStr)

	blockSize := 5

	sliceLen := len(hashStr) / blockSize

	paths := make([]string, sliceLen)

	for i := 0; i < sliceLen; i++ {
		from, to := i*blockSize, i*blockSize+blockSize
		fmt.Println("from", from, to)
		paths[i] = hashStr[from:to]
	}

	return PathKey{
		Pathname: strings.Join(paths, "/"),
		Filename: hashStr,
	}
}

type PathTransformFunc func(string) PathKey

var DefaultPathTransformFunc = func(key string) PathKey {
	return PathKey{
		Pathname: key,
		Filename: key,
	}
}

func (p PathKey) FullPath() string {
	return fmt.Sprintf("%s/%s", p.Filename, p.Pathname)
}

func (p PathKey) FirstPathName() string {
	arr := strings.Split(p.Pathname, "/")
	if len(arr) == 0 {
		return ""
	}
	firstPath := arr[0]
	return fmt.Sprintf("Here is the first path name: %s", firstPath)
}

type StoreOpts struct {
	// Root is the folder name of the root, containing all the folders/files of the system.
	Root              string
	PathTransformFunc PathTransformFunc
}
type Store struct {
	StoreOpts
}

func NewStore(opts StoreOpts) *Store {
	if opts.PathTransformFunc != nil {
		opts.PathTransformFunc = DefaultPathTransformFunc
	}
	return &Store{
		StoreOpts: opts,
	}
}

func (s *Store) Has(id string, key string) bool {
	pathKey := s.PathTransformFunc(key)
	fullPathWithRoot := fmt.Sprintf("%s/%s/%s", s.Root, id, pathKey.FullPath())

	_, err := os.Stat(fullPathWithRoot)
	return !errors.Is(err, os.ErrNotExist)
}

func (s *Store) Clear() error {
	return os.RemoveAll(s.Root)
}

func (s *Store) Delete(id string, key string) error {
	pathKey := s.PathTransformFunc(key)

	defer func() {
		log.Printf("deleted [%s] from disk", pathKey.Filename)
	}()

	firstPathNameWithRoot := fmt.Sprintf("%s/%s/%s", s.Root, id, pathKey.FirstPathName())

	return os.RemoveAll(firstPathNameWithRoot)
}


func (s *Store) Write(id string, key string, r io.Reader) (int64, error) {
	return s.WriteStream(id, key, r)
}

func (s *Store) Read(id string, key string) (int64, io.Reader, error) {
	return s.ReadStream(id, key)
}

func (s *Store) OpenFileForWriting(id string, key string) (*os.File, error) {
	// open file for wriing
	pathKey := s.PathTransformFunc(key)

	pathNameWithRoot := fmt.Sprintf("%s/%s/%s", s.Root, id, pathKey.Pathname)
	if err := os.MkdirAll(pathNameWithRoot, os.ModePerm); err != nil {
		return nil, err
	}
	fullPathWithRoot := fmt.Sprintf("%s/%s/%s", s.Root, id, pathKey.FullPath())
	fmt.Println("fullpath", fullPathWithRoot)
	return os.Create(fullPathWithRoot)
}

func (s *Store) WriteStream(id string, key string, r io.Reader) (int64, error) {
	file, err := s.OpenFileForWriting(id, key)

	if err != nil {
		return 0, err
	}
	return io.Copy(file, r)
}

func (s *Store) ReadStream(id string, key string) (int64, io.Reader, error) {
	pathKey := s.PathTransformFunc(key)
	fullPathWithRoot := fmt.Sprintf("%s/%s/%s", s.Root, id, pathKey.FullPath())
	file, err := os.Open(fullPathWithRoot)
	if err != nil {
		return 0, nil, err
	}

	fi, err := file.Stat()
	if err != nil {
		return 0, nil, err
	}

	return fi.Size(), file, nil
}