package cli

import "github.com/spf13/cobra"

import (
	"github.com/radiospiel/envy/envy"
	"log"
	"os"
	"os/exec"
	"strings"
)

func scp_binary() string {
	scp_binary, ok := os.LookupEnv("SCP")
	if ok {
		log.Printf("using scp command %s", scp_binary)
		return scp_binary
	}

	return "scp"
}

/*
 * Copy the current secret via ${SSH_BINARY:-ssh} to a target account.
 */
func do_secret_install(target_account string) {
	envy.SecretMustExist()

	target_str := target_account
	if !strings.Contains(target_str, ":") {
		target_str = target_str + ":" + envy.ENVY_BASE_NAME
	}

	/*
	 * It would be nice if we could prevent scp from overwriting the target.
	 * However, this is not possible via scp(1)
	 */
	scp_binary := scp_binary()
	log.Printf("running: %s %s %s", scp_binary, envy.SecretFile(), target_str)

	cmd := exec.Command(scp_binary, envy.SecretFile(), target_str)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	check(err)
}

func init() {
	var cmd = &cobra.Command{
		Use:   "secret:install",
		Short: "Install the current envy secret to a remote location",
		Long:  `Install the current envy secret to a remote location`,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			do_secret_install(args[0])
		},
	}

	rootCmd.AddCommand(cmd)
}
