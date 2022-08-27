package cli

import (
	"log"
	"os"

	"github.com/radiospiel/envy/envy"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "envy",
		Short: "envy handles secure environments",
		Long:  `envy handles secure environments`,
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// Exists reports whether the named file or directory exists.
func fileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func fileMustExist(path string) {
	if !fileExists(path) {
		log.Fatalf("%s: file does not exist", path)
	}
}

func fileMustNotExist(path string) {
	if fileExists(path) {
		log.Fatalf("%s: file exists already", path)
	}
}

func editFile(path string) bool {
	/*
	 * make sure we have an EDITOR env variable
	 */
	_, ok := os.LookupEnv("EDITOR")
	if !ok {
		os.Setenv("EDITOR", "vi")
	}

	/*
	 * run editor on tmpFile. The editor is determined by $EDITOR, and
	 */
	err := shell("$EDITOR " + path)
	return err == nil
}

func mustLoadConfig(path string) map[string]string {
	config, err := envy.LoadConfig(path)
	check(err)

	return config
}
