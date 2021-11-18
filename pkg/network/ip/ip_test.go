package ip

import (
	"fmt"
	"testing"
)

func TestFindIP(t *testing.T) {
	qqWry, err := NewQQwry()
	if err != nil {
		fmt.Println(err)
	}
	cc := qqWry.Find("111.60.118.155")
	fmt.Println(cc)
}
