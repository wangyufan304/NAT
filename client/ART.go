package main

import "fmt"

func art() {
	fmt.Println("  _   _              _______    _____   _        _____   ______   _   _   _______ \n | \\ | |     /\\     |__   __|  / ____| | |      |_   _| |  ____| | \\ | | |__   __|\n |  \\| |    /  \\       | |    | |      | |        | |   | |__    |  \\| |    | |   \n | . ` |   / /\\ \\      | |    | |      | |        | |   |  __|   | . ` |    | |   \n | |\\  |  / ____ \\     | |    | |____  | |____   _| |_  | |____  | |\\  |    | |   \n |_| \\_| /_/    \\_\\    |_|     \\_____| |______| |_____| |______| |_| \\_|    |_|   \n                                                                                  \n                                                                                  ")
}

func printRelationInformation() {
	fmt.Println("[Client Running Successfully!]")
	fmt.Println("[PublicAddress]", objectConfig.PublicServerAddr)
	fmt.Println("[TunnelAddress]", objectConfig.TunnelServerAddr)
	fmt.Println("[LocalAddress]", objectConfig.LocalServerAddr)
}
