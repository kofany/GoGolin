package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	irc "github.com/thoj/go-ircevent"
)

type Config map[string]string

func main() {
	//config
	config, err := ReadConfig(`config.txt`)
	if err != nil {
		fmt.Println(err)
	}
	// assign values from config file to variables
	secretChan := config["secretChan"]
	server := config["server"]
	port := config["port"]
	owner := config["owner"]
	botnick := config["botnick"]
	ident := config["ident"]
	owner = ("!" + owner)
	ircobj := irc.IRC(botnick, ident)
	ircobj.RealName = config["realname"]
	errCon := ircobj.Connect(server + ":" + port)
	if errCon != nil {
		fmt.Println("Failed connecting")
		return
	}
	ircobj.AddCallback("001", func(e *irc.Event) {
		ircobj.Join(secretChan)
	})
	//what is my bot nick on irc
	var mynick string = ircobj.GetNick()

	// ircobj.AddCallback("JOIN", func(e *irc.Event) {
	//	ircobj.Privmsg(secretChan, "Hello! I am a friendly IRC bot who will echo everything you say.")
	// })

	ircobj.AddCallback("PRIVMSG", func(e *irc.Event) {
		var result string = e.Message()
		var isOwner string = e.Source
		var sliceChan []string = strings.Split(e.Raw, " ")
		var curChan string = sliceChan[2]
		//check if it is a owner
		if strings.Contains(isOwner, owner) {
			// !op command
			if strings.Contains(result, "!op ") {

				result = strings.TrimPrefix(result, "!op ")
				var docommand string = ("MODE " + curChan + " +ooo " + result)
				ircobj.SendRawf(docommand)

			}
			// !6 command op x 6

			if strings.Contains(result, "!6 ") {

				result = strings.TrimPrefix(result, "!6 ")
				var sliceOP []string = strings.Split(result, " ")
				var op1 string = ("MODE " + curChan + " +ooo " + sliceOP[0] + " " + sliceOP[1] + " " + sliceOP[2])
				var op2 string = ("MODE " + curChan + " +ooo " + sliceOP[3] + " " + sliceOP[4] + " " + sliceOP[5])

				ircobj.SendRawf(op1)
				ircobj.SendRawf(op2)

			}
			// !deop command
			if strings.Contains(result, "!dop ") {
				result = strings.TrimPrefix(result, "!dop ")
				var docommand string = ("MODE " + curChan + " -ooo " + result)
				ircobj.SendRawf(docommand)
			}
			// !+v command
			if strings.Contains(result, "!+v ") {
				result = strings.TrimPrefix(result, "!+v ")
				var docommand string = ("MODE " + curChan + " +vvv " + result)
				ircobj.SendRawf(docommand)
			}
			// !-v command
			if strings.Contains(result, "!-v ") {
				result = strings.TrimPrefix(result, "!-v ")
				var docommand string = ("MODE " + curChan + " -vvv " + result)
				ircobj.SendRawf(docommand)

			}
			// !+b command
			if strings.Contains(result, "!+b ") {
				result = strings.TrimPrefix(result, "!+b ")
				var docommand string = ("MODE " + curChan + " +b " + result)
				ircobj.SendRawf(docommand)
			}
			// !-b command
			if strings.Contains(result, "!-b ") {
				result = strings.TrimPrefix(result, "!-b ")
				var docommand string = ("MODE " + curChan + " -b " + result)
				ircobj.SendRawf(docommand)
			}
			// !j command
			if strings.Contains(result, "!j ") {
				result = strings.Trim(result, "!j ")
				ircobj.Join(result)
			}
			// !p command
			if strings.Contains(result, "!p ") {
				result = strings.TrimPrefix(result, "!p ")
				ircobj.Part(result)
			}
			// !die command
			if strings.Contains(result, "!die ") {
				result = strings.TrimPrefix(result, "!die ")
				ircobj.QuitMessage = result
				ircobj.Quit()
				//os.Exit(1)
			}
			// !raw command
			if strings.Contains(result, "!raw") {
				ircobj.Privmsg(e.Nick, e.Raw)
				ircobj.Privmsg(e.Nick, mynick)
			}
			// !s command
			if strings.Contains(result, "!say ") {
				result = strings.TrimPrefix(result, "!say ")
				ircobj.Privmsg(curChan, result)
			}
			// !msg command
			if strings.Contains(result, "!msg ") {
				value := strings.TrimPrefix(result, "!msg ")
				target := firstWords(value, 1)
				result := strings.TrimPrefix(value, target)
				ircobj.Privmsg(target, result)

			}
			// !help command
			if strings.Contains(result, "!help") && strings.Contains(e.Raw, mynick) {
				ircobj.Privmsg(e.Nick, "My Commands: ")
				time.Sleep(1 * time.Second)
				ircobj.Privmsg(e.Nick, "!op - gives op up to 3 nicks - !op nick1 nick2 nick3")
				ircobj.Privmsg(e.Nick, "!6 - gives op up to 6 nicks - !6 nick1 nick2 nick3 nick4 nick5 nick6")
				time.Sleep(1 * time.Second)
				ircobj.Privmsg(e.Nick, "!dop - takes op up to 3 nicks - !dop nick1 nick2 nick3")
				time.Sleep(1 * time.Second)
				ircobj.Privmsg(e.Nick, "!+v - gives voice up to 3 nicks - !+v nick1 nick2 nick3")
				ircobj.Privmsg(e.Nick, "!-v - takes voice up to 3 nicks - !-v nick1 nick2 nick3")
				ircobj.Privmsg(e.Nick, "!+b - ban usermask - !+b *!ident@host")
				time.Sleep(1 * time.Second)
				ircobj.Privmsg(e.Nick, "!-b - unban usermask - !-b *!ident@host")
				ircobj.Privmsg(e.Nick, "!j #channel - joining #channel")
				time.Sleep(1 * time.Second)
				ircobj.Privmsg(e.Nick, "!p #channel - part channel ")
				time.Sleep(1 * time.Second)
				ircobj.Privmsg(e.Nick, "!die quit_text - killing bot - !die Im going... ")
				ircobj.Privmsg(e.Nick, "!msg nick text - send priv msg to nick")
			}
		}
	})
	ircobj.Loop()
}

func ReadConfig(filename string) (Config, error) {
	// init with some bogus data
	config := Config{
		"secretChan": "",
		"server":     "",
		"port":       "",
		"owner":      "",
	}
	if len(filename) == 0 {
		return config, nil
	}
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	for {
		line, err := reader.ReadString('\n')

		// check if the line has = sign
		// and process the line. Ignore the rest.
		if equal := strings.Index(line, "="); equal >= 0 {
			if key := strings.TrimSpace(line[:equal]); len(key) > 0 {
				value := ""
				if len(line) > equal {
					value = strings.TrimSpace(line[equal+1:])
				}
				// assign the config map
				config[key] = value
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
	}
	return config, nil
}
func firstWords(value string, count int) string {
	// Loop over all indexes in the string.
	for i := range value {
		// If we encounter a space, reduce the count.
		if value[i] == ' ' {
			count -= 1
			// When no more words required, return a substring.
			if count == 0 {
				return value[0:i]
			}
		}
	}
	// Return the entire string.
	return value
}
