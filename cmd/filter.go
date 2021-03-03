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
	"errors"
	"fmt"
	"path"

	"github.com/spf13/cobra"
)

var (
	lh        int
	freq      int
	bandwidth int
)

func isValidFilter(filter string) bool {
	filters := []string{"avg", "biquad", "windowedsinc", "highpass"}
	for _, v := range filters {
		if filter == v {
			return true
		}
	}
	return false
}

// filterCmd represents the filter command
var filterCmd = &cobra.Command{
	Use:   "filter",
	Short: "Apply various filters on a track.",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires a filter to be selected (avg, biquad, windowedsinc, highpass)")
		}
		if len(args) < 2 {
			return errors.New("requires an input file to be specfied")
		}
		if isValidFilter(args[0]) {
			return nil
		}
		return fmt.Errorf("invalid filter specified: %s", args[0])
	},
	Run: func(cmd *cobra.Command, args []string) {
		filter := args[0]
		file1 := args[1]

		track1 := dsp.NewWav()
		track1.ReadFile(file1)
		fmt.Printf("---------------\n%s details:\n", path.Base(file1))
		track1.DumpHeader(false)

		switch filter {
		case "avg":
			fmt.Printf("Convolving using rolling average (M=%d)...\n", bandwidth)
			track1.RollingAvgLowpass(bandwidth)

		case "windowedsinc":
			fmt.Printf("Convolving using Windowed-Sinc (fc=%d, M=%d)...\n", freq, bandwidth)
			track1.WindowedSinc(freq, bandwidth)

		case "biquad":
			fmt.Printf("Convolving using Biquad (fc=%d, lh=%d)...\n", freq, lh)
			track1.Biquad(freq, lh)

		case "highpass":
			fmt.Printf("Convolving using highpass...")
			track1.Highpass()

		case "cheb":
			// track1.chebyshev()

		default:
			fmt.Println("Please enter a valid filter for convolution")
			return
		}
		track1.WriteFile(outFile)
	},
}

func init() {
	rootCmd.AddCommand(filterCmd)
	filterCmd.Flags().IntVarP(&lh, "lh", "l", 0, "Low pass: 0, High pass: 1")
	filterCmd.Flags().IntVarP(&freq, "freq", "f", 5000, "Cut off frequency")
	filterCmd.Flags().IntVarP(&bandwidth, "bandwidth", "b", 20, "Roll off value for certain filters")
}
