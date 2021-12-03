package line

import "testing"

func TestSend(t *testing.T) {
	Init()
	tk := "G27YdjfMjtHTJU74QpF1wp6UmtpInL6LZkCdgpbqxj9"
	Send(&tk, "Hello World.")
}
