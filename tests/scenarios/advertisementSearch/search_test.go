package advertisementSearch

import (
	asserions "api-tests-template/internal/assertions/advertisementSearch"
	"api-tests-template/internal/managers/advertisementSearch"
	"api-tests-template/internal/managers/advertisementSearch/models"
	"api-tests-template/internal/utils"
	base "api-tests-template/tests"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const (
	existingAdTitle = "Холодильник Samsung от sql enjoyer"
	searchBeginning = "Холодильник"
	searchMiddle    = "Samsung"
	searchEnd       = "от sql enjoyer"
)

type TestSuite struct {
	suite.Suite
	searchData models.AdvertisementsList
}

func TestSuiteRun(t *testing.T) {
	suite.Run(t, &TestSuite{})
}

func (s *TestSuite) SetupSuite() {
	base.SetupSuite()

	base.Precondition("Проверяем, что константы заданы")
	require.NotEmpty(s.T(), existingAdTitle, "задайте existingAdTitle")
	require.NotEmpty(s.T(), searchBeginning, "задайте searchBeginning")
	require.NotEmpty(s.T(), searchMiddle, "задайте searchMiddle")
	require.NotEmpty(s.T(), searchEnd, "задайте searchEnd")
}

func (s *TestSuite) TestSearchByTitle() {
	list := advertisementSearch.GetAdvertisementsList(s.T(), searchMiddle)

	s.Run("В ответе есть хотя бы одно объявление", func() {
		require.NotEmpty(s.T(), list.Items)
	})
	s.Run("Среди найденных есть наше объявление", func() {
		asserions.ContainsAdvertisementByTitle(s.T(), list.Items, existingAdTitle)
	})
}

func (s *TestSuite) TestSearchByTitle_Substring() {
	cases := []struct {
		name   string
		search string
	}{
		{"начало", searchBeginning},
		{"середина", searchMiddle},
		{"конец", searchEnd},
		{"нижний регистр", strings.ToLower(searchMiddle)},
		{"верхний регистр", strings.ToUpper(searchMiddle)},
		{"смешанный регистр", utils.MixCase(searchMiddle)},
	}

	for _, tc := range cases {
		s.Run(tc.name, func() {
			list := advertisementSearch.GetAdvertisementsList(s.T(), tc.search)
			asserions.ContainsAdvertisementByTitle(s.T(), list.Items, existingAdTitle)
		})
	}
}
