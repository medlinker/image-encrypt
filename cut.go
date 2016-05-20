// Package imageEncrypt cut is Cutting Image
package imageEncrypt

import (
	"image"
	"io"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/sosop/imaging"
)

const (
	// DefaultPatitionX default cols
	DefaultPatitionX = 4
	// DefaultPatitionY default rows
	DefaultPatitionY = 4
)

// Cut interface
type Cut interface {
	Cutting(reader io.Reader, filename string, condition ...interface{}) (MetaCuttedImage, error)
}

// RectangleCut cutting image to litle Rectangle image
type RectangleCut struct {
	// cols
	partitionX int
	// rows
	partitionY int

	// Image Storage Interface
	storage Storage
	// Meta-information storage interface
	meta Meta
}

// NewDefaultRectangleCut constructor
func NewDefaultRectangleCut(storage Storage, meta Meta) *RectangleCut {
	return NewRectangleCut(DefaultPatitionX, DefaultPatitionY, storage, meta)
}

// NewRectangleCut constructor
func NewRectangleCut(partitionX, patitionY int, storage Storage, meta Meta) *RectangleCut {
	return &RectangleCut{partitionX: partitionX, partitionY: patitionY, storage: storage, meta: meta}
}

// Cutting implement the interface of Cut
func (r RectangleCut) Cutting(reader io.Reader, filename string, condition ...interface{}) (*MetaCuttedImage, error) {
	ext := strings.ToLower(filepath.Ext(filename))
	src, err := imaging.Decode(reader)
	if err != nil {
		return nil, err
	}

	rect := src.Bounds()
	x := rect.Max.X - rect.Min.X
	y := rect.Max.Y - rect.Min.X
	// step of x
	stepX := x / r.partitionX
	// step of y
	stepY := y / r.partitionY
	images := make([]CuttedImage, r.partitionX*r.partitionY)
	k := 0
	// goroutine save splice image
	wg := new(sync.WaitGroup)
	wg.Add(r.partitionX * r.partitionY)
	for row := 0; row < r.partitionY; row++ {
		for col := 0; col < r.partitionX; col++ {
			images[k] = CuttedImage{ID: k}
			p1 := Point{}
			p2 := Point{}
			if col > 0 {
				p1.X = images[k-1].Points[1].X
			}
			if row > 0 {
				p1.Y = images[k-r.partitionX].Points[1].Y
			}
			if col == r.partitionX-1 {
				p2.X = x
			} else {
				p2.X = p1.X + stepX
			}
			if row == r.partitionY-1 {
				p2.Y = y
			} else {
				p2.Y = p1.Y + stepY
			}
			images[k].Points = []Point{p1, p2}
			retangle := image.Rect(p1.X, p1.Y, p2.X, p2.Y)
			subImg := imaging.Crop(src, retangle)
			img := rotate(subImg, &images[k])
			go r.storage.Save(&images[k], img, strconv.Itoa(k)+filename, wg, ext)
			k++
		}
	}
	wg.Wait()
	metaImage := MetaCuttedImage{images, x, y, Rectangle, ext}
	r.meta.Save(metaImage, condition...)
	return &metaImage, nil
}
