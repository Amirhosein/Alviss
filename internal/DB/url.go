package db

import (
	"encoding/json"
	"time"
)

type UrlMapping struct {
	Original_url string    `json:"original_url"`
	Count        int       `json:"count"`
	ExpTime      time.Time `json:"exp_time"`
}

func (s UrlMapping) MarshalBinary() ([]byte, error) {
	return json.Marshal(s)
}

// UnmarshalBinary use msgpack
func (s *UrlMapping) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, s)
}
