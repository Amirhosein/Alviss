package db

import (
	"encoding/json"
	"time"
)

type UrlMapping struct {
	OriginalUrl string    `json:"OriginalUrl"`
	Count       int       `json:"count"`
	ExpTime     time.Time `json:"ExpTime"`
}

func (s UrlMapping) MarshalBinary() ([]byte, error) {
	return json.Marshal(s)
}

// UnmarshalBinary use msgpack
func (s *UrlMapping) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, s)
}
