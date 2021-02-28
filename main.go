package main

import (
	"flag"
	"fmt"
	"math"
	"math/cmplx"
	"path"

	"github.com/mjibson/go-dsp/dsputils"
)

func main() {
	operation := flag.String("op", "", "Operation: Mix, Normalize, Compress")
	outfilename := flag.String("o", "out.wav", "Output file name")
	dBFS := flag.Float64("db", -12.0, "dBFS")
	ratio := flag.Float64("R", 2.0, "Compression ratio")
	makeup := flag.Float64("m", 1.0, "Compression makeup")
	kneeWidth := flag.Float64("W", -25.0, "Compression soft knee width")
	att := flag.Float64("att", 10.0, "Compression attack time (ms)")
	rel := flag.Float64("rel", 300.0, "Compression release time (ms)")
	lh := flag.Int("lh", 0, "Low pass: 0, High pass: 1")
	freq := flag.Int("f", 5000, "Filter frequency cut off")
	filter := flag.String("filter", "biquad", "Convolution filter selection (biquad, windowedsinc, average)")
	bandwidth := flag.Int("b", 20, "Filter rolloff for certain filters")
	flag.Parse()

	switch *operation {
	case "mix":
		file1 := flag.Arg(0)
		file2 := flag.Arg(1)
		fmt.Printf("Mixing %s and %s\n", path.Base(file1), path.Base(file2))

		track1 := NewWAV()
		track1.readFile(file1)
		fmt.Printf("---------------\n%s details:\n", path.Base(file1))
		track1.dumpHeader(true)

		track2 := NewWAV()
		track2.readFile(file2)
		fmt.Printf("---------------\n%s details:\n", path.Base(file2))
		track2.dumpHeader(true)

		newTrack := NewWAV()
		newTrack.mix(track1, track2)
		newTrack.writeFile(*outfilename)

		fmt.Printf("Mixed into %s.\n", *outfilename)

	case "normalize":
		file1 := flag.Arg(0)
		desiredPeak := *dBFS
		fmt.Printf("Noramlizing %s to %f dBFS\n", path.Base(file1), desiredPeak)

		track1 := NewWAV()
		track1.readFile(file1)
		fmt.Printf("---------------\n%s details:\n", path.Base(file1))
		track1.dumpHeader(true)

		track1.normalize(desiredPeak)
		track1.writeFile(*outfilename)

		fmt.Printf("Normalized into %s.\n", *outfilename)

	case "compress":
		file1 := flag.Arg(0)
		T := *dBFS
		R := *ratio
		W := *kneeWidth
		m := *makeup
		att := *att
		rel := *rel
		fmt.Printf("Compressing %s up to %f dBFS with %f ratio\n", path.Base(file1), T, R)

		track1 := NewWAV()
		track1.readFile(file1)
		fmt.Printf("---------------\n%s details:\n", path.Base(file1))
		track1.dumpHeader(true)

		track1.compress(T, R, att, rel, 3, 1, W, m)
		track1.writeFile(*outfilename)

		fmt.Printf("Compressed into %s.\n", *outfilename)

	case "convolve":
		file1 := flag.Arg(0)
		filter := *filter
		fc := *freq
		lh := *lh
		M := *bandwidth

		track1 := NewWAV()
		track1.readFile(file1)

		switch filter {
		case "average":
			fmt.Printf("Convolving using rolling average (M=%d)...\n", M)
			track1.rollingAvgLowpass(M)

		case "windowedsinc":
			fmt.Printf("Convolving using Windowed-Sinc (fc=%d, M=%d)...\n", fc, M)
			track1.windowedSinc(fc, M)

		case "biquad":
			fmt.Printf("Convolving using Biquad (fc=%d, lh=%d)...\n", fc, lh)
			track1.biquad(fc, lh)

		case "highpass":
			fmt.Printf("Convolving using highpass...")
			track1.highpass()

		case "cheb":
			// track1.chebyshev()
		}

	case "test":
		file1 := flag.Arg(0)
		track1 := NewWAV()
		track1.readFile(file1)

		sr := int(track1.sampleRate)
		dft := track1.getDFT()
		for i := 0; i < len(dft)/2; i++ {
			r, θ := cmplx.Polar(dft[i])
			θ *= 360.0 / (2 * math.Pi)
			if dsputils.Float64Equal(r, 0) {
				θ = 0 // (When the magnitude is close to 0, the angle is meaningless)
			}
			fsr := float64(sr)
			if r > 8184015.4 {
				fmt.Printf("dft(%f) = %.1f ∠ %.1f° ||| ", (float64(i)*fsr)/float64(len(dft)), r, θ)
				fmt.Printf("%f - %f\n", real(dft[i])/fsr, imag(dft[i])/fsr)
			}
		}

		// sr := 48000
		// dt := 1 / float64(sr)
		// x := func(n float64) float64 {
		// 	wave0 := 10000 * math.Sin(2.0*math.Pi*float64(sr)*n/8.0)
		// 	// wave1 := 5000 * math.Sin(2*math.Pi*float64(n)/4.0+3.0*math.Pi/4.0)
		// 	return wave0
		// }
		// a := make([]float64, track1.NumSamples*2)
		// t := 0.0
		// for i := 0; i < int(track1.NumSamples*2); i++ {
		// 	a[i] = x(t)
		// 	t += dt
		// }
		// fmt.Println(t)
		// X := fft.FFTReal(a)

		// for i := 0; i < len(X)/2; i++ {
		// 	r, θ := cmplx.Polar(X[i])
		// 	θ *= 360.0 / (2 * math.Pi)
		// 	if dsputils.Float64Equal(r, 0) {
		// 		θ = 0 // (When the magnitude is close to 0, the angle is meaningless)
		// 	}
		// 	fsr := float64(sr)
		// 	if r > 496295833.7 {
		// 		fmt.Printf("dft(%f) = %.1f ∠ %.1f° ||| ", (float64(i)*fsr)/float64(len(X)), r, θ)
		// 		fmt.Printf("%f - %f\n", real(X[i])/fsr, imag(X[i])/fsr)
		// 	}
		// }
		// track1.sampleRate = 48000.0
		// // track1.sam
		// track1.reconSignal(getIDFT(X))
		// track1.dumpHeader(true)
		// track1.subchunk2Size = uint32(len(track1.data) * (int(track1.bitsPerSample) / 8))
		// track1.chunkSize = track1.subchunk2Size + 36

		///////TEST//////////////

		// track1.writeFile(*outfilename)

	}
}
