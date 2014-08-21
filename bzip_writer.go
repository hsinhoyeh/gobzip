package gobzip

// #include <stdio.h>
// #include <stdlib.h>
// #include <bzlib.h>
// #cgo LDFLAGS: -lbz2
import "C"

import (
	"errors"
	"fmt"
	"io"
	"os"
	"unsafe"
)

// this is a go binding to bzip2 compressor

const (
	defaultDir        = "/tmp"
	defaultPrefix     = "gobzip"
	defaultBlockSize  = 9 // from 1 to 9
	defaultVerbosity  = 0
	defaultWorkFactor = 0

	defaultBufferLen = 1024
)

type BzipWriter struct {
	fd   *C.FILE
	bzfd *C.BZFILE
	// err is the cgo operation error
	err int
	// tmpfile is a temp storage of the bzip stream
	tmpfile string
	w       io.Writer
}

// NewBzipWriter returns a BzipWriter which compresses byte data to w
func NewBzipWriter(w io.Writer) (*BzipWriter, error) {
	b := &BzipWriter{w: w}
	if err := b.bz2_bzWriteOpen(defaultBlockSize, defaultVerbosity, defaultWorkFactor); err != nil {
		return nil, err
	}
	return b, nil
}

// Write writes the byte data to the bzip writer
func (b *BzipWriter) Write(d []byte) (int, error) {
	b.bz2_bzWrite(d)
	if err := BzipError(b.err); err != nil {
		return 0, err
	}
	return len(d), nil
}

// intercept intercepts the underlying temp buffer to the w
func (b *BzipWriter) intercept() error {
	// seek the tmpfile to the beginning
	// and start to intercept it to the w (io.Writer)
	C.fseek(b.fd, C.long(0), C.int(C.SEEK_SET))
	buffer := make([]byte, defaultBufferLen)
	n := C.fread(unsafe.Pointer(&buffer[0]), C.size_t(1), C.size_t(defaultBufferLen), b.fd)
	for {
		_, err := b.w.Write(buffer[:n])
		if err != nil {
			return err
		}
		if n != C.size_t(defaultBufferLen) {
			// exit because we have reached eof
			return nil
		}
		n = C.fread(unsafe.Pointer(&buffer[0]), C.size_t(1), C.size_t(defaultBufferLen), b.fd)
	}
	return nil
}

// Close closes the bzip writer and flushes the data to the w
func (b *BzipWriter) Close() error {
	defer os.Remove(b.tmpfile)
	b.bz2_bzWriteClose()
	if err := BzipError(b.err); err != nil {
		return err
	}
	// intercept the result to w io.Writer
	if err := b.intercept(); err != nil {
		return err
	}
	if err := C.fclose(b.fd); err != C.int(0) {
		return fmt.Errorf("close file returns non-zero, file: %s", b.tmpfile)
	}
	return nil
}

// bz2_bzWriteOpen wraps C.BZ2_bzWriteOpen
func (b *BzipWriter) bz2_bzWriteOpen(blockSize int, verbosity int, workFactor int) error {
	// get a temp file for storing the bzip stream
	b.tmpfile = tempFile(defaultDir, defaultPrefix)

	filename := C.CString(b.tmpfile)
	defer C.free(unsafe.Pointer(filename))
	cmode := C.CString("w+")
	defer C.free(unsafe.Pointer(cmode))
	b.fd = C.fopen(filename, cmode)
	b.bzfd = (*C.BZFILE)(unsafe.Pointer(C.BZ2_bzWriteOpen((*C.int)(unsafe.Pointer(&b.err)), b.fd, C.int(blockSize), C.int(verbosity), C.int(workFactor))))
	if err := BzipError(b.err); err != nil {
		return fmt.Errorf("fd: %v, tmpfile:%s, errcode: %d, err: %s", b.fd, b.tmpfile, b.err, err)
	}
	return nil
}

// bz2_bzWrite wraps C.bz2_bzWrite
func (b *BzipWriter) bz2_bzWrite(buf []byte) {
	if buf == nil || len(buf) == 0 {
		return
	}
	C.BZ2_bzWrite((*C.int)(unsafe.Pointer(&b.err)), unsafe.Pointer(b.bzfd), unsafe.Pointer(&buf[0]), C.int(len(buf)))
}

// bz2_bzWriteClose wraps C.BZ2_bzWriteClose
func (b *BzipWriter) bz2_bzWriteClose() (byteIn, byteOut int) {
	abandon := 0
	C.BZ2_bzWriteClose((*C.int)(unsafe.Pointer(&b.err)), unsafe.Pointer(b.bzfd), C.int(abandon), (*C.uint)(unsafe.Pointer(&byteIn)), (*C.uint)(unsafe.Pointer(&byteOut)))
	return
}

var (
	ErrStream     = errors.New("bzip: stream")
	ErrConfig     = errors.New("bzip: config")
	ErrSequence   = errors.New("bzip: sequence")
	ErrParam      = errors.New("bzip: parameter")
	ErrMem        = errors.New("bzip: memory")
	ErrData       = errors.New("bzip: data")
	ErrDataMagic  = errors.New("bzip: data magic")
	ErrIO         = errors.New("bzip: i/o")
	ErrUnexpected = errors.New("bzip: unexpected")
	ErrOutbufFull = errors.New("bzip: output buffer full")
	ErrUnknown    = errors.New("bzip: unknown error")
)

// BzipError converts the err codes into golang's error
func BzipError(err int) error {
	switch C.int(err) {
	case C.BZ_OK:
		return nil
	case C.BZ_RUN_OK:
		return nil
	case C.BZ_FLUSH_OK:
		return nil
	case C.BZ_FINISH_OK:
		return nil
	case C.BZ_STREAM_END:
		return ErrStream
	case C.BZ_CONFIG_ERROR:
		return ErrConfig
	case C.BZ_SEQUENCE_ERROR:
		return ErrSequence
	case C.BZ_PARAM_ERROR:
		return ErrParam
	case C.BZ_MEM_ERROR:
		return ErrMem
	case C.BZ_DATA_ERROR:
		return ErrData
	case C.BZ_DATA_ERROR_MAGIC:
		return ErrDataMagic
	case C.BZ_IO_ERROR:
		return ErrIO
	case C.BZ_UNEXPECTED_EOF:
		return ErrUnexpected
	case C.BZ_OUTBUFF_FULL:
		return ErrOutbufFull
	}
	return ErrUnknown
}
