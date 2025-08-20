package domain

import (
	"time"
)

type Role struct {
	ID          uint64
	Name        string
	DisplayName string
	Description string
	Permissions []string
	IsSystem    bool
	CreatedAt   time.Time
}