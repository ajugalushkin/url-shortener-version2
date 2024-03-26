package storage

import (
	"bufio"
	"encoding/json"
	"github.com/ajugalushkin/url-shortener-version2/internal/model"
	"os"
	"path/filepath"
)

type Writer struct {
	file   *os.File
	writer *bufio.Writer
}

func NewWriter(fileName string) (*Writer, error) {
	directory, _ := filepath.Split(fileName)
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		err := os.Mkdir(directory, os.ModePerm)
		if err != nil {
			return nil, err
		}
	}

	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	return &Writer{
		file:   file,
		writer: bufio.NewWriter(file),
	}, nil
}

func (p *Writer) WriteFile(file *model.File) error {
	data, err := json.Marshal(&file)
	if err != nil {
		return err
	}

	if _, err := p.writer.Write(data); err != nil {
		return err
	}

	if err := p.writer.WriteByte('\n'); err != nil {
		return err
	}

	return p.writer.Flush()
}

func (p *Writer) Close() error {
	return p.file.Close()
}
