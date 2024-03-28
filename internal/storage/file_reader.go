package storage

import (
	"bufio"
	"github.com/ajugalushkin/url-shortener-version2/internal/model"
	"os"
)

type Reader struct {
	file    *os.File
	scanner *bufio.Scanner
}

func NewReader(fileName string) (*Reader, error) {
	file, err := os.OpenFile(fileName, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	return &Reader{
		file:    file,
		scanner: bufio.NewScanner(file),
	}, nil
}

func (c *Reader) ReadFile() (map[string]model.File, error) {
	c.scanner.Split(bufio.ScanLines)

	files := make(map[string]model.File)
	for c.scanner.Scan() {
		data := c.scanner.Bytes()

		file := model.File{}
		err := file.UnmarshalJSON(data)
		//err := json.Unmarshal(data, &file)
		if err != nil {
			return nil, err
		}
		files[file.ShortURL] = file
	}

	return files, nil
}

func (c *Reader) Close() error {
	return c.file.Close()
}
