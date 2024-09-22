package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"time"
)

// Function to simulate listening to JMS and receiving a message
func startJmsListener() {
	for {
		// Simulate pulling the message from the JMS server
		message := `<ns5:MessageCollection>...</ns5:MessageCollection>`

		// Write the message to the message.xml file
		err := ioutil.WriteFile("../xml-scripts/messages/message.xml", []byte(message), 0644)
		if err != nil {
			log.Fatalf("Failed to write file: %s", err)
		}
		fmt.Println("XML message saved to ../xml-scripts/messages/message.xml")

		// Process the message using Ruby script
		processXMLWithRuby()

		time.Sleep(10 * time.Second) // Simulate pulling messages every 10 seconds
	}
}

// Function to call the Ruby script to process the XML
func processXMLWithRuby() {
	cmd := exec.Command("ruby", "../xml-scripts/xml-details.rb", "../xml-scripts/messages/message.xml")
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Error executing Ruby script: %s\nOutput: %s", err, output)
	}
	fmt.Printf("Ruby script output:\n%s\n", output)
}

// Function to simulate the map visualization (you can add actual visualization logic here)
func visualizeData() {
	fmt.Println("Loading maps and data blocks...")

	// Add your logic here for drawing maps or handling visualization in Go
}

func main() {
	fmt.Println("Starting the visualization process...")

	// Start JMS listener in a separate goroutine
	go startJmsListener()

	// Start the visualization process
	visualizeData()

	// Keep the program running
	select {}
}
