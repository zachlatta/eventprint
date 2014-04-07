package main

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"flag"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"math"
	"net"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/adkennan/dox2go"
	"github.com/boombuler/barcode/qr"
	"github.com/gorilla/websocket"
	"github.com/zachlatta/dox2go/pdf"
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
			break
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

		if err := exec.Command("lp", path).Run(); err != nil {
			log.Fatal("Error printing badge:", err)
		}

		// Clean up.
		if err := os.Remove(path); err != nil {
			log.Fatal("Error deleting temp badge:", err)
		}

		log.Printf("Printed badge for %s %s\n", attendee.FirstName, attendee.LastName)
	}
}

func GeneratePDF(attendee model.Attendee) (string, error) {
	var b bytes.Buffer

	doc := pdf.NewPdfDoc(&b)

	pWidth, pHeight := dox2go.StandardSize(dox2go.PS_A8, dox2go.U_MM)
	page := doc.CreatePage(dox2go.U_MM, pWidth, pHeight, dox2go.PO_Portrait)

	s := page.Surface()

	qrImage, err := CreateQRCode(attendee.Barcode)
	if err != nil {
		return "", err
	}

	logoImage, err := LoadImageFromFile(logoPath)
	if err != nil {
		return "", err
	}

	nameSize := (pWidth + 25) / math.Max(float64(len(attendee.FirstName)),
		float64(len(attendee.LastName)))
	name := doc.CreateFont(pdf.FONT_Helvetica, dox2go.FS_Bold, nameSize)
	typeOfAdmissionPath := dox2go.NewPath()
	typeOfAdmission := doc.CreateFont(pdf.FONT_Helvetica, dox2go.FS_Bold, 6)
	qrcode := doc.CreateImage(qrImage)
	logo := doc.CreateImage(*logoImage)

	s.Image(logo, 18, 55, 45, 45)

	s.Bg(dox2go.RGB(0, 0, 0))
	s.Text(name, 0, 37, attendee.LastName)
	s.Text(name, 0, nameSize, attendee.FirstName)

	s.Image(qrcode, 18, 15, 50, 50)

	s.Fg(dox2go.RGB(0, 0, 0))
	typeOfAdmissionPath.Rect(0, 0, 52, 10)
	s.Fill(typeOfAdmissionPath)

	s.Bg(dox2go.RGB(255, 255, 255))
	s.Text(typeOfAdmission, 3, 3, strings.ToUpper(attendee.Title))

	doc.Close()

	h := md5.New()
	io.WriteString(h, attendee.Barcode)

	path := fmt.Sprintf("/tmp/%x.pdf", h.Sum(nil))

	f, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	_, err = b.WriteTo(f)
	if err != nil {
		return "", err
	}

	return path, nil
}

func CreateQRCode(text string) (image.Image, error) {
	qrcode, err := qr.Encode(text, qr.L, qr.Auto)
	if err != nil {
		return nil, err
	}

	return qrcode, nil
}

func LoadImageFromFile(path string) (*image.Image, error) {
	reader, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	m, _, err := image.Decode(reader)
	if err != nil {
		return nil, err
	}

	return &m, nil
}
