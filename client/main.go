package main

import (
	"crypto/md5"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/url"
	"os"
	"os/exec"
	"time"

	"github.com/gorilla/websocket"
	"github.com/zachlatta/eventprint/server/model"
)

const (
	pongWait = 60 * time.Second
)

var (
	baseUrl  string
	logoPath string
)

func main() {
	flag.StringVar(&baseUrl, "url", "http://localhost:3000", "base url of server")
	flag.StringVar(&logoPath, "logo", "lahacks.png", "path to hackathon logo")
	flag.Parse()

	url, err := url.Parse(baseUrl + "/api/ws")
	if err != nil {
		panic(err)
	}

	c, err := net.Dial("tcp", url.Host)
	if err != nil {
		log.Fatal(err)
	}

	ws, _, err := websocket.NewClient(c, url, nil, 1024, 1024)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	log.Printf("Connected to Eventprint server at %s", baseUrl)

	ws.SetPongHandler(func(string) error { ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			log.Println(err)
		}

		log.Println("Message received!")

		var attendee model.Attendee
		if err := json.Unmarshal(msg, &attendee); err != nil {
			log.Fatal("Error unmarshalling message JSON:", err)
		}

		path, err := GeneratePDF(attendee)
		if err != nil {
			log.Fatal("Error generating badge:", err)
		}

		if err := exec.Command("lp", "-o", "fit-to-page", path).Run(); err != nil {
			log.Fatal("Error printing badge:", err)
		}

		// Clean up.
		if err := os.Remove(path); err != nil {
			log.Fatal("Error deleting temp badge:", err)
		}

		log.Printf("Printed badge for %s %s\n", attendee.FirstName, attendee.LastName)
	}
}

func GeneratePDF(a model.Attendee) (string, error) {
	h := md5.New()
	io.WriteString(h, a.Barcode)

	path := fmt.Sprintf("/tmp/%x.pdf", h.Sum(nil))

	err := exec.Command("./generate_pdf.rb", a.Barcode, a.FirstName, a.LastName,
		a.Title, path).Run()
	if err != nil {
		return "", err
	}

	return path, nil
}
