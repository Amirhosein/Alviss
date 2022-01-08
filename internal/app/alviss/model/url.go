package model

import (
	"database/sql"
	"encoding/json"
	"log"
	"time"
)

type URLMapping struct {
	Original_url string    `json:"original_url"`
	Count        int       `json:"count"`
	ExpTime      time.Time `json:"exp_time"`
}

func (s URLMapping) MarshalBinary() ([]byte, error) {
	return json.Marshal(s)
}

func (s *URLMapping) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, s)
}

type URLRepo interface {
	SaveURLMapping(shortURL string, urlMapping URLMapping, expTime time.Duration) error
	GetURLMapping(shortURL string) (URLMapping, error)
	UpdateURLMapping(shortURL string, urlMapping URLMapping) error
}

type SQLURLRepo struct {
	DB *sql.DB
}

func (s SQLURLRepo) SaveURLMapping(shortURL string, urlMapping URLMapping, expTime time.Duration) error {
	if expTime == 0 {
		_, err := s.DB.Exec("INSERT INTO url_mapping(short_url, original_url, count, exp_time) VALUES($1, $2, $3, $4)", shortURL, urlMapping.Original_url, urlMapping.Count, sql.NullTime{})
		return err
	} else {
		_, err := s.DB.Exec("INSERT INTO url_mapping(short_url, original_url, count, exp_time) VALUES($1, $2, $3, $4)", shortURL, urlMapping.Original_url, urlMapping.Count, time.Now().Add(expTime))
		return err
	}
}

func (s SQLURLRepo) GetURLMapping(shortURL string) (URLMapping, error) {
	var urlMapping URLMapping
	err := s.DB.QueryRow("SELECT original_url, count FROM url_mapping WHERE short_url=$1", shortURL).Scan(&urlMapping.Original_url, &urlMapping.Count)
	if err != nil {
		log.Println(err)
	}
	err = s.DB.QueryRow("SELECT exp_time FROM url_mapping WHERE short_url=$1", shortURL).Scan(&urlMapping.ExpTime)
	return urlMapping, err
}

func (s SQLURLRepo) UpdateURLMapping(shortURL string, urlMapping URLMapping) error {
	log.Println(urlMapping)
	_, err := s.DB.Exec("UPDATE url_mapping SET count=$1, exp_time=$2 WHERE short_url=$3", urlMapping.Count, urlMapping.ExpTime, shortURL)
	return err
}
