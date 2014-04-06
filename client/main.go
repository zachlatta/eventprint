package main

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"image"
	"io"
	"log"
	"net"
	"net/url"
	"os"
	"os/exec"
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

func main() {
	url, err := url.Parse("http://localhost:3000/api/ws")
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

	ws.SetPongHandler(func(string) error { ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			break
		}

		var attendee model.Attendee
		if err := json.Unmarshal(msg, &attendee); err != nil {
			log.Fatal(err)
		}

		path, err := GeneratePDF(attendee)
		if err != nil {
			log.Fatal(err)
		}

		if err := exec.Command("pdftops", path, path+".ps").Run(); err != nil {
			log.Fatal(err)
		}

		if err := exec.Command("lp", path+".ps").Run(); err != nil {
			log.Fatal(err)
		}

		// Clean up
		if err := os.Remove(path); err != nil {
			log.Fatal(err)
		}

		if err := os.Remove(path + ".ps"); err != nil {
			log.Fatal(err)
		}

		log.Printf("Printed badge for %s %s\n", attendee.FirstName, attendee.LastName)
	}
}

func GeneratePDF(attendee model.Attendee) (string, error) {
	var b bytes.Buffer

	doc := pdf.NewPdfDoc(&b)

	pWidth, pHeight := dox2go.StandardSize(dox2go.PS_A7, dox2go.U_MM)
	page := doc.CreatePage(dox2go.U_MM, pWidth, pHeight, dox2go.PO_Portrait)

	s := page.Surface()

	qrImage, err := CreateQRCode(attendee.Barcode)
	if err != nil {
		return "", err
	}

	name := doc.CreateFont(pdf.FONT_Helvetica, dox2go.FS_Bold, 12)
	typeOfAdmissionPath := dox2go.NewPath()
	typeOfAdmission := doc.CreateFont(pdf.FONT_Helvetica, dox2go.FS_Bold, 10)
	qrcode := doc.CreateImage(qrImage)

	s.Bg(dox2go.RGB(0, 0, 0))
	s.Text(name, 17, 75, attendee.FirstName)
	s.Text(name, 0, -12, attendee.LastName)

	s.Image(qrcode, 23, 30, 75, 75)

	s.Fg(dox2go.RGB(0, 0, 0))
	typeOfAdmissionPath.Rect(5, 5, 69, 20)
	s.Fill(typeOfAdmissionPath)

	s.Bg(dox2go.RGB(255, 255, 255))
	s.Text(typeOfAdmission, 11, 9, attendee.Title)

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
	//h := md5.New()
	//io.WriteString(h, text)

	//qrImage := fmt.Sprintf("/tmp/%x.png", h.Sum(nil))

	//f, err := os.Create(qrImage)
	//if err != nil {
	//return "", err
	//}
	//defer f.Close()

	qrcode, err := qr.Encode(text, qr.L, qr.Auto)
	if err != nil {
		return nil, err
	}

	return qrcode, nil
}
