package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"time"

	irc "github.com/kofany/go-ircevent"
	"golang.org/x/exp/slices"
)

type Config map[string]string

var ownerfile string = "owner.txt"
var aopfile string = "aop.txt"
var shitfile string = "shit.txt"
var configfile string = "config.txt"

func main() {
	createConfig(configfile)
	createOwner(ownerfile)
	createAS(aopfile)
	createAS(shitfile)
	//config
	config, err := ReadConfig(`config.txt`)
	if err != nil {
		fmt.Println(err)
	}
	// assign values from config file to variables
	secretChan := config["secretChan"]
	server := config["server"]
	port := config["port"]
	botnick := config["botnick"]
	ident := config["ident"]
	myhost := config["myhost"]
	ircobj := irc.IRC(botnick, ident, myhost)
	ircobj.PingFreq = 1 * time.Minute
	ircobj.RealName = config["realname"]
	ircobj.Version = "GoGolin v 1.1 - irc client in Go"
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
	//autoop list
	ircobj.AddCallback("JOIN", func(e *irc.Event) {
		lines, err := readLines("aop.txt")
		if err != nil {
			return
		}
		isOp := lines
		var toOp string = strings.TrimPrefix(e.Source, e.Nick)
		var sliceChan []string = strings.Split(e.Raw, " ")
		var curChan string = sliceChan[2]
		curChan = strings.TrimPrefix(curChan, ":")
		tooOp := (curChan + " *" + toOp)
		if slices.Contains(isOp, tooOp) {
			ircobj.Mode(curChan, " +o "+e.Nick)
		}
	})
	//shit list
	ircobj.AddCallback("JOIN", func(e *irc.Event) {
		lines, err := readLines("shit.txt")
		if err != nil {
			return
		}
		var sliceChan []string = strings.Split(e.Raw, " ")
		var curChan string = sliceChan[2]
		curChan = strings.TrimPrefix(curChan, ":")
		isShit := lines
		var toShit string = strings.TrimPrefix(e.Source, e.Nick)
		if strings.Contains(toShit, "~") {
			toShit = strings.TrimPrefix(toShit, "!~")
			toShit = ("!" + toShit)
		}
		toShitList := (curChan + " *" + toShit)
		if slices.Contains(isShit, toShitList) {
			if strings.Contains(e.User, "~") {
				e.User = "*"
			}
			ircobj.Mode(curChan, "+b "+"*!"+e.User+"@"+e.Host)
			ircobj.Kick(e.Nick, curChan, "You are on my shitlist! Bye Bye")
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
		isOwner = ("*" + isOwner)
		var sliceChan []string = strings.Split(e.Raw, " ")
		var curChan string = sliceChan[2]
		//check if it is a owner
		if slices.Contains(owner, isOwner) {
			// !op command
			if strings.HasPrefix(result, "!+o ") {

				result = strings.TrimPrefix(result, "!+o ")
				var docommand string = ("MODE " + curChan + " +ooo " + result)
				ircobj.SendRawf(docommand)

			}
			// !6 command op x 6

			if strings.HasPrefix(result, "!6 ") {

				result = strings.TrimPrefix(result, "!6 ")
				var sliceOP []string = strings.Split(result, " ")
				var op1 string = ("MODE " + curChan + " +ooo " + sliceOP[0] + " " + sliceOP[1] + " " + sliceOP[2])
				var op2 string = ("MODE " + curChan + " +ooo " + sliceOP[3] + " " + sliceOP[4] + " " + sliceOP[5])
				ircobj.SendRawf(op1)
				ircobj.SendRawf(op2)

			}
			// !deop command
			if strings.HasPrefix(result, "!-o ") {
				result = strings.TrimPrefix(result, "!-o ")
				var docommand string = ("MODE " + curChan + " -ooo " + result)
				ircobj.SendRawf(docommand)
			}
			// !+v command
			if strings.HasPrefix(result, "!+v ") {
				result = strings.TrimPrefix(result, "!+v ")
				var docommand string = ("MODE " + curChan + " +vvv " + result)
				ircobj.SendRawf(docommand)
			}
			// !-v command
			if strings.HasPrefix(result, "!-v ") {
				result = strings.TrimPrefix(result, "!-v ")
				var docommand string = ("MODE " + curChan + " -vvv " + result)
				ircobj.SendRawf(docommand)

			}
			// !+b command
			if strings.HasPrefix(result, "!+b ") {
				result = strings.TrimPrefix(result, "!+b ")
				var docommand string = ("MODE " + curChan + " +b " + result)
				ircobj.SendRawf(docommand)
			}
			// !-b command
			if strings.HasPrefix(result, "!-b ") {
				result = strings.TrimPrefix(result, "!-b ")
				var docommand string = ("MODE " + curChan + " -b " + result)
				ircobj.SendRawf(docommand)
			}
			// !j command
			if strings.HasPrefix(result, "!j ") {
				result = strings.Trim(result, "!j ")
				ircobj.Join(result)
			}
			// !p command
			if strings.HasPrefix(result, "!p ") {
				result = strings.TrimPrefix(result, "!p ")
				ircobj.Part(result)
			}
			// !die command
			if strings.HasPrefix(result, "!die ") {
				result = strings.TrimPrefix(result, "!die ")
				ircobj.QuitMessage = result
				ircobj.Quit()
				//os.Exit(1)
			}
			// !s command
			if strings.HasPrefix(result, "!s ") {
				result = strings.TrimPrefix(result, "!s ")
				ircobj.Privmsg(curChan, result)
			}
			// !k command
			if strings.HasPrefix(result, "!k ") {
				result = strings.TrimPrefix(result, "!k ")
				ircobj.Kick(result, curChan, "Not welcome here!")
			}
			// !a command
			if strings.HasPrefix(result, "!hi") {
				ircobj.Notice(e.Nick, "Welcome my master!")
			}
			// !msg command
			if strings.HasPrefix(result, "!msg ") {
				value := strings.TrimPrefix(result, "!msg ")
				target := firstWords(value, 1)
				result := strings.TrimPrefix(value, target)
				ircobj.Privmsg(target, result)
			}
			// !+owner
			if strings.HasPrefix(result, "!+owner ") {
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
			if strings.HasPrefix(result, "!-owner ") {
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
			//owners
			if strings.HasPrefix(result, "!owners") {
				lines, err := readLines("owner.txt")
				if err != nil {
					return
				}
				for _, s := range lines {
					ircobj.Notice(e.Nick, "Owner: "+s)
				}
			}
			//!+aop
			if strings.HasPrefix(result, "!+aop ") {

				lines, err := readLines("aop.txt")
				if err != nil {
					return
				}
				isOp := lines

				result := strings.TrimPrefix(result, "!+aop ")
				if slices.Contains(isOp, result) {
					ircobj.Notice(e.Nick, "Sorry "+result+" exists on my list - not adding")
				} else {
					f, err := os.OpenFile("aop.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
					if err != nil {
						ircobj.Notice(e.Nick, "Something wrong with aop.txt")
					}
					defer f.Close()
					if _, err := f.WriteString(result + "\n"); err != nil {
						ircobj.Notice(e.Nick, "Something wrong with aop.txt")
					}
				}
				lines2, err := readLines("aop.txt")
				if err != nil {
					return
				}
				lines2 = delete_empty(lines2)
				if err := writeLines(lines2, "aop.txt"); err != nil {
					return
				}
			}
			//!-aop
			if strings.HasPrefix(result, "!-aop ") {
				lines, err := readLines("aop.txt")
				if err != nil {
					return
				}
				lines = delete_empty(lines)
				if err := writeLines(lines, "aop.txt"); err != nil {
					return
				}
				isOp := lines
				result := strings.TrimPrefix(result, "!-aop ")
				if slices.Contains(isOp, result) {

					lines, err := readLines("aop.txt")
					if err != nil {
						return
					}
					for i, v := range lines {
						if v == result {
							lines = append(lines[:i], lines[i+1:]...)
						}
					}
					lines = delete_empty(lines)
					if err := writeLines(lines, "aop.txt"); err != nil {
						return
					}
				} else {
					ircobj.Notice(e.Nick, "Sorry "+result+" not exists on my list")
				}
			}
			//!aops
			if strings.HasPrefix(result, "!aops") {
				lines, err := readLines("aop.txt")
				if err != nil {
					return
				}
				for _, s := range lines {
					ircobj.Notice(e.Nick, "Autoop: "+s)
				}
			}

			//!+shit
			if strings.HasPrefix(result, "!+shit ") {
				lines, err := readLines("shit.txt")
				if err != nil {
					return
				}

				isShit := lines

				result := strings.TrimPrefix(result, "!+shit ")
				if slices.Contains(isShit, result) {
					ircobj.Notice(e.Nick, "Sorry "+result+" exists on my list - not adding")
				} else {
					f, err := os.OpenFile("shit.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
					if err != nil {
						ircobj.Notice(e.Nick, "Something wrong with shit.txt")
					}
					defer f.Close()
					if _, err := f.WriteString(result + "\n"); err != nil {
						ircobj.Notice(e.Nick, "Something wrong with shit.txt")
					}
				}
				lines2, err := readLines("shit.txt")
				if err != nil {
					return
				}
				lines2 = delete_empty(lines2)
				if err := writeLines(lines2, "shit.txt"); err != nil {
					return
				}
			}
			//!-shit
			if strings.HasPrefix(result, "!-shit ") {
				lines, err := readLines("shit.txt")
				if err != nil {
					return
				}
				isShit := lines
				result := strings.TrimPrefix(result, "!-shit ")
				if slices.Contains(isShit, result) {

					lines, err := readLines("shit.txt")
					if err != nil {
						return
					}
					for i, v := range lines {
						if v == result {
							lines = append(lines[:i], lines[i+1:]...)
						}
					}
					lines = delete_empty(lines)
					if err := writeLines(lines, "shit.txt"); err != nil {
						return
					}
				} else {
					ircobj.Notice(e.Nick, "Sorry "+result+" not exists on my list")
				}
			}
			//!shits
			if strings.HasPrefix(result, "!shits") {
				lines, err := readLines("shit.txt")
				if err != nil {
					return
				}
				for _, s := range lines {
					ircobj.Notice(e.Nick, "Shit: "+s)
				}
			}
			// !help command
			if strings.HasPrefix(result, "!help") && strings.Contains(e.Raw, mynick) {
				ircobj.Privmsg(e.Nick, "My Commands: ")
				time.Sleep(1 * time.Second)
				ircobj.Privmsg(e.Nick, "!+o - gives op up to 3 nicks - !+o nick1 nick2 nick3")
				ircobj.Privmsg(e.Nick, "!6 - gives op up to 6 nicks - !6 nick1 nick2 nick3 nick4 nick5 nick6")
				time.Sleep(1 * time.Second)
				ircobj.Privmsg(e.Nick, "!-o - takes op up to 3 nicks - !-o nick1 nick2 nick3")
				time.Sleep(1 * time.Second)
				ircobj.Privmsg(e.Nick, "!+v - gives voice up to 3 nicks - !+v nick1 nick2 nick3")
				ircobj.Privmsg(e.Nick, "!-v - takes voice up to 3 nicks - !-v nick1 nick2 nick3")
				ircobj.Privmsg(e.Nick, "!+b - ban usermask - !+b *!ident@host")
				time.Sleep(1 * time.Second)
				ircobj.Privmsg(e.Nick, "!-b - unban usermask - !-b *!ident@host")
				ircobj.Privmsg(e.Nick, "!k nick - kick nick")
				ircobj.Privmsg(e.Nick, "!j #channel - joining #channel")
				time.Sleep(1 * time.Second)
				ircobj.Privmsg(e.Nick, "!p #channel - part channel ")
				ircobj.Privmsg(e.Nick, "!die quit_text - killing bot - !die Im going... ")
				ircobj.Privmsg(e.Nick, "!msg nick text - send priv msg to nick")
				ircobj.Privmsg(e.Nick, "!s text - say text to current channel")
				time.Sleep(1 * time.Second)
				ircobj.Privmsg(e.Nick, "!+owner add owner to bot, *!ident@host")
				ircobj.Privmsg(e.Nick, "!-owner delte owner from bot, *!ident@host")
				ircobj.Privmsg(e.Nick, "!hi - bot says hallo to You")
				ircobj.Privmsg(e.Nick, "!owners - notice owners list")
				time.Sleep(1 * time.Second)
				ircobj.Privmsg(e.Nick, "!+aop -add autoop: !+aop #channel *!ident@host")
				ircobj.Privmsg(e.Nick, "!-aop -del autoop: !-aop #channel *!ident@host")
				ircobj.Privmsg(e.Nick, "!aops - aops list")
				ircobj.Privmsg(e.Nick, "!+shit -add shit: !+shit #channel *!ident@host")
				ircobj.Privmsg(e.Nick, "!-shit -del shit: !-shit #channel *!ident@host")
				ircobj.Privmsg(e.Nick, "!shits - shit list")

			}
		}
	})
	ircobj.Loop()

}

func ReadConfig(filename string) (Config, error) {
	// init with some bogus data
	config := Config{
		"server":     "",
		"port":       "",
		"secretChan": "",
		"ident":      "",
		"botnick":    "",
		"realname":   "",
		"myhost":     "",
		"permowner":  "",
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

func isError(err error) bool {
	if err != nil {
		fmt.Println(err.Error())
	}

	return (err != nil)
}
func createAS(filepath string) {
	// check if file exists
	var _, err = os.Stat(filepath)

	// create file if not exists
	if os.IsNotExist(err) {
		var file, err = os.Create(filepath)

		if isError(err) {
			return
		}
		defer file.Close()

		var writefile, err2 = os.OpenFile(filepath, os.O_RDWR, 0644)
		if isError(err2) {
			return
		}
		defer writefile.Close()

		// Write some text line-by-line to file.
		_, err2 = writefile.WriteString("#channel *!ident@host \n")
		if isError(err2) {
			return
		}
		err = writefile.Sync()
		if isError(err) {
			return
		}
		fmt.Println("File Created Successfully", filepath)
	}
}

func writeFile(name string, text string) {
	// Open file using READ & WRITE permission.
	var file, err = os.OpenFile(name, os.O_RDWR, 0644)
	if isError(err) {
		return
	}
	defer file.Close()

	// Write some text line-by-line to file.
	_, err = file.WriteString(text)
	if isError(err) {
		return
	}

	// Save file changes.
	err = file.Sync()
	if isError(err) {
		return
	}
}

func createOwner(name string) {
	// check if file exists
	var _, err = os.Stat(name)

	// create file if not exists
	if os.IsNotExist(err) {
		var file, err = os.Create(name)
		if isError(err) {
			return
		}
		defer file.Close()
		var firstowner string = ""
		fmt.Println("Enter owner host (*!ident@host)")
		fmt.Scan(&firstowner)
		writeFile(name, firstowner)

		fmt.Println("File Created Successfully", name)
	}

}
func getIP() {
	foo, err := net.InterfaceAddrs()

	if err == nil {
		for _, v := range foo {
			fmt.Println(v)

		}
	}
}
func createConfig(name string) {
	// check if file exists
	var _, err = os.Stat(name)

	// create file if not exists
	if os.IsNotExist(err) {
		var file, err = os.Create(name)
		if isError(err) {
			return
		}
		defer file.Close()
		var cServer string = "server = "
		var cPort string = "port = "
		var cSecret string = "secretChan = "
		var cIdent string = "ident = "
		var cBotNick string = "botnick = "
		var cRealName string = "realname = "
		var cMyhost string = "myhost = "
		var iServer string = ""
		var iPort string = ""
		var iSecret string = ""
		var iIdent string = ""
		var iBotNick string = ""
		var iMyhost string = ""

		fmt.Println("I see that config.txt is not there, let's create one.")

		fmt.Println("Enter irc server ip or hostname")
		fmt.Scan(&iServer)

		fmt.Println("Enter irc server port")
		fmt.Scan(&iPort)

		fmt.Println("Enter secret channel name eg. #channel")
		fmt.Scan(&iSecret)

		fmt.Println("Enter ident (shell login)")
		fmt.Scan(&iIdent)

		fmt.Println("Enter botnick")
		fmt.Scan(&iBotNick)

		fmt.Println("Enter bot realname")
		inputReader := bufio.NewReader(os.Stdin)
		iRealName, _ := inputReader.ReadString('\n')
		time.Sleep(2 * time.Second)
		getIP()
		time.Sleep(2 * time.Second)
		fmt.Println("Enter ip to use with connection (vhost ip)")
		fmt.Scan(&iMyhost)

		var conffile, errc = os.OpenFile(name, os.O_RDWR, 0644)
		if isError(errc) {
			return
		}
		defer conffile.Close()

		// Write some text line-by-line to file.
		_, errc = conffile.WriteString(cServer + iServer + "\n")
		if isError(errc) {
			return
		}
		_, errc = conffile.WriteString(cPort + iPort + "\n")
		if isError(errc) {
			return
		}
		_, errc = conffile.WriteString(cSecret + iSecret + "\n")
		if isError(errc) {
			return
		}
		_, errc = conffile.WriteString(cIdent + iIdent + "\n")
		if isError(errc) {
			return
		}
		_, errc = conffile.WriteString(cBotNick + iBotNick + "\n")
		if isError(errc) {
			return
		}
		_, errc = conffile.WriteString(cRealName + iRealName)
		if isError(errc) {
			return
		}
		_, errc = conffile.WriteString(cMyhost + iMyhost + "\n")
		if isError(errc) {
			return
		}

		// Save file changes.
		errc = conffile.Sync()
		if isError(errc) {
			return
		}

		fmt.Println("File Created Successfully", name)

	}

}

// Split function for future use.
//func Split(r rune) bool {
//	return r == '!' || r == '@'
//
