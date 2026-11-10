package models

import (
	"api-tests-template/internal/managers/auth/models"
	"time"
)

type AdvertisementsList struct {
	Items []Item `json:"items,omitempty"`
}

type Item struct {
	Id          string      `json:"id,omitempty"`
	Title       string      `json:"title,omitempty"`
	Description string      `json:"description,omitempty"`
	Price       int32       `json:"price,omitempty"`
	Quantity    int32       `json:"quantity,omitempty"`
	UserId      string      `json:"userId,omitempty"`
	User        models.User `json:"user,omitempty"`
	CreatedAt   time.Time   `json:"created_at,omitempty"`
	UpdatedAt   time.Time   `json:"updated_at,omitempty"`
	Photos      []Photos    `json:"photos,omitempty"`
}

type Photos struct {
	Id              string    `json:"id,omitempty"`
	AdvertisementId string    `json:"advertisementId,omitempty"`
	Url             string    `json:"url,omitempty"`
	SortOrder       int32     `json:"sort_order,omitempty"`
	CreatedAt       time.Time `json:"created_at,omitempty"`
}
