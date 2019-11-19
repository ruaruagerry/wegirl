package gamecfg

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_UnitToString(t *testing.T) {

	u1 := UNumber{U: 10, V: 11.34}
	assert.Equal(t, u1.String(), "11.34J")

	u2 := UNumber{U: 320, V: 11.34}
	fmt.Println(u2.String())
	assert.Equal(t, u2.String(), "11.34LH")
}

func Test_UnitAdd(t *testing.T) {
	u11 := UNumber{U: 10, V: 11.34}
	u21 := UNumber{U: 10, V: 11.34}
	u11.Add(&u21)
	assert.Equal(t, u11.String(), "22.68J")

	u12 := UNumber{U: 10, V: 11.34}
	u22 := UNumber{U: 320, V: 11.34}
	u12.Add(&u22)
	assert.Equal(t, u12.String(), "11.34LH")

	u13 := UNumber{U: 10, V: 1134}
	u23 := UNumber{U: 11, V: 11.34}
	u13.Add(&u23)
	assert.Equal(t, u13.String(), "12.47K")

}

func Test_UnitDiv(t *testing.T) {
	u11 := UNumber{U: 10, V: 11.34}
	u21 := UNumber{U: 10, V: 11.34}
	u11.Div(&u21)
	assert.Equal(t, u11.String(), "1.00raw")

	u12 := UNumber{U: 10, V: 22.68}
	u22 := UNumber{U: 10, V: 11.34}
	u12.Div(&u22)
	assert.Equal(t, u12.String(), "2.00raw")

	u13 := UNumber{U: 9, V: 1000}
	u23 := UNumber{U: 10, V: 1}
	u13.Div(&u23)
	assert.Equal(t, u23.String(), "1.00raw")
}

func Test_UnitLe(t *testing.T) {
	u11 := UNumber{U: 10, V: 11.34}
	u21 := UNumber{U: 10, V: 11.34}
	assert.Equal(t, u11.Le(&u21), true)

	u12 := UNumber{U: 11, V: 11.34}
	u22 := UNumber{U: 10, V: 11.34}
	assert.Equal(t, u12.Le(&u22), false)

	u13 := UNumber{U: 11, V: -11.34}
	u23 := UNumber{U: 10, V: 11.34}
	assert.Equal(t, u13.Le(&u23), true)

	u14 := UNumber{U: 11, V: 11.34}
	u24 := UNumber{U: 12, V: 11.34}
	assert.Equal(t, u14.Le(&u24), true)
}

func Test_UnitGe(t *testing.T) {
	u11 := UNumber{U: 10, V: 11.34}
	u21 := UNumber{U: 10, V: 11.34}
	assert.Equal(t, u11.Ge(&u21), true)

	u12 := UNumber{U: 11, V: 11.34}
	u22 := UNumber{U: 10, V: 11.34}
	assert.Equal(t, u12.Ge(&u22), true)

	u13 := UNumber{U: 11, V: -11.34}
	u23 := UNumber{U: 10, V: 11.34}
	assert.Equal(t, u13.Ge(&u23), false)

	u14 := UNumber{U: 11, V: 11.34}
	u24 := UNumber{U: 12, V: 11.34}
	assert.Equal(t, u14.Ge(&u24), false)

}

func Test_ParseBigNumber(t *testing.T) {
	bigNumber := 4.40282346638528859811704183484516925440e+38
	un := ParseBigNumber(bigNumber)
	fmt.Println(un)
}
