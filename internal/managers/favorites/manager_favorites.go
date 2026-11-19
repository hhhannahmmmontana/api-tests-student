package favorites

import (
	"encoding/json"
	"io"
	"testing"

	client "api-tests-template/internal/client/http/favorites"
	"api-tests-template/internal/managers/favorites/models"

	"github.com/stretchr/testify/require"
)

func AddFavorite(t *testing.T, token string, id string) models.Favorite {
	raw := client.HttpCreateFavorite(t, token, id)
	defer raw.Body.Close()
	body, err := io.ReadAll(raw.Body)
	require.NoError(t, err)

	resp := models.Favorite{}
	require.NoError(t, json.Unmarshal(body, &resp))
	return resp
}

func GetLastFavorite(t *testing.T, token string) (models.Favorite, bool) {
	raw := client.HttpGet10LastFavorites(t, token)
	defer raw.Body.Close()
	body, err := io.ReadAll(raw.Body)
	require.NoError(t, err)

	resp := models.GetFavoritesResponse{}
	require.NoError(t, json.Unmarshal(body, &resp))
	if len(resp.Items) == 0 {
		return models.Favorite{}, false
	}
	return resp.Items[0], true
}

func DeleteFavorite(t *testing.T, token, id string) {
	raw := client.HttpDeleteFavorite(t, token, id)
	defer raw.Body.Close()
	_, err := io.ReadAll(raw.Body)
	require.NoError(t, err)
}
