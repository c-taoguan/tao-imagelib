package imageop

import (
	"fmt"
	"image"
	"log"
	"os"

	_ "image/jpeg"
	"image/png"

	"golang.org/x/image/draw"
)

// The functions that manipulate the images
func ResizeImage(src image.Image, dstSize image.Point) *image.RGBA {
	srcRect := src.Bounds()
	dstRect := image.Rectangle{
		Min: image.Point{0, 0},
		Max: dstSize,
	}
	dst := image.NewRGBA(dstRect)
	draw.CatmullRom.Scale(dst, dstRect, src, srcRect, draw.Over, nil)
	return dst
}

func ConvertImageToGrey(img image.Image, new_img_name string) {
	// Convert image to grayscale
	grayImg := image.NewGray(img.Bounds())
	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			grayImg.Set(x, y, img.At(x, y))
		}
	}

	// Working with grayscale image, e.g. convert to png
	f, err := os.Create(new_img_name)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	if err := png.Encode(f, grayImg); err != nil {
		log.Fatal(err)
	}
}

func CropImage(img image.Image, cropRect image.Rectangle) (cropImg image.Image, newImg bool) {
	//Interface for asserting whether `img`
	//implements SubImage or not.
	//This can be defined globally.
	type CropableImage interface {
		image.Image
		SubImage(r image.Rectangle) image.Image
	}

	if p, ok := img.(CropableImage); ok {
		// Call SubImage. This should be fast,
		// since SubImage (usually) shares underlying pixel.
		cropImg = p.SubImage(cropRect)
	} else if cropRect = cropRect.Intersect(img.Bounds()); !cropRect.Empty() {
		// If `img` does not implement `SubImage`,
		// copy (and silently convert) the image portion to RGBA image.
		rgbaImg := image.NewRGBA(cropRect)
		for y := cropRect.Min.Y; y < cropRect.Max.Y; y++ {
			for x := cropRect.Min.X; x < cropRect.Max.X; x++ {
				rgbaImg.Set(x, y, img.At(x, y))
			}
		}
		cropImg = rgbaImg
		newImg = true
	} else {
		// Return an empty RGBA image
		cropImg = &image.RGBA{}
		newImg = true
	}

	return cropImg, newImg
}

func TestImage() {

	f, err := os.Open("/Users/taoguan/Documents/MyGoProject/images/9R9A8601.jpg")
	if err != nil {
		log.Fatalf("os.Open() failed with %s\n", err)
	}
	defer f.Close()

	img, fmtName, err := image.Decode(f)
	if err != nil {
		log.Fatalf("image.Decode() failed with %s\n", err)
	}

	fmt.Println(fmtName)

	size := img.Bounds().Size()
	fmt.Printf("  size: %dx%d\n", size.X, size.Y)

	b := img.Bounds()
	width, height := b.Dx(), b.Dy()

	fmt.Println(width, height)

	ConvertImageToGrey(img, "/Users/taoguan/Documents/MyGoProject/images/grey.jpg")

}
