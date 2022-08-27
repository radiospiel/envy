package cli

import "github.com/spf13/cobra"

import (
	"github.com/radiospiel/envy/envy"
)

/*
 * Copy the current secret via ${SSH_BINARY:-ssh} to a target account.
 */
func init() {
	var cmd = &cobra.Command{
		Use:   "secret:generate",
		Short: "Generate the envy secret",
		Long:  `Generate the envy secret`,
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			envy.GenerateSecret()
		},
	}

	rootCmd.AddCommand(cmd)
}
