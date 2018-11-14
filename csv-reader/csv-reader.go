package csv_reader

import (
	"bufio"
	"os"
	"strings"
)

type CSVReader struct {
	file    *os.File
	reader  *bufio.Reader
	sep     string
	headers []string
}

// NewCSVReader creates new CSVReader
func NewCSVReader(filePath string) (*CSVReader, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	c := &CSVReader{
		file:   file,
		reader: bufio.NewReader(file),
		sep:    ",",
	}

	if err := c.readHeader(); err != nil {
		return nil, err
	}

	return c, nil
}

// SetSeparator used to set custom separator
func (c *CSVReader) SetSeparator(sep string) {
	c.sep = sep
}

// readHeader gets first line which includes headers/columns
func (c *CSVReader) readHeader() error {
	l, _, err := c.reader.ReadLine()
	if err != nil {
		return err
	}
	c.headers = strings.Split(string(l), c.sep)

	return nil
}

// Parse reads one line of file
func (c *CSVReader) Parse(data *map[string]string) error {
	l, _, err := c.reader.ReadLine()
	if err != nil {
		return err
	}
	cells := strings.Split(string(l), c.sep)
	max := len(cells)

	for i, key := range c.headers {
		if i < max {
			(*data)[key] = cells[i]
		}
	}

	return nil
}

// Close
func (c *CSVReader) Close() error {
	if c.file != nil {
		return c.file.Close()
	}
	return nil
}
