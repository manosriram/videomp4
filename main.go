package main

import (
	"encoding/binary"
	"fmt"
	"os"
)

// var position = 0

type mp4 interface{}

type Box struct {
	size int
	typ  []byte
	data []byte
}

type Iods struct {
	typ     []byte
	size    int
	flags   int
	version int
	start   int
	data    []byte
}

type MovieHeaderBox struct {
	typ               []byte
	size              int
	flags             int
	version           int
	start             int
	creation_time     int
	modification_time int
	timescale         int
	duration          int
	rate              int
	volume            int
	matrix            []byte
	next_track_id     int
}

func ByteToInt(b []byte) int {
	if lesser := len(b) < binary.MaxVarintLen64; lesser {
		b = b[:cap(b)]
	}
	return int(binary.BigEndian.Uint32(b))
}

func readIntFromByte(b []byte, position *int, byteCount int) int {
	i := ByteToInt(b[*position : *position+byteCount])
	*position += byteCount
	return i
}

func (b *Box) iods(data []byte, position int) Box {
	sizeInInt := readIntFromByte(data, &position, 4)
	typ := data[position : position+4]
	position += 4
	fmt.Println(sizeInInt, string(typ))

	// version := readIntFromByte(data, &position, 1)
	// flags := readIntFromByte(data, &position, 3)
	dtag := readIntFromByte(data, &position, 1)
	// dext := data[position : position+4]
	dext1 := readIntFromByte(data, &position, 1)
	dext2 := readIntFromByte(data, &position, 1)
	dext3 := readIntFromByte(data, &position, 1)
	ln := readIntFromByte(data, &position, 1)
	// dlen := readIntFromByte(data, &position, 1)
	// start := readIntFromByte(data, &position, 4)
	// d := data[position:sizeInInt]
	fmt.Println(dtag, dext1, dext2, dext3, ln)

	return Box{}
}

func (b *Box) mvhd(data []byte, position int) MovieHeaderBox {
	sizeInInt := readIntFromByte(data, &position, 4)

	typ := data[position : position+4]
	position += 4

	version := readIntFromByte(data, &position, 1)
	flagsInInt := readIntFromByte(data, &position, 3)
	creationTime := readIntFromByte(data, &position, 4)
	modificationTime := readIntFromByte(data, &position, 4)
	timescale := readIntFromByte(data, &position, 4)
	duration := readIntFromByte(data, &position, 4)
	rate := readIntFromByte(data, &position, 4)

	position += 10

	volume := readIntFromByte(data, &position, 2)
	mat := data[position : position+36]

	position += 36
	position += 24 // predefined

	nextTrackId := readIntFromByte(data, &position, 4)

	b.iods(data, position)

	return MovieHeaderBox{
		size:              sizeInInt,
		typ:               typ,
		version:           version,
		flags:             flagsInInt,
		creation_time:     creationTime,
		modification_time: modificationTime,
		timescale:         timescale,
		duration:          duration,
		rate:              rate,
		volume:            volume,
		matrix:            mat,
		next_track_id:     nextTrackId,
	}
}

func (b *Box) moov(data []byte, position int) Box {
	sz := data[position : position+4]
	position += 4
	sizeInInt := ByteToInt(sz)

	typ := data[position : position+4]
	position += 4

	moovData := data[position:sizeInInt]

	b.mvhd(moovData, 0)

	return Box{
		size: sizeInInt,
		typ:  typ,
		data: moovData,
	}
}

func (b *Box) ftyp(data []byte, position int) Box {
	sz := data[position : position+4]
	position += 4
	sizeInInt := ByteToInt(sz)

	typ := data[position : position+4]
	position += 4

	// fmt.Println(sizeInInt - 8)
	d := data[position:sizeInInt]
	position += sizeInInt - 8

	b.moov(data, position)

	return Box{
		size: sizeInInt,
		typ:  typ,
		data: d,
	}
}

func main() {
	f, _ := os.Open("./video.mp4")
	fileSize, _ := f.Stat()
	b := make([]byte, fileSize.Size())

	f.Read(b)

	box := Box{}
	box.ftyp(b, 0)
	// box.moov(b)
	// fmt.Println(string(moov.data))
}

// box:ftyp
// size:4
// type:4
// data:value(size)
