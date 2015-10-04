package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"gopkg.in/stomp.v1"
)

const TOPIC string = "/topic/TD_LNE_GN_SIG_AREA"

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

func run(cmd *cobra.Command, args []string) {
	// Connect to STOMP
	conn := fmt.Sprintf("%v:%v", viper.GetString("addr"), viper.GetString("port"))
	s, _ := stomp.Dial("tcp", conn, stomp.Options{
		Login:    viper.GetString("user"),
		Passcode: viper.GetString("key"),
	})

	// Subscribe to the topoc
	sub, _ := s.Subscribe(TOPIC, stomp.AckClient)

	// Get messages
	for {
		msg := <-sub.C
		fmt.Println(string(msg.Body[:]))
	}
}

func main() {
	c := cli(flags(), run)
	c.Execute()
}
