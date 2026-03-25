package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"k9s-multi-cluster-plugin/internal/generator"
)

type generateOptions struct {
	kubeconfigPath string
	templateDir    string
	overrideDir    string
	outputPath     string
}

func newGenerateCommand() *cobra.Command {
	opts := &generateOptions{}

	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate K9s plugin YAML for the active cluster",
		RunE: func(cmd *cobra.Command, args []string) error {
			request := generator.Request{
				KubeconfigPath: opts.kubeconfigPath,
				TemplateDir:    opts.templateDir,
				OverrideDir:    opts.overrideDir,
				OutputPath:     opts.outputPath,
			}

			result, err := generator.Generate(request)
			if err != nil {
				return err
			}

			_, err = fmt.Fprintf(cmd.OutOrStdout(), "generated %s for active cluster %s\n", result.OutputPath, result.ActiveCluster)
			return err
		},
	}

	cmd.Flags().StringVar(&opts.kubeconfigPath, "kubeconfig", "", "Path to the kubeconfig file")
	cmd.Flags().StringVar(&opts.templateDir, "template-dir", "", "Path to the template folder")
	cmd.Flags().StringVar(&opts.overrideDir, "override-dir", "", "Path to the override folder")
	cmd.Flags().StringVar(&opts.outputPath, "output", "", "Path to the generated K9s plugin YAML file")

	_ = cmd.MarkFlagRequired("kubeconfig")
	_ = cmd.MarkFlagRequired("template-dir")
	_ = cmd.MarkFlagRequired("override-dir")
	_ = cmd.MarkFlagRequired("output")

	return cmd
}
