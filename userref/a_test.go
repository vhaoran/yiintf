package userref

import (
	"fmt"
	"testing"

	"github.com/vhaoran/yiintf/userref/ep"
)

func Test_aaa(t *testing.T) {
	fmt.Println("abc")
	_, err := new(ep.InnerMasterInfoGetH).Call(ep.InnerMasterInfoGetIn{})

	fmt.Println(err)
}
