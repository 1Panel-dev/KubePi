package randomstring

import (
	"crypto/rand"
	"math/big"
	"strconv"
)

//Generate a random string in specified length
func Generate(length int) string {
	const base = 36
	size := big.NewInt(base)
	n := make([]byte, length)
	for i := range n {
		c, _ := rand.Int(rand.Reader, size)
		n[i] = strconv.FormatInt(c.Int64(), base)[0]
	}
	return string(n)
}
