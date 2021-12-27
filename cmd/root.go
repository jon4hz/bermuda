package cmd

import (
	"fmt"

	"github.com/jon4hz/bermuda/internal/bermuda"
	"github.com/jon4hz/bermuda/internal/config"
	"github.com/jon4hz/bermuda/internal/logger"
	"github.com/jon4hz/bermuda/internal/version"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootFlags struct {
	cfgFile string
}

var rootCmd = &cobra.Command{
	Version: version.Version,
	Use:     "bermuda",
	Short:   "Bermuda is a docker garbage collector",
	RunE: func(cmd *cobra.Command, args []string) error {
		return bermuda.Run(config.Get())
	},
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&rootFlags.cfgFile, "config", "", fmt.Sprintf("config file (default is %s)", config.DefaultCfgFile))

	rootCmd.Flags().Bool("image", false, "prune images")
	rootCmd.Flags().StringSlice("image-include", []string{}, "include images with matching labels")
	rootCmd.Flags().StringSlice("image-exclude", []string{}, "exclude images with matching labels")
	rootCmd.Flags().String("image-until", "", "prune images create before the given timestamp")
	rootCmd.Flags().Bool("image-all", false, "prune all images")

	rootCmd.Flags().Bool("container", false, "prune containers")
	rootCmd.Flags().StringSlice("container-include", []string{}, "include containers with matching labels")
	rootCmd.Flags().StringSlice("container-exclude", []string{}, "exclude containers with matching labels")
	rootCmd.Flags().String("container-until", "", "prune containers create before the given timestamp")

	viper.BindPFlag("image.active", rootCmd.Flags().Lookup("image"))
	viper.BindPFlag("image.includeLabels", rootCmd.Flags().Lookup("image-include"))
	viper.BindPFlag("image.excludeLabels", rootCmd.Flags().Lookup("image-exclude"))
	viper.BindPFlag("image.until", rootCmd.Flags().Lookup("image-until"))
	viper.BindPFlag("image.all", rootCmd.Flags().Lookup("image-all"))

	viper.BindPFlag("container.active", rootCmd.Flags().Lookup("container"))
	viper.BindPFlag("container.includeLabels", rootCmd.Flags().Lookup("container-include"))
	viper.BindPFlag("container.excludeLabels", rootCmd.Flags().Lookup("container-exclude"))
	viper.BindPFlag("container.until", rootCmd.Flags().Lookup("container-until"))

	rootCmd.AddCommand(versionCmd)
}

func initConfig() {
	config.LoadConfig(rootFlags.cfgFile)
	logger.New(config.Get().Logging)
}

func Execute() error {
	return rootCmd.Execute()
}
