package utils

import (
	"api-tests-template/internal/managers/advertisement/models"
	"bytes"
	"image"
	"image/color"
	"image/png"
	"math/rand/v2"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func RandomFloat32(min, max float32) float32 {
	return min + rand.Float32()*(max-min)
}

func RandomInt(min, max int) int {
	return min + rand.IntN(max+1-min)
}

func RandomString(length int) string {
	b := make([]byte, length)
	for i := range b {
		randCharsetIndex := rand.IntN(len(charset))
		b[i] = charset[randCharsetIndex]
	}

	return string(b)
}

func RandomPhoto(name string) models.PhotoFile {
	colorValue := color.RGBA{
		R: uint8(RandomInt(0, 255)),
		G: uint8(RandomInt(0, 255)),
		B: uint8(RandomInt(0, 255)),
		A: 255,
	}
	picture := image.NewRGBA(image.Rect(0, 0, 16, 16))
	picture.Set(0, 0, colorValue)
	var content bytes.Buffer
	png.Encode(&content, picture)
	return models.PhotoFile{
		Name:    name + ".png",
		Content: content.Bytes(),
	}
}
