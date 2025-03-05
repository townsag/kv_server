package main

import (
	bolt "go.etcd.io/bbolt"
	"log"
	"fmt"
)

func main() {
	var file_path string = "temp.db"

	db, err := bolt.Open(file_path, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// create a bucket to write key value pairs to
	// use a read - write transaction to do this
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("test-bucket"))
		if err != nil {
			return fmt.Errorf("error creating bucket: %s", err)
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}


	// read the created buckets
	err = db.View(func(tx *bolt.Tx) error {
		tx.ForEach(func(name []byte, b *bolt.Bucket) error {
			fmt.Printf("this is my bucket: %s\n", string(name))
			return nil
		})
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	// write a value to my bucket
	err = db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("test-bucket"))
		err := bucket.Put([]byte("asdf"), []byte("qwer"))
		return err
	})
	if err != nil {
		log.Fatal(err)
	}

	// read the value written to the bucket
	err = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("test-bucket"))
		value := bucket.Get([]byte("asdf"))
		if value == nil {
			fmt.Printf("unable to find a value for key: %s\n", "asdf")
			return nil
		}
		fmt.Printf("the value associated with key: %s is %s\n", "asdf", string(value))
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	// read write transaction
	err = db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("test-bucket"))
		// read the value in the database for "asdf"
		value := bucket.Get([]byte("asdf"))
		if value == nil {
			return fmt.Errorf("did not find a value associated with key %s", "asdf")
		}
		// update that value
		err := bucket.Put([]byte("asdf"), append(value, []byte("qwer")...))
		return err
	})
	if err != nil {
		log.Fatal(err)
	}
}