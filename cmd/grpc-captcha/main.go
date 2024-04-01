package main

import (
	"log"

	"gitlab.com/gtsh77-workshop/grpc-captcha/internal/app"
)

// type CaptPNG struct {
// 	*captcha.Image
// }

// func main() {
// 	var (
// 	// f   *os.File
// 	// w   *captcha.Image
// 	// err error
// 	)

// 	png := &CaptPNG{
// 		Image: captcha.NewImage("123", captcha.RandomDigits(5), 180, 80),
// 	}
// 	// if f, err = os.Create("test.png"); err != nil {
// 	// 	log.Fatal(err)
// 	// }

// 	// if _, err = w.WriteTo(f); err != nil {
// 	// 	log.Fatal(err)
// 	// }

// 	bb := new(bytes.Buffer)
// 	png.Image.WriteTo(bb)

// 	fmt.Println(hex.EncodeToString(bb.Bytes()))
// }

var name, version, compiledAt string //nolint:gochecknoglobals //lld flags

func main() {
	if _, err := app.New(name, version, compiledAt).Start(); err != nil {
		log.Fatal(err)
	}
}
