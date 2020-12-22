package line

import "testing"

func TestSend(t *testing.T) {
	tk := "G27YdjfMjtHTJU74QpF1wp6UmtpInL6LZkCdgpbqxj9"
	Init(&tk)
	Send("Hello World.")
}
