package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	hbot "github.com/whyrusleeping/hellabot"
)

var serv = flag.String("server", "irc.wolfy.me:6667", "hostname and port for irc server to connect to")
var nick = flag.String("nick", "irc_printer_bot", "nickname for the bot")

// -info trigger
var infoMessage = hbot.Trigger{
	Condition: func(bot *hbot.Bot, m *hbot.Message) bool {
		return m.Command == "PRIVMSG" && m.Content == "-info"
	},
	Action: func(irc *hbot.Bot, m *hbot.Message) bool {
		irc.Reply(m, "An IRC bot that do logging")
		return false
	},
}

// Logger
var logger = hbot.Trigger{
	Condition: func(bot *hbot.Bot, m *hbot.Message) bool {
		return m.Command == "PRIVMSG" && m.Content != "-info"
	},
	Action: func(irc *hbot.Bot, m *hbot.Message) bool {
		server := *serv
		go func(server string, channel string, user string, content string) {
			file, err := os.OpenFile(server+"."+channel[1:]+".log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
			if err != nil {
				log.Fatal(err)
			}
			log.SetOutput(file)
		
			fmt.Println("<"+user+">: "+content)
			log.Println("<"+user+">: "+content)
		}(server[:len(server)-5], m.To, m.From, m.Content)
		return false
	},
}

func main() {
	flag.Parse()

	hijackSession := func(bot *hbot.Bot) {
		bot.HijackSession = true
	}
	channels := func(bot *hbot.Bot) {
		bot.Channels = []string{"#channel1", "#channel2"}
	}
	irc, err := hbot.NewBot(*serv, *nick, hijackSession, channels)
	if err != nil {
		panic(err)
	}

	irc.AddTrigger(infoMessage)
	irc.AddTrigger(logger)

	irc.Run()
	fmt.Println("Bot shutting down.")
}
