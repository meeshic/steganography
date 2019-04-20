package main

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math"
	"os"
	"strings"
)

func init() {
	// image.RegisterFormat("jpeg", "jpeg", jpeg.Decode, jpeg.DecodeConfig)
	image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)
}

// embed message into picture by modifying pixels
/*func embedMessage(img *image.RGBA, message string) bool {
	fmt.Println(message)

	b := []byte(message)
	toWrite := len(b) * 8
	curr := 0
	setZeroes := false
	numZeroes := 8

	height := img.Bounds().Dy()
	width := img.Bounds().Dx()

	f := func(pixel *uint32, b byte, curr int, setZeroes bool) {
		if setZeroes {
			//set
			numZeroes--
		}
		bit := b & byte(math.Pow(2, float64(curr%8)))
		pixel = (pixel & 0xfe) | bit
		toWrite--
	}

	for x := 0; x < height; x++ {
		for y := 0; y < width; y++ {
			r, g, b, a := img.RGBAAt(x, y).RGBA()
			f(&r, setZeroes)
			img.Set(x, y, color.RGBA{
				R: uint8(r),
				G: uint8(g),
				B: uint8(b),
				A: uint8(a),
			})
		}
	}

	return true
}*/

func embedMessage(img *image.RGBA, message string) bool {
	byteMessage := []byte(message)
	fmt.Println(message)
	for _, b := range byteMessage {
		fmt.Println(b)
		fmt.Println(string(b))
	}
	fmt.Println("=====")

	// maxHeight := img.Bounds().Dy()
	maxWidth := img.Bounds().Dx()

	height := 0
	width := 0

	curr := 0

	var bit byte

	// 1 byte = 8 bits
	// 1 pixel = 4 bits (R,G,B,A)
	// 1 byte -> 2 pixels
	for i := 0; i < len(byteMessage); i++ {
		if width == maxWidth {
			height++
			width = 0
		}
		// fmt.Println("curr char: ", string(byteMessage[i]))
		r, g, b, a := img.RGBAAt(height, width).RGBA()
		fmt.Printf("before encoded pixel 1 - %v: %v %v %v %v\n", i+1, r, g, b, a)
		bit = byteMessage[i] & byte(math.Pow(2, float64(0)))
		// fmt.Println("bit 1: ", bit)
		r = (r & 0xfe) | uint32(bit)
		bit = byteMessage[i] & byte(math.Pow(2, float64(1)))
		// fmt.Println("bit 2: ", bit>>1)
		g = (g & 0xfe) | uint32(bit>>1)
		bit = byteMessage[i] & byte(math.Pow(2, float64(2)))
		// fmt.Println("bit 3: ", bit>>2)
		b = (b & 0xfe) | uint32(bit>>2)
		bit = byteMessage[i] & byte(math.Pow(2, float64(3)))
		// fmt.Println("bit 4: ", bit>>3)
		a = (a & 0xfe) | uint32(bit>>3)
		img.Set(height, width, color.RGBA{
			R: uint8(r),
			G: uint8(g),
			B: uint8(b),
			A: uint8(a),
		})

		width++
		if width == maxWidth {
			height++
			width = 0
		}
		fmt.Printf("encoded pixel 1 - %v: %v %v %v %v\n", i+1, r, g, b, a)
		r, g, b, a = img.RGBAAt(height, width).RGBA()
		fmt.Printf("before encoded pixel 2 - %v: %v %v %v %v\n", i+1, r, g, b, a)
		bit = byteMessage[i] & byte(math.Pow(2, float64(4)))
		// fmt.Println("bit 5: ", bit>>4)
		r = (r & 0xfe) | uint32(bit>>4)
		bit = byteMessage[i] & byte(math.Pow(2, float64(5)))
		// fmt.Println("bit 6: ", bit>>5)
		g = (g & 0xfe) | uint32(bit>>5)
		bit = byteMessage[i] & byte(math.Pow(2, float64(6)))
		// fmt.Println("bit 7: ", bit>>6)
		b = (b & 0xfe) | uint32(bit>>6)
		bit = byteMessage[i] & byte(math.Pow(2, float64(7)))
		// fmt.Println("bit 8: ", bit>>7)
		a = (a & 0xfe) | uint32(bit>>7)
		img.Set(height, width, color.RGBA{
			R: uint8(r),
			G: uint8(g),
			B: uint8(b),
			A: uint8(a),
		})
		fmt.Printf("encoded pixel 2 - %v: %v %v %v %v\n", i+1, r, g, b, a)
		width++
		curr++
	}

	if width == maxWidth {
		height++
		width = 0
	}

	// signal end of message
	img.Set(height, width, color.RGBA{
		R: uint8(0),
		G: uint8(0),
		B: uint8(0),
		A: uint8(0),
	})

	return true
}

// iterate over pixels to find embedded message
// func extractMessage(img image.RGBA) string {
// 	var message []byte

// 	// maxHeight := img.Bounds().Dy()
// 	maxWidth := img.Bounds().Dx()

// 	var bits []byte

// 	for h := 0; h < 2; h++ {
// 		for w := 0; w < maxWidth; w++ {
// 			r, g, b, a := img.RGBAAt(h, w).RGBA()
// 			fmt.Printf("extracted pixel: %v %v %v %v\n", r, g, b, a)
// 			if r == 0 && g == 0 && b == 0 && a == 0 {
// 				return string(message)
// 			}
// 			bits = append(bits, byte(r&1))
// 			bits = append(bits, byte(g&1))
// 			bits = append(bits, byte(b&1))
// 			bits = append(bits, byte(a&1))
// 			if len(bits) == 8 {
// 				var charByte byte
// 				// fmt.Println("initial char: ", charByte)
// 				for i, bit := range bits {
// 					move := uint(i)
// 					// fmt.Println("bit: ", bit<<move)
// 					charByte = charByte | (bit << move)
// 					// fmt.Println("charByte after append: ", charByte)
// 				}
// 				fmt.Println("char: ", string(charByte))
// 				message = append(message, charByte)
// 				bits = nil
// 			}
// 		}
// 	}
// 	return ""
// }

func extractMessage(img image.Image) string {
	var message []byte

	// maxHeight := img.Bounds().Dy()
	// maxWidth := img.Bounds().Dx()

	var bits []byte

	for h := 0; h < 4; h++ {
		for w := 0; w < 5; w++ {
			r, g, b, a := img.At(h, w).RGBA()
			fmt.Printf("extracted pixel: %v %v %v %v\n", r, g, b, a)
			// if r == 0 && g == 0 && b == 0 && a == 0 {
			// 	return string(message)
			// }
			bits = append(bits, byte(r&1))
			bits = append(bits, byte(g&1))
			bits = append(bits, byte(b&1))
			bits = append(bits, byte(a&1))
			if len(bits) == 8 {
				var charByte byte
				// fmt.Println("initial char: ", charByte)
				for i, bit := range bits {
					move := uint(i)
					// fmt.Println("bit: ", bit<<move)
					charByte = charByte | (bit << move)
					// fmt.Println("charByte after append: ", charByte)
				}
				fmt.Println("char: ", string(charByte))
				message = append(message, charByte)
				bits = nil
			}
		}
	}
	return string(message)
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("Enter: ")
		fmt.Println("'embed' to embed a message")
		fmt.Println("'extract' to extract a message")

		fmt.Print("> ")

		userInput, err := reader.ReadString('\n')
		userInput = strings.Replace(userInput, "\r\n", "", -1)
		userInput = strings.TrimSpace(userInput)
		for err != nil {
			fmt.Println("ERROR: ", err)
			fmt.Println("Please enter 'embed' or 'extract'")
		}
		switch userInput {
		case "embed":
			// read in picture & location
			// fmt.Print("Enter image (jpg) location: ")
			// fileLocation, err := reader.ReadString('\n')
			// for err != nil {
			// 	fmt.Printf("file location %s not valid!\n", fileLocation)
			// 	fmt.Print("Re-enter image (jpg) location: ")
			// 	fileLocation, err = reader.ReadString('\n')
			// }
			// fileLocation = strings.Replace(fileLocation, "\r\n", "", -1)
			// file, err := os.Open(fileLocation)
			file, err := os.OpenFile("C:/Users/miche/Desktop/Avatar_cat.png", os.O_RDWR, 0600)
			for err != nil {
				fmt.Printf("ERROR: Image not found - %s\n", err.Error())
				fmt.Print("Enter image (jpg) location: ")
				fileLocation, err := reader.ReadString('\n')
				for err != nil {
					fmt.Printf("file location %s not valid!\n", fileLocation)
					fmt.Print("Re-enter image (jpg) location: ")
					fileLocation, err = reader.ReadString('\n')
				}
				fileLocation = strings.Replace(fileLocation, "\r\n", "", -1)
				fileLocation = strings.TrimSpace(fileLocation)
				file, err = os.Open(fileLocation)
				if err != nil {
					fmt.Println("ERROR: ", err)
				}
			}
			// read in message to embed
			fmt.Print("Enter message to embed: ")
			message, err := reader.ReadString('\n')
			fmt.Println(message)
			if err != nil {
				fmt.Println("ERROR: ", err)
				// log.Fatal(err)
			}
			//defer file.Close()
			file.Seek(0, 0)
			source, _, err := image.Decode(file)
			extracted := extractMessage(source)
			fmt.Println("Extracted is: ", extracted)
			b := source.Bounds()
			m := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
			// Draw(dst Image, r image.Rectangle, src image.Image, sp image.Point, op Op)
			draw.Draw(m, m.Bounds(), source, b.Min, draw.Src)

			embedMessage(m, message)
			// jpeg.Encode(file, m, nil)
			err = png.Encode(file, m)
			if err != nil {
				fmt.Println("ERROR: ", err)
			}

			fmt.Println("successfully embedded message!")
			//extractedMessage := extractMessage(*m)
			file.Seek(0, 0)
			source, _, err = image.Decode(file)
			extractedMessage := extractMessage(source)
			fmt.Println("Extracted message is: ", extractedMessage)
			fmt.Println("=================================================")
			fmt.Println()

			file.Close()

			// case "extract":
			// read in picture & location
			// fmt.Print("Enter image (jpg) location: ")
			// fileLocation, err := reader.ReadString('\n')
			// for err != nil {
			// 	fmt.Printf("file location %s not valid!\n", fileLocation)
			// 	fmt.Print("Re-enter image (jpg) location: ")
			// 	fileLocation, err = reader.ReadString('\n')
			// }
			// fileLocation = strings.Replace(fileLocation, "\r\n", "", -1)
			// fileLocation = strings.TrimSpace(fileLocation)
			// file, err := os.Open(fileLocation)
			// file, err := os.Open("C:/Users/miche/Desktop/catP.png")
			// if err != nil {
			// 	fmt.Printf("ERROR: Image not found - %s\n", err.Error())
			// 	os.Exit(1)
			// }
			// // defer file.Close()
			// file.Seek(0, 0)
			// source, _, err := image.Decode(file)
			// b := source.Bounds()
			// m := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
			// draw.Draw(m, m.Bounds(), source, b.Min, draw.Src)

			// message := extractMessage(*m)
			// fmt.Println("Extracted message is: ", message)
			// file.Close()

		default:
			fmt.Println("Invalid input")
		}

		// file, err := os.Open("C:/Users/miche/Desktop/cat.jpg")
		// file, err := os.Open(fileLocation)

		// if err != nil {
		// 	fmt.Printf("ERROR: Image not found - %s\n", err.Error())
		// 	os.Exit(1)
		// }

		// defer file.Close()

		// config, format, err := image.DecodeConfig(file)
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// fmt.Println("Width:", config.Width, "Height:", config.Height, "Format:", format)

		// file.Seek(0, 0)

		// source, _, err := image.Decode(file)
		// b := source.Bounds()
		// m := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
		// draw.Draw(m, m.Bounds(), source, b.Min, draw.Src)

		// embedMessage(m, "hello there")
		// jpeg.Encode(file, m, nil)

		// message := extractMessage(*m)

		// fmt.Println("Extracted message is: ", message)
	}
}
