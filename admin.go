package client

import (
	"bytes"
)

// Lock contains the metadata for a SF system lock
type LockMetadata struct {
	Node      string `json:"node"`
	Operation string `json:"operation"`
	PID       int    `json:"pid"`
}

// Locks is a map of lock references to lock metadata
type Locks map[string]LockMetadata

// GetLocks retrieves a list of SF system locks
func (c *Client) GetLocks() (Locks, error) {
	locks := Locks{}
	err := c.doRequestJSON("admin/locks", "GET", bytes.Buffer{}, &locks)

	return locks, err
}
