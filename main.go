package main

import (
	"flag"
	"fmt"
	"os"
	"path"
)

func main() {
	operation := flag.String("op", "", "Operation: Mix, Normalize, Compress")
	outfilename := flag.String("o", "out.wav", "Output file name")
	dBFS := flag.Float64("db", 0.0, "dBFS")
	ratio := flag.Int("R", 2, "Compression ratio")
	makeup := flag.Float64("m", 1.0, "Compression makeup")
	kneeWidth := flag.Float64("W", 0.0, "Compression soft knee width")
	flag.Parse()
	fmt.Println(*outfilename)

	switch *operation {
	case "mix":
		file1dir, file1name := path.Split(flag.Arg(0))
		file2dir, file2name := path.Split(flag.Arg(1))
		fmt.Printf("Mixing %s and %s\n", file1name, file2name)

		f, err := os.Open(file1dir + file1name)
		check(err)
		var track1 wav
		readWAV(f, &track1)
		fmt.Printf("---------------\n%s details:\n", file1name)
		dumpWAVHeader(track1, true)

		f, err = os.Open(file2dir + file2name)
		var track2 wav
		readWAV(f, &track2)
		fmt.Printf("---------------\n%s details:\n", file2name)
		dumpWAVHeader(track2, false)

		new := mix(track1, track2)
		outfile, err := os.Create(*outfilename)
		check(err)
		writeWAV(outfile, new)
		fmt.Printf("Mixed into %s.\n", *outfilename)
		f.Close()

	case "normalize":
		file1dir, file1name := path.Split(flag.Arg(0))
		desiredPeak := *dBFS
		fmt.Printf("Noramlizing %s to %f dBFS\n", file1name, desiredPeak)

		f, err := os.Open(file1dir + file1name)
		check(err)
		var track1 wav
		readWAV(f, &track1)
		fmt.Printf("---------------\n%s details:\n", file1name)
		dumpWAVHeader(track1, true)

		new := normalize(track1, desiredPeak)
		outfile, err := os.Create(*outfilename)
		check(err)
		writeWAV(outfile, new)
		fmt.Printf("Normalized into %s.\n", *outfilename)

		f.Close()

	case "compress":
		file1dir, file1name := path.Split(flag.Arg(0))
		T := *dBFS
		R := *ratio
		W := *kneeWidth
		m := *makeup
		fmt.Printf("Compressing %s up to %f dBFS with %d ratio\n", file1name, T, R)

		f, err := os.Open(file1dir + file1name)
		check(err)
		var track1 wav
		readWAV(f, &track1)
		fmt.Printf("---------------\n%s details:\n", file1name)
		dumpWAVHeader(track1, true)

		new := compress(track1, T, R, m, W)
		outfile, err := os.Create(*outfilename)
		check(err)
		writeWAV(outfile, new)
		fmt.Printf("Compressed into %s.\n", *outfilename)

		f.Close()
	}
}
