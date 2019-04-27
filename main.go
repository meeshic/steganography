package main

import (
	"bufio"
	"crypto/rand"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"math"
	"math/big"
	"os"
	"strings"
)

func init() {
	// image.RegisterFormat("jpeg", "jpeg", jpeg.Decode, jpeg.DecodeConfig)
	image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)
}

// Thanks to: https://medium.com/@kpbird/golang-generate-fixed-size-random-string-dd6dbd5e63c0
func randomString(len int) string {
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		randomNum, err := rand.Int(rand.Reader, big.NewInt(26))
		if err != nil {
			log.Println("ERROR generating random number: ", err)
		}
		bytes[i] = byte(65 + randomNum.Int64()) //A=65 and Z = 65+25
	}
	return string(bytes)
}

func embedMessage(img *image.NRGBA, message string) bool {
	byteMessage := []byte(message)
	// for _, b := range byteMessage {
	// 	fmt.Print(string(b))
	// 	// fmt.Printf("% 08b \n", b)
	// }

	maxWidth := img.Bounds().Dx()

	height := 0
	width := 0
	var bit byte

	// 1 byte = 8 bits
	// 1 pixel = 4 bits (R,G,B,A)
	// 1 byte -> 2 pixels
	for i := 0; i < len(byteMessage); i++ {
		// fmt.Println("encode char: ", string(byteMessage[i]))
		if width == maxWidth {
			height++
			width = 0
		}
		pix := img.NRGBAAt(height, width)
		r, g, b, a := pix.R, pix.G, pix.B, pix.A
		bit = byteMessage[i] & byte(math.Pow(2, float64(0)))
		// fmt.Println("bit 1: ", bit)
		// fmt.Println("r before: ", strconv.FormatInt(int64(r), 2))
		r = (r & 0xfe) | bit
		// fmt.Println("r after: ", strconv.FormatInt(int64(r), 2))
		bit = byteMessage[i] & byte(math.Pow(2, float64(1)))
		// fmt.Println("bit 2: ", bit>>1)
		// fmt.Println("g before: ", strconv.FormatInt(int64(g), 2))
		g = (g & 0xfe) | (bit >> 1)
		// fmt.Println("g after: ", strconv.FormatInt(int64(g), 2))
		bit = byteMessage[i] & byte(math.Pow(2, float64(2)))
		// fmt.Println("bit 3: ", bit>>2)
		// fmt.Println("b before: ", strconv.FormatInt(int64(b), 2))
		b = (b & 0xfe) | (bit >> 2)
		// fmt.Println("b after: ", strconv.FormatInt(int64(b), 2))
		bit = byteMessage[i] & byte(math.Pow(2, float64(3)))
		// fmt.Println("bit 4: ", bit>>3)
		// fmt.Println("a before: ", strconv.FormatInt(int64(a), 2))
		a = (a & 0xfe) | (bit >> 3)
		// fmt.Println("a after: ", strconv.FormatInt(int64(a), 2))
		img.Set(height, width, color.NRGBA{
			R: r,
			G: g,
			B: b,
			A: a,
		})

		width++
		if width == maxWidth {
			height++
			width = 0
		}
		pix = img.NRGBAAt(height, width)
		r, g, b, a = pix.R, pix.G, pix.B, pix.A
		bit = byteMessage[i] & byte(math.Pow(2, float64(4)))
		// fmt.Println("bit 5: ", bit>>4)
		// fmt.Println("r before: ", strconv.FormatInt(int64(r), 2))
		r = (r & 0xfe) | (bit >> 4)
		// fmt.Println("r after: ", strconv.FormatInt(int64(r), 2))
		bit = byteMessage[i] & byte(math.Pow(2, float64(5)))
		// fmt.Println("bit 6: ", bit>>5)
		// fmt.Println("g before: ", strconv.FormatInt(int64(g), 2))
		g = (g & 0xfe) | (bit >> 5)
		// fmt.Println("g after: ", strconv.FormatInt(int64(g), 2))
		bit = byteMessage[i] & byte(math.Pow(2, float64(6)))
		// fmt.Println("bit 7: ", bit>>6)
		// fmt.Println("b before: ", strconv.FormatInt(int64(b), 2))
		b = (b & 0xfe) | (bit >> 6)
		// fmt.Println("b after: ", strconv.FormatInt(int64(b), 2))
		bit = byteMessage[i] & byte(math.Pow(2, float64(7)))
		// fmt.Println("bit 8: ", bit>>7)
		// fmt.Println("a before: ", strconv.FormatInt(int64(a), 2))
		a = (a & 0xfe) | (bit >> 7)
		// fmt.Println("a after: ", strconv.FormatInt(int64(a), 2))
		img.Set(height, width, color.NRGBA{
			R: r,
			G: g,
			B: b,
			A: a,
		})
		width++
	}

	if width == maxWidth {
		height++
		width = 0
	}

	img.Set(height, width, color.RGBA{
		R: 0,
		G: 0,
		B: 0,
		A: 0,
	})

	width++
	if width == maxWidth {
		height++
		width = 0
	}

	img.Set(height, width, color.RGBA{
		R: 0,
		G: 0,
		B: 0,
		A: 0,
	})

	return true
}

// iterate over pixels to find embedded message
func extractMessage(img image.NRGBA) string {
	var message []byte

	maxHeight := img.Bounds().Dy()
	maxWidth := img.Bounds().Dx()

	var bits []byte

	for h := 0; h < maxHeight; h++ {
		for w := 0; w < maxWidth; w++ {
			pixel := img.NRGBAAt(h, w)
			r, g, b, a := pixel.R, pixel.G, pixel.B, pixel.A
			// fmt.Printf("extracted pixel NRBGA: %v %v %v %v\n", strconv.FormatInt(int64(r), 2), strconv.FormatInt(int64(g), 2), strconv.FormatInt(int64(b), 2), strconv.FormatInt(int64(a), 2))
			if r == 0 && g == 0 && b == 0 && a == 0 {
				nextW := w + 1
				nextH := h
				if nextW == maxWidth {
					nextH = h + 1
					nextW = 0
				}
				nextPixel := img.NRGBAAt(nextH, nextW)
				nextR, nextG, nextB, nextA := nextPixel.R, nextPixel.G, nextPixel.B, nextPixel.A
				if nextR == 0 && nextG == 0 && nextB == 0 && nextA == 0 {
					return string(message)
				}

			}
			bits = append(bits, byte(r&1))
			bits = append(bits, byte(g&1))
			bits = append(bits, byte(b&1))
			bits = append(bits, byte(a&1))
			if len(bits) == 8 {
				var charByte byte
				for i, bit := range bits {
					move := byte(i)
					charByte = charByte | (bit << move)
				}
				// fmt.Printf("charByte: % 08b \n", charByte)
				// fmt.Println("char: ", string(charByte))
				message = append(message, charByte)
				// fmt.Println("message in progress: ", string(message))
				bits = []byte{}
			}
		}
	}
	return ""
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("===========================================")
		fmt.Println("Enter: ")
		fmt.Println("'embed' to embed a message")
		fmt.Println("'extract' to extract a message")
		fmt.Println("'exit' to exit program")

		fmt.Print("> ")

		userInput, err := reader.ReadString('\n')
		userInput = strings.Replace(userInput, "\r\n", "", -1)
		userInput = strings.TrimSpace(userInput)
		for err != nil {
			log.Println("ERROR: ", err)
			log.Println("Please enter 'embed' or 'extract'")
		}
		switch userInput {
		case "embed":
			// read in picture location
			fmt.Print("	Enter image (png) location: ")
			fileLocation, err := reader.ReadString('\n')
			for err != nil {
				log.Printf("file location %s not valid!\n", fileLocation)
				log.Print("Re-enter image (png) location: ")
				fileLocation, err = reader.ReadString('\n')
			}
			fileLocation = strings.TrimSpace(fileLocation)
			fileLocation = strings.Replace(fileLocation, "\r\n", "", -1)
			file, err := os.OpenFile(fileLocation, os.O_RDWR, 0600)
			// file, err := os.OpenFile("C:/Users/miche/Desktop/Avatar_cat.png", os.O_RDWR, 0600)
			for err != nil {
				log.Printf("	ERROR: Image not found - %s\n", err)
				log.Print("	Re-enter image (png) location: ")
				fileLocation, err := reader.ReadString('\n')
				for err != nil {
					log.Printf("	ERROR: file location %s not valid!\n", fileLocation)
					log.Print("	Re-enter image (png) location: ")
					fileLocation, err = reader.ReadString('\n')
				}
				fileLocation = strings.TrimSpace(fileLocation)
				fileLocation = strings.Replace(fileLocation, "\r\n", "", -1)
				file, err = os.Open(fileLocation)
				if err != nil {
					log.Printf("	ERROR: Image cannot be opened - %s\n", err)
				}
			}
			// read in message to embed
			fmt.Print("	Enter message to embed: ")
			message, err := reader.ReadString('\n')
			if err != nil {
				log.Println("	ERROR in reading message: ", err)
			}

			file.Seek(0, 0)
			source, _, err := image.Decode(file)

			b := source.Bounds()
			m := image.NewNRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
			draw.Draw(m, m.Bounds(), source, b.Min, draw.Src)

			embedMessage(m, message)

			// fmt.Println("Did message go correctly?: ", extractMessage(*m))

			fileName := strings.Split(fileLocation, ".png")
			outputFileLocation := fileName[0] + "_" + randomString(4) + ".png"
			outputFile, err := os.Create(outputFileLocation)
			// outputFile, err := os.Create("C:/Users/miche/Desktop/Avatar_cat_1.png")
			if err != nil {
				fmt.Println("	ERROR: cannot create file")
			}
			outputFile.Seek(0, 0)
			err = png.Encode(outputFile, m)
			if err != nil {
				fmt.Println("	ERROR: ", err)
			}
			outputFile.Close()

			fmt.Println()
			fmt.Printf("Embedded the message to %s :) \n", outputFileLocation)

		case "extract":
			// read in picture & location
			fmt.Print("	Enter image (png) location: ")
			fileLocation, err := reader.ReadString('\n')
			for err != nil {
				fmt.Printf("	ERROR: file location %s not valid!\n", fileLocation)
				fmt.Print("	Re-enter image (png) location: ")
				fileLocation, err = reader.ReadString('\n')
			}
			fileLocation = strings.TrimSpace(fileLocation)
			fileLocation = strings.Replace(fileLocation, "\r\n", "", -1)

			// f, err := os.Open("C:/Users/miche/Desktop/Avatar_cat_1.png")
			f, err := os.Open(fileLocation)
			if err != nil {
				fmt.Println("	ERROR: cannot open image file")
			}
			f.Seek(0, 0)
			source, err := png.Decode(f)

			b := source.Bounds()
			m := image.NewNRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
			draw.Draw(m, m.Bounds(), source, b.Min, draw.Src)

			extractedMessage := extractMessage(*m)

			fmt.Println()
			fmt.Println("Extracted message is: ", extractedMessage)
			f.Close()

		case "exit":
			fmt.Println("BYE BYE ~")
			os.Exit(0)

		default:
			fmt.Println("Invalid input :( Try again")
		}
	}
}
