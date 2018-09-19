package notice

import (
	"fmt"
	"testing"
)

func TestSend(t *testing.T) {
	r, err := Send("wenzhenglin", "hello")
	if err != nil {
		t.Errorf("send err %v\n", err)
	}
	fmt.Println("reply", r)
}
