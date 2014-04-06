package main

import (
	"bytes"
	"fmt"
	"image"
	"os"

	"github.com/adkennan/dox2go"
	"github.com/boombuler/barcode/qr"
	"github.com/zachlatta/dox2go/pdf"
)

func main() {
	var b bytes.Buffer

	doc := pdf.NewPdfDoc(&b)

	pWidth, pHeight := dox2go.StandardSize(dox2go.PS_A7, dox2go.U_MM)
	page := doc.CreatePage(dox2go.U_MM, pWidth, pHeight, dox2go.PO_Portrait)

	s := page.Surface()

	qrImage, err := CreateQRCode("Hello, World!")
	if err != nil {
		panic(err)
	}

	name := doc.CreateFont(pdf.FONT_Helvetica, dox2go.FS_Bold, 12)
	typeOfAdmissionPath := dox2go.NewPath()
	typeOfAdmission := doc.CreateFont(pdf.FONT_Helvetica, dox2go.FS_Bold, 10)
	qrcode := doc.CreateImage(qrImage)

	s.Bg(dox2go.RGB(0, 0, 0))
	s.Text(name, 17, 75, "FIRST")
	s.Text(name, 0, -12, "LAST")

	s.Image(qrcode, 23, 30, 75, 75)

	s.Fg(dox2go.RGB(0, 0, 0))
	typeOfAdmissionPath.Rect(5, 5, 69, 20)
	s.Fill(typeOfAdmissionPath)

	s.Bg(dox2go.RGB(255, 255, 255))
	s.Text(typeOfAdmission, 11, 9, "GENERAL")

	doc.Close()

	f, err := os.Create("/tmp/tmp.pdf")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	n, err := b.WriteTo(f)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Wrote %d bytes\n", n)
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
