package render

import (
	"fmt"
	"image"
	"image/color"

	"golang.org/x/image/draw"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

func pt(p fixed.Point26_6) image.Point {
	return image.Point{
		X: int(p.X+32) >> 6,
		Y: int(p.Y+32) >> 6,
	}
}

// Function to parse a hex color string into an RGBA color
func HexToRGBA(hex string) (color.RGBA, error) {
	var c color.RGBA
	if len(hex) != 7 || hex[0] != '#' {
		return c, fmt.Errorf("invalid hex color format")
	}

	_, err := fmt.Sscanf(hex, "#%02x%02x%02x", &c.R, &c.G, &c.B)
	if err != nil {
		return c, err
	}

	c.A = 255 // Set alpha to 255 (fully opaque)

	return c, nil
}
func RenderText(text string, scale float64, texture, bgTex image.Image) *image.NRGBA {
	face := basicfont.Face7x13
	stringBounds, _ := font.BoundString(face, text)

	b := image.Rectangle{pt(stringBounds.Min), pt(stringBounds.Max)}
	img := image.NewNRGBA(b)

	if bgTex != nil {
		draw.Draw(img, b, bgTex, image.Point{}, draw.Src)
	}

	d := font.Drawer{
		Dst:  img,
		Src:  texture,
		Face: face,
	}
	d.DrawString(text)

	// normalize bounds to start at 0,0
	img.Rect = img.Bounds().Sub(img.Bounds().Min)

	// scale up, as this font is quite small
	return ScaleImage(img, scale, scale, false)
}
