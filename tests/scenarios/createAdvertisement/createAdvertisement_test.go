package createAdvertisement

import (
	httpAdvertisement "api-tests-template/internal/client/http/advertisement"
	"api-tests-template/internal/managers/advertisement"
	"api-tests-template/internal/managers/advertisement/models"
	"api-tests-template/internal/managers/auth"
	authModels "api-tests-template/internal/managers/auth/models"
	"api-tests-template/internal/utils"
	base "api-tests-template/tests"
	"api-tests-template/tests/utils/advertisementCreating"
	errorUtils "api-tests-template/tests/utils/error"
	"api-tests-template/tests/utils/token"
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
	ad := advertisementCreating.CreateTestAdvertisement(&s.Suite, s.loginData.Token)
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
	advertisementCreating.CleanupTestAdvertisement(&s.Suite, s.loginData.Token, ad.Id)
}

func (s *TestSuite) TestCreateAdvertisementEmptyRequiredData() {
	exampleReq := models.CreateAdvertisementRequest{
		Title:       utils.RandomString(32),
		Description: utils.RandomString(32),
		Price:       utils.RandomFloat32(1, 1000000),
		Quantity:    int32(utils.RandomInt(0, 100)),
		Photos:      []models.PhotoFile{utils.RandomPhoto(utils.RandomString(12))},
	}
	s.Run("Проверяем, что объявление без заголовка не создастся", func() {
		req := exampleReq
		req.Title = ""
		raw := httpAdvertisement.HttpCreateAdvertisement(s.T(), s.loginData.Token, req)
		s.Equal(http.StatusBadRequest, raw.StatusCode)
		resp := errorUtils.ParseError(s.T(), raw)
		s.Equal("invalid_request", resp.Error)
		s.Equal("title is required", resp.Message)
	})
	s.Run("Проверяем, что объявление без описания не создастся", func() {
		req := exampleReq
		req.Description = ""
		raw := httpAdvertisement.HttpCreateAdvertisement(s.T(), s.loginData.Token, req)
		s.Equal(http.StatusBadRequest, raw.StatusCode)
		resp := errorUtils.ParseError(s.T(), raw)
		s.Equal("invalid_request", resp.Error)
		s.Equal("description is required", resp.Message)
	})
	s.Run("Проверяем, что объявление без фотографий не создастся", func() {
		req := exampleReq
		req.Photos = nil
		raw := httpAdvertisement.HttpCreateAdvertisement(s.T(), s.loginData.Token, req)
		s.Equal(http.StatusBadRequest, raw.StatusCode)
		resp := errorUtils.ParseError(s.T(), raw)
		s.Equal("invalid_request", resp.Error)
		s.Equal("at least 1 photo is required", resp.Message)
	})
	s.Run("Проверяем, что объявление с большим количеством фотографий не создастся", func() {
		req := exampleReq
		req.Photos = make([]models.PhotoFile, photosLimit+1)
		for i := range req.Photos {
			req.Photos[i] = utils.RandomPhoto(utils.RandomString(12))
		}
		raw := httpAdvertisement.HttpCreateAdvertisement(s.T(), s.loginData.Token, req)
		s.Equal(http.StatusBadRequest, raw.StatusCode)
		resp := errorUtils.ParseError(s.T(), raw)
		s.Equal("invalid_request", resp.Error)
		s.Equal("maximum 3 photos allowed", resp.Message)
	})
	s.Run("Проверяем, что объявление c негативной ценой не создастся", func() {
		req := exampleReq
		req.Price = -1
		raw := httpAdvertisement.HttpCreateAdvertisement(s.T(), s.loginData.Token, req)
		s.Equal(http.StatusBadRequest, raw.StatusCode)
		resp := errorUtils.ParseError(s.T(), raw)
		s.Equal("invalid_request", resp.Error)
		s.Equal("price must be between 0 and 1000000 inclusive", resp.Message)
	})
	s.Run("Проверяем, что объявление cо слишком большой ценой не создастся", func() {
		req := exampleReq
		req.Price = priceLimit + 1
		raw := httpAdvertisement.HttpCreateAdvertisement(s.T(), s.loginData.Token, req)
		s.Equal(http.StatusBadRequest, raw.StatusCode)
		resp := errorUtils.ParseError(s.T(), raw)
		s.Equal("invalid_request", resp.Error)
		s.Equal("price must be between 0 and 1000000 inclusive", resp.Message)
	})
	s.Run("Проверяем, что объявление c негативным количеством не создастся", func() {
		req := exampleReq
		req.Quantity = -1
		raw := httpAdvertisement.HttpCreateAdvertisement(s.T(), s.loginData.Token, req)
		s.Equal(http.StatusBadRequest, raw.StatusCode)
		resp := errorUtils.ParseError(s.T(), raw)
		s.Equal("invalid_request", resp.Error)
		s.Equal("quantity must be between 1 and 100 inclusive", resp.Message)
	})
	s.Run("Проверяем, что объявление cо слишком большим количеством не создастся", func() {
		req := exampleReq
		req.Quantity = quantityLimit + 1
		raw := httpAdvertisement.HttpCreateAdvertisement(s.T(), s.loginData.Token, req)
		s.Equal(http.StatusBadRequest, raw.StatusCode)
		resp := errorUtils.ParseError(s.T(), raw)
		s.Equal("invalid_request", resp.Error)
		s.Equal("quantity must be between 1 and 100 inclusive", resp.Message)
	})
	s.Run("Проверяем, что объявление cо слишком большим заголовком не создастся", func() {
		req := exampleReq
		req.Title = utils.RandomString(titleLengthLimit + 1)
		raw := httpAdvertisement.HttpCreateAdvertisement(s.T(), s.loginData.Token, req)
		s.Equal(http.StatusBadRequest, raw.StatusCode)
		resp := errorUtils.ParseError(s.T(), raw)
		s.Equal("invalid_request", resp.Error)
		s.Equal("title must not exceed 50 characters provided 51", resp.Message)
	})
	s.Run("Проверяем, что объявление cо слишком большим описанием не создастся", func() {
		req := exampleReq
		req.Description = utils.RandomString(descriptionLengthLimit + 1)
		raw := httpAdvertisement.HttpCreateAdvertisement(s.T(), s.loginData.Token, req)
		s.Equal(http.StatusBadRequest, raw.StatusCode)
		resp := errorUtils.ParseError(s.T(), raw)
		s.Equal("invalid_request", resp.Error)
		s.Equal("description must be between 1 and 500 characters", resp.Message)
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
		s.Equal(http.StatusUnauthorized, raw.StatusCode)
	})
	s.Run("Проверяем, что неверный токен выдаст ошибку авторизации", func() {
		raw := httpAdvertisement.HttpCreateAdvertisement(s.T(), token.DistortToken(s.loginData.Token), req)
		s.Equal(http.StatusUnauthorized, raw.StatusCode)
	})
}
