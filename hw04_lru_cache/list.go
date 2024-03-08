package hw04lrucache

import (
	"errors"
	"fmt"
)

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
	front  *ListItem
	back   *ListItem
	length int
}

func NewList() *list {
	return new(list)
}

func (l *list) Len() int {
	return l.length
}

func (l *list) Back() *ListItem {
	return l.back
}

func (l *list) Front() *ListItem {
	return l.front
}

func (l *list) PushFront(item interface{}) *ListItem {
	newItem := &ListItem{
		Value: item,
		Next:  l.front,
	}

	if l.front == nil {
		l.front = newItem
		l.back = l.front
	} else {
		l.front.Prev = newItem
		l.front = newItem
	}

	l.length++
	return newItem
}

func (l *list) PushBack(item interface{}) *ListItem {
	newItem := &ListItem{
		Value: item,
		Prev:  l.back,
	}
	if l.back == nil {
		l.back = newItem
		l.front = l.back
	} else {
		l.back.Next = newItem
		l.back = newItem
	}

	l.length++
	return newItem
}

// потенциально опасная операция. в item может быть передан элемент от другого списка
// но иначе - не будет О(1)
func (l *list) Remove(item *ListItem) {
	if item == nil {
		return
	}

	switch item {
	case l.front:
		l.front = item.Next
		item.Next.Prev = nil
		item.Next = nil
	case l.back:
		l.back = item.Prev
		item.Prev.Next = nil
		item.Prev = nil
	default:
		item.Prev.Next = item.Next
		item.Next.Prev = item.Prev
	}

	l.length--
}

// та же ситуация, что и с Remove()
func (l *list) MoveToFront(item *ListItem) {
	if item == nil {
		return
	}
	l.Remove(item) // удалили элемент length уменьшилась на 1
	l.front.Prev = item
	item.Next = l.front
	l.front = item
	l.length++
}

// ищем элемент с указанным значением начиная от указанного
func (l *list) SearchNext(startItem *ListItem, v interface{}) (*ListItem, error) {
	if !l.checkItem(startItem) {
		return nil, errors.New("данный элемент не принадлежит текущему списку")
	}

	for i := startItem; i != nil; i = i.Next {
		if i.Value == v {
			return i, nil
		}
	}

	return nil, nil
}

// возвращаем первый найденный по содержимому элемент списка
func (l *list) SearchFirst(v interface{}) *ListItem {
	item, _ := l.SearchNext(l.front, v)
	return item
}

// проверяем принадлежит ли указанный элемент списка текущему листу
func (l *list) checkItem(item *ListItem) bool {
	if item == nil {
		return false
	}

	for i := l.Front(); i != nil; i = i.Next {
		if i == item {
			return true
		}
	}
	return false
}

func (l list) String() string {
	res := ""
	for i := l.Front(); i != nil; i = i.Next {
		if i == l.front {
			res = fmt.Sprintf("%v", i.Value)
		} else {
			res = res + fmt.Sprintf(", %v", i.Value)
		}
	}
	return res
}
