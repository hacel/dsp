/*
Copyright © 2021 hacel <hasel@ammasa.net>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package cmd

import (
	"dsp/dsp"
	"fmt"
	"path"

	"github.com/spf13/cobra"
)

var (
	threshold float64
	ratio     float64
	att       float64
	rel       float64
	knee      float64
	makeup    bool
	gain      float64
)

// compressCmd represents the compress command
var compressCmd = &cobra.Command{
	Use:   "compress",
	Short: "Dynamic range compressor",
	Long:  `Dynamic range compressor`,
	Run: func(cmd *cobra.Command, args []string) {
		file1 := args[0]
		fmt.Printf("Compressing %s up to %f dBFS with %f ratio\n", path.Base(file1), threshold, ratio)

		track1 := dsp.NewWav()
		track1.ReadFile(file1)
		fmt.Printf("---------------\n%s details:\n", path.Base(file1))
		track1.DumpHeader(false)

		track1.Compress(threshold, ratio, att, rel, 10, knee, gain, makeup)
		track1.WriteFile(outFile)

		fmt.Printf("Compressed into %s.\n", outFile)
	},
}

func init() {
	rootCmd.AddCommand(compressCmd)
	compressCmd.Flags().Float64VarP(&threshold, "threshold", "t", -12.0, "Compression threshold in dB")
	compressCmd.Flags().Float64VarP(&ratio, "ratio", "r", 2.0, "Compression ratio")
	compressCmd.Flags().Float64VarP(&att, "attack", "a", 10.0, "Compression attack time in ms")
	compressCmd.Flags().Float64VarP(&rel, "release", "R", 300.0, "Compression release time in ms")
	compressCmd.Flags().Float64VarP(&knee, "knee", "k", -25.0, "Compression soft knee width in dB")
	compressCmd.Flags().BoolP("makeup", "m", false, "Apply makeup gain")
	compressCmd.Flags().Float64VarP(&gain, "gain", "g", -1.0, "Makeup gain in dB")
}
