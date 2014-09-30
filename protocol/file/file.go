package file

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/mjibson/mog/codec"
	"github.com/mjibson/mog/protocol"
)

func init() {
	protocol.Register("file", []string{"directory"}, List)
}

func List(params []string) (protocol.SongList, error) {
	if len(params) != 1 {
		return nil, fmt.Errorf("bad params")
	}
	songs := make(protocol.SongList)
	err := filepath.Walk(params[0], func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		f, err := os.Open(path)
		if err != nil {
			return nil
		}
		defer f.Close()
		ss, _, err := codec.Decode(fileReader(path))
		if err != nil {
			return nil
		}
		for i, s := range ss {
			id := fmt.Sprintf("%v-%v", i, path)
			songs[id] = s
		}
		return nil
	})
	return songs, err
}

func fileReader(path string) codec.Reader {
	return func() (io.ReadCloser, error) {
		return os.Open(path)
	}
}
