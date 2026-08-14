package main

import (
	"fmt"
	"image"
	"image/color"

	"github.com/disintegration/gift"
	"github.com/disintegration/imaging"
)

func colorLerp(a, b, t float32) float32 {
	return a + t*(b-a)
}

func colorize(src image.Image, bg, fg, mid color.Color) image.Image {
	cBg := color.NRGBAModel.Convert(bg).(color.NRGBA)
	cFg := color.NRGBAModel.Convert(fg).(color.NRGBA)
	cMid := color.NRGBAModel.Convert(mid).(color.NRGBA)

	bgR, bgG, bgB := float32(cBg.R)/255.0, float32(cBg.G)/255.0, float32(cBg.B)/255.0
	fgR, fgG, fgB := float32(cFg.R)/255.0, float32(cFg.G)/255.0, float32(cFg.B)/255.0
	midR, midG, midB := float32(cMid.R)/255.0, float32(cMid.G)/255.0, float32(cMid.B)/255.0

	filter := gift.ColorFunc(func(r0, g0, b0, a0 float32) (float32, float32, float32, float32) {
		gray := 0.299*r0 + 0.587*g0 + 0.114*b0
		var r, g, b float32
		if gray < 0.5 {
			t := gray / 0.5
			r = colorLerp(bgR, midR, t)
			g = colorLerp(bgG, midG, t)
			b = colorLerp(bgB, midB, t)
		} else {
			t := (gray - 0.5) / 0.5
			r = colorLerp(midR, fgR, t)
			g = colorLerp(midG, fgG, t)
			b = colorLerp(midB, fgB, t)
		}
		return r, g, b, a0
	})

	g := gift.New(filter)

	dst := image.NewRGBA(g.Bounds(src.Bounds()))
	g.Draw(dst, src)
	return dst
}

func main() {
	draculaColors := map[string]color.RGBA{
		"bg":    {40, 42, 54, 255},
		"mid":   {68, 71, 90, 255},
		"fg":    {248, 248, 242, 255},
		"red":   {255, 85, 85, 255},
		"green": {80, 250, 123, 255},
		"cyan":  {139, 233, 253, 255},
	}

	img, err := imaging.Open("image.jpg")
	if err != nil {
		fmt.Println("Error opening image:", err)
		return
	}

	result := colorize(img, draculaColors["bg"], draculaColors["fg"], draculaColors["mid"])
	err = imaging.Save(result, "output.jpg")
	if err != nil {
		fmt.Println("error while saving:", err)
		return
	}
	fmt.Println("saved successfully.")
}
