package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strings"
)

func decodeChunk(byteChunk []byte, n int) ([]byte, error) {
	if len(byteChunk) < n {
		log.Fatalln("Len(byteChunk) < ", n)
	}
	if len(byteChunk) != 4 {
		return nil, errors.New("string is not multiple of 4")
	}

	data := string(byteChunk)

	var tmp [4]byte
	var out []byte
	var x uint32
	padIdx := len(data)
	for i, letter := range data {
		if padIdx == len(data) && string(letter) == Pad {
			padIdx = i

		}
		x = x | Decoding64[string(letter)]<<(6*(3-i))
	}
	binary.BigEndian.PutUint32(tmp[:], x)
	out = tmp[1:padIdx]

	return out, nil
}

func encodeChunk(byteChunk []byte, n int) string {
	if len(byteChunk) < n || n > 3 {
		log.Fatalln("Len(byteChunk) < ", n, "|| n > 3")
	}

	var out string
	var tmp [4]byte
	copy(tmp[1:], byteChunk)

	chunk := binary.BigEndian.Uint32(tmp[:])
	numChars := int(math.Ceil((float64(n) * 8.0) / 6.0))
	for i := 0; i < numChars; i++ {
		x := (chunk >> (6 * (3 - i))) & (0x40 - 1)
		out += Encoding64[int(x)]
	}

	for i := len(out); i <= 3; i++ {
		out += Pad
	}

	return out
}

func decodeData(in io.Reader, out io.Writer) error {
	for {
		buff := make([]byte, 4)

		n, err := in.Read(buff)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		chunk, err := decodeChunk(buff, n)
		if err != nil {
			return err
		}

		io.Copy(out, bytes.NewReader(chunk))
	}

	return nil
}

func encodeData(in io.Reader, out io.Writer) error {
	for {
		buff := make([]byte, 3)

		n, err := in.Read(buff)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		chunk := encodeChunk(buff, n)
		io.Copy(out, strings.NewReader(chunk))
	}

	return nil
}

func encodeString(data string) string {
	r := strings.NewReader(data)
	var buf bytes.Buffer

	err := encodeData(r, &buf)
	if err != nil {
		log.Fatalln(err)
	}
	return buf.String()
}

func decodeString(data string) string {
	r := strings.NewReader(data)
	var buf bytes.Buffer

	err := decodeData(r, &buf)
	if err != nil {
		log.Fatalln(err)
	}
	return buf.String()
}

func main() {

	decode := flag.Bool("d", false, "Decode base64-file")
	flag.Parse()
	files := flag.Args()

	if len(files) != 1 {
		fmt.Println("Usage: gobase64 <file>")
		os.Exit(1)
	}

	file, err := os.Open(files[0])
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	if *decode {
		err = decodeData(file, os.Stdout)
	} else {
		err = encodeData(file, os.Stdout)
	}
	if err != nil {
		log.Fatalln(err)
	}
}
