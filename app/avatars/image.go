package avatars

import (
	"bytes"
	"io"

	"github.com/disintegration/imaging"
)

func ResizeImageToAvatar(file io.Reader) ([]byte, error) {
	img, err := imaging.Decode(file)
	if err != nil {
		return nil, err
	}

	resized := imaging.Resize(img, 256, 256, imaging.Lanczos)

	buf := new(bytes.Buffer)
	if err := imaging.Encode(buf, resized, imaging.JPEG); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
