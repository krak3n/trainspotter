package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func flags() *pflag.FlagSet {
	// Create flagset
	set := pflag.NewFlagSet(os.Args[0], pflag.ExitOnError)

	// Add arguments
	set.StringP("user", "u", "", "User")
	set.StringP("key", "k", "", "Key")
	set.StringP("addr", "a", "", "Address")
	set.StringP("port", "p", "", "Port")

	return set
}

func cli(flags *pflag.FlagSet, run func(cmd *cobra.Command, args []string)) *cobra.Command {
	// Create command
	cmd := &cobra.Command{
		Use:   "trainsporter",
		Short: "Trainspotter",
		Long:  "Trainspotter",
		Run:   run,
	}

	// Register flags with command
	cmd.Flags().AddFlagSet(flags)

	// Bind to viper
	viper.BindPFlags(flags)

	return cmd
}

func main() {
	// Create CLI
	c := cli(flags(), func(cmd *cobra.Command, args []string) { // Method called when Tamzin is invoked
		fmt.Println("Hello World")
	})

	// Invoke CLI
	c.Execute()
}
