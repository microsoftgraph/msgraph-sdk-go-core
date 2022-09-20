package internal

type Person struct {
	DisplayName *string
	callRecord  *CallRecord
}

func NewPerson() *Person {
	return &Person{}
}

type CallRecord struct {
	Entity
}

func NewCallRecord() *CallRecord {
	return &CallRecord{}
}

func (u *Person) SetCallRecord(record *CallRecord) {
	u.callRecord = record
}

func (u *Person) GetCallRecord() *CallRecord {
	return u.callRecord
}
