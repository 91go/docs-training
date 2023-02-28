package dir

import (
	"fmt"
	"testing"
)

func TestDir_Xz(t *testing.T) {
	d := NewDir("arch")
	d.Xz()
	fmt.Println(d)
}
