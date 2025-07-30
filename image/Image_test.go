package image_test

import (
	"fmt"
	"testing"

	"github.com/wsva/lib_go/image"
)

func TestImage(T *testing.T) {
	img, err := image.GetImage("testimage.png")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(image.GetImageSize(img))
	img, err = image.CropImageHeight(img, 100)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(image.GetImageSize(img))
	img, err = image.CropImageWidth(img, 100)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(image.GetImageSize(img))
}
