package main

import (
	"../client"
	"fmt"
	"log"
	"os"
	"text/tabwriter"
)

func main() {
	if len(os.Args) != 3 {
		usage()
		os.Exit(-1)
	}
	if os.Args[1] == "" {
		usage()
		os.Exit(-1)
	}

	c, err := client.New("https://" + os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	if !c.Login(os.Args[2], os.Getenv("UCMPASSWORD")) {
		log.Fatal("Can't login please verify username and password")
	}

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 0, '\t', 0)
	fmt.Fprint(w, "Profile\tTrunk\tUser\tType\tStatus\n")
	for _, trunk := range c.ListTrunks() {
		fmt.Fprintf(w, "Trunk\t%s\t%s\t%s\t%s\n",
			trunk.TrunkName,
			trunk.Username,
			trunk.Type,
			trunk.Status)
	}
	for _, account := range c.ListAccounts() {
		fmt.Fprintf(w, "Account\t%s\t%s\t%s\t%s\n",
			account.Extension,
			"-",
			account.Type,
			account.Status)
	}
	w.Flush()
}

func usage() {
	fmt.Println(`ucmstatus: <host:port> <user>
	env UCMPASSWORD if need password.
	example usage: UCMPASSWORD=mipassword ./ucmstatus localhost:8089 admin`)
}
