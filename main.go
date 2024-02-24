package main

import (
	"github.com/gofiber/fiber/v3"
	"image"
	"image/draw"
	"image/gif"
	"image/png"
	"os"
	"strconv"
)

func hello(c fiber.Ctx) error {
	return c.SendString("Hello, World 👋!")
}

func upload(c fiber.Ctx) error {
	img, err := c.FormFile("image")
	if err != nil {
		return c.SendString(err.Error())
	}
	pic, err := img.Open()
	if err != nil {
		return c.SendString(err.Error())
	}
	imgDecoded, err := gif.DecodeAll(pic)
	if err != nil {
		return c.SendString(err.Error())
	}
	gifFile, err := os.Create("./asset/gambar.gif")
	if err != nil {
		panic(err.Error() + "Error creating GIF file")
	}
	defer gifFile.Close()
	err = gif.EncodeAll(gifFile, imgDecoded)

	if err != nil {
		panic(err.Error() + "Error encoding image to PNG")
	}
	return c.SendString("File uploaded successfully")
}

func Extract(c fiber.Ctx) error {
	file, err := os.Open("./asset/gambar.gif")
	if err != nil {
		return c.SendString(err.Error())

	}
	defer file.Close()
	gifData, err := gif.DecodeAll(file)
	if err != nil {
		return c.SendString(err.Error())
	}
	overpaintImage := image.NewRGBA(image.Rect(0, 0, gifData.Config.Width, gifData.Config.Height))
	draw.Draw(overpaintImage, overpaintImage.Bounds(), gifData.Image[0], image.Point{gifData.Config.Width, gifData.Config.Height}, draw.Src)

	for i, frame := range gifData.Image {
		draw.Draw(overpaintImage, overpaintImage.Bounds(), frame, image.ZP, draw.Over)
		out, err := os.Create("./output/frame" + strconv.Itoa(i) + ".png")
		if err != nil {
			return c.SendString(err.Error())
		}
		defer out.Close()
		err = png.Encode(out, overpaintImage)
		if err != nil {
			return c.SendString(err.Error())
		}
	}
	return c.SendString("Extract Succes")
}

func main() {
	app := fiber.New()
	app.Get("/", hello)
	app.Post("/upload", upload)
	app.Post("/extract", Extract)
	app.Static("/", "./asset")
	app.Static("/", "./output")

	app.Listen(":3000")

}
