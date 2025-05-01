package main

import "fmt"

func WriteUint16(buf []byte, n uint16) []byte {
	return append(buf, byte((n>>8)&0xff),byte(n&0xff))
}

func WriteString(buf []byte, s string) []byte{
	if len(s) > 0xffff{
		panic("WriteString is too big")
	}
	buf = WriteUint16(buf, uint16(len(s)))
	return append(buf, []byte(s)...)
}

func WriteStringSlice(buf []byte, s []string) []byte{
	if len(s) > 0xffff{
		panic("WriteStringSlice is too big")
	}
	buf = WriteUint16(buf, uint16(len(s)))
	
	for _, elt := range s {
		buf = WriteString(buf, elt)
	}
	return buf
}

func ReadUint16(buf []byte) (uint16, []byte, error){
	if len(buf) < 2 {
		return 0, nil, fmt.Errorf("ReadUint16 has reached end of input")
	}
	return uint16(buf[0])<<8 | uint16(buf[1]), buf [2:], nil
}

func ReadString(buf []byte) (string, []byte, error){
	var size uint16
	size, buf, err := ReadUint16(buf)
	if err != nil {
		return "", nil, err
	}

	//Checks the size
	if len(buf) < int(size) {
		return "", nil, fmt.Errorf("ReadString given insufficient bytes")
	}

	return string(buf[:size]), buf[size:], nil
}

func ReadStringSlice(buf []byte) ([]string, []byte, error){
	var size uint16
	size, buf, err := ReadUint16(buf)
	if err != nil {
		return nil, nil, err
	}

	result := make([]string, size)
	for i := range result{
		result[i], buf, err = ReadString(buf)
		if err != nil {
			return nil, nil, err
		}
	}
	return result, buf, nil
}