package main

import (
	"fmt"
	"os"

	"github.com/jessie846/eram/custom_map"
	"github.com/jessie846/eram/file_list"
	"github.com/jessie846/eram/flight"
	"github.com/jessie846/eram/message_receiver"
	"github.com/jessie846/eram/window"
)

func main() {
	mapNames := []string{
		"Boundary2.geojson",
		"High Event Split 2.geojson",
		"ATCTLow Boundary.geojson.geojson",
	}

	var maps []custom_map.Map // or []*map.Map if map.FromFile returns pointers
	for _, mapName := range mapNames {
		mapPath := fmt.Sprintf("../ZJX Maps/%s", mapName)
		mapData, err := customMap.FromFile(mapPath)
		if err != nil {
			fmt.Printf("Error loading map: %v\n", err)
			return
		}
		// Append mapData to maps based on its return type (use mapData or *mapData)
		maps = append(maps, mapData) // Or maps = append(maps, *mapData) if necessary
	}

	args := os.Args
	if len(args) < 3 {
		fmt.Printf("Usage: %s <facility> <sector>\n", args[0])
		return
	}

	facility := args[1]
	sector := args[2]
	currentPosition := flight.Owner{Facility: facility, Sector: sector}

	if len(args) == 4 && args[3] == "--files" {
		fileList := file_list.FromGlob("../xml-scripts/messages/*.xml")
		messageReceiver := message_receiver.NewFileListMessageReceiver(fileList)
		if err := window.Show(&currentPosition, &maps, messageReceiver); err != nil {
			fmt.Printf("Error showing window: %v\n", err)
		}
	} else {
		messageReceiver, err := message_receiver.NewRabbitMQMessageReceiver()
		if err != nil {
			fmt.Printf("Error creating RabbitMQ message receiver: %v\n", err)
			return
		}
		if err := window.Show(&currentPosition, &maps, messageReceiver); err != nil {
			fmt.Printf("Error showing window: %v\n", err)
		}
	}
}
