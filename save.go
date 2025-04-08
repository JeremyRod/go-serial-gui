package main

import (
	"os"
)

type File struct {
	file *os.File
}

func NewFile(path string) (File, error) {
	file, err := os.Create(path)
	if err != nil {
		return File{nil}, err
	}
	return File{file}, nil

}
func (f *File) Close() error {
	err := f.file.Close()
	if err != nil {
		return err
	}
	return nil
}

func (f *File) Save(data []byte) {
	f.file.Write(data)
}

func (f *File) Load() string {
	data, err := os.ReadFile(f.file.Name())
	if err != nil {
		return ""
	}
	return string(data)
}
