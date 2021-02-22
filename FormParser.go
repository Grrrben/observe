package observe

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// max upload of 10 MB files, increase if needed...
const maxFileSize = 10 << 20

type FormParser struct {
	r *http.Request
}

func NewFormParser(r *http.Request) *FormParser {
	fp := new(FormParser)
	fp.r = r

	return fp
}

func (fp *FormParser) ParseFieldAsString(key string) (t time.Time, err error) {
	dt := fp.r.FormValue(key)
	if dt == "" {
		return t, errors.New("required date missing")
	}

	format := "02-01-2006 15:04"
	t, err = time.Parse(format, dt)
	if err != nil {
		return t, errors.New("required date has invalid format; should be [dd-mm-yyyy hh:mm]")
	}

	return
}

func (fp *FormParser) GetFileFromField(key string) (b []byte, err error) {
	err = fp.r.ParseMultipartForm(maxFileSize)
	if err != nil {
		return b, fmt.Errorf("unable to parse/upload image (max 10MB); %s", err)
	}

	f, _, err := fp.r.FormFile(key)
	defer f.Close()
	if err != nil {
		return b, fmt.Errorf("error retrieving the image file; %s", err)
	}

	b, err = ioutil.ReadAll(f)
	if err != nil {
		return b, fmt.Errorf("error retrieving the image file; %s", err)
	}

	return
}
