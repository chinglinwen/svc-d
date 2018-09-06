package check

import (
	"fmt"
	"testing"
	"time"
)

func TestSimpleCheck(t *testing.T) {
	fmt.Println(SimpleCheck("172.28.137.221", "8000"))
}

func TestCheckIPWithConfig(t *testing.T) {
	fmt.Println(CheckIPWithConfig("ops_test", "172.28.40.251", "3000"))
}

func TestSimpleCheckLonger(t *testing.T) {
	fmt.Println(CheckLonger("ops_test", "172.28.137.221", "8000", 10*time.Second))
}
