package main

import (
	// "github.com/automatico/jato/cli"
	// "github.com/automatico/jato/ssh"
	"github.com/automatico/jato/telnet"
)

func main() {

	// cliParams := cli.CLI()
	// if cliParams.NoOp != true {
	// 	ssh.SSH(cliParams)
	// }
	telnet.Telnet()

}
