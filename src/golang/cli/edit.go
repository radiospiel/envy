package cli

import "github.com/spf13/cobra"

import (
	"github.com/radiospiel/envy/envy"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
)

func copyEnvyFile(srcPath string, dest *os.File, do_encrypt bool) error {
	return envy.ParseFile(srcPath, func(mode envy.Mode, pt1 string, pt2 string) {
		switch mode {
		case envy.Mode_Value:
			fmt.Fprintf(dest, "%s=%s\n", pt1, pt2)
		case envy.Mode_Secured_Value:
			var value string
			if do_encrypt {
				value = envy.EncryptSecuredValue(pt2)
			} else {
				value = envy.DecryptSecuredValue(pt2)
			}

			fmt.Fprintf(dest, "%s=%s\n", pt1, value)
		case envy.Mode_Line:
			fmt.Fprintf(dest, "%s\n", pt1)
		}
	})
}

func copyAndDecryptEnvyFile(srcPath string, dest *os.File) error {
	return copyEnvyFile(srcPath, dest, false)
}

func copyAndEncryptEnvyFile(srcPath string, dest *os.File) error {
	return copyEnvyFile(srcPath, dest, true)
}

func shell(str string) error {
	cmd := exec.Command("/bin/sh", "-c", str) // editor, tmpFile.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func do_edit(path string) {
	/*
	 * create a temp file, and decrypt src into the temp file.
	 */
	tmpFile, err := ioutil.TempFile("", "envy")
	defer os.Remove(tmpFile.Name())

	check(err)

	err = copyAndDecryptEnvyFile(path, tmpFile)
	check(err)

	tmpFile.Close()

	if !editFile(tmpFile.Name()) {
		return
	}

	/*
	 * encrypt temp file into the original path.
	 */

	dest, err := os.Create(path)
	check(err)
	defer dest.Close()

	copyAndEncryptEnvyFile(tmpFile.Name(), dest)
	check(err)
}

func init() {
	var cmd = &cobra.Command{
		Use:   "edit",
		Short: "Edit an envy file",
		Long:  `Edit an envy file in your $EDITOR`,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			do_edit(args[0])
		},
	}

	rootCmd.AddCommand(cmd)
}
