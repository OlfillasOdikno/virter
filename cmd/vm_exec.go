package cmd

import (
	"context"
	"io/ioutil"
	"log"

	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func vmExecCommand() *cobra.Command {
	execCmd := &cobra.Command{
		Use:   "exec vm_name docker_image",
		Short: "Run a Docker container against a VM",
		Long:  `Run a Docker container on the host with a connection to a VM.`,
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			vmName := args[0]
			dockerImageName := args[1]

			dockerTimeout := viper.GetDuration("time.docker_timeout")
			ctx, cancel := context.WithTimeout(context.Background(), dockerTimeout)
			defer cancel()

			v, err := VirterConnect()
			if err != nil {
				log.Fatal(err)
			}

			docker, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
			if err != nil {
				log.Fatalf("could not connect to Docker %v", err)
			}

			privateKeyPath := viper.GetString("auth.virter_private_key_path")
			privateKey, err := ioutil.ReadFile(privateKeyPath)
			if err != nil {
				log.Fatalf("failed to load private key from '%s': %v", privateKeyPath, err)
			}

			err = v.VMExec(ctx, docker, vmName, dockerImageName, privateKey)
			if err != nil {
				log.Fatal(err)
			}
		},
	}

	return execCmd
}
