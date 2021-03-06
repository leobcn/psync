package lib

import (
	"os"
	"io"
	"bufio"
	"path/filepath"
)

func ParseHashList(r io.Reader, fn func(Checksum)) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			break
		}
		fn(Checksum(line))
	}
}

type HashList []Checksum

func NewHashList(r io.Reader) *HashList {
	h := HashList{}
	ParseHashList(r, func(s Checksum) {
		h = append(h, s)
	})
	return &h
}

func (h *HashList) Resolve(fn func (io.Reader, error) error) {
	root := BlocksDir()
	for _, checksum := range *h {
		f, err := os.Open(filepath.Join(root, string(checksum)))
		if err != nil {
			fn(nil, err)
			break
		}
		err = fn(f, nil)
		f.Close()
		if err != nil {
			break
		}
	}
}

func (h *HashList) Missing() []Checksum {
	xs := []Checksum{}
	root := BlocksDir()
	for _, checksum := range *h {
		_, err := os.Stat(filepath.Join(root, string(checksum)))
		if err != nil {
			xs = append(xs, checksum)
		}
	}
	return xs
}
