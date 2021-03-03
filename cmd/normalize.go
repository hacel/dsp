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
	peak float64
)

// normalizeCmd represents the normalize command
var normalizeCmd = &cobra.Command{
	Use:   "normalize",
	Short: "Normalizes the given tracks amplitude",
	Long:  `Normalizes the given tracks amplitude`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		file1 := args[0]
		fmt.Printf("Noramlizing %s to %f dBFS\n", path.Base(file1), peak)

		track1 := dsp.NewWav()
		track1.ReadFile(file1)
		fmt.Printf("---------------\n%s details:\n", path.Base(file1))
		track1.DumpHeader(false)

		track1.Normalize(peak)
		track1.WriteFile(outFile)

		fmt.Printf("Normalized into %s.\n", outFile)
	},
}

func init() {
	rootCmd.AddCommand(normalizeCmd)
	normalizeCmd.Flags().Float64VarP(&peak, "peak", "p", -1.0, "Desired peak in dB")
}
