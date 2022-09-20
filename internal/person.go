package internal

type Person struct {
	DisplayName *string
	callRecord  *CallRecord
	callRecords []*CallRecord
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

func (u *Person) SetCallRecords(records []*CallRecord) {
	u.callRecords = records
}

func (u *Person) GetCallRecords() []*CallRecord {
	return u.callRecords
}
