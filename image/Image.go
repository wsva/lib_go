package image

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"net/http"
	"os"
)

func GetImageType(filename string) (string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()

	// why 512 bytes ?
	// see http://golang.org/pkg/net/http/#DetectContentType
	buff := make([]byte, 512)
	_, err = f.Read(buff)
	if err != nil {
		return "", err
	}
	return http.DetectContentType(buff), nil
}

func GetImageAndType(filename string) (image.Image, string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, "", err
	}
	defer f.Close()
	return image.Decode(f)
}

func GetImage(filename string) (image.Image, error) {
	imgtype, err := GetImageType(filename)
	if err != nil {
		return nil, err
	}

	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	switch imgtype {
	case "image/jpeg", "image/jpg":
		return jpeg.Decode(f)
	case "image/gif":
		return gif.Decode(f)
	case "image/png":
		return png.Decode(f)
	default:
		return nil, fmt.Errorf("unsupported image type: %v", imgtype)
	}
}

// width, height
func GetImageSize(img image.Image) (int, int) {
	min, max := img.Bounds().Min, img.Bounds().Max
	return max.X - min.X, max.Y - min.Y
}

func cropImage(img image.Image, rect image.Rectangle) (image.Image, error) {
	switch img := img.(type) {
	case *image.YCbCr:
		return img.SubImage(rect), nil
	case *image.NRGBA:
		return img.SubImage(rect), nil
	case *image.RGBA:
		return img.SubImage(rect), nil
	case *image.Paletted:
		return img.SubImage(rect), nil
	default:
		return nil, errors.New("unsupported image type")
	}
}

func CropImageHeight(img image.Image, offset int) (image.Image, error) {
	min, max := img.Bounds().Min, img.Bounds().Max
	if max.Y-min.Y < offset {
		return nil, fmt.Errorf("offset too large")
	}
	return cropImage(img, image.Rect(
		min.X, min.Y+offset/2, max.X, max.Y-offset/2))
}

func CropImageWidth(img image.Image, offset int) (image.Image, error) {
	min, max := img.Bounds().Min, img.Bounds().Max
	if max.X-min.X < offset {
		return nil, fmt.Errorf("offset too large")
	}
	return cropImage(img, image.Rect(
		min.X+offset/2, min.Y, max.X-offset/2, max.Y))
}

func WriteImagePng(img image.Image, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	return png.Encode(file, img)
}

func GetImageHash(img image.Image) (string, error) {
	buffer := new(bytes.Buffer)
	err := png.Encode(buffer, img)
	if err != nil {
		return "", err
	}
	hash := sha256.Sum256(buffer.Bytes())
	return base64.URLEncoding.EncodeToString(hash[:]), nil
}
