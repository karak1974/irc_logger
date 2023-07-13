package main

import (
	"flag"
	"fmt"
	"net"

	hbot "github.com/whyrusleeping/hellabot"
	log "gopkg.in/inconshreveable/log15.v2"
)

var serv = flag.String("server", "irc.ircnet.com:6667", "hostname and port for irc server to connect to")
var nick = flag.String("nick", "camp_printer", "nickname for the bot")
var enters = flag.Int("enters", 7, "New lines before messages")

func printMessage(msg string) {
	conn, err := net.Dial("tcp", "172.16.0.28:9100")
	if err != nil {
		fmt.Println("Error connecting:", err)
		return
	}
	defer conn.Close()

	// Print the message
	_, err = fmt.Fprintf(conn, "%s\n", msg)
	if err != nil {
		fmt.Println("Error sending data:", err)
		return
	}
	// This is necessarry to see the message on the printer
	for i := 0; i < *enters; i++ {
		_, err = fmt.Fprintf(conn, "\n")
		if err != nil {
			fmt.Println("Error sending new line:", err)
			return
		}
	}
}

func main() {
	flag.Parse()

	hijackSession := func(bot *hbot.Bot) {
		bot.HijackSession = true
	}
	channels := func(bot *hbot.Bot) {
		bot.Channels = []string{"#camp++"}
	}
	irc, err := hbot.NewBot(*serv, *nick, hijackSession, channels)
	if err != nil {
		panic(err)
	}

	irc.AddTrigger(sayInfoMessage)
	irc.AddTrigger(logMessage)
	irc.Logger.SetHandler(log.StdoutHandler)

	irc.Run()
	fmt.Println("Bot shutting down.")
}

var sayInfoMessage = hbot.Trigger{
	Condition: func(bot *hbot.Bot, m *hbot.Message) bool {
		return m.Command == "PRIVMSG" && m.Content == "-info"
	},
	Action: func(irc *hbot.Bot, m *hbot.Message) bool {
		// Reply, but also logged
		irc.Reply(m, "I'm a printer in the hackcenter, Wolfy is my creator")
		return false
	},
}

var logMessage = hbot.Trigger{
	Condition: func(bot *hbot.Bot, m *hbot.Message) bool {
		return m.Command == "PRIVMSG"
	},
	Action: func(irc *hbot.Bot, m *hbot.Message) bool {
		var msg string
		if m.To == *nick {
			msg = fmt.Sprintf("<%s><%s>: %s", "Private", m.From, m.Content)
		} else {
			msg = fmt.Sprintf("<%s><%s>: %s", m.To, m.From, m.Content)
		}

		fmt.Println(msg)
		printMessage(msg)
		return false
	},
}
