/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:              "docs-training",
	Short:            "A brief description of your application",
	PersistentPreRun: checkEnv,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.docs-training.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.PersistentFlags().BoolP("toggle", "t", false, "Help message for toggle")

	// number of questions
	rootCmd.PersistentFlags().IntP("num", "n", 30, "number of questions")

	// dirs and files
	rootCmd.PersistentFlags().StringSliceP("wf", "w", nil, "dirs and files")

	// exclude files
	rootCmd.PersistentFlags().StringSliceP("exclude", "e", nil, "exclude specified files")
}

func checkEnv(cmd *cobra.Command, args []string) {
	var EnvVar = "BaseURL"
	if value := os.Getenv(EnvVar); value == "" {
		fmt.Printf("环境变量 %s 不存在\n", EnvVar)
		err := os.Setenv(EnvVar, "https://blog.wrss.top/")
		if err != nil {
			os.Exit(1)
		}
	}
}
