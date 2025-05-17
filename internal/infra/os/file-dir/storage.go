package filedir

import (
	"errors"
	"io"
	"os"
	"time"

	"github.com/gabriel-vasile/mimetype"
)

type FileClient struct {
	sync *fileDirectorySync
	root *os.Root
}

func NewFileClient(dir string) (*FileClient, error) {
	d, err := NewFileDirectorySync()
	if err != nil {
		return nil, err
	}
	root, err := os.OpenRoot(dir)
	if err != nil {
		return nil, err
	}
	return &FileClient{
		sync: d,
		root: root,
	}, nil
}

const GetBufferSize = 4096

func (fcl *FileClient) Get(path string) (*FileReader, error) {
	fcl.sync.RLock(path)
	stat, err := fcl.root.Stat(path)
	fcl.sync.RUnlock(path)

	if errors.Is(err, os.ErrNotExist) {
		return nil, ErrFileNotFound
	} else if err != nil {
		return nil, err
	}
	if stat.IsDir() {
		return nil, ErrFileIsDir
	}

	reqCh := make(chan readRequest)
	resCh := make(chan readResponse)

	go func() {
		defer close(resCh)
		fcl.sync.RLock(path)
		defer fcl.sync.RUnlock(path)
		f, err := fcl.root.Open(path)
		if err != nil {
			resCh <- readResponse{
				err: err,
			}
			return
		}
		defer f.Close()

		buf := make([]byte, GetBufferSize)
		for req := range reqCh {
			if req.offset == -1 {
				n, err := f.Read(buf)
				if err != nil && err != io.EOF {
					resCh <- readResponse{
						err: err,
					}
					continue
				} else if err == io.EOF {
					resCh <- readResponse{
						data: buf[:n],
						n:    n,
						err:  io.EOF,
					}
					continue
				}
				resCh <- readResponse{
					data: buf[:n],
					n:    n,
				}
				continue
			} else if req.length == -1 {
				_, err := f.Seek(req.offset, 0)
				if err != nil {
					resCh <- readResponse{
						err: err,
					}
					continue
				}
				resCh <- readResponse{
					n: 0,
				}
				continue
			}
			n, err := f.ReadAt(buf, req.offset)
			if err != nil && err != io.EOF {
				resCh <- readResponse{
					err: err,
				}
				continue
			} else if err == io.EOF {
				resCh <- readResponse{
					data: buf[:n],
					n:    n,
					err:  io.EOF,
				}
				continue

			}
			resCh <- readResponse{
				data: buf[:n],
				n:    n,
			}
		}
	}()

	fr := NewFileReader(reqCh, resCh)
	return fr, nil
}

func (fcl *FileClient) Put(path string, reader io.Reader) error {
	fcl.sync.Lock(path)
	defer fcl.sync.Unlock(path)

	f, err := fcl.root.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, reader)
	if err != nil {
		return err
	}

	return nil
}

func (fcl *FileClient) Delete(path string) error {
	fcl.sync.RLock(path)
	stat, err := fcl.root.Stat(path)
	fcl.sync.RUnlock(path)

	if errors.Is(err, os.ErrNotExist) {
		return ErrFileNotFound
	} else if err != nil {
		return err
	}

	if stat.IsDir() {
		return ErrFileIsDir
	}

	fcl.sync.Lock(path)
	defer fcl.sync.Unlock(path)

	return fcl.root.Remove(path)
}

type FileInfo struct {
	Name        string
	Size        int64
	Mode        os.FileMode
	ModTime     time.Time
	ContentType string
}

func (fcl *FileClient) Stat(path string) (*FileInfo, error) {
	fcl.sync.RLock(path)
	defer fcl.sync.RUnlock(path)

	stat, err := fcl.root.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil, ErrFileNotFound
	} else if err != nil {
		return nil, err
	}
	if stat.IsDir() {
		return nil, ErrFileIsDir
	}

	cntType, err := fcl.getContentType(path)
	if err != nil {
		return nil, err
	}

	return &FileInfo{
		Name:        stat.Name(),
		Size:        stat.Size(),
		Mode:        stat.Mode(),
		ModTime:     stat.ModTime(),
		ContentType: cntType,
	}, nil
}

func (fcl *FileClient) Mkdir(path string) error {
	return fcl.root.Mkdir(path, 0755)
}

func (fcl *FileClient) getContentType(path string) (string, error) {
	fcl.sync.RLock(path)
	defer fcl.sync.RUnlock(path)

	f, err := fcl.root.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	mimetype, err := mimetype.DetectReader(f)
	if err != nil {
		return "", err
	}

	return mimetype.String(), nil

}
