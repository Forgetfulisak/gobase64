package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"flag"
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

	data := byteChunk

	var uintBuf [4]byte
	var x uint32
	padIdx := len(data)
	for i, letter := range data {
		if padIdx == len(data) && letter == Pad {
			padIdx = i

		}
		x = x | Decoding64[letter]<<(6*(3-i))
	}
	binary.BigEndian.PutUint32(uintBuf[:], x)
	out := uintBuf[1:padIdx]

	return out, nil
}

func encodeChunk(byteChunk []byte) []byte {
	if len(byteChunk) > 3 {
		panic("len(byteChunk) must be > 3")
	}

	n := len(byteChunk)
	out := make([]byte, 0, 4)

	var uintBuf [4]byte
	copy(uintBuf[1:], byteChunk)

	chunk := binary.BigEndian.Uint32(uintBuf[:])
	numChars := int(math.Ceil((float64(n) * 8.0) / 6.0))
	for i := 0; i < numChars; i++ {
		x := (chunk >> (6 * (3 - i))) & (0x40 - 1)
		out = append(out, Encoding64[int(x)])
	}

	for i := len(out); i <= 3; i++ {
		out = append(out, Pad)
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
		_, err = out.Write(chunk)
		if err != nil {
			return err
		}
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

		chunk := encodeChunk(buff[:n])
		_, err = out.Write(chunk)
		if err != nil {
			return err
		}
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
	var input io.Reader
	var err error

	decode := flag.Bool("d", false, "Decode base64-file")
	flag.Parse()
	files := flag.Args()

	if len(files) != 1 {
		input = os.Stdin
	} else {

		file, err := os.Open(files[0])
		if err != nil {
			log.Fatalln(err)
		}
		defer file.Close()
		input = file
	}

	if *decode {
		err = decodeData(input, os.Stdout)
	} else {
		err = encodeData(input, os.Stdout)
	}
	if err != nil {
		log.Fatalln(err)
	}
}
