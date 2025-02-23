package testserver

import "time"

type (
	Namespace string
	Table     map[uint]Model
	State     map[Namespace]Table
)

const (
	UserNamespace     Namespace = "users"
	ResourceNamespace Namespace = "resources"
)

func NewState() *State {
	return &State{
		UserNamespace: {
			1: &User{ID: 1, Name: "Alice Johnson", Birthdate: time.Date(1995, time.June, 15, 0, 0, 0, 0, time.UTC)},
			2: &User{ID: 2, Name: "Bob Williams", Birthdate: time.Date(1990, time.December, 5, 0, 0, 0, 0, time.UTC)},
			3: &User{ID: 3, Name: "Charlie Davis", Birthdate: time.Date(1988, time.March, 22, 0, 0, 0, 0, time.UTC)},
			4: &User{ID: 4, Name: "Diana Evans", Birthdate: time.Date(1999, time.April, 10, 0, 0, 0, 0, time.UTC)},
			5: &User{ID: 5, Name: "Ethan Moore", Birthdate: time.Date(1985, time.September, 30, 0, 0, 0, 0, time.UTC)},
		},
		ResourceNamespace: {
			1: &Resource{ID: 1, Type: "Image", Data: "image_123.png", Timestamp: time.Now().Add(-48 * time.Hour)},
			2: &Resource{ID: 2, Type: "Video", Data: "video_456.mp4", Timestamp: time.Now().Add(-24 * time.Hour)},
			3: &Resource{ID: 3, Type: "Document", Data: "doc_789.pdf", Timestamp: time.Now().Add(-72 * time.Hour)},
			4: &Resource{ID: 4, Type: "Audio", Data: "audio_321.mp3", Timestamp: time.Now().Add(-6 * time.Hour)},
			5: &Resource{ID: 5, Type: "Text", Data: "note_654.txt", Timestamp: time.Now().Add(-12 * time.Hour)},
		},
	}
}
