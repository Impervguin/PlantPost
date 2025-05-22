package filedir

import (
	"errors"
	"io"
	"os"
)

type readRequest struct {
	offset int64
	length int64
}

type readResponse struct {
	data []byte
	n    int
	err  error
}

type FileReader struct {
	reqCh chan<- readRequest
	resCh <-chan readResponse
}

func NewFileReader(reqCh chan<- readRequest, resCh <-chan readResponse) *FileReader {
	return &FileReader{
		reqCh: reqCh,
		resCh: resCh,
	}
}

func (fr *FileReader) Read(p []byte) (n int, err error) {
	req := readRequest{
		offset: -1,
		length: int64(len(p)),
	}
	select {
	case fr.reqCh <- req:
	case <-fr.resCh:
		return 0, os.ErrClosed
	}

	response := <-fr.resCh
	if errors.Is(response.err, EOF) {
		if response.n != 0 {
			copy(p, response.data)
		}
		return response.n, io.EOF
	}
	if response.err != nil {
		return 0, response.err
	}
	return copy(p, response.data), nil
}

func (fr *FileReader) Close() error {
	close(fr.reqCh)
	return nil
}

func (fr *FileReader) Seek(offset int64, whence int) (int64, error) {
	req := readRequest{
		offset: offset,
		length: -1,
	}
	select {
	case fr.reqCh <- req:
	case <-fr.resCh:
		return 0, os.ErrClosed
	}
	response := <-fr.resCh
	if response.err != nil {
		return 0, response.err
	}
	return offset, nil
}

func (fr *FileReader) ReadAt(p []byte, offset int64) (n int, err error) {
	req := readRequest{
		offset: offset,
		length: int64(len(p)),
	}
	select {
	case fr.reqCh <- req:
	case <-fr.resCh:
		return 0, os.ErrClosed
	}
	response := <-fr.resCh
	if response.err != nil {
		return 0, response.err
	}
	if response.n == 0 {
		return 0, EOF
	}
	return copy(p, response.data), nil
}
