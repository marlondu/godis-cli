package core

import (
	"fmt"
	"testing"
)

// created: 2018/9/6

func TestEncrypt(t *testing.T) {
	src := "I love you so how much you will never know"
	enc := Encrypt(src)
	fmt.Println(enc)
	t.Logf("Ecrypt:%s\n", enc)
}

func TestDecrypt(t *testing.T) {
	src := "dn59CF4QjkBhwuTYuPnnt2o/o14Vz3vR3pl2/waF6mY+BEQTUgDM3tV/jQiu0Bf+"
	dec := Decrypt(src)
	if dec != "" {
		t.Logf("decrypted to:%s", dec)
	} else {
		t.Errorf("decrypted error")
	}
}
