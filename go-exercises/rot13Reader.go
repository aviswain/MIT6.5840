package main

import (
	"io"
	"os"
	"strings"
)

type rot13Reader struct {
	r io.Reader
}

func rot13(b byte) byte {
    switch {
    case b >= 'A' && b <= 'Z':
        if b <= 'M' {
            return b + 13
        }
        return b - 13
    case b >= 'a' && b <= 'z':
        if b <= 'm' {
            return b + 13
        }
        return b - 13
    default:
        return b
    }
}

func (reader *rot13Reader) Read(b []byte) (int, error) {
	n, err := reader.r.Read(b)
	if err != nil {
		return n, err
	}
	
	for i := 0; i < n; i++ {
		b[i] = rot13(b[i])
	}

	return n, err
}

func main() {
	s := strings.NewReader("Lbh penpxrq gur pbqr!")
	r := rot13Reader{s}
	io.Copy(os.Stdout, &r)
}
