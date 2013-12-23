package io

import (
	"io"
	"os"
)

type OnProgressFunc func(file_size, uploaded int64)

type UpFile struct {
	uploaded   int64
	file       *os.File
	tag        bool
	fsize      int64
	onProgress func(file_size, uploaded int64)
}

func OpenUpFile(name string, onProgress OnProgressFunc) (pfile *UpFile, err error) {
	f, err := os.Open(name)
	pfile = new(UpFile)
	pfile.file = f
	fi, err := pfile.file.Stat()
	if err != nil {
		return
	}
	pfile.onProgress = onProgress
	pfile.fsize = fi.Size()
	return
}

func (p *UpFile) Size() int64 {
	return p.fsize
}

func (pfile *UpFile) ReadAt(p []byte, off int64) (n int, err error) {
	n, err = pfile.file.ReadAt(p, off)
	if err == io.EOF {
		return
	} else if err != nil {
		return
	}
	if !pfile.tag {
		pfile.tag = true
		return
	}
	go pfile.onProgress(pfile.fsize, pfile.uploaded)
	pfile.uploaded += int64(n)
	return
}

func (pfile *UpFile) Close() {
	pfile.file.Close()
}

func (pfile *UpFile) Read(b []byte) (n int, err error) {

	n, err = pfile.file.Read(b)
	if err == io.EOF {
		//almost finished
		return
	} else if err != nil {
		return
	}

	if !pfile.tag {
		pfile.tag = true
		return
	}
	go pfile.onProgress(pfile.fsize, pfile.uploaded)
	pfile.uploaded += int64(n)
	return
}
