package server

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"encoding/json"
	"io"
	"net"
)

type BasePacket struct {
	Name string `json:"name"`
}

func readSize(conn net.Conn) (uint32, error) {
	var size_packet uint32
	buf_size_raw := make([]byte, 4)
	_, err := conn.Read(buf_size_raw)
	if err != nil {
		return 0, err
	}

	reader_size := bytes.NewReader(buf_size_raw)
	err = binary.Read(reader_size, binary.BigEndian, &size_packet)
	if err != nil {
		return 0, err
	}

	return size_packet, nil
}

func encodeSize(data []byte) ([]byte, error) {
	buffer := new(bytes.Buffer)
	err := binary.Write(buffer, binary.BigEndian, uint32(len(data)))
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func readPacket(conn net.Conn, size uint32) ([]byte, error) {
	buf := make([]byte, size)
	_, err := conn.Read(buf)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func compressPacket(data string) ([]byte, error) {
	var buffer bytes.Buffer
	writer := zlib.NewWriter(&buffer)
	_, err := writer.Write([]byte(data))
	if err != nil {
		return nil, err
	}

	err = writer.Close()
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func deflatePacket(data []byte) (string, error) {
	buffer := bytes.NewReader(data)
	reader, err := zlib.NewReader(buffer)
	if err != nil {
		return "", err
	}

	result, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}

	return bytes.NewBuffer(result).String(), nil
}

func getPacketName(data string) (string, error) {
	var packet BasePacket
	err := json.Unmarshal([]byte(data), &packet)
	if err != nil {
		return "", err
	}
	return packet.Name, nil
}
