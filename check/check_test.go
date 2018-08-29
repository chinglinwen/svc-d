package check

import (
	"fmt"
	"testing"
	"time"
)

func TestSimpleCheck(t *testing.T) {
	fmt.Println(SimpleCheck("172.28.137.221", "8000"))
}

func TestSimpleCheckLonger(t *testing.T) {
	fmt.Println(SimpleCheckLonger("172.28.137.221", "8000", 10*time.Second))
}
