package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"gopkg.in/stomp.v2"
	"gopkg.in/stomp.v2/frame"
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

func feed(addr string, user string, key string) (*stomp.Subscription, error) {
	// Connect to Stomp Service
	conn, err := stomp.Dial("tcp", addr,
		stomp.ConnOpt.Login(user, key),
		stomp.ConnOpt.AcceptVersion(stomp.V11),
		stomp.ConnOpt.AcceptVersion(stomp.V12),
		stomp.ConnOpt.Header(frame.NewHeader("client-id", user)))
	if err != nil {
		return nil, err
	}

	// Subscribe to Feed
	subscriptionName := frame.NewHeader("activemq.subscriptionName", "trainspotter-td-gn")
	sub, err := conn.Subscribe(TOPIC, stomp.AckClient, stomp.SubscribeOpt.Header(subscriptionName))
	if err != nil {
		return nil, err
	}

	return sub, nil
}

func run(cmd *cobra.Command, args []string) {
	// Message hub
	hub := NewHub()
	go hub.Run()

	// Websocket Service
	ws := NewWSService(hub)

	// Connect to STOMP
	addr := fmt.Sprintf("%v:%v", viper.GetString("addr"), viper.GetString("port"))
	sub, err := feed(addr, viper.GetString("user"), viper.GetString("key"))
	if err != nil {
		panic(err)
	}

	// Get messages
	go func() {
		for {
			msg := <-sub.C

			if err := msg.Conn.Ack(msg); err != nil{
				fmt.Println("Faild to ACK Message")
			}

			fmt.Println("-----------------------------")
			if msg.Body != nil {
				fmt.Println(string(msg.Body[:]))
				Process(msg.Body, hub)
			}
			fmt.Println("-----------------------------")
		}
	}()

	http.ListenAndServe(":5000", ws)
}

func main() {
	c := cli(flags(), run)
	c.Execute()
}
