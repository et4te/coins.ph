package main

import (
	"errors"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
)

//------------------------------------------------------------------------------
// Persistence
//------------------------------------------------------------------------------

type BoltState struct {
	*bolt.DB
}

func OpenDB(path string) *BoltState {
	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		log.Fatalf("Error creating database: %s", err)
	}
	return &BoltState{db}
}

// Ensures a given bucket exists in the database.
func (db *BoltState) SyncBucket(bucket []byte) error {
	return db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(bucket)
		if err != nil {
			return err
		}
		return nil
	})
}

// Get a record from a bucket.
func (db *BoltState) GetRecord(bucket []byte, k []byte) (*Record, error) {
	fmt.Println(fmt.Sprintf("(%s) lookup up %s", string(bucket), string(k)))
	if len(k) <= 0 {
		return nil, errors.New("Key was 0 bytes.")
	}

	var r []byte
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		if x := b.Get(k); x != nil {
			r = make([]byte, len(x))
			copy(r, x)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	fmt.Println(fmt.Sprintf("(%s) found %s", string(bucket), string(r)))

	if r != nil {
		return &Record{Key: k, Value: r}, nil
	} else {
		return nil, nil
	}
}

// Fetches all records within a bucket.
func (db *BoltState) GetRecords(bucket []byte) ([]*Record, error) {
	lim, err := db.CountBucket(bucket)
	if err != nil {
		return nil, err
	}

	keys, vals, err := db.getKVs(bucket, lim)
	if err != nil {
		return nil, err
	}

	var records []*Record
	for i := 0; i < len(keys); i++ {
		records = append([]*Record{&Record{Key: keys[i], Value: vals[i]}}, records...)
	}

	return records, nil
}

// Puts a record in a bucket.
func (db *BoltState) PutRecord(bucket []byte, record *Record) error {
	return db.putKV(bucket, record.Key, record.Value)
}

func (db *BoltState) CountBucket(bucket []byte) (uint64, error) {
	fmt.Println(fmt.Sprintf("(%s) counting ...", string(bucket)))
	count := make([]uint64, 1)
	err := db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(bucket)
		if err != nil {
			return err
		}
		var i uint64 = 0
		c := b.Cursor()
		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			i += 1
			copy(count, []uint64{i})
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	fmt.Println(fmt.Sprintf("(%s) found %v", bucket, count[0]))
	return count[0], nil
}

//------------------------------------------------------------------------------
// Internal
//------------------------------------------------------------------------------

// Fetch key / value pairs within a bucket up to lim.
func (db *BoltState) getKVs(bucket []byte, lim uint64) ([][]byte, [][]byte, error) {
	fmt.Println(fmt.Sprintf("(%s) fetching kvs where count = %v", string(bucket), lim))
	keys := make([]([]byte), lim)
	vals := make([]([]byte), lim)
	err := db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(bucket)
		if err != nil {
			return err
		}
		var i uint64 = 0
		b.ForEach(func(k, v []byte) error {
			keys[i] = make([]byte, len(k))
			vals[i] = make([]byte, len(v))
			copy(keys[i], k)
			copy(vals[i], v)
			i = i + 1
			return nil
		})
		return nil
	})
	if err != nil {
		return nil, nil, err
	}
	return keys, vals, nil
}

// Put a key / value pair within bucket.
func (db *BoltState) putKV(bucket []byte, k []byte, v []byte) error {
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		if err := b.Put(k, v); err != nil {
			return err
		}
		return nil
	})
}
