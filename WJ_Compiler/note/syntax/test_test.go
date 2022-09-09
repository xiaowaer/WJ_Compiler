package syntax

import (
	"fmt"
	"testing"
)


func TestA(t *testing.T) {
    var hh Mytest
    hh.init()
    a := hh.hhh()
   fmt.Printf("%d", uint64(a)) 
   want := a
   if(want <= 0 ){
	t.Errorf("%d", want)
   }
}
