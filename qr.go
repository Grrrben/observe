package observe

import (
	"errors"
	"fmt"
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"image/png"
	"os"
)

type Qr struct{}

func (q *Qr) Create(subject string) (barcode.Barcode, error) {
	qrCode, err := qr.Encode(subject, qr.L, qr.Auto)
	if err != nil {
		return nil, fmt.Errorf("could not create QR code for string %s", subject)
	}
	bc, err := barcode.Scale(qrCode, 512, 512)
	if err != nil {
		return nil, errors.New("could not create barcode")
	}

	return bc, nil
}

func (q *Qr) Save(filename string, qrCode barcode.Barcode) error {
	qrFile, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("could not create file %s", filename)
	}

	encoder := png.Encoder{CompressionLevel: png.BestCompression}
	err = encoder.Encode(qrFile, qrCode)
	if err != nil {
		return fmt.Errorf("could not create QR code; %s", err)
	}

	return nil
}
