package testserver

import (
	"fmt"
	"time"
)

var (
	_ Model = (*User)(nil)
	_ Model = (*Resource)(nil)
)

type (
	Model interface {
		SetID(id uint)
		GetID() uint
		String() string
	}

	User struct {
		ID        uint      `json:"id"`
		Name      string    `json:"name"`
		Birthdate time.Time `json:"birthdate"`
	}

	Resource struct {
		ID        uint      `json:"id"`
		Type      string    `json:"type"`
		Data      string    `json:"data"`
		Timestamp time.Time `json:"timestamp"`
	}
)

// GetID implements Model.
func (u *User) GetID() uint {
	return u.ID
}

// SetID implements Model.
func (u *User) SetID(id uint) {
	u.ID = id
}

// String implements Model.
func (u *User) String() string {
	return fmt.Sprintf(
		"User{ID:%d Name:%s Birthdate:%s}",
		u.ID,
		u.Name,
		u.Birthdate.Format(time.RFC3339),
	)
}

// GetID implements Model.
func (r *Resource) GetID() uint {
	return r.ID
}

// SetID implements Model.
func (r *Resource) SetID(id uint) {
	r.ID = id
}

// String implements Model.
func (r *Resource) String() string {
	return fmt.Sprintf(
		"Resource{ID:%d Type:%s Data:%s Timestamp:%s}",
		r.ID,
		r.Type,
		r.Data,
		r.Timestamp.Format(time.RFC3339),
	)
}
