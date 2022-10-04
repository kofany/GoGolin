package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	irc "github.com/thoj/go-ircevent"
	"golang.org/x/exp/slices"
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
	//owner := config["owner"]
	botnick := config["botnick"]
	ident := config["ident"]
	ircobj := irc.IRC(botnick, ident)
	ircobj.RealName = config["realname"]
	ircobj.Version = "GoGolin v 0.3 - irc client in Go"
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

	ircobj.AddCallback("JOIN", func(e *irc.Event) {
		//owner
		lines, err := readLines("owner.txt")
		if err != nil {
			return
		}
		owner := lines
		var isOwner string = strings.TrimPrefix(e.Source, e.Nick)
		var sliceChan []string = strings.Split(e.Raw, " ")
		var curChan string = sliceChan[2]
		curChan = strings.TrimPrefix(curChan, ":")
		if slices.Contains(owner, isOwner) {
			var docommand string = ("MODE " + curChan + " +o " + e.Nick)
			ircobj.SendRawf(docommand)

		}
	})
	ircobj.AddCallback("PRIVMSG", func(e *irc.Event) {
		//owner
		lines, err := readLines("owner.txt")
		if err != nil {
			return
		}

		owner := lines
		var result string = e.Message()
		var isOwner string = strings.TrimPrefix(e.Source, e.Nick)
		var sliceChan []string = strings.Split(e.Raw, " ")
		var curChan string = sliceChan[2]
		//check if it is a owner
		if slices.Contains(owner, isOwner) {
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
			// !s command
			if strings.Contains(result, "!say ") {
				result = strings.TrimPrefix(result, "!say ")
				ircobj.Privmsg(curChan, result)
			}
			// !a command
			if strings.Contains(result, "!a ") {
				ircobj.Notice(e.Nick, "Welcom my master!")
			}
			// !msg command
			if strings.Contains(result, "!msg ") {
				value := strings.TrimPrefix(result, "!msg ")
				target := firstWords(value, 1)
				result := strings.TrimPrefix(value, target)
				ircobj.Privmsg(target, result)

			}
			// !+owner
			if strings.Contains(result, "!+owner ") {
				result := strings.TrimPrefix(result, "!+owner ")
				if slices.Contains(owner, result) {
					ircobj.Notice(e.Nick, "Sorry "+result+" exists on my list - not adding")
				} else {
					f, err := os.OpenFile("owner.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
					if err != nil {
						ircobj.Notice(e.Nick, "Something wrong with owner.text")
					}
					defer f.Close()
					if _, err := f.WriteString("\n" + result); err != nil {
						ircobj.Notice(e.Nick, "Something wrong with owner.text")
					}
				}
			}
			// !-owner
			if strings.Contains(result, "!-owner ") {
				result := strings.TrimPrefix(result, "!-owner ")
				if slices.Contains(owner, result) {

					lines, err := readLines("owner.txt")
					if err != nil {
						return
					}
					for i, v := range lines {
						if v == result {
							lines = append(lines[:i], lines[i+1:]...)
						}
					}
					lines = delete_empty(lines)
					if err := writeLines(lines, "owner.txt"); err != nil {
						return
					}
				} else {
					ircobj.Notice(e.Nick, "Sorry "+result+" not exists on my list")
				}
			}
			if strings.Contains(result, "!owners") {
				lines, err := readLines("owner.txt")
				if err != nil {
					return
				}
				for _, s := range lines {
					ircobj.Notice(e.Nick, "Owner: "+s)
				}
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
				time.Sleep(1 * time.Second)
				ircobj.Privmsg(e.Nick, "!+owner add owner to bot")
				ircobj.Privmsg(e.Nick, "!-owner delte owner from bot")
				ircobj.Privmsg(e.Nick, "!a - bot says hallo to You")
				ircobj.Privmsg(e.Nick, "!owners - notice owners list")

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
func readLines(path string) (lines []string, err error) {
	var (
		file   *os.File
		part   []byte
		prefix bool
	)
	if file, err = os.Open(path); err != nil {
		return
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	buffer := bytes.NewBuffer(make([]byte, 0))
	for {
		if part, prefix, err = reader.ReadLine(); err != nil {
			break
		}
		buffer.Write(part)
		if !prefix {
			lines = append(lines, buffer.String())
			buffer.Reset()
		}
	}
	if err == io.EOF {
		err = nil
	}
	return
}

func writeLines(lines []string, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for _, line := range lines {
		fmt.Fprintln(w, line)
	}
	return w.Flush()
}
func delete_empty(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}
