package dsp

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"os"
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

func rms(val []float64) float64 {
	var sum float64
	for _, s := range val {
		sum += math.Pow(s, 2)
	}
	return math.Sqrt(sum / float64(len(val)))
}

// Wav is a struct to hold wav data.
type Wav struct {
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
	// Derived fields
	NumSamples uint32
	SampleSize uint16
	Duration   float64
}

// NewWav returns a new Wav struct.
func NewWav() *Wav {
	return &Wav{}
}

// Read reads binary data from an io.Reader into Wav.
func (w *Wav) Read(r io.Reader) {
	binary.Read(r, binary.BigEndian, &w.chunkID)
	binary.Read(r, binary.LittleEndian, &w.chunkSize)
	binary.Read(r, binary.BigEndian, &w.format)
	binary.Read(r, binary.BigEndian, &w.subchunk1ID)
	binary.Read(r, binary.LittleEndian, &w.subchunk1Size)
	binary.Read(r, binary.LittleEndian, &w.audioFormat)
	binary.Read(r, binary.LittleEndian, &w.numChannels)
	binary.Read(r, binary.LittleEndian, &w.sampleRate)
	binary.Read(r, binary.LittleEndian, &w.byteRate)
	binary.Read(r, binary.LittleEndian, &w.blockAlign)
	binary.Read(r, binary.LittleEndian, &w.bitsPerSample)
	binary.Read(r, binary.BigEndian, &w.subchunk2ID)
	binary.Read(r, binary.LittleEndian, &w.subchunk2Size)
	w.NumSamples = (8 * w.subchunk2Size) / uint32(w.bitsPerSample)
	w.SampleSize = (w.numChannels * w.bitsPerSample) / 8
	w.Duration = float64(w.subchunk2Size) / float64(w.byteRate)
	for i := 0; i < int(w.NumSamples); i++ {
		x := make([]byte, w.bitsPerSample/8)
		binary.Read(r, binary.LittleEndian, &x)
		smp := float64(int16(binary.LittleEndian.Uint16(x)))
		w.data = append(w.data, smp)
	}
}

// ReadFile opens the given file string and passes it to Read.
func (w *Wav) ReadFile(path string) {
	f, err := os.Open(path)
	check(err)
	defer f.Close()
	w.Read(f)
}

// Write writes Wav data into an io.Writer as binary.
func (w *Wav) Write(r io.Writer) {
	binary.Write(r, binary.BigEndian, w.chunkID)
	binary.Write(r, binary.LittleEndian, w.chunkSize)
	binary.Write(r, binary.BigEndian, w.format)
	binary.Write(r, binary.BigEndian, w.subchunk1ID)
	binary.Write(r, binary.LittleEndian, w.subchunk1Size)
	binary.Write(r, binary.LittleEndian, w.audioFormat)
	binary.Write(r, binary.LittleEndian, w.numChannels)
	binary.Write(r, binary.LittleEndian, w.sampleRate)
	binary.Write(r, binary.LittleEndian, w.byteRate)
	binary.Write(r, binary.LittleEndian, w.blockAlign)
	binary.Write(r, binary.LittleEndian, w.bitsPerSample)
	binary.Write(r, binary.BigEndian, w.subchunk2ID)
	binary.Write(r, binary.LittleEndian, w.subchunk2Size)
	for i := 0; i < int(w.NumSamples); i++ {
		signal := make([]byte, w.bitsPerSample/8)
		binary.LittleEndian.PutUint16(signal, uint16(w.data[i]))
		binary.Write(r, binary.LittleEndian, signal)
	}
}

// WriteFile opens the given file string and passes it to Write
func (w *Wav) WriteFile(path string) {
	f, err := os.Create(path)
	check(err)
	defer f.Close()
	w.Write(f)
}

// DumpHeader prints Wav header information.
func (w *Wav) DumpHeader(more bool) {
	fmt.Printf("%-14s %.2fKB\n", "File size:", float64(w.chunkSize)/1000)
	fmt.Printf("%-14s %.2fs\n", "Duration:", w.Duration)
	fmt.Printf("%-14s %d\n", "Sample rate:", w.sampleRate)
	if more {
		fmt.Printf("Size of each sample: %d bytes\n", w.SampleSize)
		fmt.Printf("Number of samples: %d\n", w.NumSamples)
		fmt.Printf("%-14s %s\n", "chunkID:", w.chunkID)
		fmt.Printf("%-14s %d\n", "chunkSize:", w.chunkSize)
		fmt.Printf("%-14s %s\n", "format:", w.format)
		fmt.Printf("%-14s %s\n", "subchunk1ID:", w.subchunk1ID)
		fmt.Printf("%-14s %d\n", "subchunk1Size:", w.subchunk1Size)
		fmt.Printf("%-14s %d\n", "audioFormat:", w.audioFormat)
		fmt.Printf("%-14s %d\n", "numChannels:", w.numChannels)
		fmt.Printf("%-14s %d\n", "byteRate:", w.byteRate)
		fmt.Printf("%-14s %d\n", "blockAlign:", w.blockAlign)
		fmt.Printf("%-14s %d\n", "bitsPerSample:", w.bitsPerSample)
		fmt.Printf("%-14s %s\n", "subchunk2ID:", w.subchunk2ID)
		fmt.Printf("%-14s %d\n", "subchunk2Size:", w.subchunk2Size)
	}
}

// ReconSignal reconstructs signal data into Wav from a given inverse DFT.
func (w *Wav) ReconSignal(idft []complex128) {
	if len(w.data) < len(idft) {
		i := 0
		for ; i < len(w.data); i++ {
			w.data[i] = real(idft[i])
		}
		for ; i < len(idft); i++ {
			w.data = append(w.data, real(idft[i]))
		}
	} else {
		for i := 0; i < len(w.data); i++ {
			w.data[i] = real(idft[i])
		}
	}
}

// Mix mixes two tracks into one.
func (w *Wav) Mix(t1 *Wav, t2 *Wav) {
	var longerTrack, shorterTrack *Wav
	if t1.NumSamples >= t2.NumSamples {
		longerTrack = t1
		shorterTrack = t2
	} else {
		longerTrack = t2
		shorterTrack = t1
	}
	*w = *longerTrack
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
		w.data[i] = x
	}
}

// Normalize normalizes a track according to the desired peak in dBFS.
func (w *Wav) Normalize(desiredPeak float64) {
	base := math.Pow(2, float64(w.bitsPerSample-1)) * math.Pow(10, (desiredPeak/20))
	var peak float64 = 0
	for i := 0; i < int(w.NumSamples); i++ {
		x := math.Abs(w.data[i])
		if x > peak {
			peak = x
		}
	}
	normNum := base / peak
	for i := 0; i < int(w.NumSamples); i++ {
		x := w.data[i]
		x *= normNum
		w.data[i] = x
	}
}

// Compress is a dynamic range compressor.
func (w *Wav) Compress(threshold, ratio, tatt, trel, tla, knee, gain float64, makeup bool) {
	threshold = math.Pow(2, float64(w.bitsPerSample-1)) * math.Pow(10, threshold/20)
	sr := float64(w.sampleRate)
	tatt *= math.Pow(10, -3) // attack time
	trel *= math.Pow(10, -3) // release time
	tla *= math.Pow(10, -3)  // lookahead
	knee = math.Pow(2, float64(w.bitsPerSample-1)) * math.Pow(10, (knee/20))
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
	nla := sr * tla

	for i := 0; i < int(w.NumSamples); i++ {
		summ := 0.0
		for j := 0; j < int(nla); j++ {
			var smp float64
			if i+j >= len(w.data) {
				smp = 0.0
			} else {
				smp = w.data[i+j]
			}
			summ += smp
		}

		peak := summ / nla
		var theta float64
		if peak > env {
			theta = att
		} else {
			theta = rel
		}
		env = ((1.0-theta)*peak + theta*env)

		var gain float64
		if env-threshold < -knee/2 {
			gain = 1.0
		} else if math.Abs(env-threshold) <= knee/2 {
			gain = (env + ((1/ratio-1)*math.Pow(env-threshold+knee/2, 2))/(knee*2)) / env
		} else if env-threshold > knee/2 {
			gain = (threshold + (env-threshold)/ratio) / env
		}

		x := w.data[i]
		x *= gain

		w.data[i] = x
	}
	if makeup {
		fmt.Printf("Normalizing...\n")
		w.Normalize(gain)
	}
}

// RollingAvgLowpass is a low pass filter using rolling average.
func (w *Wav) RollingAvgLowpass(bandwidth int) {
	var period []float64
	for i := 0; i < int(w.NumSamples)-5; i++ {
		x := w.data[i]
		if len(period) == bandwidth {
			period = period[1:]
		}
		period = append(period, x)
		avg := avg(period)
		w.data[i] = avg
	}
}

// Biquad is an implementation of the Biquad filter
func (w *Wav) Biquad(fc, lh int) {
	r := math.Sqrt(2) // Rez
	sr := float64(w.sampleRate)
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
	for i := 0; i < int(w.NumSamples); i++ {
		y := 0.0
		x0 := w.data[i]
		if len(period) == 2 {
			x1 := period[1]
			x2 := period[0]
			y1 := w.data[i-1]
			y2 := w.data[i-2]
			y = a1*x0 + a2*x1 + a3*x2 - b1*y1 - b2*y2
			period = period[1:]
		} else {
			y = x0
		}
		period = append(period, x0)
		w.data[i] = y
	}
}

// WindowedSinc is a Hamming windowed-sinc  low pass filter
func (w *Wav) WindowedSinc(cutoff, bandwidth int) {
	if cutoff > int(w.sampleRate)/2 {
		panic("Cutoff frequency too high.")
	}
	FC := float64(cutoff) / float64(w.sampleRate) // Cut off (freq/sample rate)
	M := bandwidth                                // Filter roll off
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
	for j := M; j < int(w.NumSamples); j++ {
		y := 0.0
		x := 0.0
		for i := range kernel {
			x = w.data[j-i]
			y += x * kernel[i]
		}
		filteredData = append(filteredData, y)
	}
	w.data = filteredData
}

// Highpass is a basic highpass filter
func (w *Wav) Highpass() {
	var period []float64
	for i := 0; i < int(w.NumSamples); i++ {
		var y float64
		x := w.data[i]
		if len(period) == 2 {
			y = period[0] + -2*period[1] + x
			period = period[1:]
		}
		period = append(period, x)
		w.data[i] = y
	}
}

// Chebyshev sub routine
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

// Chebyshev is an implementation of the Chebyshev filter
// WIP
func (w *Wav) Chebyshev() {
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
	for i := 0; i < int(w.NumSamples); i++ {
		x := w.data[i]
		x *= GAIN
		w.data[i] = x
	}
}
