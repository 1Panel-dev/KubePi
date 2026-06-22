package ip

import (
	"os"
	"testing"
)

func TestFindIP(t *testing.T) {
	data, err := os.ReadFile("../../cmd/server/helper/ip/qqwry.dat")
	if err != nil {
		t.Skipf("qqwry.dat is not available: %v", err)
	}
	IpCommonDictionary = data

	qqWry, err := NewQQwry()
	if err != nil {
		t.Fatal(err)
	}
	cc := qqWry.Find("111.60.118.155")
	if cc.IP != "111.60.118.155" {
		t.Fatalf("unexpected ip result: %+v", cc)
	}
	t.Log(cc)
}
