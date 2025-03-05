package kv_store

/*
This is a kv store backed by bbolt
Considerations:
	- bbolt has an exclusive write lock so only one thread can write at a time
		- Individual transactions and all objects created from them (e.g. buckets, keys) are not thread safe. To work with data in multiple goroutines you must start a transaction for each one or use locking to ensure only one goroutine accesses a transaction at a time. Creating transaction from the DB is thread safe.
	- use contexts to enforce timeout
		- need to read more about the context package
			- do contexts govern the request lifetime?
		- need to understand how contexts interact with middleware
		- need to understand how contexts interact with bbolt transactions
	- do I want to implement the WAL here?
		- bbolt is not write optimized, applying the operations in a WAL using the batch update api is much more performant
		- is that unnecessary complexity?
		- how do I make writing to the WAL thread safe?
			- will there be multiple goroutines trying to access the store concurrently? probably right
			- are request handlers called in their own goroutine?
*/