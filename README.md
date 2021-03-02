# dsp
Basic 16-bit PCM digital signal processor in Go

## Status
| Func | Status  | Description | Notes |
| --- |--------|--------| -----|
| getDFT() | Working | Returns the Discrete Fourier Transform of a track | |
| getIDFT() | Working | Returns the Inverse Discrete Fourier Transform of a track | |
| reconSignal() | Working | Reconstructs signal data from IDFT output | |
| mix()      | Working | Adds two tracks together | |
| normalize() | Working | Normalizes track amplitude | |
| compress() | Working | Dynamic range compressor |Controls are not as impactful as they should be. Add noise floor. Peak or RMS?|
| rollingAvgLowpass() | Working | LP filter using rolling average |  No real controls. |
| biquad() | Working | LP/HP filter using Biquad |Best lowpass/highpass filter so far |
| windowedSinc() | Working | LP filter using Hamming Windowed-Sinc  |Can't filter above SR/2?|
| highpass() | Working | Very basic high pass filter with no controls | |
| chebyshev() | Not working | Chebyshev filter | WIP |

#### TODO
- Error checking, log package
- Stereo compatibility