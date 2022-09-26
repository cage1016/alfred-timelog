/*
Copyright © 2022 KAI CHU CHUNG <cage.chung@gmail.com>

*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/spf13/cobra"

	"github.com/cage1016/alfred-timelog/alfred"
)

var (
	snre = regexp.MustCompile(`(?m)^\d{4}w\d{2}.xlsx$`)
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list timelog sheets",
	Run:   runListCmd,
}

func runListCmd(cmd *cobra.Command, args []string) {
	CheckForUpdate()

	dir := alfred.GetBaseFolder(wf)
	if dir == "" {
		dir = wf.DataDir()
	}

	res := []string{}
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return err
		}
		if snre.MatchString(info.Name()) {
			res = append(res, info.Name())
		}
		return nil
	})
	if err != nil {
		ErrorHandle(fmt.Errorf("failed list timelog sheets: %v", err))
	}

	for _, v := range res {
		ni := wf.NewItem(v).
			Subtitle("⇧ ⌥, ↩ Open file with default application").
			Quicklook(filepath.Join(dir, v)).
			Arg(filepath.Join(dir, v)).
			Valid(true)

		ni.Opt().
			Subtitle("↩ Reveal file in Finder").
			Arg(filepath.Join(dir, v)).
			Valid(true)
	}

	if args[0] != "" {
		wf.Filter(args[0])
	}

	wf.SendFeedback()
}

func init() {
	rootCmd.AddCommand(listCmd)
}
