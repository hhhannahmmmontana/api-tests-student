package advertisementCreating

import (
	"api-tests-template/internal/managers/advertisement"
	"api-tests-template/internal/managers/advertisement/models"
	"api-tests-template/internal/utils"

	"github.com/stretchr/testify/suite"
)

func CreateTestAdvertisement(s *suite.Suite, token string) models.Advertisement {
	var ad models.Advertisement
	s.Run("Создаем объявление", func() {
		req := models.CreateAdvertisementRequest{
			Title:       utils.RandomString(32),
			Description: utils.RandomString(32),
			Price:       utils.RandomFloat32(1, 1000000),
			Quantity:    int32(utils.RandomInt(0, 100)),
			Photos:      []models.PhotoFile{utils.RandomPhoto(utils.RandomString(12))},
		}
		ad = advertisement.CreateAdvertisement(s.T(), token, req)
	})
	return ad
}

func DeleteAdvertisement(s *suite.Suite, token string, id string) {
	s.Run("Удаляем объявление", func() {
		advertisement.DeleteAdvertisement(s.T(), token, id)
	})
}

func CleanupTestAdvertisement(s *suite.Suite, token string, id string) {
	s.T().Cleanup(func() {
		if id != "" {
			DeleteAdvertisement(s, token, id)
		}
	})
}
