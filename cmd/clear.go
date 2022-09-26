/*
Copyright © 2022 KAI CHU CHUNG <cage.chung@gmail.com>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cage1016/alfred-timelog/alfred"
)

// clearCmd represents the clear command
var clearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear cache",
	Run:   runClearCmd,
}

func runClearCmd(cmd *cobra.Command, args []string) {
	err := alfred.StoreOngoingTimelog(wf, alfred.Timelog{})
	if err != nil {
		ErrorHandle(fmt.Errorf("failed clear cache: %v", err))
		return
	}

	av.Var("msg", "clear cache success")
	if err := av.Send(); err != nil {
		wf.Fatalf("failed to send args to Alfred: %v", err)
	}
}

func init() {
	rootCmd.AddCommand(clearCmd)
}
