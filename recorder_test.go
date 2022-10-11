package recoder_test

import (
	"fmt"
	"io"
	"os"
	"testing"

	recoder "github.com/WelcomerTeam/Recoder"
)

const (
	testGIFLargeLocation      = "test_large.gif"
	testGIFMediumLocation     = "test_medium.gif"
	testGIFSmallLocation      = "test_small.gif"
	testGIFExtraSmallLocation = "test_extrasmall.gif"
)

func runBenchmark(b *testing.B, attr recoder.QuantizationAttributes, loc string) {
	r, err := os.OpenFile(loc, os.O_RDONLY, 0444)
	if err != nil {
		fmt.Println(fmt.Errorf("os.OpenFile(%s): %w", loc, err))

		return
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err = recoder.RecodeImage(r, attr)

		if err != nil {
			b.Fatal(err)
		}

		b.StopTimer()

		_, _ = r.Seek(0, io.SeekStart)

		b.StartTimer()
	}
}

func fastAttrs() recoder.QuantizationAttributes {
	attrs := recoder.NewQuantizationAttributes()
	attrs.Speed = recoder.SpeedFastest

	return attrs
}

func defaultAttrs() recoder.QuantizationAttributes {
	attrs := recoder.NewQuantizationAttributes()
	attrs.Speed = recoder.SpeedDefault

	return attrs
}
func slowAttrs() recoder.QuantizationAttributes {
	attrs := recoder.NewQuantizationAttributes()
	attrs.Speed = recoder.SpeedSlowest

	return attrs
}

// Large Benchmark.
func BenchmarkRecoderLargeFast(b *testing.B) { runBenchmark(b, fastAttrs(), testGIFLargeLocation) }
func BenchmarkRecoderLargeDefault(b *testing.B) {
	runBenchmark(b, defaultAttrs(), testGIFLargeLocation)
}
func BenchmarkRecoderLargeSlow(b *testing.B) { runBenchmark(b, slowAttrs(), testGIFLargeLocation) }

// Medium Benchmark.
func BenchmarkRecoderMediumFast(b *testing.B) { runBenchmark(b, fastAttrs(), testGIFMediumLocation) }
func BenchmarkRecoderMediumDefault(b *testing.B) {
	runBenchmark(b, defaultAttrs(), testGIFMediumLocation)
}
func BenchmarkMRecoderMediumSlow(b *testing.B) { runBenchmark(b, slowAttrs(), testGIFMediumLocation) }

// Small Benchmark.
func BenchmarkRecoderSmallFast(b *testing.B) { runBenchmark(b, fastAttrs(), testGIFSmallLocation) }
func BenchmarkRecoderSmallDefault(b *testing.B) {
	runBenchmark(b, defaultAttrs(), testGIFSmallLocation)
}
func BenchmarkRecoderSmallSlow(b *testing.B) { runBenchmark(b, slowAttrs(), testGIFSmallLocation) }

// ExtraSmall Benchmark.
func BenchmarkRecoderExtraSmallFast(b *testing.B) {
	runBenchmark(b, fastAttrs(), testGIFExtraSmallLocation)
}
func BenchmarkRecoderExtraSmallDefault(b *testing.B) {
	runBenchmark(b, defaultAttrs(), testGIFExtraSmallLocation)
}
func BenchmarkExtraRecoderExtraSmallSlow(b *testing.B) {
	runBenchmark(b, slowAttrs(), testGIFExtraSmallLocation)
}
