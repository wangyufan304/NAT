package utils

import (
	"fmt"
	"testing"
)

func TestParseFile(t *testing.T) {
	pf := ParseFile("server.yml")
	name := pf.ViperInstance.GetString("Server.Name")
	addr := pf.ViperInstance.GetString("Server.Addr")
	port := pf.ViperInstance.GetString("Server.Port")
	fmt.Println("[Name]", name, "[Addr]", addr, "[Port]", port)
}
