package db

import (
	"encoding/csv"
	"os"
)

type CloseFunc func() error

func newDBReader(file string) (*csv.Reader, CloseFunc) {
	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}

	r := csv.NewReader(f)
	r.Comma = ';'
	return r, f.Close
}

func newDBWriter(file string) (*csv.Writer, CloseFunc) {
	f, err := os.Create(file)
	if err != nil {
		panic(err)
	}

	w := csv.NewWriter(f)
	w.Comma = ';'
	return w, f.Close
}


