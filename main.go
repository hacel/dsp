package main

import (
	"flag"
	"fmt"
	"os"
	"path"
)

var (
	mixCmd    = flag.NewFlagSet("mix", flag.ExitOnError)
	mFilename = mixCmd.String("o", "out.wav", "Output file name")

	normalizeCmd = flag.NewFlagSet("normalize", flag.ExitOnError)
	nFilename    = normalizeCmd.String("o", "out.wav", "Output file name")
	ndBFS        = normalizeCmd.Float64("db", -1.0, "dBFS")

	compressCmd = flag.NewFlagSet("compress", flag.ExitOnError)
	cFilename   = compressCmd.String("o", "out.wav", "Output file name")
	cdBFS       = compressCmd.Float64("db", -12.0, "dBFS")
	cRatio      = compressCmd.Float64("R", 2.0, "Compression ratio")
	cMakeup     = compressCmd.Float64("m", 1.0, "Compression makeup")
	cKneeWidth  = compressCmd.Float64("W", -25.0, "Compression soft knee width")
	cAtt        = compressCmd.Float64("att", 10.0, "Compression attack time (ms)")
	cRel        = compressCmd.Float64("rel", 300.0, "Compression release time (ms)")

	convolveCmd = flag.NewFlagSet("convolve", flag.ExitOnError)
	vFilename   = convolveCmd.String("o", "out.wav", "Output file name")
	vFilter     = convolveCmd.String("filter", "biquad", "Convolution filter selection (biquad, windowedsinc, average)")
	vLh         = convolveCmd.Int("lh", 0, "Low pass: 0, High pass: 1")
	vFreq       = convolveCmd.Int("f", 5000, "Filter frequency cut off")
	vBandwidth  = convolveCmd.Int("b", 20, "Filter rolloff for certain filters")
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage:\n\t. (mix, normalize, compress, convolve) [options] <file(s)>")
		return
	}

	switch os.Args[1] {
	case "mix":
		mixCmd.Parse(os.Args[2:])
		file1 := mixCmd.Arg(0)
		file2 := mixCmd.Arg(1)
		if file1 == "" || file2 == "" {
			fmt.Println("Specify two files to be mixed.")
			return
		}
		fmt.Printf("Mixing %s and %s\n", path.Base(file1), path.Base(file2))

		track1 := NewWAV()
		track1.ReadFile(file1)
		fmt.Printf("---------------\n%s details:\n", path.Base(file1))
		track1.dumpHeader(false)

		track2 := NewWAV()
		track2.ReadFile(file2)
		fmt.Printf("---------------\n%s details:\n", path.Base(file2))
		track2.dumpHeader(false)

		newTrack := NewWAV()
		newTrack.mix(track1, track2)
		newTrack.WriteFile(*mFilename)

		fmt.Printf("Mixed into %s.\n", *mFilename)

	case "normalize":
		normalizeCmd.Parse(os.Args[2:])
		file1 := normalizeCmd.Arg(0)
		if file1 == "" {
			fmt.Println("Specify a file to be normalized.")
			return
		}
		fmt.Printf("Noramlizing %s to %f dBFS\n", path.Base(file1), *ndBFS)

		track1 := NewWAV()
		track1.ReadFile(file1)
		fmt.Printf("---------------\n%s details:\n", path.Base(file1))
		track1.dumpHeader(false)

		track1.normalize(*ndBFS)
		track1.WriteFile(*nFilename)

		fmt.Printf("Normalized into %s.\n", *nFilename)

	case "compress":
		compressCmd.Parse(os.Args[2:])
		file1 := compressCmd.Arg(0)
		if file1 == "" {
			fmt.Println("Specify a file to be compressed.")
			return
		}
		fmt.Printf("Compressing %s up to %f dBFS with %f ratio\n", path.Base(file1), *cdBFS, *cRatio)

		track1 := NewWAV()
		track1.ReadFile(file1)
		fmt.Printf("---------------\n%s details:\n", path.Base(file1))
		track1.dumpHeader(false)

		track1.compress(*cdBFS, *cRatio, *cAtt, *cRel, 10, *cKneeWidth, *cMakeup)
		track1.WriteFile(*cFilename)

		fmt.Printf("Compressed into %s.\n", *cFilename)

	case "convolve":
		convolveCmd.Parse(os.Args[2:])
		file1 := convolveCmd.Arg(0)
		if file1 == "" {
			fmt.Println("Specify a file to be convolved.")
			return
		}

		track1 := NewWAV()
		track1.ReadFile(file1)
		fmt.Printf("---------------\n%s details:\n", path.Base(file1))
		track1.dumpHeader(false)

		switch *vFilter {
		case "average":
			fmt.Printf("Convolving using rolling average (M=%d)...\n", *vBandwidth)
			track1.rollingAvgLowpass(*vBandwidth)

		case "windowedsinc":
			fmt.Printf("Convolving using Windowed-Sinc (fc=%d, M=%d)...\n", *vFreq, *vBandwidth)
			track1.windowedSinc(*vFreq, *vBandwidth)

		case "biquad":
			fmt.Printf("Convolving using Biquad (fc=%d, lh=%d)...\n", *vFreq, *vLh)
			track1.biquad(*vFreq, *vLh)

		case "highpass":
			fmt.Printf("Convolving using highpass...")
			track1.highpass()

		case "cheb":
			// track1.chebyshev()
		default:
			fmt.Println("Please enter a valid filter for convolution")
			return
		}
		track1.WriteFile(*vFilename)
	}
}
