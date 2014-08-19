package gobzip

import (
	"bytes"
	"compress/bzip2"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBzip(t *testing.T) {
	message := "At Castle Black, news arrives of Eddard's arrest and imprisonment." +
		"Jon Snow is unable to do anything about it, to his frustration." +
		"Alliser Thorne taunts Jon that his father is a traitor." +
		"Jon threatens him with a knife and is confined to quarters for his trouble. Meanwhile, the bodies of several men from Benjen Stark's patrol have been found," +
		"but there is no sign of Benjen himself." +
		"Samwell Tarly notes that the bodies do not smell like they've been dead for weeks and there is no evidence of decomposition." +
		"Jon and several other Sworn Brothers urge Lord Commander Jeor Mormont to burn the bodies, but he refuses, wanting Maester Aemon to examine them."

	cBuf := &bytes.Buffer{}
	w, err := NewBzipWriter(cBuf)
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
