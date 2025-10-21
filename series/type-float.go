package series

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

type floatElement struct {
	e         float64
	precision uint8
	nan       bool
}

// force floatElement struct to implement Element interface
var _ Element = (*floatElement)(nil)

func (e *floatElement) Set(value interface{}) {
	e.nan = false
	switch val := value.(type) {
	case string:
		if val == "NaN" {
			e.nan = true
			return
		}
		v := strings.TrimSpace(value.(string))
		v = strings.ReplaceAll(v, ",", "")
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			e.nan = true
			return
		}
		e.e = f
		e.precision = uint8(len(v) - strings.LastIndex(v, ".") - 1)
	case int:
		e.e = float64(val)
	case float64:
		e.e = float64(val)
		e.precision = 255
	case bool:
		b := val
		if b {
			e.e = 1
		} else {
			e.e = 0
		}
	case Element:
		if fe, ok := value.(*floatElement); ok {
			e.e = fe.e
			e.precision = fe.precision
			e.nan = fe.nan
			return
		}
		e.e = val.Float()
	default:
		e.nan = true
		return
	}
}

func (e floatElement) Copy() Element {
	return &floatElement{e.e, e.precision, e.nan}
}

func (e floatElement) IsNA() bool {
	if e.nan || math.IsNaN(e.e) {
		return true
	}
	return false
}

func (e floatElement) Type() Type {
	return Float
}

func (e floatElement) Val() ElementValue {
	if e.IsNA() {
		return nil
	}
	return convertToDecimalPlaces(e.e, e.precision)
}

func (e floatElement) String() string {
	if e.IsNA() {
		return "NaN"
	}
	return fmt.Sprintf("%f", e.e)
}

func (e floatElement) Int() (int, error) {
	if e.IsNA() {
		return 0, fmt.Errorf("can't convert NaN to int")
	}
	f := e.e
	if math.IsInf(f, 1) || math.IsInf(f, -1) {
		return 0, fmt.Errorf("can't convert Inf to int")
	}
	if math.IsNaN(f) {
		return 0, fmt.Errorf("can't convert NaN to int")
	}
	return int(f), nil
}

func (e floatElement) Float() float64 {
	if e.IsNA() {
		return math.NaN()
	}
	return float64(e.e)
}

func (e floatElement) Bool() (bool, error) {
	if e.IsNA() {
		return false, fmt.Errorf("can't convert NaN to bool")
	}
	switch e.e {
	case 1:
		return true, nil
	case 0:
		return false, nil
	}
	return false, fmt.Errorf("can't convert Float \"%v\" to bool", e.e)
}

func (e floatElement) IntElement() *intElement {
	el := &intElement{}
	el.Set(e.Val())
	return el
}

func (e floatElement) FloatElement() *floatElement {
	return &e
}

func (e floatElement) StringElement() *stringElement {
	el := &stringElement{}
	el.Set(e.Val())
	return el
}

func (e floatElement) BoolElement() *boolElement {
	el := &boolElement{}
	el.Set(e.Val())
	return el
}

func (e floatElement) Eq(elem Element) bool {
	if e.IsNA() || elem.IsNA() {
		return false
	}

	switch fe := elem.(type) {
	case *floatElement:
		return convertToDecimalPlaces(e.e, e.precision) == convertToDecimalPlaces(fe.e, fe.precision)
	default:
		f := fe.Float()
		return e.e == f
	}
}

func (e floatElement) Neq(elem Element) bool {
	f := elem.Float()
	if e.IsNA() || math.IsNaN(f) {
		return false
	}
	return e.e != f
}

func (e floatElement) Less(elem Element) bool {
	f := elem.Float()
	if e.IsNA() || math.IsNaN(f) {
		return false
	}
	return e.e < f
}

func (e floatElement) LessEq(elem Element) bool {
	f := elem.Float()
	if e.IsNA() || math.IsNaN(f) {
		return false
	}
	return e.e <= f
}

func (e floatElement) Greater(elem Element) bool {
	f := elem.Float()
	if e.IsNA() || math.IsNaN(f) {
		return false
	}
	return e.e > f
}

func (e floatElement) GreaterEq(elem Element) bool {
	f := elem.Float()
	if e.IsNA() || math.IsNaN(f) {
		return false
	}
	return e.e >= f
}

// WithPrecision creates a new floatElement with the specified precision
func (e floatElement) WithPrecision(precision uint8) *floatElement {
	if e.precision == precision {
		return &e
	}

	if e.IsNA() {
		return &floatElement{nan: true}
	}
	return &floatElement{e.e, precision, false}
}

func convertToDecimalPlaces(f float64, precision uint8) float64 {
	p := math.Pow(10, float64(precision))
	a := math.Round(f*p) / p
	return a
}
