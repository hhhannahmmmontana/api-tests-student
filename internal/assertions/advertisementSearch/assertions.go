package advertisementSearch

import (
	"api-tests-template/internal/managers/advertisementSearch/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func ContainsAdvertisementByTitle(t *testing.T, items []models.Item, title string) {
	t.Helper()

	assert.Condition(t, func() bool {
		for _, item := range items {
			if item.Title == title {
				return true
			}
		}
		return false
	}, "Объявление с названием %q не найдено", title)
}
