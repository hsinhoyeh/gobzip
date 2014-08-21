package gobzip

import (
	"bytes"
	"compress/bzip2"
	"io/ioutil"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

var message = "At Castle Black, news arrives of Eddard's arrest and imprisonment." +
	"Jon Snow is unable to do anything about it, to his frustration." +
	"Alliser Thorne taunts Jon that his father is a traitor." +
	"Jon threatens him with a knife and is confined to quarters for his trouble. Meanwhile, the bodies of several men from Benjen Stark's patrol have been found," +
	"but there is no sign of Benjen himself." +
	"Samwell Tarly notes that the bodies do not smell like they've been dead for weeks and there is no evidence of decomposition." +
	"Jon and several other Sworn Brothers urge Lord Commander Jeor Mormont to burn the bodies, but he refuses, wanting Maester Aemon to examine them."

func TestBzip(t *testing.T) {
	cBuf := &bytes.Buffer{}
	w, err := NewBzipWriter(cBuf)
	if err != nil {
		t.Fatal("should not be nil, err :%s", err)
	}
	assert.NoError(t, err, "no error")
	n, err := w.Write([]byte(message))
	assert.NoError(t, err, "no error")
	assert.Equal(t, len(message), int(n), "should be equal")
	err = w.Close()
	assert.NoError(t, err, "no error")

	r := bzip2.NewReader(cBuf)
	b, err := ioutil.ReadAll(r)
	assert.NoError(t, err, "no error")
	assert.True(t, bytes.Compare([]byte(message), b) == 0, "exp: %s, got:%s", message, string(b))
}

func TestMultiThreadingBzip(t *testing.T) {
	numRoutine := 100
	numTest := 100
	errChan := make(chan error, 1)
	wg := &sync.WaitGroup{}
	go func() {
		for {
			err := <-errChan
			t.Errorf("fail: %s", err)
		}
	}()

	for j := 0; j < numRoutine; j++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < numTest; i++ {
				err := func() error {
					cBuf := &bytes.Buffer{}
					w, err := NewBzipWriter(cBuf)
					if err != nil {
						return err
					}
					_, err = w.Write([]byte(message))
					if err != nil {
						return err
					}
					err = w.Close()
					if err != nil {
						return err
					}
					return nil
				}()
				if err != nil {
					errChan <- err
				}
			}
		}()
	}
	wg.Wait()
}
