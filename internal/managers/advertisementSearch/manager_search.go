package advertisementSearch

import (
	"api-tests-template/internal/client/http/advertisementSearch"
	"api-tests-template/internal/managers/advertisementSearch/models"
	"encoding/json"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func GetAdvertisementsList(t *testing.T, text string) models.AdvertisementsList {
	raw := advertisementSearch.HttpGetAdvertisements(t, text)
	defer raw.Body.Close()
	body, err := io.ReadAll(raw.Body)
	require.NoError(t, err)

	resp := models.AdvertisementsList{}
	require.NoError(t, json.Unmarshal(body, &resp))
	return resp
}
