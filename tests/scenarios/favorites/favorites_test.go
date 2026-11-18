package favorites

import (
	"os"
	"testing"

	"api-tests-template/internal/managers/auth"
	authModels "api-tests-template/internal/managers/auth/models"
	base "api-tests-template/tests"

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
