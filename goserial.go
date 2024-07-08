package main

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
	"go.bug.st/serial"
)

func main() {
	println("starting main")

	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("yaml")   // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(".")      // optionally look for config in the working directory
	err := viper.ReadInConfig()   // Find and read the config file
	if err != nil {               // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	err = checkAllConfig()
	if err != nil {
		panic(err.Error())
	}

	mode := &serial.Mode{
		BaudRate: viper.GetInt("baud"),
		Parity:   serial.EvenParity,
		DataBits: 7,
		StopBits: serial.OneStopBit,
	}

	listenerPort := viper.GetString("listener")

	forwarderPort := viper.GetString("forwarder")

	tportPort := viper.GetString("tport")

	listener, err := serial.Open(listenerPort, mode)
	listener.SetReadTimeout(time.Millisecond * 10)

	if err != nil {
		println("error opening port: %v", err)
	}

	fowarder, err := serial.Open(forwarderPort, mode)
	fowarder.SetReadTimeout(time.Millisecond * 10)
	if err != nil {
		fmt.Printf("error opening port: %v \n", err)
	}

	tport, err := serial.Open(tportPort, mode)

	if err != nil {
		fmt.Printf("error opening port: %v \n", err)
	}

	go func() {
		println("running listener loop")
		for {
			buff := make([]byte, viper.GetInt("buffer"))
			n, err := listener.Read(buff)
			if err != nil {
				fmt.Printf("error reading listener %v: %v \n", listenerPort, err)
			}
			if n != 0 {
				fowarder.Write(buff[:n])
				fmt.Printf("writing to forwarder (drip feed) %v %s, %s \n", n, forwarderPort, string(buff[:n]))
				tportCopy := buff[:n]
				go tport.Write(tportCopy)
			}

		}
	}()

	go func() {
		println("running forwarder loop")

		for {
			buff := make([]byte, viper.GetInt("buffer"))
			n, err := fowarder.Read(buff)
			if err != nil {
				fmt.Printf("error reading forwarder %v: %v \n", forwarderPort, err)
			}
			if n != 0 {
				listener.Write(buff[:n])
				fmt.Printf("writing to listener (cnc) %v  %s , %s \n", n, listenerPort, string(buff[:n]))
				tportCopy := buff[:n]
				go tport.Write(tportCopy)
			}

		}
	}()

	print("running forwarder with T port sniffer \n")

	for {

	}
}

func checkAllConfig() error {

	keys := []string{
		"baud",
		"listener",
		"forwarder",
		"tport",
		"buffer",
	}

	for _, key := range keys {

		isSet := viper.IsSet(key)
		if !isSet {
			fmt.Printf("unable to locate key %v", key)
			return fmt.Errorf("unable to locate key %v", key)
		}
	}
	return nil

}
