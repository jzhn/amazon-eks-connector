package main

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s.io/klog/v2"

	"golang.a2z.com/eks-connector/pkg/config"
	"golang.a2z.com/eks-connector/pkg/proxy"
	"golang.a2z.com/eks-connector/pkg/server"
	"golang.a2z.com/eks-connector/pkg/serviceaccount"
)

var serverCmd = &cobra.Command{
	Use:     "server",
	Short:   "Run EKS connector proxy server",
	Example: "",
	Run: func(cmd *cobra.Command, args []string) {
		configProvider := config.NewProvider()
		configuration, err := configProvider.Get()
		if err != nil {
			klog.Fatalf("failed to load configuration: %v", err)
		}

		secretProvider := serviceaccount.NewProvider()

		server := server.Server{
			ProxyConfig:  configuration.ProxyConfig,
			ProxyHandler: proxy.NewProxyHandler(configuration.ProxyConfig, secretProvider),
		}

		server.Run()
	},
}

func init() {
	serverCmd.Flags().String("proxy.socketType",
		"unix",
		"The socket type of proxy. Can be 'unix' or 'tcp'")
	serverCmd.Flags().String("proxy.socketAddr",
		"/var/eks/shared/connector.sock",
		"The address of proxy, should be a FS path or network address depending on socket type")
	serverCmd.Flags().String("proxy.targetHost",
		"kubernetes.default.svc:443",
		"The target of the proxy, should be api server's address")
	serverCmd.Flags().String("proxy.targetProtocol",
		"https",
		"The target protocol of the proxy. Can be 'https' or 'http'")

	err := viper.BindPFlags(serverCmd.Flags())
	if err != nil {
		klog.Fatal("failed to bind cmd flags: %v", err)
	}

	rootCmd.AddCommand(serverCmd)
}
