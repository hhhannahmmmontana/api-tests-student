package models

import (
	"api-tests-template/internal/managers/auth/models"
	"time"
)

type AdvertisementsList struct {
	Items []Advertisement `json:"items,omitempty"`
}

type Advertisement struct {
	Id          string      `json:"id,omitempty"`
	Title       string      `json:"title,omitempty"`
	Description string      `json:"description,omitempty"`
	Price       float64     `json:"price,omitempty"`
	Quantity    int32       `json:"quantity,omitempty"`
	UserId      string      `json:"user_id,omitempty"`
	User        models.User `json:"user,omitempty"`
	CreatedAt   time.Time   `json:"created_at,omitempty"`
	UpdatedAt   time.Time   `json:"updated_at,omitempty"`
	Photos      []Photo     `json:"photos,omitempty"`
}

type Photo struct {
	Id              string    `json:"id,omitempty"`
	AdvertisementId string    `json:"advertisement_id,omitempty"`
	Url             string    `json:"url,omitempty"`
	SortOrder       int32     `json:"sort_order,omitempty"`
	CreatedAt       time.Time `json:"created_at,omitempty"`
}
