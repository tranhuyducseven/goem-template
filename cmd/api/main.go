package main

import (
	"fmt"

	"github.com/tranhuyducseven/goem-template/cmd/api/server"
	"github.com/tranhuyducseven/goem-template/config"
)

func main() {
	s := server.New()
	host := "0.0.0.0"

	err := s.Start(fmt.Sprintf("%s:%s", host, config.Env.PORT))
	defer s.Close()
	if err != nil {
		panic(err)
	}
}
