package main

import (
	"fmt"
	"time"

	"go.bug.st/serial"
)

func main() {
	println("starting main")
	mode := &serial.Mode{
		BaudRate: 9600,
		Parity:   serial.EvenParity,
		DataBits: 7,
		StopBits: serial.OneStopBit,
	}
	listenerPort := "COM5"

	forwarderPort := "COM8"

	listener, err := serial.Open(listenerPort, mode)

	if err != nil {
		println("error opening port: %v", err)
	}

	listener.SetReadTimeout(time.Millisecond * 10)

	fowarder, err := serial.Open(forwarderPort, mode)
	fowarder.SetReadTimeout(time.Millisecond * 10)
	if err != nil {
		fmt.Printf("error opening port: %v \n", err)
	}
	go func() {
		println("running listener loop")
		for {
			buff := make([]byte, 100)
			n, err := listener.Read(buff)
			if err != nil {
				fmt.Printf("error reading listener %v: %v \n", listenerPort, err)
			}
			if n != 0 {
				fowarder.Write(buff[:n])
				fmt.Printf("writing %v %s, %s \n", n, forwarderPort, string(buff[:n]))
			}

		}
	}()

	go func() {
		println("running forw arder loop")

		for {
			buff := make([]byte, 100)
			n, err := fowarder.Read(buff)
			if err != nil {
				fmt.Printf("error reading forwarder %v: %v \n", forwarderPort, err)
			}
			if n != 0 {
				listener.Write(buff[:n])
				fmt.Printf("writing %v  %s , %s \n", n, listenerPort, string(buff[:n]))
			}

		}
	}()

	go func() {
		println("test func")
	}()

	for {
	}
}
