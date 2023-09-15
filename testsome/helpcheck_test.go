package testsome

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"
)
func getRand(length int, isDigit bool) string {
	if length < 1 {
	 
		return ""
	}
	var char string
	if isDigit {
		char = "0123456789"
	} else {
		char = "abcdefg0123456789"
	}
	charArr := strings.Split(char, "")
	charlen := len(charArr)
	ran := rand.New(rand.NewSource(time.Now().Unix()))
	var rchar string = ""
	for i := 1; i <= length; i++ {
		rchar = rchar + charArr[ran.Intn(charlen)]
	}
	return rchar
}

func TestFunc(t *testing.T) {
	fmt.Println(getRand(20, false))
}