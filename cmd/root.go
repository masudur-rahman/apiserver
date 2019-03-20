package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "apiserver",
	Short: "It's a server containing workers of appscode",
	Long: "All the worker profile of the AppsCode Ltd." +
		" is included in this server",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome.........!!!")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
