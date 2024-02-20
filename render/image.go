package render

import (
	"fmt"
	"image"
	"image/color"
	colorpkg "image/color"
	"image/jpeg"
	"image/png"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"golang.org/x/image/draw"
)

func OnlyColor(width int, height int, color []string) ([]string, error) {
	// Construct message list
	var messages []string
	for _, col := range color {
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				messages = append(messages, fmt.Sprintf("PX %d %d %s\n", x, y, col))
			}
		}
	}
	return messages, nil
}

func GetPixels(filePath string, startingX int, startingY int) ([]string, error) {
	// Open the image file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Determine the image format
	format := strings.ToLower(filepath.Ext(filePath))
	var img image.Image
	switch format {
	case ".jpg", ".jpeg":
		img, err = jpeg.Decode(file)
	case ".png":
		img, err = png.Decode(file)
	default:
		return nil, fmt.Errorf("unsupported image format")
	}
	if err != nil {
		return nil, err
	}

	// Get dimensions of the image
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	// Define the background color to filter out (e.g., white)
	backgroundColor := colorpkg.RGBA{255, 255, 255, 255} // White color in RGBA format

	// Iterate over each pixel and collect their information
	var pixels []string
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// Get the color of the pixel at (x, y)
			rgba := img.At(x, y)
			r, g, b, a := rgba.RGBA()
			if a == 0 || rgba == backgroundColor {
				continue
			}
			// Format the pixel information and append to the slice
			pixels = append(pixels, fmt.Sprintf("PX %d %d %02x%02x%02x\n", startingX+x, startingY+y, r>>8, g>>8, b>>8))
		}
	}

	return pixels, nil
}

// From Hochwasser

func ScaleImage(img image.Image, factorX, factorY float64, highQuality bool) (scaled *image.NRGBA) {
	b := img.Bounds()
	newX := int(math.Ceil(factorX * float64(b.Max.X)))
	newY := int(math.Ceil(factorY * float64(b.Max.Y)))
	scaledBounds := image.Rect(0, 0, newX, newY)
	scaledImg := image.NewNRGBA(scaledBounds)
	scaler := draw.NearestNeighbor
	if highQuality {
		scaler = draw.CatmullRom
	}
	scaler.Scale(scaledImg, scaledBounds, img, b, draw.Src, nil)
	return scaledImg
}

type Commands [][]byte

type RenderOrder uint8

func (t RenderOrder) String() string   { return []string{"→", "↓", "←", "↑", "random"}[t] }
func (t RenderOrder) IsVertical() bool { return t&0b01 != 0 }
func (t RenderOrder) IsReverse() bool  { return t&0b10 != 0 }
func NewOrder(v string) RenderOrder {
	switch v {
	case "ltr", "l", "→":
		return LeftToRight
	case "rtl", "r", "←":
		return RightToLeft
	case "ttb", "t", "↓":
		return TopToBottom
	case "btt", "b", "↑":
		return BottomToTop
	default:
		return Shuffle
	}
}

const (
	LeftToRight = 0b000
	TopToBottom = 0b001
	RightToLeft = 0b010
	BottomToTop = 0b011
	Shuffle     = 0b100
)

// Shuffle reorders commands randomly, in place.
func (c Commands) Shuffle() {
	for i := range c {
		j := rand.Intn(i + 1)
		c[i], c[j] = c[j], c[i]
	}
}
func (c Commands) ToString() []string {
	var output []string
	for _, i := range c {
		output = append(output, string(i))
	}
	return output
}

func ReadImage(path string) (*image.NRGBA, error) {
	reader, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	img, _, err := image.Decode(reader)
	if err != nil {
		return nil, err
	}
	return imgToNRGBA(img), nil
}

func imgToNRGBA(img image.Image) *image.NRGBA {
	b := img.Bounds()
	r := image.NewNRGBA(b)
	for x := b.Min.X; x < b.Max.X; x++ {
		for y := b.Min.Y; y < b.Max.Y; y++ {
			r.Set(x, y, img.At(x, y))
		}
	}
	return r
}

// CommandsFromImage converts an image to the respective pixelflut Commands
func CommandsFromImage(img *image.NRGBA, order RenderOrder, offset image.Point) (cmds Commands) {
	b := img.Bounds()
	cmds = make([][]byte, b.Size().X*b.Size().Y)
	numCmds := 0

	max1 := b.Max.X
	max2 := b.Max.Y
	min1 := b.Min.X
	min2 := b.Min.Y
	dir := 1
	if order.IsVertical() {
		max1, max2 = max2, max1
		min1, min2 = min2, min1
	}
	if order.IsReverse() {
		min1, max1 = max1, min1
		min2, max2 = max2, min2
		dir *= -1
	}

	for i1 := min1; i1 != max1; i1 += dir {
		for i2 := min2; i2 != max2; i2 += dir {
			x := i1
			y := i2
			if order.IsVertical() {
				x, y = y, x
			}

			c := img.NRGBAAt(x, y)
			if c.A == 0 {
				continue
			}

			var cmd []byte
			cmd = append(cmd, []byte("PX ")...)
			cmd = strconv.AppendUint(cmd, uint64(x+offset.X), 10)
			cmd = append(cmd, ' ')
			cmd = strconv.AppendUint(cmd, uint64(y+offset.Y), 10)
			cmd = append(cmd, ' ')
			appendColor(&cmd, c)
			cmd = append(cmd, '\n')
			cmds[numCmds] = cmd
			numCmds++
		}
	}

	cmds = cmds[:numCmds]

	if order == Shuffle {
		cmds.Shuffle()
	}

	return
}

func appendColor(buf *[]byte, c color.NRGBA) {
	var mask uint32 = 0xf0000000
	// merge into uint32
	var col = uint32(c.R)<<24 + uint32(c.G)<<16 + uint32(c.B)<<8 + uint32(c.A)
	// if alpha is ff, drop it.
	if 0xff&col == 0xff {
		col = col >> 8
		mask = mask >> 8
	}
	// add leading zeros if needed
	for mask > 0xf {
		if mask&col == 0 {
			*buf = append(*buf, '0')
		} else {
			break
		}
		mask = mask >> 4
	}
	*buf = strconv.AppendUint(*buf, uint64(col), 16)
}
