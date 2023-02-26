package vsv

import (
	"fmt"
	"io"
	"os"
	"sync"
)

const CHUNK_SIZE = 4096

type File struct {
	mutex sync.Mutex 

	path string
	raw *os.File
	chunks uint64

	currentFrame uint64
	frameChunks uint64
	
	header *Header
}

func New(path string) (*File, error) {
	stat, err := os.Stat(path)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("file '%s' does not exist", path)
	}
	if err != nil {
		return nil, fmt.Errorf("file '%s' error %w", path, err)
	}
	if stat.IsDir() {
		return nil, fmt.Errorf("file '%s' is a directory", path)
	}

	return &File{
		chunks: uint64(stat.Size() / CHUNK_SIZE),
		currentFrame: 0,
		header: nil,
		path: path,
		raw: nil,
	}, nil
}

func (f *File) Open() error {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	raw, err := os.Open(f.path)
	if err != nil {
		return fmt.Errorf("file '%s' open error %w", f.path, err)
	}

	f.raw = raw
	f.currentFrame = 0

	headerData, err := f.readChunks(0, 1)
	if err != nil {
		return err
	}

	header, err := NewHeader(headerData)
	if err != nil {
		return fmt.Errorf("file '%s' error %w", f.path, err)
	}
	f.header = header
	f.frameChunks = f.header.FrameChunks()

	return nil
}

func (f *File) Close() error {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if err := f.Close(); err != nil {
		return fmt.Errorf("file '%s' close error %w", f.path, err)
	}

	f.raw = nil
	f.currentFrame = 0
	f.frameChunks = 0
	f.header = nil

	return nil
}

func (f *File) Path() string {
	return f.path
}

func (f *File) Header() *Header {
	if f.header == nil {
		return nil
	}
	return f.header.Clone()
}

func (f *File) FrameCount() uint64 {
	return (f.chunks - 1) / f.frameChunks
}

func (f *File) Frame(id uint64) (*Frame, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	offset := f.frameChunks * id + 1 // Skip header chunk
	if offset >= f.chunks {
		return nil, io.EOF
	}

	frameData, err := f.readChunks(offset, f.frameChunks)
	if err != nil {
		return nil, err	
	}

	frame, err := NewFrame(frameData, f.header.Bounds(), f.header.ColorMode())
	if err != nil {
		return nil, fmt.Errorf("file '%s' error %w", f.path, err)
	}

	f.currentFrame = id

	return frame, nil
}

func (f *File) seek(offset uint64) error {
	if _, err := f.raw.Seek(int64(offset), 0); err != nil {
		return fmt.Errorf("file '%s' seek to %d error %w", f.path, offset, err)
	}
	return nil
}

func (f *File) read(size uint64) ([]byte, error) {
	data := make([]byte, size)
	if _, err := f.raw.Read(data); err != nil {
		return nil, fmt.Errorf("file '%s' read error %w", f.path, err)
	}
	return data, nil
}

func (f *File) readChunks(offset uint64, count uint64) ([]byte, error) {
	err := f.seek(offset * CHUNK_SIZE)
	if err != nil {
		return nil, err
	}

	return f.read(count * CHUNK_SIZE)
}
