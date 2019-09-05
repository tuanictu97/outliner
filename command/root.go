package command

import (
	"os/user"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	ol "github.com/jyny/outliner/pkg/outliner"
	"github.com/jyny/outliner/pkg/util"

	// deployer agnet
	"github.com/jyny/outliner/pkg/deployer/ssh"

	// cloud provider
	"github.com/jyny/outliner/pkg/cloud/linode"
	//"github.com/jyny/outliner/pkg/cloud/vultr"
	//"github.com/jyny/outliner/pkg/cloud/digitalocean"
)

// Persistent Flags
var cfgFile string
var version = ""

// Persistent for commends
var cloud = ol.NewCloud()
var deployer = ol.NewDeployer()

// RootCmd commands
var RootCmd = &cobra.Command{
	Use:   "outliner",
	Short: "Auto setup & deploy tool for outline VPN server",
	Long:  "Auto setup & deploy tool for outline VPN server",
}

// Execute entry of commandline
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		panic(err)
	}
}

func init() {
	cobra.OnInitialize(initOutliner)
	RootCmd.PersistentFlags().StringVarP(&cfgFile, "file", "F", "", "config file (default is $HOME/.outliner/.env)")

	u, err := user.Current()
	if err != nil {
		panic(err)
	}

	// `.env` as config file name
	viper.SetConfigType("env")
	viper.SetConfigName("")

	// search from possible paths
	viper.AddConfigPath(filepath.Join(u.HomeDir, "/.outliner/"))
	viper.AddConfigPath(u.HomeDir)
	viper.AddConfigPath(".")

	if cfgFile != "" {
		// top precedence order of paths
		viper.SetConfigFile(cfgFile)
	} else {
		// set flag to load config from $ENV
		viper.AutomaticEnv()
	}

	// load config file
	viper.ReadInConfig()
}

func initOutliner() {
	// register deployer agent
	deployer.RegisterAgent(
		ssh.NewAgent(),
	)

	// Activate & register cloud providers
	err := cloud.RegisterProvider(
		util.Validater,
		linode.Activator{},
		//vultr.Activator{},
		//digitalocean.Activator{},
	)
	if err != nil {
		panic(err)
	}
}