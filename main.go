package main

import (
	"log"
	"spy-cat-agency/cmd"
	_ "spy-cat-agency/docs"
)

//	@title			Spy Cat Agency API
//	@version		1.0
//	@description	A REST API for managing spy cats and their missions
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		localhost:8080
//	@BasePath	/api/v1

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}