package pkg

type Cache interface {
	// Get tries to fetch the value for the given key
	// It will always return second argument as false, in cases where first
	// value is expired or not found
	Get(key string) (string, bool)
	// Set writes the value for the given key
	// Expiration duration is hard coded as 30 minutes
	// If a value for the specified key already exists, then it will be overwritten
	Set(key, value string)
}
