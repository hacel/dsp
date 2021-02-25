package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"os"

	"github.com/mjibson/go-dsp/fft"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func avg(val []float64) float64 {
	var sum float64
	for _, s := range val {
		sum += s
	}
	return sum / float64(len(val))
}

func rms(val []float64) float64 { // THIS IS AVG?
	var sum float64
	for _, s := range val {
		sum += math.Pow(s, 2)
	}
	return math.Sqrt(sum / float64(len(val)))
}

// WAV is a struct to hold wave file format data
type WAV struct {
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
	data          []float64
	// Calculated fields
	NumSamples uint32
	SampleSize uint16
	Duration   float64
}

// NewWAV creates new WAV object and returns a pointer to it
func NewWAV() *WAV {
	return &WAV{}
}

func (object *WAV) read(r io.Reader) {
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
		x := make([]byte, int(object.SampleSize))
		binary.Read(r, binary.LittleEndian, &x)
		smp := float64(int16(binary.LittleEndian.Uint16(x)))
		object.data = append(object.data, smp)
	}
}

func (object *WAV) readFile(path string) {
	f, err := os.Open(path)
	check(err)
	defer f.Close()
	object.read(f)
}

func (object *WAV) write(r io.Writer) {
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
		signal := make([]byte, object.SampleSize)
		binary.LittleEndian.PutUint16(signal, uint16(object.data[i]))
		binary.Write(r, binary.LittleEndian, signal)
	}
}

func (object *WAV) writeFile(path string) {
	f, err := os.Create(path)
	check(err)
	defer f.Close()
	object.write(f)
}

func (object *WAV) dumpHeader(more bool) {
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

func (object *WAV) getDFT() []complex128 {
	return fft.FFTReal(object.data)
}

func getIDFT(dft []complex128) []complex128 {
	return fft.IFFT(dft)
}

func (object *WAV) reconSignal(idft []complex128) {
	for i := range object.data {
		object.data[i] = real(idft[i])
	}
}

func (object *WAV) mix(t1 *WAV, t2 *WAV) {
	var longerTrack, shorterTrack *WAV
	if t1.NumSamples >= t2.NumSamples {
		longerTrack = t1
		shorterTrack = t2
	} else {
		longerTrack = t2
		shorterTrack = t1
	}
	*object = *longerTrack
	for i := 0; i < int(longerTrack.NumSamples); i++ {
		var x float64
		if i < int(shorterTrack.NumSamples) {
			x = longerTrack.data[i] + shorterTrack.data[i]
		} else {
			x = longerTrack.data[i]
		}
		if x > 32767 {
			x = 32767
		} else if x < -32768 {
			x = -32768
		}
		object.data[i] = x
	}
}

func (object *WAV) normalize(desiredPeak float64) {
	base := math.Pow(2, float64(object.bitsPerSample-1)) * math.Pow(10, (desiredPeak/20))
	var peak float64 = 0
	for i := 0; i < int(object.NumSamples); i++ {
		x := math.Abs(object.data[i])
		if x > peak {
			peak = x
		}
	}
	normNum := base / peak
	for i := 0; i < int(object.NumSamples); i++ {
		x := object.data[i]
		x *= normNum
		object.data[i] = x
	}
}

func (object *WAV) compress(threshold, ratio, tatt, trel, tla, twnd, W, makeup float64) {
	threshold = math.Pow(2, float64(object.bitsPerSample-1)) * math.Pow(10, threshold/20)
	sr := float64(object.sampleRate)
	tatt *= math.Pow(10, -3) // attack time
	trel *= math.Pow(10, -3) // release time
	tla *= math.Pow(10, -3)  // rms lookahead offset
	// twnd *= math.Pow(10, -3) // rms window size
	W = math.Pow(2, float64(object.bitsPerSample-1)) * math.Pow(10, (W/20))
	var att, rel float64
	if tatt == 0 {
		att = 0.0
	} else {
		att = math.Exp(-1.0 / (sr * tatt))
	}
	if trel == 0 {
		rel = 0.0
	} else {
		rel = math.Exp(-1.0 / (sr * trel))
	}
	env := 0.0
	// lhsmp := sr * tla // sample offset in lookahead
	// nrms := sr * twnd // sample count in window
	nla := sr * tla

	for i := 0; i < int(object.NumSamples); i++ {
		summ := 0.0
		// for j := 0; j < int(nrms); j++ {
		// 	lki := i + j + int(lhsmp)
		// 	var smp float64
		// 	if lki >= len(object.data) {
		// 		smp = 0.0
		// 	} else {
		// 		smp = object.data[lki]
		// 	}
		// 	summ += math.Pow(smp, 2)
		// }
		// rms := math.Sqrt(summ / nrms)

		for j := 0; j < int(nla); j++ {
			var smp float64
			if i+j >= len(object.data) {
				smp = 0.0
			} else {
				smp = object.data[i+j]
			}
			summ += smp
		}

		peak := summ / nla
		var theta float64
		// if rms > env {
		if peak > env {
			theta = att
		} else {
			theta = rel
		}
		// env = ((1.0-theta)*rms + theta*env)
		env = ((1.0-theta)*peak + theta*env)

		var gain float64
		// if env > threshold {
		// 	gain = (threshold + (env-threshold)/ratio) / env
		// } else {
		// 	gain = 1.0
		// }
		if env-threshold < -W/2 {
			gain = 1.0
		} else if math.Abs(env-threshold) <= W/2 {
			// env = env + ((1/ratio-1)*math.Pow(env-threshold+W/2, 2))/(W*2)
			gain = (env + ((1/ratio-1)*math.Pow(env-threshold+W/2, 2))/(W*2)) / env
		} else if env-threshold > W/2 {
			// env = threshold + (env-threshold)/ratio
			gain = (threshold + (env-threshold)/ratio) / env
		}

		x := object.data[i]
		x *= gain

		object.data[i] = x
	}
	if makeup != 1.0 {
		fmt.Printf("Normalizing...\n")
		object.normalize(makeup)
	}
}

func (object *WAV) rollingAvgLowpass(bandwidth int) {
	var period []float64
	for i := 0; i < int(object.NumSamples)-5; i++ {
		x := object.data[i]
		// period = append(period, v)
		// if len(period) == 5 {
		// 	period = period[1:]
		// }
		// period = append(period, object.data[i+5])

		if len(period) == bandwidth {
			period = period[1:]
		}
		period = append(period, x)
		avg := avg(period)
		object.data[i] = avg
	}
}

func (object *WAV) biquad(fc, lh int) {
	r := math.Sqrt(2) // Rez
	sr := float64(object.sampleRate)
	var c, a1, a2, a3, b1, b2 float64
	if lh == 0 { // Low pass
		c = 1.0 / math.Tan(math.Pi*float64(fc)/sr)
		a1 = 1.0 / (1.0 + r*c + c*c)
		a2 = 2 * a1
		a3 = a1
		b1 = 2.0 * (1.0 - c*c) * a1
		b2 = (1.0 - r*c + c*c) * a1
	} else { // High pass
		c = math.Tan(math.Pi * float64(fc) / sr)
		a1 = 1.0 / (1.0 + r*c + c*c)
		a2 = -2 * a1
		a3 = a1
		b1 = 2.0 * (c*c - 1.0) * a1
		b2 = (1.0 - r*c + c*c) * a1
	}
	var period []float64
	for i := 0; i < int(object.NumSamples); i++ {
		y := 0.0
		x0 := object.data[i]
		if len(period) == 2 {
			x1 := period[1]
			x2 := period[0]
			y1 := object.data[i-1]
			y2 := object.data[i-2]
			y = a1*x0 + a2*x1 + a3*x2 - b1*y1 - b2*y2
			period = period[1:]
		} else {
			y = x0
		}
		period = append(period, x0)
		object.data[i] = y
	}
}

func (object *WAV) windowedSinc(cutoff, bandwidth int) {
	if cutoff > int(object.sampleRate)/2 {
		panic("Cutoff frequency too high.")
	}
	FC := float64(cutoff) / float64(object.sampleRate) // Cut off (freq/sample rate)
	M := bandwidth                                     // Filter roll off
	kernel := make([]float64, M)
	for i := range kernel {
		if i-M/2 == 0 {
			kernel[i] = 2 * math.Pi * FC
		} else {
			kernel[i] = math.Sin(2*math.Pi*FC*float64(i-M/2)) / float64(i-M/2)
		}
		kernel[i] = kernel[i] * (0.54 - 0.46*math.Cos(2*math.Pi*float64(i)/float64(M)))
	}
	var sum float64 = 0
	for i := range kernel {
		sum += kernel[i]
	}
	for i := range kernel {
		kernel[i] /= sum
	}
	var filteredData []float64
	for j := M; j < int(object.NumSamples); j++ {
		y := 0.0
		x := 0.0
		for i := range kernel {
			x = object.data[j-i]
			y += x * kernel[i]
		}
		filteredData = append(filteredData, y)
	}
	object.data = filteredData
}

func (object *WAV) highpass() {
	var period []float64
	for i := 0; i < int(object.NumSamples); i++ {
		var y float64
		x := object.data[i]
		if len(period) == 2 {
			y = period[0] + -2*period[1] + x
			period = period[1:]
		}
		period = append(period, x)
		object.data[i] = y
	}
}

func _cheb(FC, PR, LH, NP, P float64) (float64, float64, float64, float64, float64) {
	RP := -math.Cos(math.Pi/(NP*2) + (P-1)*math.Pi/NP)
	IP := math.Sin(math.Pi/(NP*2) + (P-1)*math.Pi/NP)
	if PR != 0 {
		ES := math.Sqrt(math.Pow((100/(100-PR)), 2) - 1)
		VX := (1/NP)*math.Log(1/ES) + math.Sqrt(math.Pow(1/ES, 2)+1)
		KX := (1/NP)*math.Log(1/ES) + math.Sqrt(math.Pow(1/ES, 2)-1)
		KX = (math.Exp(KX) + math.Exp(-KX)) / 2
		RP = RP * ((math.Exp(VX) - math.Exp(-VX)) / 2) / KX
		IP = IP * ((math.Exp(VX) + math.Exp(-VX)) / 2) / KX
	}
	// fmt.Printf("RP=%f, IP=%f\n", RP, IP)

	T := 2 * math.Tan(0.5)
	W := 2 * math.Pi * FC
	M := math.Pow(RP, 2) + math.Pow(IP, 2)
	D := 4 - 4*RP*T + M*math.Pow(T, 2)
	X0 := math.Pow(T, 2) / D
	X1 := 2 * math.Pow(T, 2) / D
	X2 := math.Pow(T, 2) / D
	Y1 := (8 - 2*M*math.Pow(T, 2)) / D
	Y2 := (-4 - 4*RP*T - M*math.Pow(T, 2)) / D
	// fmt.Printf("T=%f, W=%f, M=%f, D=%f\n, X0=%f, X1=%f, X2=%f, Y1=%f, Y2=%f\n", T, W, M, D, X0, X1, X2, Y1, Y2)

	K := 0.0
	if LH == 1 {
		K = -math.Cos(W/2+0.5) / math.Cos(W/2-0.5)
	} else {
		K = math.Sin(0.5-W/2) / math.Sin(0.5+W/2)
	}
	D = 1 + Y1*K - Y2*math.Pow(K, 2)

	A0 := (X0 - X1*K + X2*math.Pow(K, 2)) / D
	A1 := (-2*X0*K + X1 + X1*math.Pow(K, 2) - 2*X2*K) / D
	A2 := (X0*math.Pow(K, 2) - X1*K + X2) / D
	B1 := (2*K + Y1 + Y1*math.Pow(K, 2) - 2*Y2*K) / D
	B2 := (-(math.Pow(K, 2)) - Y1*K + Y2) / D

	if LH == 1 {
		A1 = -A1
		B1 = -B1
	}
	// fmt.Printf("A0=%f, A1=%f, A2=%f, B1=%f, B2=%f\n, K=%f, D=%f\n", A0, A1, A2, B1, B2, K, D)
	return A0, A1, A2, B1, B2
}

func (object *WAV) chebyshev() {
	var A [22]float64
	var B [22]float64
	var TA [22]float64
	var TB [22]float64

	A[2] = 1.0
	B[2] = 1.0

	FC := 0.1 // Cut off
	LH := 0.0 // 0: LP, 1: HP
	PR := 0.0 // Percent ripple
	NP := 4.0 // Number of poles

	for P := 2; float64(P) < NP/2; P++ {
		A0, A1, A2, B1, B2 := _cheb(FC, PR, LH, NP, float64(P))
		for I := 0; I < 22; I++ {
			TA[I] = A[I]
			TB[I] = B[I]
		}
		for I := 2; I < 22; I++ {
			A[I] = A0*TA[I] + A1*TA[I-1] + A2*TA[I-2]
			B[I] = TB[I] - B1*TB[I-1] - B2*TB[I-2]
		}
	}

	B[2] = 0
	for I := 0; I < 20; I++ {
		A[I] = A[I+2]
		B[I] = -B[I+2]
	}

	SA := 0.0
	SB := 0.0
	for I := 0; I < 20; I++ {
		if LH == 1 {
			SA = SA + A[I]*math.Pow(-1, float64(I))
			SB = SB + B[I]*math.Pow(-1, float64(I))
			fmt.Printf("%f, SB=%f\n", SA, SB)
		} else {
			SA = SA + A[I]
			SB = SB + B[I]
		}
	}

	GAIN := SA / (1 - SB)
	for I := 0; I < 20; I++ {
		A[I] = A[I] / GAIN
	}
	for i := 0; i < int(object.NumSamples); i++ {
		x := object.data[i]
		x *= GAIN
		object.data[i] = x
	}
}
