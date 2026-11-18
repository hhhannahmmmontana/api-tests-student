package favorites

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"api-tests-template/internal/constants/path"
	apiRunner "api-tests-template/internal/helpers/api-runner"
)

func HttpCreateFavorite(t *testing.T, token string, id string) *http.Response {
	request := struct {
		AdvertisementID string `json:"advertisement_id"`
	}{
		AdvertisementID: id,
	}
	return apiRunner.GetRunner().Auth(token).Create().Post(path.FavoritesPath).
		ContentType("application/json").
		JSON(request).
		Expect(t).
		End().Response
}

func HttpGet10LastFavorites(t *testing.T, token string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().Get(path.FavoritesPath).
		Query("page", "1").
		Query("page_size", "10").
		Query("sort_by", "added_at").
		Query("sort_order", "desc").
		Expect(t).
		End().Response
}

func HttpDeleteFavorite(t *testing.T, token, id string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().
		Delete(fmt.Sprintf(path.FavoritePath, url.PathEscape(id))).
		Expect(t).
		End().Response
}
