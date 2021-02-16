package observe

import (
	"encoding/base64"
	"fmt"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

const thumb = "thumb"
const small = "small"
const medium = "medium"
const large = "large"
const Raw = "raw"

type ImageHandler struct {
	*Scaler
}

func (h *ImageHandler) SavePng(newFileName string, m *image.Image) error {
	f, err := os.Create(newFileName)
	if err != nil {
		return fmt.Errorf("could not create file %s; %s", newFileName, err)
	}

	encoder := png.Encoder{CompressionLevel: png.BestCompression}
	err = encoder.Encode(f, *m)
	if err != nil {
		return fmt.Errorf("could not create image %s", err)
	}

	return nil
}

func (h *ImageHandler) Base64ToImage(base64image string) (*image.Image, error) {
	// data:webImage/png;base64,iVBORw0KGgoAAAANSUhEUgAAAoAAAAHgCAYAAAA10dzkAAAgAElEQVR4nJTcZ3db55no/XyF86zz [etc]
	// data can be split on the comma ,
	i := strings.Index(base64image, ",")
	if i < 0 {
		i = 0
	}

	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(base64image[i+1:]))
	m, _, err := image.Decode(reader)
	if err != nil {
		return nil, fmt.Errorf("Error decoding base64 string: %s ", err)
	}

	return &m, nil
}

func (h *ImageHandler) GetImgFromPath(path string) (image.Image, error) {
	r, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	img, _, err := image.Decode(r)

	return img, err
}

// Scales an Image file to all sizes and saves it as a PNG.
// If filename is a path, only the basename of the file is used.
func (h *ImageHandler) ScaleImageAllSizes(img image.Image, filename string) error {
	var wg sync.WaitGroup
	sizes := []string{thumb, small, medium, large}

	wg.Add(len(sizes))
	for _, size := range sizes {
		m, err := h.Scale(img, size)
		if err != nil {
			return err
		}
		p, _ := h.GetPath(size)
		b := filepath.Base(filename)
		f := fmt.Sprintf("%s/%s", p, b)
		err = h.SavePng(f, &m)
		if err != nil {
			return err
		}
		wg.Done()
	}
	wg.Wait()

	return nil
}

func (h *ImageHandler) GetPath(size string) (string, error) {
	switch size {
	case Raw:
		return "./static/images/observation/raw", nil
	case thumb:
		return "./static/images/observation/thumb", nil
	case small:
		return "./static/images/observation/small", nil
	case medium:
		return "./static/images/observation/medium", nil
	case large:
		return "./static/images/observation/large", nil
	default:
		return "", fmt.Errorf("unable to create path for unknown size %s", size)
	}
}
func (h *ImageHandler) getFullFilename(size string, filename string) (string, error) {
	p, err := h.GetPath(size)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/%s", p, filename), nil
}
