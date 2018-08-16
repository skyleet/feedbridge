package store

// Storage defines the interface for a data store backend
type StorageRepository interface {
	// Save stores an item in the store
	Save(string, string) error

	// Get item from store
	Get(string) (string, error)
}