// Package cmd implements the socketcli CLI command.
package cmd

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/antoniomika/socketcli/entry"
	"github.com/antoniomika/socketcli/utils"
	"github.com/fsnotify/fsnotify"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	// Version describes the version of the current build.
	Version = "dev"

	// Commit describes the commit of the current build.
	Commit = "none"

	// Date describes the date of the current build.
	Date = "unknown"

	// configFile holds the location of the config file from CLI flags.
	configFile string

	// rootCmd is the root cobra command.
	rootCmd = &cobra.Command{
		Use:     "socketcli",
		Short:   "The socketcli command",
		Long:    "The socketcli command",
		Run:     runCommand,
		Version: Version,
	}
)

// init initializes flags used by the root command.
func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.SetVersionTemplate(fmt.Sprintf("Version: %v\nCommit: %v\nDate: %v\n", Version, Commit, Date))

	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "config.yml", "Config file")

	rootCmd.PersistentFlags().StringP("time-format", "", "2006/01/02 - 15:04:05", "The time format to use for general log messages")
	rootCmd.PersistentFlags().StringP("log-to-file-path", "", "/tmp/socketcli.log", "The file to write log output to")
	rootCmd.PersistentFlags().StringP("websocket-address", "w", "wss://chaotic-stream.herokuapp.com", "The websocket address to connect to")

	rootCmd.PersistentFlags().BoolP("debug", "", false, "Enable debugging information")
	rootCmd.PersistentFlags().BoolP("log-to-stdout", "", true, "Enable writing log output to stdout")
	rootCmd.PersistentFlags().BoolP("log-to-file", "", false, "Enable writing log output to file, specified by log-to-file-path")
	rootCmd.PersistentFlags().BoolP("log-to-file-compress", "", false, "Enable compressing log output files")
	rootCmd.PersistentFlags().BoolP("stop-after", "", false, "Enable stopping the program after stop-after-time")

	rootCmd.PersistentFlags().IntP("log-to-file-max-size", "", 500, "The maximum size of outputed log files in megabytes")
	rootCmd.PersistentFlags().IntP("log-to-file-max-backups", "", 3, "The maxium number of rotated logs files to keep")
	rootCmd.PersistentFlags().IntP("log-to-file-max-age", "", 28, "The maxium number of days to store log output in a file")

	rootCmd.PersistentFlags().DurationP("stop-after-time", "", 1*time.Minute, "The duration to stop the program after if stop-after is set")
}

// initConfig initializes the configuration and loads needed
// values. It initializes logging and other vars.
func initConfig() {
	viper.SetConfigFile(configFile)

	err := viper.BindPFlags(rootCmd.PersistentFlags())
	if err != nil {
		logrus.Println("Unable to bind pflags:", err)
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		logrus.Println("Using config file:", viper.ConfigFileUsed())
	}

	viper.WatchConfig()

	writers := []io.Writer{}

	if viper.GetBool("log-to-stdout") {
		writers = append(writers, os.Stdout)
	}

	if viper.GetBool("log-to-file") {
		writers = append(writers, &lumberjack.Logger{
			Filename:   viper.GetString("log-to-file-path"),
			MaxSize:    viper.GetInt("log-to-file-max-size"),
			MaxBackups: viper.GetInt("log-to-file-max-backups"),
			MaxAge:     viper.GetInt("log-to-file-max-age"),
			Compress:   viper.GetBool("log-to-file-compress"),
		})
	}

	multiWriter := io.MultiWriter(writers...)

	viper.OnConfigChange(func(e fsnotify.Event) {
		logrus.Println("Reloaded configuration file.")

		log.SetFlags(0)
		log.SetOutput(utils.LogWriter{
			TimeFmt:     viper.GetString("time-format"),
			MultiWriter: multiWriter,
		})

		if viper.GetBool("debug") {
			logrus.SetLevel(logrus.DebugLevel)
		}
	})

	log.SetFlags(0)
	log.SetOutput(utils.LogWriter{
		TimeFmt:     viper.GetString("time-format"),
		MultiWriter: multiWriter,
	})

	if viper.GetBool("debug") {
		logrus.SetLevel(logrus.DebugLevel)
	}

	logrus.SetOutput(multiWriter)
}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

// runCommand is used to start the root muxer.
func runCommand(cmd *cobra.Command, args []string) {
	entry.Start()
}
