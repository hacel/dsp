package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"strconv"
)

func main() {
	outfilename := flag.String("o", "out.wav", "usage")
	flag.Parse()

	operation := flag.Arg(0)
	switch operation {
	case "mix":
		file1dir, file1name := path.Split(flag.Arg(1))
		file2dir, file2name := path.Split(flag.Arg(2))
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
		file1dir, file1name := path.Split(flag.Arg(1))
		dBFS, _ := strconv.Atoi(flag.Arg(2))
		fmt.Printf("Noramlizing %s to %d dBFS\n", file1name, dBFS)

		f, err := os.Open(file1dir + file1name)
		check(err)
		var track1 wav
		readWAV(f, &track1)
		fmt.Printf("---------------\n%s details:\n", file1name)
		dumpWAVHeader(track1, true)

		new := normalize(track1, float64(dBFS))
		outfile, err := os.Create(*outfilename)
		check(err)
		writeWAV(outfile, new)
		fmt.Printf("Normalized into %s.\n", *outfilename)

		f.Close()
	case "compress":
		file1dir, file1name := path.Split(flag.Arg(1))
		dBFS, _ := strconv.Atoi(flag.Arg(2))
		ratio, _ := strconv.Atoi(flag.Arg(3))
		fmt.Printf("Compressing %s up to %d dBFS with %d ratio\n", file1name, dBFS, ratio)

		f, err := os.Open(file1dir + file1name)
		check(err)
		var track1 wav
		readWAV(f, &track1)
		fmt.Printf("---------------\n%s details:\n", file1name)
		dumpWAVHeader(track1, true)

		new := compress(track1, float64(dBFS), ratio)
		outfile, err := os.Create(*outfilename)
		check(err)
		writeWAV(outfile, new)
		fmt.Printf("Compressed into %s.\n", *outfilename)

		f.Close()
	}
}
