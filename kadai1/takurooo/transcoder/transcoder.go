package transcoder

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const (
	JPEG string = "jpeg"
	JPG  string = "jpg"
	PNG  string = "png"
)

type Transcoder struct {
	inFormat  string
	outFormat string
}

func IsSupported(format string) bool {
	switch format {
	case JPEG:
	case JPG:
	case PNG:
	default:
		return false
	}
	return true
}

func NewTranscoder(inFormat, outFormat string) *Transcoder {
	return &Transcoder{inFormat, outFormat}
}

func (t *Transcoder) decode(r io.Reader, format string) (img image.Image, err error) {
	switch format {
	case PNG:
		img, err = png.Decode(r)
	case JPEG, JPG:
		img, err = jpeg.Decode(r)
	default:
		return nil, fmt.Errorf("transcoder.decode: invalid format %s", format)
	}
	if err != nil {
		return nil, err
	}
	return img, nil
}

func (t *Transcoder) encode(w io.Writer, img image.Image, format string) (err error) {
	switch format {
	case PNG:
		err = png.Encode(w, img)
	case JPEG, JPG:
		err = jpeg.Encode(w, img, nil)
	default:
		return fmt.Errorf("transcoder.encode: invalid format %s", format)
	}
	if err != nil {
		return err
	}
	return nil
}

// CanTrans はパス拡張子からトランスコード可能な入力フォーマットか判定する.
func (t *Transcoder) CanTrans(path string) bool {
	ext := filepath.Ext(path)
	format := strings.TrimLeft(ext, ".")

	if t.inFormat == JPG || t.inFormat == JPEG {
		return format == JPG || format == JPEG
	}

	return format == t.inFormat
}

// Do 指定されたパスのファイルを指定されたフォーマットにトランスコードする.
// トランスコードされたファイルは入力ファイルと同じ階層に拡張子が変換された状態で保存される。
func (t *Transcoder) Do(inPath string) (err error) {
	// -------------------------
	// decode phase
	// -------------------------
	rf, err := os.Open(inPath)
	if err != nil {
		return err
	}
	defer func() {
		if derr := rf.Close(); derr != nil {
			err = fmt.Errorf("transcoder.Do: read file close err: %v, Do err: %v", derr, err)
		}
	}()

	var img image.Image

	img, err = t.decode(rf, t.inFormat)
	if err != nil {
		return err
	}

	// -------------------------
	// encode phase
	// -------------------------
	dir := filepath.Dir(inPath)
	base := filepath.Base(strings.TrimSuffix(inPath, filepath.Ext(inPath)))
	outPath := filepath.Join(dir, base+"."+t.outFormat)

	wf, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer func() {
		if derr := wf.Close(); derr != nil {
			err = fmt.Errorf("transcoder.Do: write file close err: %v, Do err: %v", derr, err)
		}
	}()

	err = t.encode(wf, img, t.outFormat)
	if err != nil {
		return err
	}

	return nil
}
