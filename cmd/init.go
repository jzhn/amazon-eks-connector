package main

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s.io/klog/v2"

	"golang.a2z.com/eks-connector/pkg/agent"
	"golang.a2z.com/eks-connector/pkg/config"
	"golang.a2z.com/eks-connector/pkg/initializer"
	"golang.a2z.com/eks-connector/pkg/k8s"
	"golang.a2z.com/eks-connector/pkg/ssm"
	"golang.a2z.com/eks-connector/pkg/state"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize EKS connector",
	Run: func(cmd *cobra.Command, args []string) {
		configProvider := config.NewProvider()
		configuration, err := configProvider.Get()
		if err != nil {
			klog.Fatalf("failed to load configuration: %v", err)
		}

		ssmService := ssm.NewClient(configuration.AgentConfig)
		registration := agent.NewRegistration(ssmService, configuration.ActivationConfig)
		fsPersistence := state.NewFileSystemPersistence(configuration.StateConfig)
		secret, err := k8s.NewSecretInCluster(configuration.StateConfig)
		if err != nil {
			klog.Fatalf("failed to initiate kubernetes client: %v", err)
		}
		secretPersistence := state.NewSecretPersistence(secret)

		initer := initializer.NewInitializer(
			secretPersistence,
			fsPersistence,
			registration,
		)

		if err = initer.Initialize(); err != nil {
			klog.Fatalf("failed to initiate eks-connector: %v", err)
		}
	},
}

func init() {
	initCmd.Flags().String("agent.region",
		"us-west-2",
		"The AWS region that EKS connector agent communicates to")
	initCmd.Flags().String("agent.endpoint",
		"",
		"The SSM endpoint that EKS connector agent communicates to")
	initCmd.Flags().String("activation.id",
		"",
		"EKS connector activationId, as provided by RegisterCluster API")
	initCmd.Flags().String("activation.code",
		"",
		"EKS connector activationCode, as provided by RegisterCluster API")
	initCmd.Flags().String("state.baseDir",
		state.DirSsmVault,
		"The vault folder of ssm agent container")
	initCmd.Flags().String("state.secretNamePrefix",
		"eks-connector-state",
		"Prefix of Kubernetes Secret name used to persist eks-connector state")
	initCmd.Flags().String("state.secretNamespace",
		"eks-connector",
		"Kubernetes namespace of the Secret used to persist eks-connector state")
	_ = initCmd.MarkFlagRequired("activation.id")
	_ = initCmd.MarkFlagRequired("activation.code")

	err := viper.BindPFlags(initCmd.Flags())
	if err != nil {
		klog.Fatal("failed to bind cmd flags: %v", err)
	}

	rootCmd.AddCommand(initCmd)
}
