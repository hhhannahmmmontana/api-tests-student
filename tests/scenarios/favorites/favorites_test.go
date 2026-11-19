package favorites

import (
	httpFavorites "api-tests-template/internal/client/http/favorites"
	"api-tests-template/internal/managers/auth"
	authModels "api-tests-template/internal/managers/auth/models"
	"api-tests-template/internal/managers/favorites"
	"api-tests-template/internal/managers/favorites/models"
	base "api-tests-template/tests"
	"api-tests-template/tests/utils/advertisementCreating"
	errorUtils "api-tests-template/tests/utils/error"
	"api-tests-template/tests/utils/token"
	"net/http"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite
	loginData authModels.LoginOkResponse
	ownerData authModels.LoginOkResponse
}

func TestSuiteRun(t *testing.T) { suite.Run(t, &TestSuite{}) }

func (s *TestSuite) SetupSuite() {
	base.SetupSuite()
	base.Precondition("Авторизация пользователя избранного и владельца тестовых объявлений")
	s.loginData = auth.Login(s.T(), os.Getenv("TEST_FAVORITES_LOGIN"), os.Getenv("TEST_FAVORITES_PASSWORD"))
	s.ownerData = auth.Login(s.T(), os.Getenv("TEST_LOGIN"), os.Getenv("TEST_PASSWORD"))
	require.NotEmpty(s.T(), s.loginData.User.Id)
	require.NotEmpty(s.T(), s.ownerData.User.Id)
	require.NotEqual(s.T(), s.loginData.User.Id, s.ownerData.User.Id, "Нужны два разных пользователя")
}

func (s *TestSuite) TearDownSuite() { base.TearDownSuite() }

func CreateTestFavourite(s *TestSuite) models.Favorite {
	ad := advertisementCreating.CreateTestAdvertisement(&s.Suite, s.ownerData.Token)
	var favorite models.Favorite
	s.Run("Добавляем созданное объявление в избранное", func() {
		favorite = favorites.AddFavorite(s.T(), s.loginData.Token, ad.Id)
		s.Equal(ad.UserId, favorite.Advertisement.User.Id)
		ad.User = favorite.Advertisement.User
		s.Equal(ad, favorite.Advertisement)
	})
	s.Run("Проверяем, что объявление находится в избранных", func() {
		foundFavorite, isFound := favorites.GetLastFavorite(s.T(), s.loginData.Token)
		s.True(isFound, "Список избранного пуст")
		s.Equal(favorite.Advertisement.UserId, foundFavorite.Advertisement.User.Id)
		foundFavorite.Advertisement.User = favorite.Advertisement.User
		s.Equal(favorite, foundFavorite)
	})
	return favorite
}

func (s *TestSuite) TestCreateAndDeleteFavoritePositive() {
	favorite := CreateTestFavourite(s)
	s.Run("Удаляем объявление из избранного", func() {
		favorites.DeleteFavorite(s.T(), s.loginData.Token, favorite.Advertisement.Id)
		foundFavorite, isFound := favorites.GetLastFavorite(s.T(), s.loginData.Token)
		if isFound {
			s.NotEqual(favorite.FavouriteId, foundFavorite.FavouriteId)
		}
	})
	advertisementCreating.CleanupTestAdvertisement(&s.Suite, s.ownerData.Token, favorite.Advertisement.Id)
}

func (s *TestSuite) TestCreateFavoriteAndDeleteAdvertisementPositive() {
	favorite := CreateTestFavourite(s)
	s.Run("Удаляем объявление", func() {
		advertisementCreating.DeleteAdvertisement(&s.Suite, s.ownerData.Token, favorite.Advertisement.Id)
		foundFavorite, isFound := favorites.GetLastFavorite(s.T(), s.loginData.Token)
		if isFound {
			s.NotEqual(favorite.FavouriteId, foundFavorite.FavouriteId)
		}
	})
	advertisementCreating.CleanupTestAdvertisement(&s.Suite, s.ownerData.Token, favorite.Advertisement.Id)
}

func (s *TestSuite) TestAddUnexistingAdvertisementToFavorites() {
	s.Run("Создаем избранное с несуществующим объявлением", func() {
		raw := httpFavorites.HttpCreateFavorite(s.T(), s.loginData.Token, uuid.New().String())
		s.Equal(http.StatusNotFound, raw.StatusCode)
		resp := errorUtils.ParseError(s.T(), raw)
		s.Equal("not_found", resp.Error)
		s.Equal("Advertisement not found", resp.Message)
	})
	s.Run("Удаляем из избранного несуществующее объявление", func() {
		raw := httpFavorites.HttpDeleteFavorite(s.T(), s.loginData.Token, uuid.New().String())
		s.Equal(http.StatusNotFound, raw.StatusCode)
		resp := errorUtils.ParseError(s.T(), raw)
		s.Equal("not_found", resp.Error)
		s.Equal("Favorite entry not found", resp.Message)
	})
}

func (s *TestSuite) TestUnauthorizedFavorite() {
	ad := advertisementCreating.CreateTestAdvertisement(&s.Suite, s.ownerData.Token)
	s.Run("Добавляем в избранное, не войдя в систему", func() {
		raw := httpFavorites.HttpCreateFavorite(s.T(), "", ad.Id)
		s.Equal(http.StatusUnauthorized, raw.StatusCode)
	})
	s.Run("Добавляем в избранное с некорректным токеном", func() {
		raw := httpFavorites.HttpCreateFavorite(s.T(), token.DistortToken(s.loginData.Token), ad.Id)
		s.Equal(http.StatusUnauthorized, raw.StatusCode)
	})
	advertisementCreating.CleanupTestAdvertisement(&s.Suite, s.ownerData.Token, ad.Id)
}

func (s *TestSuite) TestAddAdvertisementTwiceToFavorites() {
	favorite := CreateTestFavourite(s)
	s.Run("Добавляем объявление в избранное снова", func() {
		raw := httpFavorites.HttpCreateFavorite(s.T(), s.loginData.Token, favorite.Advertisement.Id)
		s.Equal(http.StatusConflict, raw.StatusCode)
		resp := errorUtils.ParseError(s.T(), raw)
		s.Equal("already_exists", resp.Error)
		s.Equal("Advertisement is already in favorites", resp.Message)
	})
	advertisementCreating.CleanupTestAdvertisement(&s.Suite, s.ownerData.Token, favorite.Advertisement.Id)
}

func (s *TestSuite) TestDeleteAdvertisementTwiceFromFavorites() {
	favorite := CreateTestFavourite(s)
	s.Run("Удаляем объявление из избранного", func() {
		favorites.DeleteFavorite(s.T(), s.loginData.Token, favorite.Advertisement.Id)
	})
	s.Run("Удаляем объявление из избранного снова", func() {
		raw := httpFavorites.HttpDeleteFavorite(s.T(), s.loginData.Token, favorite.Advertisement.Id)
		s.Equal(http.StatusNotFound, raw.StatusCode)
		resp := errorUtils.ParseError(s.T(), raw)
		s.Equal("not_found", resp.Error)
		s.Equal("Favorite entry not found", resp.Message)
	})
	advertisementCreating.CleanupTestAdvertisement(&s.Suite, s.ownerData.Token, favorite.Advertisement.Id)
}

func (s *TestSuite) TestAddOwnedAdvertisementToFavorites() {
	ad := advertisementCreating.CreateTestAdvertisement(&s.Suite, s.ownerData.Token)
	s.Run("Добавляем созданное объявление в избранное", func() {
		raw := httpFavorites.HttpCreateFavorite(s.T(), s.ownerData.Token, ad.Id)
		s.Equal(http.StatusForbidden, raw.StatusCode)
		resp := errorUtils.ParseError(s.T(), raw)
		s.Equal("forbidden", resp.Error)
		s.Equal("Cannot add own advertisement to favorites", resp.Message)
	})
	advertisementCreating.CleanupTestAdvertisement(&s.Suite, s.ownerData.Token, ad.Id)
}
