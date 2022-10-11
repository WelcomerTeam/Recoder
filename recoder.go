package recoder

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"image/gif"
	"io"

	"github.com/savsgio/gotils/strconv"
	"github.com/ultimate-guitar/go-imagequant"
)

// VERSION respects semantic versioning.
const VERSION = "1.0.0"

const (
	ColorsMin = 2
	ColorsMax = 256

	QualityMin int = 0
	QualityMax int = 100

	SpeedSlowest int = 1
	SpeedDefault int = 3
	SpeedFastest int = 10
)

// QuantizationAttributes represents all attributes provided to imagequant.
type QuantizationAttributes struct {
	// Specifies maximum number of colors to use. The default is 256.
	// Instead of setting a fixed limit it's better to use MinQuality and MaxQuality
	MaxColors int
	// Quality is in range 0 (worst) to 100 (best) and values are analoguous to JPEG quality
	// (i.e. 80 is usually good enough). Quantization will attempt to use the lowest number
	// of colors needed to achieve maximum quality. maximum value of 100 is the default and
	// means conversion as good as possible. If it's not possible to convert the image with
	// at least minimum quality (i.e. 256 colors is not enough to meet the minimum quality),
	// then Image.Quantize() will fail. The default minimum is 0 (proceeds regardless of quality).
	//
	// Features dependent on speed:
	// speed 1-5: Noise-sensitive dithering
	// speed 8-10 or if image has more than million colors: Forced posterization
	// speed 1-7 or if minimum quality is set: Quantization error known
	// speed 1-6: Additional quantization techniques
	MinQuality int
	MaxQuality int
	// Higher speed levels disable expensive algorithms and reduce quantization precision. The
	// default speed is 3. Speed 1 gives marginally better quality at significant CPU cost.
	// Speed 10 has usually 5% lower quality, but is 8 times faster than the default. High
	// speeds combined with Quality will use more colors than necessary and will be less likely
	// to meet minimum required quality.
	Speed int
}

// NewQuantizationAttributes returns the default quantization attributes.
func NewQuantizationAttributes() QuantizationAttributes {
	return QuantizationAttributes{
		MaxColors:  ColorsMax,
		MinQuality: QualityMin,
		MaxQuality: QualityMax,
		Speed:      SpeedDefault,
	}
}

// RecodeImage handles GIF re-encoding.
func RecodeImage(r io.Reader, qa QuantizationAttributes) (dst io.Reader, err error) {
	attrs := attributesToImageQuant(qa)
	defer attrs.Release()

	src, err := gif.DecodeAll(r)
	if err != nil {
		return nil, fmt.Errorf("failed to decode gif: %w", err)
	}

	imOverlay := image.NewRGBA(image.Rect(0, 0, src.Config.Width, src.Config.Height))
	bounds := imOverlay.Bounds()

	draw.Draw(imOverlay, bounds, src.Image[0], image.Point{0, 0}, draw.Src)

	out := &gif.GIF{
		Image:           make([]*image.Paletted, len(src.Image)),
		Delay:           src.Delay,
		LoopCount:       src.LoopCount,
		Disposal:        nil,
		Config:          src.Config,
		BackgroundIndex: src.BackgroundIndex,
	}

	out.Config.ColorModel = nil

	for index, frame := range src.Image {
		draw.Draw(imOverlay, bounds, frame, image.Point{0, 0}, draw.Over)

		out.Image[index] = quantizeImage(imOverlay.Pix, bounds, attrs)
		if out.Image[index] == nil {
			return nil, fmt.Errorf("failed to recode image: %w", err)
		}
	}

	buf := &bytes.Buffer{}

	err = gif.EncodeAll(buf, out)
	if err != nil {
		return nil, fmt.Errorf("failed to encode gif: %w", err)
	}

	return buf, nil
}

func attributesToImageQuant(qa QuantizationAttributes) (attrs *imagequant.Attributes) {
	attrs, _ = imagequant.NewAttributes()
	_ = attrs.SetMaxColors(qa.MaxColors)
	_ = attrs.SetQuality(qa.MinQuality, qa.MaxQuality)
	_ = attrs.SetSpeed(qa.Speed)

	return
}

func quantizeImage(pix []uint8, bounds image.Rectangle, attrs *imagequant.Attributes) *image.Paletted {
	qsrc, _ := imagequant.NewImage(attrs, strconv.B2S(pix), bounds.Dx(), bounds.Dy(), 1)
	res, _ := qsrc.Quantize(attrs)

	// WriteRemappedImage returns a list of bytes pointing directly to
	// the palette indexes. We can directly copy over this to the dst.
	rmap, err := res.WriteRemappedImage()
	if err != nil {
		return nil
	}

	dst := &image.Paletted{
		Pix:     rmap,
		Stride:  1 * bounds.Dx(),
		Rect:    bounds,
		Palette: res.GetPalette(),
	}

	qsrc.Release()
	res.Release()

	return dst
}
