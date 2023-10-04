package bbolt

import "go.etcd.io/bbolt"

func Get(db *bbolt.DB, bucket []byte, key string) (s string) {
	_ = db.View(func(tx *bbolt.Tx) error {
		s = string(tx.Bucket(bucket).Get([]byte(key)))
		return nil
	})

	return
}

func Put(db *bbolt.DB, bucket []byte, key, value string) error {
	return db.Update(func(tx *bbolt.Tx) error {
		return tx.Bucket(bucket).Put([]byte(key), []byte(value))
	})
}
