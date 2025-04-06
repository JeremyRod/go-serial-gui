package main

import (
	"os"
)

type File struct {
	file *os.File
}

func (f *File) Close() error {
	err := f.file.Close()
	if err != nil {
		return err
	}
	return nil
}

func (f *File) Save(data string) {
	f.file.Write([]byte(data))
}

func (f *File) Load() string {
	data, err := os.ReadFile(f.file.Name())
	if err != nil {
		return ""
	}
	return string(data)
}
