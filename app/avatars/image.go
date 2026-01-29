package avatars

// ------------------------------------------------------------------------------------------------
// Imports
// ------------------------------------------------------------------------------------------------

import (
	"bytes"
	"io"

	"github.com/disintegration/imaging"
)

// ------------------------------------------------------------------------------------------------
// Functions
// ------------------------------------------------------------------------------------------------

func ResizeImageToAvatar(file io.Reader) ([]byte, error) {
	img, err := imaging.Decode(file, imaging.AutoOrientation(true))
	if err != nil {
		return nil, err
	}

	resized := imaging.Fill(img, 256, 256, imaging.Center, imaging.Lanczos)

	buf := new(bytes.Buffer)
	err = imaging.Encode(buf, resized, imaging.JPEG, imaging.JPEGQuality(85))
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// ------------------------------------------------------------------------------------------------
