package common

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
)

type Image struct {
	Id        int    `json:"id" gorm:"column:id;"`
	Url       string `json:"url" gorm:"column:url;"`
	Width     int    `json:"width" gorm:"column:width;"`
	Height    int    `json:"height" gorm:"column:height;"`
	CloudName string `json:"cloud_name,omitempty" gorm:"-"`
	Extension string `json:"extension,omitempty" gorm:"-"`
}

func (Image) TableName() string { return "images" }

func (j *Image) Fulfill(domain string) {
	domain = strings.TrimRight(domain, "/")
	url := strings.TrimLeft(j.Url, "/")

	if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
		j.Url = url
		return
	}

	j.Url = fmt.Sprintf("%s/%s", domain, url)
}

// Scan implements sql.Scanner interface
func (j *Image) Scan(value interface{}) error {
	if value == nil {
		*j = Image{}
		return nil
	}

	var data []byte
	switch v := value.(type) {
	case []byte:
		data = v
	case string:
		data = []byte(v)
	default:
		return fmt.Errorf("cannot scan %T into Image", value)
	}

	if len(data) == 0 {
		*j = Image{}
		return nil
	}

	var img Image
	if err := json.Unmarshal(data, &img); err != nil {
		return err
	}
	*j = img
	return nil
}

// Value implements driver.Valuer interface
func (j *Image) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}
