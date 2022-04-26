package cmd

import (
	"github.com/dfds/kiam2irsa/pkg/k8s/sa"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var saCmd = &cobra.Command{
	Use:   "sa",
	Short: "Find Kubernetes ServiceAccounts with certain annotations",
	Run:   func(cmd *cobra.Command, args []string) {},
}

func saInit() {
	saCmd.PersistentFlags().StringP("kubeconfig", "f", "", "Full path to the kubeconfig file")
}

func CheckSA(cmd *cobra.Command, args []string) {
	logger, _ := zap.NewDevelopment()
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
		}
	}(logger)
	sugar := logger.Sugar()

	kubeconfig, err := saCmd.Flags().GetString("kubeconfig")
	if err != nil {
		sugar.Error(err.Error())
		return
	}

	sa.GetSA(kubeconfig)
}
