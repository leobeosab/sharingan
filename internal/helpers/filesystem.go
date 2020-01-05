package helpers

import (
	"bufio"
	"io"
	"os"
)

func GetNumberOfLinesInFile(f *os.File) int {
	c := 0
	s := bufio.NewScanner(f)

	for s.Scan() {
		c++
	}

	f.Seek(0, io.SeekStart)
	return c
}
