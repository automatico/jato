package main

import (
	"github.com/automatico/jato/cli"
	"github.com/automatico/jato/ssh"
)

func main() {

	cliParams := cli.CLI()
	if cliParams.NoOp != true {
		ssh.SSH(cliParams)
	}

}
