# GoGolin
Simple irc bot reacting to !commands and written in Go

# Info
owner.txt - here bot stores owners hosts
shit.txt - here bot stores shit list
aop.txt - here bot stores aop list
config.txt - here bot stores config 

# Config

# GoGolin v1.0 config file

 server = irc.server.com
 
 port = 6667
 
 secretChan = #yourchannel
 
 ident = botident
 
 botnick = botnick
 
 realname = botrealname 
 
 myhost = ip_witch_want_to_use

# Commands

!help

My Commands:

!+o - gives op up to 3 nicks - !+o nick1 nick2 nick3

!6 - gives op up to 6 nicks - !6 nick1 nick2 nick3 nick4 nick5 nick6

!-o - takes op up to 3 nicks - !-o nick1 nick2 nick3

!+v - gives voice up to 3 nicks - !+v nick1 nick2 nick3

!-v - takes voice up to 3 nicks - !-v nick1 nick2 nick3

!+b - ban usermask - !+b *!ident@host

!-b - unban usermask - !-b *!ident@host

!k nick - kick nick

!j #channel - joining #channel

!p #channel - part channel

!die quit_text - killing bot - !die Im going...

!msg nick text - send priv msg to nick

!s text - say text to current channel

!+owner add owner to bot, *!ident@host

!-owner delte owner from bot, *!ident@host

!hi - bot says hallo to You

!owners - notice owners list

!+aop -add autoop: !+aop #channel *!ident@host

!-aop -del autoop: !-aop #channel *!ident@host

!aops - aops list

!+shit -add shit: !+shit #channel *!ident@host

!-shit -del shit: !-shit #channel *!ident@host

!shits - shit list

!help - command list (works only as private message to bot)
