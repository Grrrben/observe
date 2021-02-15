package observe

import (
	"fmt"
	"github.com/nfnt/resize"
	"image"
	_ "image/png"
)

const heightThumb = uint(64)
const heightSmall = uint(128)
const heightMedium = uint(512)
const heightLarge = uint(2048)

type Scaler struct{}

func (s *Scaler) Scale(img image.Image, size string) (image.Image, error) {
	var height uint

	switch size {
	case thumb:
		height = heightThumb
	case small:
		height = heightSmall
	case medium:
		height = heightMedium
	case large:
		height = heightLarge
	default:
		return nil, fmt.Errorf("unknown size %s", size)
	}
	mw := uint(img.Bounds().Max.X)

	return resize.Thumbnail(mw, height, img, resize.Lanczos3), nil
}
