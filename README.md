# go-dsp

## TODO
- Error checking
- Use log package
Mixing works but clipping issues



normalization:
1. find baseline signal amplitude from bitspersample and desired peak. Ex: 32767*10^(-1/20)
2. find peak in track
3. divide baseline amplitude by tracks peak amplitude
4. multiply every single in the track by the result