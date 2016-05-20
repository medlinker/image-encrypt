package imageEncrypt

import (
	"os"
	"testing"
)

func TestCutting(t *testing.T) {
	f, err := os.Open("test-asserts/test1.png")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	m := NewMetaByRedis("127.0.0.1:6379", "test")
	s := NewFileStorage("test-asserts/")
	// c := NewDefaultRectangleCut(s, m)
	c := NewRectangleCut(8, 7, s, m)
	c.Cutting(f, "test1.png", "test1")
	if err != nil {
		t.Fatal(err)
	}
}
