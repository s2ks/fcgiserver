package util

import (
	"testing"
	"github.com/s2ks/fcgiserver/logger"
)

func Test0(t *testing.T) {
	logger.LogLevel(logger.LogLevelDebug)
}

func TestSub1(t *testing.T) {
	str, err := Substitute("%WORD% World", "WORD", "Hello", "%")
	if str != "Hello World" {
		t.Errorf("TestSub1 failed %v", err)
	}
}

func TestSub2(t *testing.T) {
	keyval := make(map[string]string)

	keyval["KEY"] = "Hello"
	keyval["KEY2"] = "And goodbye"
	str, err := SubstituteMap("%KEY%, %KEY%, %KEY%, world. %KEY2%.", keyval, "%")

	if str != "Hello, Hello, Hello, world. And goodbye." {
		t.Errorf("TestSub2 failed %v", err)
	}
}

func TestSub3(t *testing.T) {
	keyval := make(map[string]string)

	keyval["KEY"] = "Hello"
	keyval["KEY2"] = "And goodbye"
	str, err := SubstituteMap("%KEY3%, %kEY%, %KEY% KEY% %KEY, world. %KEY2% KEY3.", keyval, "%")

	if str != "%KEY3%, %kEY%, Hello KEY% %KEY, world. And goodbye KEY3." {
		t.Errorf("TestSub3 failed %v", err)
	}
}

func TestSub4(t *testing.T) {
	keyval := make(map[string]string)

	raw := "%KEY%, world. %KEY2%"

	keyval["KEY6"] = "Hello"
	keyval["KEY5"] = "And goodbye"
	str, err := SubstituteMap(raw, keyval, "%")

	if raw != str && err != nil {
		t.Errorf("TestSub4 failed")
	}
}

func TestSub5(t *testing.T) {
	keyval := make(map[string]string)

	keyval["KEY"] = "key"

	dest, err := SubstituteMap("%KEY% %KEY% %KEY% %KEY%", keyval, "%")

	if err != nil {
		t.Error(err)
		return
	}

	if string(dest) != "key key key key" {
		t.Error("TestSub5 failed")
	}
}

func TestBuffer1(t *testing.T) {
	var buf *Buffer

	buf = new(Buffer)

	buf.Write([]byte("Hello world"))

	if string(buf.Bytes()) != "Hello world" {
		t.Error("TestBuffer1 failed (stage 1)")
	}

	buf.Write([]byte(", I hope you are doing well"))

	if string(buf.Bytes()) != "Hello world, I hope you are doing well" {
		t.Error("TestBuffer1 failed (stage 2)")
	}
}
