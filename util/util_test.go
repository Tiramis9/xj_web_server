package util

import (
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"math/rand"
	"testing"
)

func TestMD5(t *testing.T) {

	//fmt.Printf("%x \n", 0x23a&0x0f/0x1)
	//
	//fmt.Print(23 / 5,23%5)
	//
	//return

	//var (
	//	in       = "1"
	//	expected = "c4ca4238a0b923820dcc509a6f75849b"
	//)
	//actual := MD5(in)
	//if actual != expected {
	//	t.Errorf("MD5(%s) = %s; expected %s", in, actual, expected)
	//}

	//var x uint64 = 1
	//bytes , err:= IntToBytes(x)
	//fmt.Print(bytes,err)

	//i, err := BytesToInt16([]byte{0,1})
	//fmt.Println(i,err)

	//var x = []byte{1,2,4,5}
	//fmt.Print(x[3:4])
	var json = `{
    "app_version": "1",
    "app_name": "1",
    "num_register_origin": 1,
    "machine_id": "11",
    "user_id":0,
    "page": 1,
    "size":50
}`
	xx := VerifySign(base64.StdEncoding.EncodeToString([]byte(json)), DesKey)
	fmt.Println(xx)

	//fmt.Println(MD5(json))
	//z := 11101010101010110
	//fmt.Println(z % 10) //个位
	//fmt.Println((z % 100) / 10) //十位
	//fmt.Println((z % 1000) / 100) //十位

	m := []byte{1, 2, 3}
	newY := make([]byte, 4+len(m))
	// 默认使用小端序
	binary.BigEndian.PutUint16(newY, uint16(len(m)+2))
	binary.BigEndian.PutUint16(newY[2:], uint16(3))
	copy(newY[4:], m)
	fmt.Println(newY)
	fmt.Println(rand.Int())
	//binary.BigEndian.PutUint16(m[2:], uint16(2))
	var x map[string]string
	fmt.Println(x)
	x["1"] = "0"
	delete(x, "1")
}
