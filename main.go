package main

import (
	"os"
	"path/filepath"

	"github.com/dfds/kiam2irsa/pkg/k8s/sa"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewDevelopment()
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
		}
	}(logger)
	sugar := logger.Sugar()
	home, err := os.UserHomeDir()
	if err != nil {
		sugar.Error(err.Error())
		return
	}

	// TODO: Replace this dirty hack with the cobra CLI
	kubeconfig := filepath.Join(home, ".kube", "config")
	args := os.Args
	if len(args) >= 2 {
		kubeconfig = args[1]
	}

	sa.GetSA(kubeconfig)
}
