package main

import (
	"bytes"
	"golang.org/x/crypto/blake2b"
)

//------------------------------------------------------------------------------
// Persistable Records
//------------------------------------------------------------------------------

type Persistent interface {
	AsRecord() (*Record, error)
}

// A record is a ([]byte hash, []byte value) pair.
type Record struct {
	Key   []byte
	Value []byte
}

// Verifies the integrity of a stored record.
func (r *Record) Verify(k []byte) bool {
	h, err := HashBlake2b(k, r.Value)
	if err != nil {
		return false
	}

	if bytes.Compare(h, r.Key) == 0 {
		return true
	} else {
		return false
	}
}

// Takes the Blake2b hash of a record given an initial key 'k' which
// can be anything (e.g: 'account').
func HashBlake2b(k []byte, b []byte) ([]byte, error) {
	h, err := blake2b.New256(k)
	if err != nil {
		return nil, err
	}
	h.Write(b)
	return h.Sum(nil), nil
}
