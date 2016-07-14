package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/pkg/errors"
	"github.com/sbani/gcr/config"
	"github.com/sbani/gcr/pkg"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var c = new(config.Config)
var mutex = &sync.RWMutex{}

// RootCmd is the cobra command setup
var RootCmd = &cobra.Command{
	Use:   "gcr",
	Short: "Content repo",
}

// Execute the root cmd and print error if exists
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gcr.yaml)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	mutex.Lock()
	if cfgFile != "" {
		// enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	path := absPathify("$HOME")
	if _, err := os.Stat(filepath.Join(path, ".gcr.yml")); err != nil {
		_, _ = os.Create(filepath.Join(path, ".gcr.yml"))
	}

	viper.SetConfigType("yaml")
	viper.SetConfigName(".")     // name of config file (without extension)
	viper.AddConfigPath("$HOME") // adding home directory as first search path
	viper.AutomaticEnv()         // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf(`Config file not found because "%s"`, err)
		fmt.Println("")
	}

	if err := viper.Unmarshal(c); err != nil {
		pkg.Fatal(errors.Wrap(err, "Could not read config"))
	}

	mutex.Unlock()
}

func absPathify(inPath string) string {
	if strings.HasPrefix(inPath, "$HOME") {
		inPath = userHomeDir() + inPath[5:]
	}

	if strings.HasPrefix(inPath, "$") {
		end := strings.Index(inPath, string(os.PathSeparator))
		inPath = os.Getenv(inPath[1:end]) + inPath[end:]
	}

	if filepath.IsAbs(inPath) {
		return filepath.Clean(inPath)
	}

	p, err := filepath.Abs(inPath)
	if err == nil {
		return filepath.Clean(p)
	}
	return ""
}

func userHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
}
