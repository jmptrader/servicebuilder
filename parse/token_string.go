// generated by stringer -type=Token; DO NOT EDIT

package parse

import "fmt"

const _Token_name = "ILLEGALEOFWSIDENTLEFTBRACERIGHTBRACELEFTSQBRACERIGHTSQBRACECOLONCOMMANUMERICMODELFIELDSPAGINATIONACTIONSSTRINGINTDOUBLEDATEDATETIME"

var _Token_index = [...]uint8{7, 10, 12, 17, 26, 36, 47, 59, 64, 69, 76, 81, 87, 97, 104, 110, 113, 119, 123, 131}

func (i Token) String() string {
	if i < 0 || i >= Token(len(_Token_index)) {
		return fmt.Sprintf("Token(%d)", i)
	}
	hi := _Token_index[i]
	lo := uint8(0)
	if i > 0 {
		lo = _Token_index[i-1]
	}
	return _Token_name[lo:hi]
}
