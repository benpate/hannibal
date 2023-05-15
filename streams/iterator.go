package streams

type Iterator interface {
	Count() int
	HasNext() bool
	GetNext() Document
}
