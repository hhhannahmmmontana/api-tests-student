package createAdvertisement

import (
	httpAdvertisement "api-tests-template/internal/client/http/advertisement"
	"api-tests-template/internal/managers/advertisement"
	"api-tests-template/internal/managers/advertisement/models"
	"api-tests-template/internal/managers/auth"
	authModels "api-tests-template/internal/managers/auth/models"
	"api-tests-template/internal/utils"
	base "api-tests-template/tests"
	"net/http"
	"os"
	"slices"
	"testing"

	"github.com/stretchr/testify/suite"
)

const (
	titleLengthLimit       = 50
	descriptionLengthLimit = 500
	priceLimit             = 1000000
	quantityLimit          = 100
	photosLimit            = 3
)

type TestSuite struct {
	suite.Suite
	loginData authModels.LoginOkResponse
}

func TestSuiteRun(t *testing.T) {
	suite.Run(t, &TestSuite{})
}

func (s *TestSuite) SetupSuite() {
	base.SetupSuite()

	base.Precondition("Авторизация пользователя с кредами из переменных окружения и получение его параметров")
	s.loginData = auth.Login(s.T(), os.Getenv("TEST_LOGIN"), os.Getenv("TEST_PASSWORD"))
}

func (s *TestSuite) TearDownSuite() {
	base.TearDownSuite()
}

func (s *TestSuite) TestCreateAdvertisementPositive() {
	var ad models.Advertisement
	s.Run("Создаем объявление", func() {
		req := models.CreateAdvertisementRequest{
			Title:       utils.RandomString(32),
			Description: utils.RandomString(32),
			Price:       utils.RandomFloat32(1, 1000000),
			Quantity:    int32(utils.RandomInt(0, 100)),
			Photos:      []models.PhotoFile{utils.RandomPhoto(utils.RandomString(12))},
		}
		ad = advertisement.CreateAdvertisement(s.T(), s.loginData.Token, req)
	})
	s.Run("Проверяем, что при GET, полученное объявление совпадает с созданным", func() {
		gotAd := advertisement.GetAdvertisement(s.T(), ad.Id)
		s.Equal(ad.UserId, gotAd.User.Id)
		gotAd.User = ad.User
		s.Equal(ad, gotAd)
	})
	s.Run("Проверяем наличие объявления в поиске", func() {
		ads := advertisement.SearchAdvertisements(s.T(), ad.Title)
		found := slices.ContainsFunc(ads.Items, func(item models.Advertisement) bool {
			return item.Id == ad.Id
		})
		s.True(found, "Не получено объявления с ID, равном созданному")
	})
	s.Run("Проверяем наличие фотографий", func() {
		photos := advertisement.GetAdvertisementPhotos(s.T(), s.loginData.Token, ad.Id)
		s.Require().NotEmpty(photos)
		photo := photos[0]
		s.Require().NotEmpty(photo.Url)

		resp, err := http.Get(photo.Url)
		s.Require().NoError(err)
		defer resp.Body.Close()
		s.Equal(http.StatusOK, resp.StatusCode)
	})

	s.T().Cleanup(func() {
		if ad.Id != "" {
			advertisement.DeleteAdvertisement(s.T(), s.loginData.Token, ad.Id)
		}
	})
}

func (s *TestSuite) TestCreateAdvertisementEmptyRequiredData() {
	exampleReq := models.CreateAdvertisementRequest{
		Description: utils.RandomString(32),
		Price:       utils.RandomFloat32(1, 1000000),
		Quantity:    int32(utils.RandomInt(0, 100)),
		Photos:      []models.PhotoFile{utils.RandomPhoto(utils.RandomString(12))},
	}
	s.Run("Проверяем, что объявление без заголовка не создастся", func() {
		req := exampleReq
		req.Title = ""
		raw := httpAdvertisement.HttpCreateAdvertisement(s.T(), s.loginData.Token, req)
		defer raw.Body.Close()
		s.Equal(http.StatusBadRequest, raw.StatusCode)
	})
	s.Run("Проверяем, что объявление без описания не создастся", func() {
		req := exampleReq
		req.Description = ""
		raw := httpAdvertisement.HttpCreateAdvertisement(s.T(), s.loginData.Token, req)
		defer raw.Body.Close()
		s.Equal(http.StatusBadRequest, raw.StatusCode)
	})
	s.Run("Проверяем, что объявление без фотографий не создастся", func() {
		req := exampleReq
		req.Photos = nil
		raw := httpAdvertisement.HttpCreateAdvertisement(s.T(), s.loginData.Token, req)
		defer raw.Body.Close()
		s.Equal(http.StatusBadRequest, raw.StatusCode)
	})
	s.Run("Проверяем, что объявление с большим количеством фотографий не создастся", func() {
		req := exampleReq
		req.Photos = make([]models.PhotoFile, photosLimit+1)
		for i := range req.Photos {
			req.Photos[i] = utils.RandomPhoto(utils.RandomString(12))
		}
		raw := httpAdvertisement.HttpCreateAdvertisement(s.T(), s.loginData.Token, req)
		defer raw.Body.Close()
		s.Equal(http.StatusBadRequest, raw.StatusCode)
	})
	s.Run("Проверяем, что объявление c негативной ценой не создастся", func() {
		req := exampleReq
		req.Price = -1
		raw := httpAdvertisement.HttpCreateAdvertisement(s.T(), s.loginData.Token, req)
		defer raw.Body.Close()
		s.Equal(http.StatusBadRequest, raw.StatusCode)
	})
	s.Run("Проверяем, что объявление cо слишком большой ценой не создастся", func() {
		req := exampleReq
		req.Price = priceLimit + 1
		raw := httpAdvertisement.HttpCreateAdvertisement(s.T(), s.loginData.Token, req)
		defer raw.Body.Close()
		s.Equal(http.StatusBadRequest, raw.StatusCode)
	})
	s.Run("Проверяем, что объявление c негативным количеством не создастся", func() {
		req := exampleReq
		req.Quantity = -1
		raw := httpAdvertisement.HttpCreateAdvertisement(s.T(), s.loginData.Token, req)
		defer raw.Body.Close()
		s.Equal(http.StatusBadRequest, raw.StatusCode)
	})
	s.Run("Проверяем, что объявление cо слишком большим количеством не создастся", func() {
		req := exampleReq
		req.Quantity = quantityLimit + 1
		raw := httpAdvertisement.HttpCreateAdvertisement(s.T(), s.loginData.Token, req)
		defer raw.Body.Close()
		s.Equal(http.StatusBadRequest, raw.StatusCode)
	})
	s.Run("Проверяем, что объявление cо слишком большим заголовком не создастся", func() {
		req := exampleReq
		req.Title = utils.RandomString(titleLengthLimit + 1)
		raw := httpAdvertisement.HttpCreateAdvertisement(s.T(), s.loginData.Token, req)
		defer raw.Body.Close()
		s.Equal(http.StatusBadRequest, raw.StatusCode)
	})
	s.Run("Проверяем, что объявление cо слишком большим описанием не создастся", func() {
		req := exampleReq
		req.Title = utils.RandomString(descriptionLengthLimit + 1)
		raw := httpAdvertisement.HttpCreateAdvertisement(s.T(), s.loginData.Token, req)
		defer raw.Body.Close()
		s.Equal(http.StatusBadRequest, raw.StatusCode)
	})
}

func (s *TestSuite) TestCreateAdvertisementUnauthorized() {
	req := models.CreateAdvertisementRequest{
		Title:       utils.RandomString(titleLengthLimit + 1),
		Description: utils.RandomString(32),
		Price:       utils.RandomFloat32(1, 1000000),
		Quantity:    int32(utils.RandomInt(0, 100)),
		Photos:      []models.PhotoFile{utils.RandomPhoto(utils.RandomString(12))},
	}
	s.Run("Проверяем, что пустой токен выдаст ошибку авторизации", func() {
		raw := httpAdvertisement.HttpCreateAdvertisement(s.T(), "", req)
		defer raw.Body.Close()
		s.Equal(http.StatusUnauthorized, raw.StatusCode)
	})
	s.Run("Проверяем, что неверный токен выдаст ошибку авторизации", func() {
		token := s.loginData.Token
		lastChar := token[len(token)-1]
		replacement := "a"
		if lastChar == 'a' {
			replacement = "b"
		}
		invalidToken := token[:len(token)-1] + replacement
		raw := httpAdvertisement.HttpCreateAdvertisement(s.T(), invalidToken, req)
		defer raw.Body.Close()
		s.Equal(http.StatusUnauthorized, raw.StatusCode)
	})
}
