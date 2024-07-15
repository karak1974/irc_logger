package main

import (
	"flag"
	"fmt"
	"net"
	"regexp"

	hbot "github.com/whyrusleeping/hellabot"
	log "gopkg.in/inconshreveable/log15.v2"
)

const BOTNAME = "hsbot"

var (
	serv       = flag.String("server", "irc.atw-inter.net:6667", "hostname and port for irc server to connect to")
	nick       = flag.String("nick", "printer", "nickname for the bot")
	enters     = flag.Int("enters", 0, "New lines before messages")
	lastAuthor = ""
)

func extractAuthor(input string) string {
	pattern := `<([^>]*)>`
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(input)

	return matches[1]
}

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
	// How many new line do we want print between messages
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
		bot.Channels = []string{"#camp++printer"}
	}
	irc, err := hbot.NewBot(*serv, *nick, hijackSession, channels)
	if err != nil {
		panic(err)
	}

	irc.AddTrigger(logMessage)
	irc.Logger.SetHandler(log.StdoutHandler)

	irc.Run()
	fmt.Println("Bot shutting down.")
}

var logMessage = hbot.Trigger{
	Condition: func(bot *hbot.Bot, m *hbot.Message) bool {
		return m.Command == "PRIVMSG"
	},
	Action: func(irc *hbot.Bot, m *hbot.Message) bool {
		var msg string
		if m.From == BOTNAME {
			// From bridge
			author := extractAuthor(m.Content)
			if lastAuthor == author {
				msg = m.Content[(len(author) + 3):]
			} else {
				msg = m.Content
			}
			lastAuthor = author
		} else {
			// Directly on irc
			if lastAuthor == m.From {
				msg = m.Content
			} else {
				msg = fmt.Sprintf("<%s> %s", m.From, m.Content)
			}
			lastAuthor = m.From
		}

		printMessage(msg)
		return false
	},
}
