package FileUtils

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

type FileIO struct {
}

func (this *FileIO) ReadStrLines(filePath string) ([]string, error) {
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(b), "\n")
	return lines, nil
}

func (this *FileIO) ReadByteLines(filePath string) ([][]byte, error) {
	var line bytes.Buffer
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	buf := bufio.NewReader(f)
	lines := make([][]byte, 0)
	for {
		data, err := buf.ReadBytes('\n')
		data = bytes.TrimSuffix(data, []byte("\n"))
		if err != nil {
			if bufio.ErrBufferFull == err {
				line.Write(data)
				continue
			}
		}
		if line.Len() > 0 {
			line.Write(data)
			lines = append(lines, line.Bytes())
			line.Reset()
		} else {
			lines = append(lines, data)
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
	}

	return lines, nil
}

func (this *FileIO) WriteFile(path string, content []byte, appendFlag bool) (int, error) {
	var err error
	var f *os.File
	if appendFlag {
		f, err = os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	} else {
		f, err = os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	}

	if err != nil {
		return 0, err
	}
	defer f.Close()
	writerLenth, err := f.Write(content)
	if err != nil {
		return writerLenth, err
	}
	return writerLenth, nil
}

func (this *FileIO) WriteFileHead(path string, content []byte) error {

	var buf bytes.Buffer
	fileContent, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	buf.Write(content)
	buf.Write(fileContent)
	_, err = this.WriteFile(path, buf.Bytes(), false)
	if err != nil {
		return err
	}
	return nil

}

func (this *FileIO) UpdateFileByLine(path string, content []byte, lineNumber int) error {
	var err error
	var f *os.File
	var buf bytes.Buffer
	byteSlice, err := this.ReadByteLines(path)
	if err != nil {
		return err
	}

	defer f.Close()
	if (0 < lineNumber) && (lineNumber < len(content)-1) {
		byteSlice[lineNumber-1] = content
	} else {
		return errors.New("the number of line is wrong")
	}

	for _, b := range byteSlice {
		buf.Write(b)
		buf.WriteByte('\n')
	}

	f, err = os.OpenFile(path, os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	_, err = f.Write(buf.Bytes())
	if err != nil {
		return err
	}
	return nil

}

func (this *FileIO) FindWithPrefix(content []byte, prefix, end string) string {

	e := bytes.Index(content, []byte("package"))
	contentSub := content[:e]
	s := bytes.Index(contentSub, []byte("// +build"))
	if s <= 0 {
		return ""
	}
	contentSub = content[s:e]
	e = bytes.IndexByte(contentSub, '\n')
	return (string(contentSub[:e]))
}
