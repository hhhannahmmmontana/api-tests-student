package advertisement

import (
	"encoding/json"
	"io"
	"testing"

	client "api-tests-template/internal/client/http/advertisement"
	"api-tests-template/internal/managers/advertisement/models"

	"github.com/stretchr/testify/require"
)

func CreateAdvertisement(t *testing.T, token string, request models.CreateAdvertisementRequest) models.Advertisement {
	raw := client.HttpCreateAdvertisement(t, token, request)
	defer raw.Body.Close()
	body, err := io.ReadAll(raw.Body)
	require.NoError(t, err)

	resp := models.Advertisement{}
	require.NoError(t, json.Unmarshal(body, &resp))
	return resp
}

func GetAdvertisement(t *testing.T, id string) models.Advertisement {
	raw := client.HttpGetAdvertisement(t, id)
	defer raw.Body.Close()
	body, err := io.ReadAll(raw.Body)
	require.NoError(t, err)

	resp := models.Advertisement{}
	require.NoError(t, json.Unmarshal(body, &resp))
	return resp
}

func DeleteAdvertisement(t *testing.T, token, id string) {
	raw := client.HttpDeleteAdvertisement(t, token, id)
	defer raw.Body.Close()
	_, err := io.ReadAll(raw.Body)
	require.NoError(t, err)
}

func GetAdvertisementPhotos(t *testing.T, token, id string) []models.Photo {
	raw := client.HttpGetAdvertisementPhotos(t, token, id)
	defer raw.Body.Close()
	body, err := io.ReadAll(raw.Body)
	require.NoError(t, err)

	var resp []models.Photo
	require.NoError(t, json.Unmarshal(body, &resp))
	return resp
}

func SearchAdvertisements(t *testing.T, text string) models.AdvertisementsList {
	raw := client.HttpSearchAdvertisements(t, text)
	defer raw.Body.Close()
	body, err := io.ReadAll(raw.Body)
	require.NoError(t, err)

	resp := models.AdvertisementsList{}
	require.NoError(t, json.Unmarshal(body, &resp))
	return resp
}

func GetMyAdvertisements(t *testing.T, token string, expectedStatusCode int, limit int) string {
	myAdvertisementsRawResponse := client.HttpGetMyAdvertisements(t, token, limit)
	defer myAdvertisementsRawResponse.Body.Close()

	require.Equalf(t, expectedStatusCode, myAdvertisementsRawResponse.StatusCode,
		"HTTP status code должен быть %d", expectedStatusCode)

	myAdvertisementsResponse, errResponse := io.ReadAll(myAdvertisementsRawResponse.Body)
	require.NoError(t, errResponse)

	return string(myAdvertisementsResponse)
}
