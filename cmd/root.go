package cmd

import (
  "fmt"
  "os"
  "strings"

  homedir "github.com/mitchellh/go-homedir"
  "github.com/spf13/cobra"
  "github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
  Use: "kld",
  Short: "Command Line Tool for Kaleido resources management",
  Run: func(cmd *cobra.Command, args []string) {
    fmt.Printf("API URL: %s\n", viper.Get("api.url"))
  },
}

var cfgFile string

func Execute() {
  if err := rootCmd.Execute(); err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
}

func init() {
  cobra.OnInitialize(initConfig)

  // all environment variables for the "kld" command will have the "KLD" prefix
  // e.g "KLD_API_URL"
  viper.SetEnvPrefix("kld")
  // allows code to access env variables by name only, without the prefix
  viper.AutomaticEnv()
  // allows using "." to access env variables with "_"
  // e.g viper.Get('api.url') for value of "KLD_API_URL"
  viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

  // config files capture defaults that can be overwritten by env variables and flags
  rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "Config file that captures re-usable settings such as API URl, API Key, etc. (default is $HOME/.kld.yaml)")
  viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))

  rootCmd.AddCommand(newCreateCmd())
  rootCmd.AddCommand(newDeleteCmd())
}

func initConfig() {
  // Don't forget to read config either from cfgFile or from home directory!
  if cfgFile != "" {
    // Use config file from the flag.
    viper.SetConfigFile(cfgFile)
  } else {
    // Find home directory.
    home, err := homedir.Dir()
    if err != nil {
      fmt.Println(err)
      os.Exit(1)
    }

    // Search config in home directory with name ".kld" (without extension).
    viper.AddConfigPath(home)
    viper.SetConfigName(".kld")
  }

  if err := viper.ReadInConfig(); err != nil {
    fmt.Printf("\nCan't read config: %v, will rely on environment variables for required configurations\n", err)
  }
}