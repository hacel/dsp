package dsp

import (
	"fmt"
	"math"
	"testing"

	"github.com/mjibson/go-dsp/fft"
)

const (
	tolerance = 1e-8
)

func float64Equal(a, b float64) bool {
	fmt.Println(a, b, tolerance)
	return math.Abs(a-b) <= tolerance || math.Abs(1-a/b) <= tolerance
}

func floatSliceEqual(a, b []float64) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !float64Equal(a[i], b[i]) {
			return false
		}
	}
	return true
}

func TestDFTReconstruction(t *testing.T) {
	track := NewWav()
	track.data = []float64{1.0, 2.0, 3.0}
	dft := fft.FFTReal(track.data)
	idft := fft.IFFT(dft)
	track.ReconSignal(idft)

	got := track.data
	want := []float64{1.0, 2.0, 3.0}
	if !floatSliceEqual(got, want) {
		t.Errorf("got %f, wanted %f", got, want)
	}
}

func BenchmarkDFTReconstruction(b *testing.B) {
	track := NewWav()
	track.data = []float64{1.0, 2.0, 3.0}
	for i := 0; i < b.N; i++ {
		dft := fft.FFTReal(track.data)
		idft := fft.IFFT(dft)
		track.ReconSignal(idft)
	}
}
