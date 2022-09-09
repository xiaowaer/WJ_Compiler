package syntax

import (
	"fmt"
)
type Mytest struct{
fafas string
asdf string
aaa[64] token
tok token
}
var q[128] token


func myhash(s []byte) uint {
	return (uint(s[0])<<4 ^ uint(s[1]) + uint(len(s))) & uint(len(q)-1)
}

// func init(){
// 	// populate keywordMap
// 	for tok := _Break; tok <= _Var; tok++ {
// 		h := myhash([]byte(tok.String()))
// 		if q[h] != 0 {
// 			panic("imperfect hash")
// 		}
// 		q[h] = tok
//	}
//}

func (m *Mytest) init(){
    	// populate keywordMap
	for tok := _Break; tok <= _Var; tok++ {
		h := myhash([]byte(tok.String()))
		if q[h] != 0 {
			panic("imperfect hash")
		}
		q[h] = tok
	}
}
    
func (s *Mytest) hhh() uint{
var str string = "if"
var data []byte = []byte(str)
       a :=myhash(data) 
fmt.Println(a)
return a
}




