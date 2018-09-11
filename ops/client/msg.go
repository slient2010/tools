package client

import (
	"errors"
	"fmt"
	"net"
	"time"
)

const (
	headerByteCount = 4
)

const (
	JOB_ID_UPDATE_CLIENT_STAT = 0xffff
)

func Assemble(data []byte, v int) []byte {
	return append(data, byte(v))
}

func assemble2(data []byte, v int) []byte {
	return append(data, byte(v>>8), byte(v))
}

// func AssembleUtfStr(data []byte, str string) []byte {
// 	strBytes := []byte(str)
// 	return Assemble2(append(data, strBytes...), len(strBytes))
// }

func write(c net.Conn, bs []byte) error {
	l := len(bs)
	i := 0
	for {
		n, err := c.Write(bs[i:])
		if err != nil {
			return err
		}
		i += n
		if i >= l {
			break
		} else {
			fmt.Println("net.Conn.Write not write all data, tcp sendBuf overflow?", i, l)
			time.Sleep(time.Second)
		}
	}
	return nil
}

func writeBytes(c net.Conn, data []byte) error {
	bc, v := headerByteCount, len(data)

	bs := make([]byte, bc)
	for i, j := bc-1, 0; i >= 0; i-- {
		bs[j] = byte(v >> uint(i*8))
		j += 1
	}

	return write(c, append(bs, data...))
}

func writeString(c net.Conn, s string) error {
	return writeBytes(c, []byte(s))
}

func readBool(c net.Conn) (bool, error) {
	r, err := read(c, 1)
	if err != nil {
		return false, err
	}
	return r[0] == 1, nil
}

func split2(data []byte) ([]byte, int) {
	l := len(data) - 2
	return data[:l], int(data[l])<<8 | int(data[l+1])
}

func read(c net.Conn, bc int) ([]byte, error) {
	bs := make([]byte, bc)
	ix := 0
	for ix < bc {
		n, err := c.Read(bs[ix:])
		if err != nil {
			return nil, err
		}
		ix += n
	}
	return bs, nil
}

func readInt(c net.Conn, bc int) (int, error) {
	if bc > 4 {
		return -1, errors.New("ReadInt byteCount too large!")
	}

	bs, err := read(c, bc)
	if err != nil {
		return -1, err
	}

	r := 0
	for i, j := bc-1, 0; i >= 0; i-- {
		r |= int(bs[j]) << uint(i*8)
		j += 1
	}

	return r, nil
}

func readBytes(c net.Conn) ([]byte, error) {
	bc, err := readInt(c, headerByteCount)
	if err != nil {
		return nil, err
	}

	bs, err := read(c, bc)
	if err != nil {
		return nil, err
	}

	return bs, nil
}
