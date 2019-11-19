package gamecfg

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

// Unit aa-bb-cc-... unit
type Unit int

const (
	// RAW base
	RAW Unit = iota

	// UnitOffest unit进位
	UnitOffest = 1000
)

var (
	// UnitString unit string
	UnitString = [...]string{
		"raw", "A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N",
		"O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z",
		"AA", "AB", "AC", "AD", "AE", "AF", "AG", "AH", "AI", "AJ", "AK", "AL", "AM", "AN", "AO", "AP", "AQ", "AR",
		"AS", "AT", "AU", "AV", "AW", "AX", "AY", "AZ",
		"BA", "BB", "BC", "BD", "BE", "BF", "BG", "BH", "BI", "BJ", "BK", "BL", "BM", "BN", "BO", "BP", "BQ", "BR",
		"BS", "BT", "BU", "BV", "BW", "BX", "BY", "BZ",
		"CA", "CB", "CC", "CD", "CE", "CF", "CG", "CH", "CI", "CJ", "CK", "CL", "CM", "CN", "CO", "CP", "CQ", "CR",
		"CS", "CT", "CU", "CV", "CW", "CX", "CY", "CZ",
		"DA", "DB", "DC", "DD", "DE", "DF", "DG", "DH", "DI", "DJ", "DK", "DL", "DM", "DN", "DO", "DP", "DQ", "DR",
		"DS", "DT", "DU", "DV", "DW", "DX", "DY", "DZ",
		"EA", "EB", "EC", "ED", "EE", "EF", "EG", "EH", "EI", "EJ", "EK", "EL", "EM", "EN", "EO", "EP", "EQ", "ER",
		"ES", "ET", "EU", "EV", "EW", "EX", "EY", "EZ",
		"FA", "FB", "FC", "FD", "FE", "FF", "FG", "FH", "FI", "FJ", "FK", "FL", "FM", "FN", "FO", "FP", "FQ", "FR",
		"FS", "FT", "FU", "FV", "FW", "FX", "FY", "FZ",
		"GA", "GB", "GC", "GD", "GE", "GF", "GG", "GH", "GI", "GJ", "GK", "GL", "GM", "GN", "GO", "GP", "GQ", "GR",
		"GS", "GT", "GU", "GV", "GW", "GX", "GY", "GZ",
		"HA", "HB", "HC", "HD", "HE", "HF", "HG", "HH", "HI", "HJ", "HK", "HL", "HM", "HN", "HO", "HP", "HQ", "HR",
		"HS", "HT", "HU", "HV", "HW", "HX", "HY", "HZ",
		"IA", "IB", "IC", "ID", "IE", "IF", "IG", "IH", "II", "IJ", "IK", "IL", "IM", "IN", "IO", "IP", "IQ", "IR",
		"IS", "IT", "IU", "IV", "IW", "IX", "IY", "IZ",
		"JA", "JB", "JC", "JD", "JE", "JF", "JG", "JH", "JI", "JJ", "JK", "JL", "JM", "JN", "JO", "JP", "JQ", "JR",
		"JS", "JT", "JU", "JV", "JW", "JX", "JY", "JZ",
		"KA", "KB", "KC", "KD", "KE", "KF", "KG", "KH", "KI", "KJ", "KK", "KL", "KM", "KN", "KO", "KP", "KQ", "KR",
		"KS", "KT", "KU", "KV", "KW", "KX", "KY", "KZ",
		"LA", "LB", "LC", "LD", "LE", "LF", "LG", "LH", "LI", "LJ", "LK", "LL", "LM", "LN", "LO", "LP", "LQ", "LR",
		"LS", "LT", "LU", "LV", "LW", "LX", "LY", "LZ",
	}
	// UNumberZero 零值
	UNumberZero   = UNumber{U: RAW, V: 0}
	unitStringMap map[string]Unit
)

// UNumber number with unit
type UNumber struct {
	U Unit    `json:"u"`
	V float32 `json:"v"`
}

// NewUNmberByInt 通过int值生成
func NewUNmberByInt(U Unit, V int) *UNumber {
	return &UNumber{
		V: float32(V),
		U: U,
	}
}

// NewUNmberByUint32 通过uint值生成
func NewUNmberByUint32(U Unit, V uint32) *UNumber {
	return &UNumber{
		V: float32(V),
		U: U,
	}
}

// Add 相加
func (n *UNumber) Add(un2 *UNumber) {
	new := Add(n, un2)
	if new != nil {
		n.V = new.V
		n.U = new.U
	}
}

// Sub 相减
func (n *UNumber) Sub(un2 *UNumber) {
	if n.Lt(un2) {
		log.Errorf("u1:(%v) less than u2:(%v)", n, un2)
		return
	}
	new := Sub(n, un2)
	if new != nil {
		n.V = new.V
		n.U = new.U
	}
}

// Mul 相乘
func (n *UNumber) Mul(un2 *UNumber) {
	new := Mul(n, un2)
	if new != nil {
		n.V = new.V
		n.U = new.U
	}
}

// Div 相除
func (n *UNumber) Div(un2 *UNumber) {
	new := Div(n, un2)
	if new != nil {
		n.V = new.V
		n.U = new.U
	}
}

// Equal 判断是否相等
func (n *UNumber) Equal(un2 *UNumber) bool {
	if un2 == nil {
		return false
	}
	n.normalize()
	un2.normalize()
	return absFloat32(n.V-un2.V) < 0.000001 && n.U == un2.U
}

// Ge 大于等于
func (n *UNumber) Ge(un2 *UNumber) bool {
	if un2 == nil {
		return false
	}
	// 异号
	if n.V > 0 && n.V < 0 {
		return true
	}
	// 异号
	if n.V < 0 && n.V > 0 {
		return false
	}

	// 以下代码处理两数同号
	n.normalize()
	un2.normalize()

	diff := n.U - un2.U
	if diff > 1 {
		if n.V >= 0 {
			return true
		}
		return false
	}
	if diff < -1 {
		if n.V >= 0 {
			return false
		}
		return true
	}
	value := un2.valueByUnit(n.U)

	return n.V >= value
}

// Gt 大于
func (n *UNumber) Gt(un2 *UNumber) bool {
	if un2 == nil {
		return false
	}
	// 异号
	if n.V > 0 && n.V < 0 {
		return true
	}
	// 异号
	if n.V < 0 && n.V > 0 {
		return false
	}

	// 以下代码处理两数同号
	n.normalize()
	un2.normalize()

	diff := n.U - un2.U
	if diff > 1 {
		if n.V >= 0 {
			return true
		}

		return false
	}

	if diff < -1 {
		if n.V >= 0 {
			return false
		}
		return true
	}

	value := un2.valueByUnit(n.U)

	return n.V > value
}

// Le 小于等于
func (n *UNumber) Le(un2 *UNumber) bool {
	return !n.Gt(un2)
}

// Lt 小于
func (n *UNumber) Lt(un2 *UNumber) bool {
	return !n.Ge(un2)
}

// RoundedUp 去掉小数
func (n *UNumber) RoundedUp() {
	n.V = float32(math.Ceil(float64(n.V)))
}

// isRaw 单位是否为0
func (n *UNumber) isRaw() bool {
	return n.U == RAW
}

// ToFloat32 unmber转float32
func (n *UNumber) ToFloat32() float32 {
	value := n.V
	unit := n.U
	for {
		if unit <= 0 || unit >= 4 {
			break
		}
		value = value * UnitOffest
		unit--
	}
	return value
}

func (n *UNumber) valueByUnit(unit Unit) float32 {
	offset := n.U - unit
	value := n.V
	if offset < 0 {
		offset = -offset
		for i := 0; i < int(offset); i++ {
			value = value / (UnitOffest)
		}
	} else {
		for i := 0; i < int(offset); i++ {
			value = value * (UnitOffest)
		}
	}
	return value
}

// normalize 格式化
func (n *UNumber) normalize() {
	for {
		if absFloat32(n.V) > UnitOffest && n.U < string2Unit("LZ") {
			n.V = n.V / UnitOffest
			n.U = n.U + 1
		} else {
			break
		}
	}
	for {
		if absFloat32(n.V) < 1 && n.U > string2Unit("raw") {
			n.V = n.V * UnitOffest
			n.U = n.U - 1
		} else {
			break
		}
	}
}

// fixValue value保留2位小数
func (n *UNumber) fixValue() {
	f64, err := strconv.ParseFloat(fmt.Sprintf("%.2f", n.V), 32)
	if err != nil {
		log.Errorf("fix:(%v) to float32 err:(%v)", n.V, err)
		return
	}
	n.V = (float32(f64))
}

// String UNumber转string
func (n *UNumber) String() string {
	return fmt.Sprintf("%0.2f%s", n.V, UnitString[n.U])
}

// ParseBigNumber parse  unumber object from bigNumber
func ParseBigNumber(bigNumber float64) UNumber {
	u := RAW
	for {
		if bigNumber < float64(math.MaxFloat32) {
			break
		}
		bigNumber = bigNumber / UnitOffest
		u++
	}
	un := UNumber{
		V: float32(bigNumber),
		U: u,
	}
	un.normalize()
	return un
}

// ParseUNumber parse unumber object from string
func ParseUNumber(text string) (UNumber, error) {
	un := UNumber{}
	// 306.5A
	bytes := []byte(text)
	idx := -1
	for i, b := range bytes {
		if (b >= 65 && b <= 90) || (b >= 97 && b <= 122) {
			idx = i
			break
		}
	}

	floatText := text
	unitText := ""
	if idx != -1 {
		floatText = text[0:idx]
		unitText = text[idx:]
	}

	f64, err := strconv.ParseFloat(floatText, 32)
	if err != nil {
		return un, err
	}

	un.V = float32(f64)
	un.U = string2Unit(unitText)

	return un, nil
}

func string2Unit(unitText string) Unit {
	unitText = strings.ToUpper(unitText)
	unit, ok := unitStringMap[unitText]
	if ok {
		return unit
	}

	return RAW
}

func absUnit(u Unit) Unit {
	if u < 0 {
		return 0 - u
	}
	return u
}

func absFloat32(v float32) float32 {
	if v < 0 {
		return 0 - v
	}
	return v
}

// Add 两个UNmber相加
func Add(un1 *UNumber, un2 *UNumber) *UNumber {
	if un1 == nil || un2 == nil {
		return nil
	}

	new := &UNumber{U: 0, V: 0}

	var adjust *UNumber

	if un1.U > un2.U {
		new.V = un1.V
		new.U = un1.U
		adjust = un2
	} else {
		new.V = un2.V
		new.U = un2.U
		adjust = un1
	}

	offest := absUnit(un1.U - un2.U)
	// 如果单位相差超过2级，也即是6个0，则忽略
	if offest > 3 {
		return new
	}

	div := 1
	for i := 0; i < int(offest); i++ {
		div *= UnitOffest
	}

	new.V = new.V + adjust.V/float32(div)
	new.normalize()

	return new
}

// Sub 两个UNmber相减
func Sub(un1 *UNumber, un2 *UNumber) *UNumber {
	if un1 == nil || un2 == nil {
		return nil
	}
	un22 := &UNumber{U: un2.U, V: -un2.V}
	newUnm := Add(un1, un22)
	return newUnm
}

// Mul 两个UNmber相乘
func Mul(un1 *UNumber, un2 *UNumber) *UNumber {
	if un1 == nil || un2 == nil {
		return nil
	}
	un22 := &UNumber{U: un1.U + un2.U, V: un1.V * un2.V}
	un22.normalize()
	return un22
}

// Div 两个UNmber相除
func Div(un1 *UNumber, un2 *UNumber) *UNumber {
	if un1 == nil || un2 == nil {
		return nil
	}
	uniX := un1.U - un2.U
	value := un1.V / un2.V

	if uniX < 0 {
		uniX = -uniX
		for i := 0; i < int(uniX); i++ {
			value = value / UnitOffest
		}
		uniX = 0
	}
	newum := &UNumber{U: uniX, V: value}
	newum.normalize()
	return newum
}

func init() {
	unitStringMap = make(map[string]Unit)

	for i, s := range UnitString {
		unitStringMap[s] = Unit(i)
	}
}
