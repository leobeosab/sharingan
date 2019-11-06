package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli"
)

func main() {
	sharingan := cli.NewApp()

	sharingan.Name = "Sharingan"
	sharingan.Usage = "Wrapper and analyzer for offensive security recon tools"

	sharingan.Action = func(c *cli.Context) error {
		fmt.Printf("Hello %q", c.Args().Get(0))
		return nil
	}

	err := sharingan.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}
