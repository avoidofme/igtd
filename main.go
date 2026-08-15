package main

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"os"
	"path/filepath"
	"strings"

	"github.com/disintegration/gift"
	"github.com/disintegration/imaging"
	"github.com/spf13/cobra"
)

var draculaColors = map[string]color.RGBA{
	"bg":     {40, 42, 54, 255},
	"mid":    {68, 71, 90, 255},
	"fg":     {248, 248, 242, 255},
	"cyan":   {139, 233, 253, 255},
	"green":  {80, 250, 123, 255},
	"orange": {255, 184, 108, 255},
	"pink":   {255, 121, 198, 255},
	"purple": {189, 147, 249, 255},
	"red":    {255, 85, 85, 255},
}

func colorLerp(a, b, t float32) float32 {
	return a + t*(b-a)
}

func colorize(src image.Image, c1, c2, c3 color.RGBA) *image.RGBA {
	c1R, c1G, c1B := float32(c1.R)/255.0, float32(c1.G)/255.0, float32(c1.B)/255.0
	c2R, c2G, c2B := float32(c2.R)/255.0, float32(c2.G)/255.0, float32(c2.B)/255.0
	c3R, c3G, c3B := float32(c3.R)/255.0, float32(c3.G)/255.0, float32(c3.B)/255.0

	filter := gift.ColorFunc(func(r0, g0, b0, a0 float32) (float32, float32, float32, float32) {
		gray := 0.299*r0 + 0.587*g0 + 0.114*b0
		var r, g, b float32
		if gray < 0.5 {
			t := gray / 0.5
			r = colorLerp(c1R, c3R, t)
			g = colorLerp(c1G, c3G, t)
			b = colorLerp(c1B, c3B, t)
		} else {
			t := (gray - 0.5) / 0.5
			r = colorLerp(c3R, c2R, t)
			g = colorLerp(c3G, c2G, t)
			b = colorLerp(c3B, c2B, t)
		}
		return r, g, b, a0
	})

	g := gift.New(filter)
	dst := image.NewRGBA(g.Bounds(src.Bounds()))
	g.Draw(dst, src)
	return dst
}

func composite(fg, bg image.Image, mask *image.Gray) *image.RGBA {
	bounds := fg.Bounds()
	dst := image.NewRGBA(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			fColor := fg.At(x, y)
			bColor := bg.At(x, y)
			mVal, _, _, _ := mask.At(x, y).RGBA()
			alpha := float64(mVal>>8) / 255.0

			fr, fg_, fb, fa := fColor.RGBA()
			br, bg_, bb, ba := bColor.RGBA()

			// Linear interpolation based on mask intensity (0.0 to 1.0)
			outR := uint8((float64(fr>>8)*alpha + float64(br>>8)*(1-alpha)))
			outG := uint8((float64(fg_>>8)*alpha + float64(bg_>>8)*(1-alpha)))
			outB := uint8((float64(fb>>8)*alpha + float64(bb>>8)*(1-alpha)))
			outA := uint8((float64(fa>>8)*alpha + float64(ba>>8)*(1-alpha)))

			dst.Set(x, y, color.RGBA{outR, outG, outB, outA})
		}
	}
	return dst
}

func applyDracula(src image.Image) image.Image {
	bounds := src.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	bnwImage := colorize(src, draculaColors["bg"], draculaColors["fg"], draculaColors["mid"])
	redImage := colorize(src, draculaColors["pink"], draculaColors["red"], draculaColors["red"])
	greenImage := colorize(src, draculaColors["orange"], draculaColors["green"], draculaColors["green"])
	blueImage := colorize(src, draculaColors["purple"], draculaColors["cyan"], draculaColors["cyan"])

	redMask := image.NewGray(image.Rect(0, 0, width, height))
	greenMask := image.NewGray(image.Rect(0, 0, width, height))
	blueMask := image.NewGray(image.Rect(0, 0, width, height))

	for y := range height {
		for x := range width {
			c := src.At(x, y)
			r32, g32, b32, _ := c.RGBA()
			r, g, b := int(r32>>8), int(g32>>8), int(b32>>8)

			redDominance := int(math.Min(float64(r-g), float64(r-b)))
			greenDominance := int(math.Min(float64(g-r), float64(g-b)))
			blueDominance := int(math.Min(float64(b-r), float64(b-g)))

			if redDominance >= 20 && r >= 30 {
				val := int(float64(redDominance) * 1.2)
				redMask.SetGray(x, y, color.Gray{Y: uint8(math.Min(float64(val), 255))})
			}
			if greenDominance >= 20 && g >= 30 {
				val := int(float64(greenDominance) * 1.2)
				greenMask.SetGray(x, y, color.Gray{Y: uint8(math.Min(float64(val), 255))})
			}
			if blueDominance >= 20 && b >= 30 {
				val := int(float64(blueDominance) * 1.2)
				blueMask.SetGray(x, y, color.Gray{Y: uint8(math.Min(float64(val), 255))})
			}
		}
	}

	blurG := gift.New(gift.GaussianBlur(3.0))
	blurredRedMask := image.NewGray(bounds)
	blurredGreenMask := image.NewGray(bounds)
	blurredBlueMask := image.NewGray(bounds)

	blurG.Draw(blurredRedMask, redMask)
	blurG.Draw(blurredGreenMask, greenMask)
	blurG.Draw(blurredBlueMask, blueMask)

	res1 := composite(redImage, bnwImage, blurredRedMask)
	res2 := composite(greenImage, res1, blurredGreenMask)
	result := composite(blueImage, res2, blurredBlueMask)

	return result
}

var rootCmd = &cobra.Command{
	Use:     "./igtd <image path>",
	Example: "./igtd mycoolimage.jpg\n./igtd mycoolimage.jpg -o output.jpg",
	Short:   "Dracula theme color applier",
	Args:    cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		inputPath := args[0]
		outputPath, err := cmd.Flags().GetString("output")
		if err != nil {
			return fmt.Errorf("error parsing flag value: %v", err)
		}
		file, err := os.Stat(inputPath)
		if os.IsNotExist(err) {
			return err
		}
		if file.IsDir() {
			return fmt.Errorf("%v is a dir", inputPath)
		}
		output, err := os.Stat(outputPath)
		if err == nil && output.IsDir() {
			basename := filepath.Base(inputPath)
			outputPath = filepath.Join(outputPath, "igtd_"+basename)
		} else if os.IsNotExist(err) {
			ext := filepath.Ext(outputPath)
			if ext == "" || !strings.Contains(".jpg.png.jpeg", strings.ToLower(ext)) {
				if err := os.MkdirAll(outputPath, 0o755); err != nil {
					return fmt.Errorf("error while creating parent directory: %v", err)
				}
				basename := filepath.Base(inputPath)
				outputPath = filepath.Join(outputPath, "igtd_"+basename)
			}
		}
		src, err := imaging.Open(inputPath)
		if err != nil {
			return fmt.Errorf("error opening image: %v", err)
		}
		result := applyDracula(src)
		if err := imaging.Save(result, outputPath); err != nil {
			return fmt.Errorf("error saving image: %v", err)
		}
		fmt.Println("Saved to", outputPath)
		return nil
	},
}

func main() {
	rootCmd.Flags().StringP("output", "o", "./", "output path")

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
