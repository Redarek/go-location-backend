package socket

import (
	socketio "github.com/googollee/go-socket.io"
	"location-backend/internal/db"
)

// TODO: implement websocket server routes

type Server struct {
	App *socketio.Server
	db  db.Service
}

func New(db db.Service) *Server {
	server := &Server{
		App: socketio.NewServer(nil),
		db:  db,
	}

	return server
	//server := socketio.NewServer(nil)

	//server.App.OnConnect("/", func(s socketio.Conn) error {
	//	s.SetContext("")
	//	fmt.Println("connected:", s.ID())
	//	return nil
	//})
	//
	//server.App.OnEvent("/", "notice", func(s socketio.Conn, msg string) {
	//	fmt.Println("notice:", msg)
	//	s.Emit("reply", "have "+msg)
	//})
	//
	//server.App.OnEvent("/chat", "msg", func(s socketio.Conn, msg string) string {
	//	s.SetContext(msg)
	//	return "recv " + msg
	//})
	//
	//server.App.OnEvent("/", "bye", func(s socketio.Conn) string {
	//	last := s.Context().(string)
	//	s.Emit("bye", last)
	//	s.Close()
	//	return last
	//})
	//
	//server.App.OnError("/", func(s socketio.Conn, e error) {
	//	// server.Remove(s.ID())
	//	fmt.Println("meet error:", e)
	//})
	//
	//server.App.OnDisconnect("/", func(s socketio.Conn, reason string) {
	//	// Add the Remove session id. Fixed the connection & mem leak
	//	//server.Remove(s.ID())
	//	fmt.Println("closed", reason)
	//})
	//
	//go server.App.Serve()
	//defer server.App.Close()
	//
	//http.Handle("/socket.io/", server.App)
	//http.Handle("/", http.FileServer(http.Dir("./asset")))
	//log.Println("Serving at localhost:8000...")
	//log.Fatal(http.ListenAndServe(":8000", nil))
}
