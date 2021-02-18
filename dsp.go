package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
)

type wav struct {
	chunkID       [4]byte
	chunkSize     uint32
	format        [4]byte
	subchunk1ID   [4]byte
	subchunk1Size uint32
	audioFormat   uint16
	numChannels   uint16
	sampleRate    uint32
	byteRate      uint32
	blockAlign    uint16
	bitsPerSample uint16
	subchunk2ID   [4]byte
	subchunk2Size uint32
	data          [][]byte
	// Calculated fields
	NumSamples uint32
	SampleSize uint16
	Duration   float64
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func rms(val []float64) float64 {
	var sum float64
	for _, s := range val {
		sum += s
	}
	return sum / float64(len(val))
}

func readWAV(r io.Reader, object *wav) {
	binary.Read(r, binary.BigEndian, &object.chunkID)
	binary.Read(r, binary.LittleEndian, &object.chunkSize)
	binary.Read(r, binary.BigEndian, &object.format)
	binary.Read(r, binary.BigEndian, &object.subchunk1ID)
	binary.Read(r, binary.LittleEndian, &object.subchunk1Size)
	binary.Read(r, binary.LittleEndian, &object.audioFormat)
	binary.Read(r, binary.LittleEndian, &object.numChannels)
	binary.Read(r, binary.LittleEndian, &object.sampleRate)
	binary.Read(r, binary.LittleEndian, &object.byteRate)
	binary.Read(r, binary.LittleEndian, &object.blockAlign)
	binary.Read(r, binary.LittleEndian, &object.bitsPerSample)
	binary.Read(r, binary.BigEndian, &object.subchunk2ID)
	binary.Read(r, binary.LittleEndian, &object.subchunk2Size)
	object.NumSamples = (8 * object.subchunk2Size) / uint32((object.numChannels * object.bitsPerSample))
	object.SampleSize = (object.numChannels * object.bitsPerSample) / 8
	object.Duration = float64(object.subchunk2Size) / float64(object.byteRate)
	for i := 0; i < int(object.NumSamples); i++ {
		sample := make([]byte, int(object.SampleSize))
		binary.Read(r, binary.LittleEndian, &sample)
		object.data = append(object.data, sample)
	}
}

func writeWAV(r io.Writer, object wav) {
	binary.Write(r, binary.BigEndian, object.chunkID)
	binary.Write(r, binary.LittleEndian, object.chunkSize)
	binary.Write(r, binary.BigEndian, object.format)
	binary.Write(r, binary.BigEndian, object.subchunk1ID)
	binary.Write(r, binary.LittleEndian, object.subchunk1Size)
	binary.Write(r, binary.LittleEndian, object.audioFormat)
	binary.Write(r, binary.LittleEndian, object.numChannels)
	binary.Write(r, binary.LittleEndian, object.sampleRate)
	binary.Write(r, binary.LittleEndian, object.byteRate)
	binary.Write(r, binary.LittleEndian, object.blockAlign)
	binary.Write(r, binary.LittleEndian, object.bitsPerSample)
	binary.Write(r, binary.BigEndian, object.subchunk2ID)
	binary.Write(r, binary.LittleEndian, object.subchunk2Size)
	for i := 0; i < len(object.data); i++ {
		binary.Write(r, binary.LittleEndian, object.data[i])
	}
}

func dumpWAVHeader(object wav, more bool) {
	fmt.Printf("File size: %.2fKB\n", float64(object.chunkSize)/1000)
	fmt.Printf("Number of samples: %d\n", object.NumSamples)
	fmt.Printf("Size of each sample: %d bytes\n", object.SampleSize)
	fmt.Printf("Duration of file: %fs\n", object.Duration)
	if more {
		fmt.Printf("%-14s %s\n", "chunkID:", object.chunkID)
		fmt.Printf("%-14s %d\n", "chunkSize:", object.chunkSize)
		fmt.Printf("%-14s %s\n", "format:", object.format)
		fmt.Printf("%-14s %s\n", "subchunk1ID:", object.subchunk1ID)
		fmt.Printf("%-14s %d\n", "subchunk1Size:", object.subchunk1Size)
		fmt.Printf("%-14s %d\n", "audioFormat:", object.audioFormat)
		fmt.Printf("%-14s %d\n", "numChannels:", object.numChannels)
		fmt.Printf("%-14s %d\n", "sampleRate:", object.sampleRate)
		fmt.Printf("%-14s %d\n", "byteRate:", object.byteRate)
		fmt.Printf("%-14s %d\n", "blockAlign:", object.blockAlign)
		fmt.Printf("%-14s %d\n", "bitsPerSample:", object.bitsPerSample)
		fmt.Printf("%-14s %s\n", "subchunk2ID:", object.subchunk2ID)
		fmt.Printf("%-14s %d\n", "subchunk2Size:", object.subchunk2Size)
	}
}

func mix(t1 wav, t2 wav) wav {
	var track, longerTrack, shorterTrack wav
	if t1.NumSamples >= t2.NumSamples {
		longerTrack = t1
		shorterTrack = t2
	} else {
		longerTrack = t2
		shorterTrack = t1
	}
	track = longerTrack
	for i := 0; i < int(longerTrack.NumSamples); i++ {
		signal := make([]byte, track.SampleSize)
		var sample int32
		if i < int(shorterTrack.NumSamples) {
			sample = int32(int16(binary.LittleEndian.Uint16(longerTrack.data[i]))) + int32(int16(binary.LittleEndian.Uint16(shorterTrack.data[i])))
		} else {
			sample = int32(int16(binary.LittleEndian.Uint16(longerTrack.data[i])))
		}
		if sample > 32767 {
			sample = 32767
		} else if sample < -32768 {
			sample = -32768
		}
		binary.LittleEndian.PutUint16(signal, uint16(sample))
		track.data[i] = signal
	}
	return track
}

func normalize(track wav, desiredPeak float64) wav {
	base := math.Pow(2, float64(track.bitsPerSample-1)) * math.Pow(10, (desiredPeak/20))
	var peak float64 = 0
	for i := 0; i < int(track.NumSamples); i++ {
		sample := math.Abs(float64(int16(binary.LittleEndian.Uint16(track.data[i]))))
		if sample > peak {
			peak = sample
		}
	}
	normNum := base / peak
	for i := 0; i < int(track.NumSamples); i++ {
		signal := make([]byte, track.SampleSize)
		sample := float64(int16(binary.LittleEndian.Uint16(track.data[i])))
		sample *= normNum
		binary.LittleEndian.PutUint16(signal, uint16(sample))
		track.data[i] = signal
	}
	return track
}

func compress(track wav, thresh float64, ratio int, makeup float64) wav {
	sigThresh := math.Pow(2, float64(track.bitsPerSample-1)) * math.Pow(10, (thresh/20))
	var period []float64
	for i := 0; i < int(track.NumSamples); i++ {
		signal := make([]byte, track.SampleSize)
		sample := float64(int16(binary.LittleEndian.Uint16(track.data[i])))
		if len(period) == 500 {
			period = period[1:]
		}
		period = append(period, sample)
		sigRMS := rms(period)
		if sigRMS > sigThresh {
			sample = sigThresh + (sample-sigThresh)/float64(ratio)
		} else if sigRMS < -sigThresh {
			sample = -sigThresh + (sample+sigThresh)/float64(ratio)
		}
		binary.LittleEndian.PutUint16(signal, uint16(sample))
		track.data[i] = signal
	}
	if makeup != 1.0 {
		fmt.Printf("Normalizing...\n")
		track = normalize(track, makeup)
	}
	return track
}
