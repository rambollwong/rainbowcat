package task

type DataStore interface {
	AddData(data Data)
	GetData(dataType Type) Data
	RemoveData(dataId uint64)
	ExistData(dataType Type) bool
}
