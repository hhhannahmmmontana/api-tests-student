package models

import (
	"api-tests-template/internal/managers/advertisement/models"
	"time"
)

type Favorite struct {
	AddedAt       time.Time            `json:"added_at"`
	Advertisement models.Advertisement `json:"createAdvertisement"`
	FavouriteId   string               `json:"favourite_id"`
}

type GetFavoritesResponse struct {
	Items      []Favorite `json:"items"`
	Page       int32      `json:"page"`
	PageSize   int32      `json:"page_size"`
	Total      int32      `json:"total"`
	TotalPages int32      `json:"total_pages"`
}
