package cmd

import "github.com/spf13/cobra"

func NewRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:           "k9s-multi-cluster-plugin",
		Short:         "Generate K9s plugin YAML from templates and overrides",
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	rootCmd.AddCommand(newGenerateCommand())
	return rootCmd
}
