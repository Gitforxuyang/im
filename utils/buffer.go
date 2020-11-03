package utils

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

func Encode(buf []byte) ([]byte,error){
	var len int32= int32(len(buf))
	fmt.Println(len)
	pkg:=new(bytes.Buffer)
	err:=binary.Write(pkg,binary.BigEndian,len)
	Must(err)
	err=binary.Write(pkg,binary.BigEndian,buf)
	Must(err)
	return pkg.Bytes(),nil
}
