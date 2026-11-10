package advertisementSearch

import (
	"api-tests-template/internal/constants/path"
	apiRunner "api-tests-template/internal/helpers/api-runner"
	"net/http"
	"testing"
)

func HttpGetAdvertisements(t *testing.T, textForSearch string) *http.Response {
	return apiRunner.GetRunner().Create().
		Get(path.GetAdvertisementsPath).
		Query("search", textForSearch).
		Expect(t).End().Response
}
