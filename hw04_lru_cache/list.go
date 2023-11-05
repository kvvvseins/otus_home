package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	FirstItem *ListItem
	LastItem  *ListItem
	count     int
}

func (l *list) Len() int {
	return l.count
}

func (l *list) Front() *ListItem {
	if l.Len() == 0 {
		return nil
	}

	return l.FirstItem
}

func (l *list) Back() *ListItem {
	if l.Len() == 0 {
		return nil
	}

	return l.LastItem
}

func (l *list) PushFront(v interface{}) *ListItem {
	defer func() {
		l.count++
	}()

	currentListItem := &ListItem{
		Value: v,
	}

	if l.Len() == 0 {
		l.FirstItem = currentListItem
		l.LastItem = currentListItem

		return currentListItem
	}

	currentListItem.Next = l.FirstItem

	l.FirstItem.Prev = currentListItem
	l.FirstItem = currentListItem

	if l.LastItem.Prev == nil {
		l.LastItem.Prev = l.FirstItem
	}

	return currentListItem
}

func (l *list) PushBack(v interface{}) *ListItem {
	if l.Len() == 0 {
		return l.PushFront(v)
	}

	defer func() {
		l.count++
	}()

	currentListItem := &ListItem{
		Value: v,
		Next:  nil,
		Prev:  l.LastItem,
	}

	l.LastItem.Next = currentListItem
	l.LastItem = currentListItem

	if l.FirstItem.Next == nil {
		l.FirstItem.Next = l.LastItem
	}

	return currentListItem
}

func (l *list) Remove(i *ListItem) {
	if l.Len() == 0 {
		return
	}

	defer func() {
		l.count--
	}()

	if i.Next == nil && i.Prev == nil {
		l.FirstItem = nil
		l.LastItem = nil

		return
	}

	if i.Prev == nil {
		l.FirstItem = i.Next
	} else {
		i.Prev.Next = i.Next
	}

	if i.Next == nil {
		l.LastItem = i.Prev
	} else {
		i.Next.Prev = i.Prev
	}

	i.Next = nil
	i.Prev = nil
}

func (l *list) MoveToFront(i *ListItem) {
	if l.Len() == 0 {
		l.PushFront(i)

		return
	}

	if l.Len() == 1 || i.Prev == nil {
		return
	}

	i.Prev.Next = i.Next

	if i.Next == nil {
		l.LastItem = i.Prev
	} else {
		i.Next.Prev = i.Prev
	}

	i.Next = l.FirstItem
	i.Prev = nil

	l.FirstItem.Prev = i
	l.FirstItem = i
}

func NewList() List {
	return new(list)
}
