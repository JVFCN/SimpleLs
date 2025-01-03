package cmd

import (
	"github.com/spf13/cobra"
	"ls/TraverseDIR"
	"os"
	"regexp"
)

var rootCmd = &cobra.Command{
	Use:   "ls",
	Short: "Linux ls command for Windows",
	Long:  `Use Golang to implement Linux's LS command on Windows`,
	Run: func(cmd *cobra.Command, args []string) {
		isList, err := cmd.Flags().GetBool("list")
		if err != nil {
			os.Exit(1)
		}
		showTime, err := cmd.Flags().GetBool("time")
		mod, err := cmd.Flags().GetBool("mod")
		re, err := cmd.Flags().GetString("regular")
		all, err := cmd.Flags().GetBool("all")
		regex, err := regexp.Compile(re)
		if !isList {
			TraverseDIR.Default(showTime, regex, all)
		} else {
			TraverseDIR.List(showTime, mod, regex, all)
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("list", "l", false, "print as list")
	rootCmd.Flags().StringP("regular", "r", ".", "use regular expressions")
	rootCmd.PersistentFlags().BoolP("time", "t", false, "show run time")
	rootCmd.Flags().BoolP("mod", "m", false, "show modified time")
	rootCmd.Flags().BoolP("all", "a", false, "show all file")
}
