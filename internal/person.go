package internal

type Person struct {
	displayName    *string
	callRecord     *CallRecord
	callRecords    []*CallRecord
	status         *PersonStatus
	previousStatus []*PersonStatus
	cardNumbers    []int
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
func (u *Person) SetDisplayName(name *string) {
	u.displayName = name
}

func (u *Person) GetDisplayName() *string {
	return u.displayName
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

func (u *Person) SetStatus(personStatus *PersonStatus) {
	u.status = personStatus
}

func (u *Person) GetStatus() *PersonStatus {
	return u.status
}

func (u *Person) SetPreviousStatus(previousStatus []*PersonStatus) {
	u.previousStatus = previousStatus
}

func (u *Person) GetPreviousStatus() []*PersonStatus {
	return u.previousStatus
}

func (u *Person) SetCardNumbers(numbers []int) {
	u.cardNumbers = numbers
}

func (u *Person) GetCardNumbers() []int {
	return u.cardNumbers
}
