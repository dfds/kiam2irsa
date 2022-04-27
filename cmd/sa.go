package cmd

import (
	"fmt"
	"os"

	"github.com/dfds/kiam2irsa/pkg/k8s/sa"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var saCmd = &cobra.Command{
	Use:   "sa",
	Short: "Find Kubernetes ServiceAccounts with certain annotations",
	Run: func(cmd *cobra.Command, args []string) {
		sa.GetSA(cmd, args)
	},
}

func saInit() {
	logger, _ := zap.NewDevelopment()
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
		}
	}(logger)
	sugar := logger.Sugar()

	homeDir, err := os.UserHomeDir()
	if err != nil {
		sugar.Error(err.Error())
		return
	}

	kubeconfig := fmt.Sprintf("%s/.kube/config", homeDir) // Default value

	saCmd.PersistentFlags().StringP("kubeconfig", "f", kubeconfig, "Full path to the kubeconfig file")
}
