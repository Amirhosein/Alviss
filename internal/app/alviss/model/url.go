package model

import (
	"database/sql"
	"encoding/json"
	"log"
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

func (s *UrlMapping) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, s)
}

type UrlRepo interface {
	SaveUrlMapping(shortUrl string, urlMapping UrlMapping, expTime time.Duration) error
	GetUrlMapping(shortUrl string) (UrlMapping, error)
	UpdateUrlMapping(shortUrl string, urlMapping UrlMapping) error
}

type SQLUrlRepo struct {
	DB *sql.DB
}

func (s SQLUrlRepo) SaveUrlMapping(shortUrl string, urlMapping UrlMapping, expTime time.Duration) error {
	if expTime == 0 {
		_, err := s.DB.Exec("INSERT INTO url_mapping(short_url, original_url, count, exp_time) VALUES($1, $2, $3, $4)", shortUrl, urlMapping.Original_url, urlMapping.Count, sql.NullTime{})
		return err
	} else {
		_, err := s.DB.Exec("INSERT INTO url_mapping(short_url, original_url, count, exp_time) VALUES($1, $2, $3, $4)", shortUrl, urlMapping.Original_url, urlMapping.Count, time.Now().Add(expTime))
		return err
	}
}

func (s SQLUrlRepo) GetUrlMapping(shortUrl string) (UrlMapping, error) {
	var urlMapping UrlMapping
	err := s.DB.QueryRow("SELECT original_url, count FROM url_mapping WHERE short_url=$1", shortUrl).Scan(&urlMapping.Original_url, &urlMapping.Count)
	if err != nil {
		log.Println(err)
	}
	err = s.DB.QueryRow("SELECT exp_time FROM url_mapping WHERE short_url=$1", shortUrl).Scan(&urlMapping.ExpTime)
	return urlMapping, err
}

func (s SQLUrlRepo) UpdateUrlMapping(shortUrl string, urlMapping UrlMapping) error {
	log.Println(urlMapping)
	_, err := s.DB.Exec("UPDATE url_mapping SET count=$1, exp_time=$2 WHERE short_url=$3", urlMapping.Count, urlMapping.ExpTime, shortUrl)
	return err
}
