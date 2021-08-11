package utils

import (
	"fmt"
	"testing"
)

func TestIsIP(t *testing.T) {
	res := IsIP("10.0.0.1", true)
	fmt.Println(res)
	res = IsIP("10.0.0", true)
	fmt.Println(res)
	res = IsIP("10.0.0", false)
	fmt.Println(res)
}
