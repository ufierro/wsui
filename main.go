package main

import (
	"net/url"

	"fmt"
	"log"
	"net/http"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"github.com/Lyca0n/wsui/model"
	"github.com/Lyca0n/wsui/widgets"

	"github.com/gorilla/websocket"
)

var connList []model.Bookmark = []model.Bookmark{
	{Name: "Home", Url: url.URL{Scheme: "ws", Host: "localhost", Path: ""}},
	{Name: "Store 0020", Url: url.URL{Scheme: "ws", Host: "192.168.0.120", Path: ""}},
	{Name: "Home", Url: url.URL{Scheme: "ws", Host: "localhost", Path: ""}},
	{Name: "Store 0020", Url: url.URL{Scheme: "ws", Host: "192.168.0.120", Path: ""}},
	{Name: "Home", Url: url.URL{Scheme: "ws", Host: "localhost", Path: ""}},
	{Name: "Store 0020", Url: url.URL{Scheme: "ws", Host: "192.168.0.120", Path: ""}},
	{Name: "Store 0020", Url: url.URL{Scheme: "ws", Host: "192.168.0.120", Path: ""}},
}

func main() {
	//To be removed
	go LanuchServer()
	connectionWidget := widgets.ConnectionWidget{}
	messages := widgets.MessageWidget{}
	a := app.New()
	w := a.NewWindow("WSUI")
	w.Resize(fyne.NewSize(960, 660))
	w.SetContent(
		container.NewGridWithColumns(3, connectionWidget.MakeConnectionWidget(connList), messages.MakeMessageWidget()),
	)
	w.ShowAndRun()

	messages.NewMessage("hello bro")
}

//server stuff

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
} // use default options

func socketHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade our raw HTTP connection to a websocket based one
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("Error during connection upgradation:", err)
		return
	}
	defer conn.Close()

	// The event loop
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error during message reading:", err)
			break
		}
		log.Printf("Received: %s", message)
		err = conn.WriteMessage(messageType, message)
		if err != nil {
			log.Println("Error during message writing:", err)
			break
		}
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Index Page")
}

func LanuchServer() {
	http.HandleFunc("/socket", socketHandler)
	http.HandleFunc("/", home)
	log.Fatal(http.ListenAndServe("localhost:9090", nil))
}