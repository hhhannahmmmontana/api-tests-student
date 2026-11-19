package error

import (
	"api-tests-template/internal/managers/error/models"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func ParseError(t *testing.T, raw *http.Response) models.ErrorResponse {
	defer raw.Body.Close()
	body, err := io.ReadAll(raw.Body)
	require.NoError(t, err)
	var resp models.ErrorResponse
	require.NoError(t, json.Unmarshal(body, &resp))
	return resp
}
