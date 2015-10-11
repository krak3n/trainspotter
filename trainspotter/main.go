package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"gopkg.in/stomp.v1"
)

const TOPIC string = "/topic/TD_LNE_GN_SIG_AREA"

// Store in memory current berth state
var berths map[string]string

type Messages []Type

type Type struct {
	CA Message `json:"CA_MSG"`
	CB Message `json:"CB_MSG"`
	CC Message `json:"CC_MSG"`
	CT Message `json:"CT_MSG"`
}

type Message struct {
	AreaID  string `json:"area_id"`
	Descr   string `json:"descr"`
	From    string `json:"from"`
	MsgType string `json:"msg_type"`
	Time    string `json:"time"`
	To      string `json:"to"`
}

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
	// Message hub
	hub := NewHub()
	go hub.Run()

	// Websocket Service
	ws := NewWSService(hub)

	// Connect to STOMP
	conn := fmt.Sprintf("%v:%v", viper.GetString("addr"), viper.GetString("port"))
	s, _ := stomp.Dial("tcp", conn, stomp.Options{
		Login:    viper.GetString("user"),
		Passcode: viper.GetString("key"),
	})

	// Subscribe to the topoc
	sub, _ := s.Subscribe(TOPIC, stomp.AckClient)

	// Get messages
	go func() {
		for {
			msg := <-sub.C
			messages := &Messages{}

			json.Unmarshal(msg.Body, messages)

			fmt.Printf("%+v\n", messages)
		}
	}()

	http.ListenAndServe(":5000", ws)
}

func main() {
	c := cli(flags(), run)
	c.Execute()
}
