package main

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"log"
	"math"
	"net/http"
	"os"
	"sync"

	"github.com/gin-gonic/gin"
)

var digitFiles [10]string

func main() {

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.GET("/counter", func(c *gin.Context) {
		img, err := createDigis()

		if err != nil {
			log.Println(err)
			c.JSON(http.StatusOK, gin.H{
				"message": "Failed to generate of image.",
			})
			return
		}

		extraHeaders := map[string]string{
			"Content-Disposition": `attachment; filename="counter.png"`,
		}

		buffer := new(bytes.Buffer)
		png.Encode(buffer, img)

		c.DataFromReader(http.StatusOK, int64(buffer.Len()), "", buffer, extraHeaders)
	})

	r.Run()
}

func init() {
	digitFiles[0] = "img/1/0.png"
	digitFiles[1] = "img/1/1.png"
	digitFiles[2] = "./img/1/2.png"
	digitFiles[3] = "./img/1/3.png"
	digitFiles[4] = "./img/1/4.png"
	digitFiles[5] = "./img/1/5.png"
	digitFiles[6] = "./img/1/6.png"
	digitFiles[7] = "./img/1/7.png"
	digitFiles[8] = "./img/1/8.png"
	digitFiles[9] = "./img/1/9.png"
}

func loadImages() ([10]image.Image, error) {

	var digitsImg [10]image.Image
	for i := 0; i < 10; i++ {
		file, err := os.Open(digitFiles[i])
		if err != nil {
			return digitsImg, errors.New(fmt.Sprintf("%s path:%s", err.Error(), digitFiles[i]))
		}

		img, err := png.Decode(file)
		if err != nil {
			return digitsImg, err
		}
		digitsImg[i] = img
	}
	return digitsImg, nil
}

func createImage(count uint64, zPadLen int) (image.Image, error) {

	digitsImg, err := loadImages()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	genImg := image.NewRGBA(image.Rect(0, 0, digitsImg[0].Bounds().Dx()*zPadLen, digitsImg[0].Bounds().Dy()))

	for i := 0; i < zPadLen; i++ {
		dig := int((count / uint64(math.Pow10(zPadLen-i-1))) % 10)
		img := digitsImg[dig]
		rect := image.Rect(img.Bounds().Dx()*i, 0, img.Bounds().Dx()*(i+1), img.Bounds().Dy())
		draw.Draw(genImg, rect, img, image.Point{0, 0}, draw.Over)
	}

	return genImg, nil
}

func createDigis() (image.Image, error) {
	var mu sync.Mutex
	var filename = "counter.dat"

	mu.Lock()

	count, err := readCount(filename)
	if err != nil {
		count = 1
	}

	count++

	err = saveCount(filename, count)
	if err != nil {
		return nil, err
	}

	img, err := createImage(count, 10)

	mu.Unlock()

	return img, err
}

func readCount(filepath string) (uint64, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	var count uint64
	_, err = fmt.Fscanf(file, "%d", &count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func saveCount(filepath string, count uint64) error {
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = fmt.Fprintf(file, "%d", count)
	if err != nil {
		return err
	}
	return nil
}
