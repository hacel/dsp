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

// mixCmd represents the mix command
var mixCmd = &cobra.Command{
	Use:   "mix",
	Short: "Mixes two tracks into one.",
	Long:  `Mixes two tracks into one.`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		file1 := args[0]
		file2 := args[1]
		fmt.Printf("Mixing %s and %s\n", path.Base(file1), path.Base(file2))

		track1 := dsp.NewWav()
		track1.ReadFile(file1)
		fmt.Printf("---------------\n%s details:\n", path.Base(file1))
		track1.DumpHeader(false)

		track2 := dsp.NewWav()
		track2.ReadFile(file2)
		fmt.Printf("---------------\n%s details:\n", path.Base(file2))
		track2.DumpHeader(false)

		newTrack := dsp.NewWav()
		newTrack.Mix(track1, track2)
		newTrack.WriteFile(outFile)

		fmt.Printf("Mixed into %s.\n", outFile)
	},
}

func init() {
	rootCmd.AddCommand(mixCmd)
}
