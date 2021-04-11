package main

import (
	"bytes"
	"flag"
	"image"
	"image/draw"
	"image/gif"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/savsgio/gotils/strconv"
	"github.com/ultimate-guitar/go-imagequant"
	"gopkg.in/fsnotify.v1"
)

// VERSION respects semantic versioning.
const VERSION = "0.1+110420210202"

// Extention files must end with to be detected by recoder
const matchPath = ".recode"

var (
	attr, _ = imagequant.NewAttributes()

	sWatchFolder     *string
	sResultDirectory *string
)

func main() {
	sWatchFolder = flag.String("watch", "", "Folder to watch for files ending with .recode")
	sResultDirectory = flag.String("resulting_directory", "", "Folder where to output the files. Leave default to use same folder as watch")
	sSpeed := flag.Int("speed", 3, "Imagequant quantization speed. Speed 1 gives marginally better quality at significant CPU cost. Speed 10 has usually 5% lower quality, but is 8 times faster than the default")
	sMinQuality := flag.Int("min_quality", 0, "Minimum quality for Imagequant")
	sMaxQuality := flag.Int("max_quality", 100, "Maximum quality for Imagequant")
	flag.Parse()

	if *sWatchFolder == "" {
		println("missing required argument: watch")
		os.Exit(2)
	}

	if *sResultDirectory == "" {
		*sResultDirectory = *sWatchFolder
	}

	attr.SetSpeed(*sSpeed)
	attr.SetQuality(*sMinQuality, *sMaxQuality)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		panic(err)
	}
	defer watcher.Close()

	err = watcher.Add(*sWatchFolder)
	if err != nil {
		panic(err)
	}

	println("Watching", *sWatchFolder, "and serving to", *sResultDirectory, "...")

	files, err := ioutil.ReadDir(*sWatchFolder)
	if err != nil {
		panic(err)
	}

	// We will check for any files with .recode incase
	// any existed during downtime.

	recoded := make([]string, 0)
	for _, file := range files {
		name := file.Name()
		if strings.HasSuffix(name, matchPath) {
			recoded = append(recoded, name)
		}
	}

	if len(recoded) > 0 {
		println("Found", len(recoded), "files to recode")
		for _, file := range recoded {
			go recodeFile(path.Join(*sWatchFolder, file))
		}
	}

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if event.Op&fsnotify.Create == fsnotify.Create {
				if strings.HasSuffix(event.Name, matchPath) {
					go recodeFile(event.Name)
				}
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Println("error:", err)
		}
	}
}

// recodeFile opens a file that ends with .recode and will attempt
// to parse it as a GIF. It it cannot, it will remove the .record and
// leave the file alone. If it can parse it it will then attempt
// to remove any Disposal the GIF may have had. Once this is done, it
// will then quantize the new frames then save the new file and remove
// the .recode
func recodeFile(path string) {
	var err error
	var outputPath string

	start := time.Now()
	outputPath = filepath.Join(
		*sResultDirectory,
		filepath.Base(strings.ReplaceAll(path, matchPath, "")),
	)

	defer func() {
		println("Recoded " + path + " in " + time.Since(start).String())
		if r := recover(); r != nil {
			println("! recovered during recoding of file", path)
		}
		if err != nil {
			println("! " + err.Error())
		}
	}()

	f, err := os.Open(path)
	if err != nil {
		return
	}

	src, err := gif.DecodeAll(f)
	if err != nil {
		return
	}

	f.Seek(0, io.SeekStart)

	config, err := gif.DecodeConfig(f)
	if err != nil {
		return
	}

	overpaintImage := image.NewRGBA(image.Rect(0, 0, config.Width, config.Height))
	draw.Draw(overpaintImage, overpaintImage.Bounds(), src.Image[0], image.Point{}, draw.Src)

	for i, frame := range src.Image {
		draw.Draw(overpaintImage, overpaintImage.Bounds(), frame, image.Point{}, draw.Over)

		quant, err := quantizeImage(overpaintImage)
		if err != nil {
			return
		}

		src.Image[i] = quant
	}

	src.Disposal = nil
	src.Config.ColorModel = nil

	var b bytes.Buffer

	err = gif.EncodeAll(&b, src)
	if err != nil {
		return
	}

	f, err = os.Create(outputPath)
	if err != nil {
		return
	}

	_, err = b.WriteTo(f)
	if err != nil {
		return
	}

	err = f.Close()
	if err != nil {
		return
	}

	err = os.Remove(path)
	if err != nil {
		return
	}
}

// quantizeImage converts an image.Image to image.Paletted via imagequant
func quantizeImage(src image.Image) (*image.Paletted, error) {
	b := src.Bounds()

	qimg, err := imagequant.NewImage(attr, strconv.B2S(imagequant.ImageToRgba32(src)), b.Dx(), b.Dy(), 1)
	if err != nil {
		panic(err)
	}

	pm, err := qimg.Quantize(attr)
	if err != nil {
		panic(err)
	}

	dst := image.NewPaletted(src.Bounds(), pm.GetPalette())

	// WriteRemappedImage returns a list of bytes pointing to direct
	// palette indexes so we can just copy it over and it will be
	// using the optimimal indexes.
	rmap, err := pm.WriteRemappedImage()
	if err != nil {
		return dst, err
	}

	dst.Pix = rmap

	pm.Release()
	qimg.Release()

	return dst, nil
}
