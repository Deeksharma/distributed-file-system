package main

import (
	"distributed-file-system/p2p"
	"log"
)

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

func main() {
	transport := p2p.NewTCPTransport(":3000")
	//log.Fatal(transport.ListenAndAccept())
	if err := transport.ListenAndAccept(); err != nil {
		log.Fatal(err)
	}
	select {}
}
