package cmd

import (
	"fmt"
	"os"

	"github.com/dfds/kiam2irsa/pkg/k8s"
	"github.com/dfds/kiam2irsa/pkg/logging"

	"github.com/spf13/cobra"
)

var podCmd = &cobra.Command{
	Use:   "pods",
	Short: "Find Kubernetes Pods that is still not migrated away from KIAM",
	Run: func(cmd *cobra.Command, args []string) {
		k8s.CheckPodsMigrationStatus(cmd)
	},
}

func podInit() {
	sugar := logging.SugarLogger()

	// Setting a default value for kubeconfig
	homeDir, err := os.UserHomeDir()
	if err != nil {
		sugar.Error(err.Error())
		return
	}
	kubeconfig, exist := os.LookupEnv("KUBECONFIG")
	if !exist {
		kubeconfig = fmt.Sprintf("%s/.kube/config", homeDir)
	}

	podCmd.PersistentFlags().StringP("kubeconfig", "f", kubeconfig, "Full path to the kubeconfig file")
	podCmd.PersistentFlags().StringP("status", "s", "KIAM", "Migration status supports: KIAM, IRSA, BOTH")
	podCmd.PersistentFlags().StringP("output", "o", "text", "Output format supports: TEXT,CSV")
}
