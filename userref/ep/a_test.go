package ep

import (
	"fmt"
	"testing"
)

func Test_aaa(t *testing.T) {
	fmt.Println("abc")
	_, err := new(InnerMasterInfoGetH).Call(&InnerMasterInfoGetIn{})

	fmt.Println(err)
}
