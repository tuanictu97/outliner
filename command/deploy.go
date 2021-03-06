package command

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	deployCmd.Flags().StringP("ip", "i", "", "IP address of Server (required)")
	deployCmd.MarkFlagRequired("ip")
	viper.BindPFlag("deploy_ip", deployCmd.Flags().Lookup("ip"))
	RootCmd.AddCommand(deployCmd)
}

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "deploy outliner to Server",
	Long:  `deploy outliner to Server`,
	PreRun: func(cmd *cobra.Command, args []string) {
		printDeployInstanceStart()
	},
	Run: func(cmd *cobra.Command, args []string) {
		ip := viper.GetString("deploy_ip")
		err := deployer.DeployService(ip)
		if err != nil {
			panic(err)
		}
		printDeployInstanceWait()
		err = deployer.WaitService(ip)
		if err != nil {
			panic(err)
		}
		printDeployInstanceDone()
		viper.Set("inspect_ip", ip)

	},
	PostRun: func(cmd *cobra.Command, args []string) {
		inspectCmd.PreRun(inspectCmd, []string{})
		inspectCmd.Run(inspectCmd, []string{})
	},
}
