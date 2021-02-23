package main

import (
	"flag"
	"fmt"
	"path"
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

		// track1.compress(T, R, m, W)
		// track1.attackComp(T, R, m, W)
		track1.attackComp(T, R, att, rel, 3, 1, W, m)
		track1.writeFile(*outfilename)

		fmt.Printf("Compressed into %s.\n", *outfilename)

	case "convolve":
		file1 := flag.Arg(0)
		fmt.Printf("Convolving...\n")

		track1 := NewWAV()
		track1.readFile(file1)

		// track1.lowpass()
		// track1.highpass()
		// track1.windowedSinc()
		track1.chebyshev()

		track1.writeFile(*outfilename)
	}
}
