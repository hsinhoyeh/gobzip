#gobzip

package gobzip uses cgo to wrap the bzlip to provide bzip compressor in golang.
Because golang native library only provides bzip uncompressor (see <http://golang.org/pkg/compress/bzip2/>, they have a feature request for implementing bzip compressor, but haven't finished it yet. see the issue ticker here: <https://code.google.com/p/go/issues/detail?id=4828>). Thus BzipWriter is introduced.

####bzlib:

1. make sure you have bzlib.h, if not, please build the source code from <http://www.bzip.org>
2. we pushed built-in Darwin/Linux amd64 lib in bzip/lib. You can use them directly if you are one of them.
add them into cgo
NOTE: the part and the platform you used

```
export CGO_CFLAGS=-Ibzip/include
export CGO_LDFLAGS=-Lbzip/Darwin_amd64/lib/
```

#### import


```
go get github.com/hsinhoyeh/gobzip
```

#### example


```
    import (
        "github.com/hsihoyeh/gobzip"
    )

    func main(){

        buf := &bytes.Buffer{}
        w, _ := NewBzipWriter(cBuf)
        w.Write([]byte("I am a plain text message"))
        w.Close()
        fmt.Println(buf.Bytes())
    }
```



