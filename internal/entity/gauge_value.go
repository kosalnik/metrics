package entity

type GaugeValue struct {
	Name  string
	Value float64
}

//func NewGaugeValue(name string, value float64) (*GaugeValue, error) {
//	//v, err := strconv.ParseFloat(value, 64)
//	//if err != nil {
//	//	return nil, err
//	//}
//	return &GaugeValue{name, value}, nil
//}

func (v *GaugeValue) Set(val float64) {
	v.Value = val
}

func (v *GaugeValue) Inc(val float64) {
	v.Value += val
}
