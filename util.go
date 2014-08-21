package gobzip

import (
	"fmt"
	"hash/fnv"
	"path/filepath"

	"github.com/satori/go.uuid"
)

// use uuid to generate file name
func tempFile(dir, prefix string) string {
	fnv1a := fnv.New64a()
	fnv1a.Write(uuid.NewV1().Bytes())
	return filepath.Join(fmt.Sprintf("%s", dir), fmt.Sprintf("%s-%d", prefix, uint64(fnv1a.Sum64())))
}
