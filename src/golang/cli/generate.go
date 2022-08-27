package cli

import "github.com/spf13/cobra"

import (
	"github.com/radiospiel/envy/envy"
	"io/ioutil"
	"os"
)

const generate_template = `#
# This is an envy(3) file.
#
# Use "envy edit path-to-file" to edit this file.
#

#
# A non-secured part. Note that part names are only here for documentation purposes.
#
[http]
HTTP_PORT=80

#
# A secure block: every entry in a block named [secure] or [something.secure]
# will be encrypted.
[secure]
MY_PASSWORD=This is my password

#
# Another non-secured block
[database]
DATABASE_POOL_SIZE=10

#
# Another secured block
[database.secure]
DATABASE_URL=postgres://pg_user:pg_password/server:5432/database/schema
`

func do_generate(path string) {
	fileMustNotExist(path)

	/*
	 * Verify that we have a secret
	 */
	envy.SecretMustExist()

	/*
	 * create a temp file, and decrypt src into the temp file.
	 */
	tmpFile, err := ioutil.TempFile("", "envy")
	defer os.Remove(tmpFile.Name())

	check(err)

	_, err = tmpFile.WriteString(generate_template)
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
		Use:   "generate",
		Short: "Generate an envy file",
		Long:  `Generate and edit an envy file`,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			do_generate(args[0])
		},
	}

	rootCmd.AddCommand(cmd)
}
