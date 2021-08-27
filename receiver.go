package gocoder

type Receiver interface {
	Type
	GetName() string
	GetType() Type
	IsReceiver()
	GetValue() Value
}

var _ Receiver = (*tReceiver)(nil)

type tReceiver struct {
	Type
	ReName string
}

func (t *tReceiver) WriteCode(w Writer) {
	w.WriteCode(t)
}

func (t *tReceiver) GetName() string {
	return t.ReName
}

func (t *tReceiver) GetType() Type {
	return t.Type
}

func (t *tReceiver) GetValue() Value {
	return NewValue(t.GetName(), t.GetType())
}

func (t *tReceiver) IsReceiver() {
}
