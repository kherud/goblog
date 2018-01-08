package main

import (
	"flag"
	"fmt"
	"os"
	"bufio"
	"github.com/kherud/goblog/config"
	"github.com/kherud/goblog/backend"
	"github.com/kherud/goblog/webserver"
)

/**
Starting point that parses possible flags, ensures an user exists and creates one if not. Then starts the web server.
Also ensures the storage directory and certificate files exist.
 */
func main() {
	time := flag.Int("t", 15, "Minutes until an authentication session expires")
	port := flag.String("p", "8080", "Port that is used for the webserver")
	flag.Parse()
	config.SESSION_TIME = *time
	config.DEFAULT_PORT = *port
	_, certErr := os.Stat(config.CERT_FILE)
	_, keyErr := os.Stat(config.CERT_FILE)
	if certErr != nil || keyErr != nil {
		fmt.Println("HTTPS certificate or key file could not be found.\nPlease ensure they are at the right directory.")
	} else {
		os.MkdirAll(config.DATA_PATH, os.ModePerm)
		fmt.Println("Starting webserver on port", config.DEFAULT_PORT)
		fmt.Println("Session expiration time:", config.SESSION_TIME, "minutes")
		backend.EnsureUserExists(bufio.NewReader(os.Stdin)) // inject dependency for proper testing
		fmt.Println("Server is now running...")
		webserver.StartServer()
	}
}
