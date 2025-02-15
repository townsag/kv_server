package kv_store

import (
	"testing"
)

func TestSetAndGetHappyPath(t *testing.T) {
	var store *MemoryStore = NewMemoryStore()
	err := store.Set("asdf", "asdf")
	if err != nil {
		t.Errorf("Set operation threw error: %s", err.Error())
	}
	value, err := store.Get("asdf")
	if err != nil {
		t.Errorf("Get operation threw error: %s", err.Error())
	}
	if value != "asdf" {
		t.Errorf("got %v, wanted 'asdf'", value)
	}
}

func TestGetMissingKey(t *testing.T) {
	var store *MemoryStore = NewMemoryStore()
	_, err := store.Get("asdf")
	if err == nil {
		t.Fatalf("called get on a nonexistent value and did not get an error")
	}
}

func TestDeleteHappyPath(t *testing.T) {
	var store *MemoryStore = NewMemoryStore()
	err := store.Set("asdf", "asdf")
	if err != nil {
		t.Errorf("failed to set a value with error: %s", err.Error())
	}
	
	err = store.Delete("asdf")
	if err != nil {
		t.Errorf("failed to delete key 'asdf' with message %s", err.Error())
	}

	_, err = store.Get("asdf")
	if err == nil {
		t.Fatalf("was able to retrieve a deleted key without an error")
	}
}

func TestDeleteMissingKey(t *testing.T) {
	var store *MemoryStore = NewMemoryStore()
	err := store.Delete("asdf")
	if err == nil {
		t.Fatal("did not get an error after deleting a missing key")
	}
}