package main

import (
	"github.com/seatbelt/pkg/hook"
	"github.com/seatbelt/pkg/scanner"
	"github.com/spf13/cobra"
)

var action string

var scanSecretsCmd = &cobra.Command{
	Use:   "scan-secrets",
	Short: "Scan prompt for secrets before submission",
	Long:  `Scans the prompt for secrets. If any secrets are found, the behavior depends on the action specified.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		input, err := hook.ReadInput()
		if err != nil {
			return err
		}

		registry := scanner.NewRegistry(
			scanner.NewAPIKeyDetector(),
			scanner.NewDBConnectionDetector(),
			scanner.NewPrivateKeyDetector(),
			scanner.NewPasswordDetector(),
			scanner.NewEnvVarDetector(),
			scanner.NewDotEnvDetector(),
		)

		output := hook.ProcessInput(input, registry)
		return hook.WriteOutput(output)
	},
}

func init() {
	rootCmd.AddCommand(scanSecretsCmd)
	scanSecretsCmd.Flags().StringVar(&action, "action", "block", "action to take when secrets are found (block)")
	scanSecretsCmd.MarkFlagRequired("action")
}