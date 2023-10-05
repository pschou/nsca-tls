package main

import (
	"errors"
	"io"
)

func sliceLine(out []byte, in io.Reader) (slice []byte, err error) {
	for i := 0; i < len(out); i++ {
		_, err = in.Read(out[i : i+1])
		//fmt.Println(string(out[i:i+1]), out[i])
		if out[i] == '\n' {
			return out[:i+1], err
		}
		if err != nil {
			return
		}
	}
	return nil, errors.New("command too long")
}
