package jwt

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

func TestEasyToken_GetToken(t *testing.T) {
	token, _ := EasyToken{
		Username: strconv.Itoa(29),
		Expires:  time.Now().Unix() + 3600*24*30, //Segundos

	}.GetToken()
	fmt.Println(token)
	b, s, e := EasyToken{}.ValidateToken(token)
	fmt.Println(b, s, e)
	var tokenStr = `eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NzYwNTI5OTEsImlzcyI6IjEwIiwibmJmIjoxNTczNDYwOTkxfQ.wMOzjsLb-vq7mhiXSt3bykzwhP9XaL24kFOR9ZNqlgn7bFCGHwo8vyPo5iojs38C1-z4_MOWhyexXo7XDoGw-BLpFT0uQsVPKDR57UlVHQ8D7XCkhHgPnmwyejdll_8vXk9gS27KzfcPPMDl5DbKlE04Vo1FvP3_gdV7D8zn5BmumEr6a9eJZ07eIaH_i2ElxsNsU6eZuJRLRmiqDXYw37d07EWCfqWU__qw423rk3O5eNmpnfsWEztPbc6vDN8shVxbV8RaDbED32EUdyXZm2v0P8lg1Mki7lHlXeyriQvJuJY-PWeVRs60IC-yQhopgAlMsie9kJ9JyC7KyDtm4A`
	b1, s1, e1 := EasyToken{}.ValidateToken(tokenStr)
	fmt.Println(b1, s1, e1)

}
