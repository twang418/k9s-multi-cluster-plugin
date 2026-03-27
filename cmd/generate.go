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
	installToK9s   bool
	k9sDataDir     string
}

func newGenerateCommand() *cobra.Command {
	opts := &generateOptions{}

	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate K9s plugin YAML for the active cluster",
		RunE: func(cmd *cobra.Command, args []string) error {
			outputMode := generator.OutputModeFile
			if opts.installToK9s {
				outputMode = generator.OutputModeK9s
			}

			request := generator.Request{
				KubeconfigPath: opts.kubeconfigPath,
				TemplateDir:    opts.templateDir,
				OverrideDir:    opts.overrideDir,
				OutputPath:     opts.outputPath,
				OutputMode:     outputMode,
				K9sDataDir:     opts.k9sDataDir,
			}

			result, err := generator.Generate(request)
			if err != nil {
				return err
			}

			if outputMode == generator.OutputModeK9s {
				_, err = fmt.Fprintf(cmd.OutOrStdout(), "generated %d K9s plugin file(s) for active cluster %s under %s\n", len(result.OutputPaths), result.ActiveCluster, result.InstallRoot)
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
	cmd.Flags().BoolVar(&opts.installToK9s, "install-to-k9s", false, "Write merged plugins into K9s context-specific plugin files")
	cmd.Flags().StringVar(&opts.k9sDataDir, "k9s-data-dir", "", "Override the K9s data-home root used for install mode")

	_ = cmd.MarkFlagRequired("kubeconfig")
	_ = cmd.MarkFlagRequired("template-dir")
	_ = cmd.MarkFlagRequired("override-dir")

	return cmd
}
