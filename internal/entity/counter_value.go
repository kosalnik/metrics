package entity

type CounterValue struct {
	Name  string
	Value int64
}

//func NewCounterValue(name, value string) (*CounterValue, error) {
//	//v, err := strconv.ParseInt(value, 10, 64)
//	//if err != nil {
//	//	return nil, err
//	//}
//	return &CounterValue{name, value}, nil
//}

func (v *CounterValue) Set(val int64) {
	v.Value = val
}

func (v *CounterValue) Inc(val int64) {
	v.Value += val
}
