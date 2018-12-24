package utils

import (
	"bufio"
	"encoding/csv"
	"github.com/pkg/errors"
	"io"
	"os"
	"path/filepath"
)

// Opens the specified file
func OpenFile(filename string) *os.File {
	pathToFile, err := filepath.Abs(filename)
	if err != nil {
		panic(errors.Wrapf(err, "failed to get absolute path of %s", filename))
	}
	file, err := os.Open(pathToFile)
	if err != nil {
		panic(errors.Wrapf(err, "failed to open file %s", filename))
	}
	return file
}

// Reads the file and processes the file with the provided function
func ReadCSVFile(file *os.File, processLine func(record []string, line int)) {
	ReadCSVFileWithDelimiter(file, ',', processLine)
}

func ReadCSVFileWithDelimiter(file *os.File, delimiter rune, processLine func(record []string, line int)) {
	reader := csv.NewReader(bufio.NewReader(file))
	reader.Comma = delimiter
	line := 0
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(errors.Wrapf(err, "failed to read line $d", line))
		}
		processLine(record, line)
		line++
	}
}

func ReadTXTFile(file *os.File, lines chan string) {
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines <- scanner.Text()
	}
	close(lines)
}

// Closes the file and panics if an error occurs
func CloseFile(file *os.File) {
	err := file.Close()
	if err != nil {
		panic(errors.Wrapf(err, "failed to close %s", file.Name()))
	}
}
