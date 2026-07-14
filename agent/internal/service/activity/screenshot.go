package activity

import (
	"image"
	"image/png"
	"log"
	"os"
	"strings"

	"github.com/kbinani/screenshot"
)

func TakeScreenshot() (image.Image, error) {
	img, err := screenshot.CaptureDisplay(0)
	if err != nil {
		log.Println("error capturing display:", err)
		return nil, err
	}
	return img, nil
}

func SaveScreenshot(img image.Image, path string) error {
	s := strings.Split(path, ".")
	extension := s[len(s)-1]
	if extension != "png" {
		path = path + ".png"
	}
	file, err := os.Create(path)
	if err != nil {
		log.Println("error creating file:", err)
		return err
	}
	defer file.Close()
	err = png.Encode(file, img)
	if err != nil {
		log.Println("error encoding image:", err)
		return err
	}
	return nil
}
