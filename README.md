# go-dsp

## TODO
- Error checking
- Use log package
Mixing works but clipping issues
- Test for stereo

- refactor "object oriented"
 - use path.base
 - new instance using :=
 defer
 - maybe writewavfile/readwavfile(path)
- defer fileclose


compression:
out = (in - thresh) / ratio + thresh

- too much distortion

threshold
ratio
normalize to 0db
rms & peak
soft knee

attack time
release time
noise floor



normalization:
1. find baseline signal amplitude from bitspersample and desired peak. Ex: 32767*10^(-1/20)
2. find peak in track
3. divide baseline amplitude by tracks peak amplitude
4. multiply every single in the track by the result

mixing:
1. add two signals and watch for overflow
2. if result exceeds max then set to max