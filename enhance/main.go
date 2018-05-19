package main

import (
	"bytes"
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"
)

func main() {

	flag.Parse()
	path := *flag.String("path", ".", "the path to enhance")

	stat, err := os.Stat(path)
	if err != nil {
		panic(err)
	}

	if !stat.IsDir() {
		fileHandler(path)
		return
	}

	err = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if err != nil {
			return err
		}
		fileHandler(path)
		return nil
	})
	if err != nil {
		panic(err)
	}
}

func fileHandler(filepath string) {
	src, _ := ioutil.ReadFile(filepath)
	file := ParseFile("", src)
	if len(file.Funs) > 0 {
		codeGenerate(filepath, bytes.Split(src, []byte{'\n'}), file)
	}
}

var newLine = []byte{'\n'}
var indent = []byte{'\t'}

func codeGenerate(filepath string, lines [][]byte, f *File) {
	fd, err := os.OpenFile(filepath, os.O_RDWR|os.O_TRUNC, os.ModeType)
	if err != nil {
		panic(err)
	}
	defer fd.Close()
	pos := f.ImportPos
	at := 1
	for at < pos {
		fd.Write(lines[at-1])
		fd.Write(newLine)
		at++
	}
	if f.ImportKey {
		fd.Write([]byte(`import enhancer "enhance"`))
		fd.Write(newLine)
	} else {
		fd.Write(indent)
		fd.Write([]byte(`enhancer "enhance"`))
		fd.Write(newLine)
	}
	funs := f.Funs
	for i := range funs {
		fun := &funs[i]
		pos = fun.Lbrace
		for at < pos {
			fd.Write(lines[at-1])
			fd.Write(newLine)
			at++
		}
		define := lines[at-1]
		define = rename(fun.Name, define)
		fd.Write(define)
		fd.Write(newLine)
		at++
		pos = fun.Rbrace + 1
		for at < pos {
			fd.Write(lines[at-1])
			fd.Write(newLine)
			at++
		}
		temp.ExecuteTemplate(fd, "tmpl", fun)
	}
	for at <= len(lines) {
		fd.Write(lines[at-1])
		fd.Write(newLine)
		at++
	}

}

func rename(name string, define []byte) []byte {
	old := []byte(name)
	new := make([]byte, 0, len(old)+1)
	new = append(new, '_')
	new = append(new, old...)
	return bytes.Replace(define, old, new, 1)
}
