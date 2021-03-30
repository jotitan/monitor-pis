package agent

type Metric interface {
	//GetValue return value of metric and name
	GetValue()(float32,string,error)
}

