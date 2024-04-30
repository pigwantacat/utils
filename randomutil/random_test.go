package randomutil

import (
	"fmt"
	"testing"
)

func TestRandomIntArray(t *testing.T) {
	for i := 0; i < 100; i++ {
		fmt.Println(RandomIntArray(-100, 100, 10))
	}
}

func TestRandomString(t *testing.T) {
	for i := 0; i < 100; i++ {
		fmt.Println(RandomString(10))
	}
}
