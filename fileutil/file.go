package fileutil

import (
	"fmt"
	"os"
)

func CheckFile(filename string) {
	_, err := os.Stat(filename)
	if err == nil {
		os.Remove(filename)
	}
}

func CreateFileAndWrite(filename string, data []byte) {
	CheckFile(filename)
	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("create: %s err: %s \n", filename, err)
		return
	}
	file.Write(data)
	file.Close()
}
