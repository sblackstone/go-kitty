/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/sblackstone/go-kitty/kitty"
	"github.com/spf13/cobra"
)

var (
	snakeCount          int
	snakeMaxLen         int
	snakeInitialDelayMax int
	stringCount         int
	stringMinLen        int
	stringMaxLen        int
	stringInitialDelayMax int
	butterflyCount      int
	butterflyInitialDelayMax int
	laserCount          int
	laserInitialDelayMax int
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "go-kitty",
	Short: "Cat Entertainment",
	Long:  `A way to entertain a cat looking at a terminal window`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := kitty.DefaultKittyConfig()
		cfg.SnakeCount = snakeCount
		cfg.SnakeConfig.MaxLen = snakeMaxLen
		cfg.SnakeConfig.InitialDelayMax = snakeInitialDelayMax
		cfg.SwayStringCount = stringCount
		cfg.SwayStringConfig.MinLen = stringMinLen
		cfg.SwayStringConfig.MaxLen = stringMaxLen
		cfg.SwayStringConfig.InitialDelayMax = stringInitialDelayMax
		cfg.ButterflyCount = butterflyCount
		cfg.ButterflyConfig.InitialDelayMax = butterflyInitialDelayMax
		cfg.LaserCount = laserCount
		cfg.LaserConfig.InitialDelayMax = laserInitialDelayMax

		k, err := kitty.New(cfg)
		if err != nil {
			os.Exit(1)
		}
		k.Start(cmd.Context())
	},
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

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.go-kitty.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	defaults := kitty.DefaultKittyConfig()
	rootCmd.Flags().IntVar(&snakeCount, "snakes", defaults.SnakeCount, "Number of snakes")
	rootCmd.Flags().IntVar(&snakeMaxLen, "snake-max-len", defaults.SnakeConfig.MaxLen, "Snake max length")
	rootCmd.Flags().IntVar(&snakeInitialDelayMax, "snake-initial-delay-max", defaults.SnakeConfig.InitialDelayMax, "Max initial delay (ticks) for snakes")
	rootCmd.Flags().IntVar(&stringCount, "strings", defaults.SwayStringCount, "Number of sway strings")
	rootCmd.Flags().IntVar(&stringMinLen, "string-min-len", defaults.SwayStringConfig.MinLen, "Sway string min length")
	rootCmd.Flags().IntVar(&stringMaxLen, "string-max-len", defaults.SwayStringConfig.MaxLen, "Sway string max length")
	rootCmd.Flags().IntVar(&stringInitialDelayMax, "string-initial-delay-max", defaults.SwayStringConfig.InitialDelayMax, "Max initial delay (ticks) for sway strings")
	rootCmd.Flags().IntVar(&butterflyCount, "butterflies", defaults.ButterflyCount, "Number of butterflies")
	rootCmd.Flags().IntVar(&butterflyInitialDelayMax, "butterfly-initial-delay-max", defaults.ButterflyConfig.InitialDelayMax, "Max initial delay (ticks) for butterflies")
	rootCmd.Flags().IntVar(&laserCount, "lasers", defaults.LaserCount, "Number of laser pointers")
	rootCmd.Flags().IntVar(&laserInitialDelayMax, "laser-initial-delay-max", defaults.LaserConfig.InitialDelayMax, "Max initial delay (ticks) for lasers")
}
