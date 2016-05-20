// Package imageEncrypt models
package imageEncrypt

import (
	"image"
	"math/rand"
	"time"

	"github.com/sosop/imaging"
)

// shape of cutting
const (
	Rectangle = iota
	RightTriangle

	// Degree0
	Degree0 = iota
	// Degree90
	Degree90
	// Degree180
	Degree180
	// Degree270
	Degree270
)

// image's formats
var (
	formats = map[string]imaging.Format{
		".jpg":  imaging.JPEG,
		".jpeg": imaging.JPEG,
		".png":  imaging.PNG,
		".tif":  imaging.TIFF,
		".tiff": imaging.TIFF,
		".bmp":  imaging.BMP,
		".gif":  imaging.GIF,
	}
)

// Point
type Point struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// MetaCuttedImage meta information
type MetaCuttedImage struct {
	Images []CuttedImage `json:"images"`
	MaxX   int           `json:"maxX"`
	MaxY   int           `json:"maxY"`
	Shape  int           `json:"shape"`
	Ext    string        `json:"ext"`
}

// CuttedImage slice image
type CuttedImage struct {
	ID       int     `json:"id"`
	Location string  `json:"location"`
	Points   []Point `json:"points"`
	Rotate   int     `json:"rotate"`
}

func rotate(img *image.NRGBA, cuttedImage *CuttedImage) *image.NRGBA {
	switch randRotate() {
	case Degree90:
		cuttedImage.Rotate = Degree90
		return imaging.Rotate90(img)
	case Degree180:
		cuttedImage.Rotate = Degree180
		return imaging.Rotate180(img)
	case Degree270:
		cuttedImage.Rotate = Degree270
		return imaging.Rotate270(img)
	default:
		cuttedImage.Rotate = Degree0
	}
	return img
}

func inverseRotate(img image.Image, rotateDegree int) image.Image {
	switch rotateDegree {
	case Degree90:
		return imaging.Rotate270(img)
	case Degree180:
		return imaging.Rotate180(img)
	case Degree270:
		return imaging.Rotate90(img)
	}
	return img
}

func randRotate() int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(4)
}
