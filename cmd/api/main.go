package main

import (
	"location-backend/internal/app"
)

//func main() {
//
//	server := server.New()
//
//	server.RegisterFiberRoutes()
//	port, _ := strconv.Atoi(os.Getenv("PORT"))
//	err := server.Listen(fmt.Sprintf(":%d", port))
//	if err != nil {
//		panic(fmt.Sprintf("cannot start server: %s", err))
//	}
//}

func main() {
	s := app.New()
	s.Run()
}
