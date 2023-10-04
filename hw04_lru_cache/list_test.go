package hw04lrucache

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestList(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		l := NewList()

		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})

	t.Run("complex", func(t *testing.T) {
		l := NewList()

		l.PushFront(10) // [10]
		l.PushBack(20)  // [10, 20]
		l.PushBack(30)  // [10, 20, 30]
		require.Equal(t, 3, l.Len())

		middle := l.Front().Next // 20
		l.Remove(middle)         // [10, 30]
		require.Equal(t, 2, l.Len())

		for i, v := range [...]int{40, 50, 60, 70, 80} {
			if i%2 == 0 {
				l.PushFront(v)
			} else {
				l.PushBack(v)
			}
		} // [80, 60, 40, 10, 30, 50, 70]

		require.Equal(t, 7, l.Len())
		require.Equal(t, 80, l.Front().Value)
		require.Equal(t, 70, l.Back().Value)

		l.MoveToFront(l.Front()) // [80, 60, 40, 10, 30, 50, 70]
		l.MoveToFront(l.Back())  // [70, 80, 60, 40, 10, 30, 50]

		elems := make([]int, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			elems = append(elems, i.Value.(int))
		}
		require.Equal(t, []int{70, 80, 60, 40, 10, 30, 50}, elems)
	})

	t.Run("removeTest", func(t *testing.T) {
		l := NewList()

		// пробуем удалить из пустого списка
		l.Remove(&ListItem{Value: 11})
		require.Equal(t, 0, l.Len())

		l.PushFront(11)

		// пробуем удалить из списка одного элемента
		l.Remove(&ListItem{Value: 11})
		require.Equal(t, 0, l.Len())

		l.PushFront(10)
		l.PushBack(20)
		l.PushBack(35)
		l.PushBack(30)

		// пробуем удалить первый элемент списка
		l.Remove(l.Front())
		require.Equal(t, 3, l.Len())
		require.Equal(t, 20, l.Front().Value)
		require.Equal(t, 30, l.Back().Value)

		l.PushFront(10)

		// пробуем удалить последний элемент списка
		l.Remove(l.Back())
		require.Equal(t, 3, l.Len())
		require.Equal(t, 10, l.Front().Value)
		require.Equal(t, 35, l.Back().Value)
	})

	t.Run("MoveToFront test", func(t *testing.T) {
		l := NewList()

		item0 := l.PushFront(0)
		item1 := l.PushFront(1)
		item2 := l.PushFront(2)
		item3 := l.PushFront(3)
		item4 := l.PushFront(4)
		item5 := l.PushFront(5)

		l.MoveToFront(item2)

		require.Equal(t, 2, l.Front().Value)
		require.Equal(t, (*ListItem)(nil), l.Front().Prev)

		require.Equal(t, 0, l.Back().Value)
		require.Equal(t, (*ListItem)(nil), l.Back().Next)

		require.Equal(t, item5.Value, item2.Next.Value)
		require.Equal(t, (*ListItem)(nil), item2.Prev)

		require.Equal(t, item4.Value, item5.Next.Value)
		require.Equal(t, item2.Value, item5.Prev.Value)

		require.Equal(t, item3.Value, item4.Next.Value)
		require.Equal(t, item5.Value, item4.Prev.Value)

		require.Equal(t, item1.Value, item3.Next.Value)
		require.Equal(t, item4.Value, item3.Prev.Value)

		require.Equal(t, item0.Value, item1.Next.Value)
		require.Equal(t, item3.Value, item1.Prev.Value)

		require.Equal(t, (*ListItem)(nil), item0.Next)
		require.Equal(t, item1.Value, item0.Prev.Value)
	})
}
