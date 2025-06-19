package entrylist

import (
	"sync"
	"sync/atomic"
)

type EntryList struct {
	limit uint16
	data  sync.Map // ключ - string, значение - uint16
	state uint32   // 1 - доступен, 0 - недоступен
}

func NewEntryList(limit uint16) *EntryList {
	return &EntryList{
		limit: limit,
		state: 1,
	}
}

func (elist *EntryList) checkState() bool {
	return atomic.LoadUint32(&elist.state) == 1
}

func (elist *EntryList) IncrementFor(entry string) error {
	if !elist.checkState() {
		return ErrDownStated
	}

	newVal, err := elist.data.LoadOrStore(entry, uint16(1))
	if err {
		for {
			current := newVal.(uint16)
			if current >= elist.limit {
				return ErrPersonalOverflow
			}

			if elist.data.CompareAndSwap(entry, current, current+1) {
				return nil
			}

			newVal, _ = elist.data.Load(entry)
		}
	}

	return nil
}

func (elist *EntryList) Reset() error {
	if !atomic.CompareAndSwapUint32(&elist.state, 1, 0) {
		panic("Occasionally died, incorrect state set.")
	}

	elist.data = sync.Map{}

	if !atomic.CompareAndSwapUint32(&elist.state, 0, 1) {
		panic("Occasionally died, incorrect state set.")
	}
	return nil
}

func (elist *EntryList) GetCount(entry string) (uint16, error) {
	if !elist.checkState() {
		return 0, ErrDownStated
	}

	val, ok := elist.data.Load(entry)
	if !ok {
		return 0, nil
	}
	return val.(uint16), nil
}

func (elist *EntryList) IsLimitReached(entry string) (bool, error) {
	if !elist.checkState() {
		return false, ErrDownStated
	}

	val, ok := elist.data.Load(entry)
	if !ok {
		return false, nil
	}
	return val.(uint16) >= elist.limit, nil
}
