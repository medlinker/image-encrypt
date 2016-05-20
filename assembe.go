// Package imageEncrypt assembe interface is restoring image
package imageEncrypt

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	"image/draw"
	"sync"

	"github.com/sosop/imaging"
)

// Assembe it's a interface.
// Implement this interface
type Assembe interface {
	// assembing function do Specific work
	Assembing(condition ...interface{}) ([]byte, string, error)
}

// FileSystemAssembe Read slice image from the file system and restore
type FileSystemAssembe struct {
	// Image Storage Interface
	s Storage
	// Meta-information storage interface
	m Meta
}

// NewFileSystemAssembe constructor
func NewFileSystemAssembe(s Storage, m Meta) *FileSystemAssembe {
	return &FileSystemAssembe{s, m}
}

// Implement the interface of Assembe
func (a *FileSystemAssembe) Assembing(condition ...interface{}) ([]byte, string, error) {
	metaImage, err := a.m.Get(condition...)
	if err != nil {
		return nil, "", err
	}
	// create old image
	full := image.NewNRGBA(image.Rect(0, 0, metaImage.MaxX, metaImage.MaxY))
	n := len(metaImage.Images)
	wg := new(sync.WaitGroup)
	wg.Add(n)
	flag := true
	for _, cuttedImage := range metaImage.Images {
		go drawIt(a.s, cuttedImage, full, &flag, wg)
	}
	wg.Wait()
	if !flag {
		return nil, "", errors.New("加载失败")
	}
	// save old image on the file system
	// imaging.Save(full, fmt.Sprint("test", metaImage.Ext))

	// return image bytes
	buf := bytes.NewBuffer(nil)
	f, _ := formats[metaImage.Ext]
	err = imaging.Encode(buf, full, f)
	if err != nil {
		return nil, "", err
	}
	// return image string
	// data := base64.StdEncoding.EncodeToString(buf.Bytes())
	return buf.Bytes(), metaImage.Ext, nil
}

func (a *FileSystemAssembe) AssebingBase64(condition ...interface{}) (string, error) {
	buf, ext, err := a.Assembing(condition...)
	if err != nil {
		return "", err
	}
	imgBase64 := fmt.Sprint("data:image/", ext[1:], ";base64,", base64.StdEncoding.EncodeToString(buf))
	return imgBase64, nil
}

// draw the old image from splice image
func drawIt(s Storage, cuttedImage CuttedImage, bg *image.NRGBA, flag *bool, wg *sync.WaitGroup) {
	defer wg.Done()
	rc, err := s.Get(cuttedImage.Location)
	if err != nil {
		*flag = false
		return
	}
	defer rc.Close()
	img, err := imaging.Decode(rc)
	if err != nil {
		*flag = false
		return
	}
	invImg := inverseRotate(img, cuttedImage.Rotate)
	draw.Draw(bg, image.Rect(cuttedImage.Points[0].X, cuttedImage.Points[0].Y, cuttedImage.Points[1].X, cuttedImage.Points[1].Y), invImg, image.Pt(0, 0), draw.Src)
}
