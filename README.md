# go-dsp

## TODO
- Error checking
- Use log package
- Compatibilty with stereo tracks

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



#### Notes
normalization:
1. find baseline signal amplitude from bitspersample and desired peak. Ex: 32767*10^(-1/20)
2. find peak in track
3. divide baseline amplitude by tracks peak amplitude
4. multiply every single in the track by the result

mixing:
1. add two signals and watch for overflow
2. if result exceeds max then set to max