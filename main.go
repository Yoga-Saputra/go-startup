package main

import "startup/routes"

func main() {
	routes.InitApi()
	// overseer.Run(overseer.Config{
	// 	Program: gracefulStart,
	// 	Address: ":4000",
	// 	Fetcher: &fetcher.HTTP{
	// 		URL:      "http://localhost",
	// 		Interval: 1 * time.Second,
	// 	},
	// 	Debug: true,
	// })
}

// func gracefulStart(state overseer.State) {
// 	routes.InitApi(state)
// }
