package cmd

import (
	"github.com/dfds/kiam2irsa/pkg/k8s/sa"
	"github.com/spf13/cobra"
)

var saCmd = &cobra.Command{
	Use:   "sa",
	Short: "Find Kubernetes ServiceAccounts with certain annotations",
	Run: func(cmd *cobra.Command, args []string) {
		sa.GetSA(cmd, args)
	},
}

func saInit() {
	saCmd.PersistentFlags().StringP("kubeconfig", "f", "", "Full path to the kubeconfig file")
}
