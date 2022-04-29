package cmd

import (
	"fmt"
	"github.com/dfds/kiam2irsa/pkg/k8s"
	"github.com/dfds/kiam2irsa/pkg/logging"
	"os"

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
	podCmd.PersistentFlags().BoolP("parallelism", "p", false, "Use goroutines to make requests in parallel?")
}
