package pq_test

import (
	"sort"
	"testing"

	lib "github.com/quantime/go-priorityqueue/lib/pq"
)

type TestItem struct {
	val    string
	prio   float64
	second float64
}

var counter int

func (t TestItem) Clean() {
	counter++
}

func TestPriorityQueue(t *testing.T) {
	pq := lib.New()
	elements := []TestItem{
		TestItem{"A", 5, 0},
		TestItem{"B", 3, 0},
		TestItem{"C", 7, 0},
		TestItem{"D", 8, 0},
		TestItem{"E", 6, 0},
		TestItem{"F", 2, 0},
		TestItem{"G", 9, 0},
	}

	for _, e := range elements {
		pq.Insert(e, e.prio, e.second)
	}

	sort.SliceStable(elements, func(i, j int) bool {
		return elements[i].prio < elements[j].prio
	})

	for _, e := range elements {
		item, err := pq.Pop()
		if err != nil {
			t.Fatalf(err.Error())
		}

		i := item.(TestItem)
		if e != i {
			t.Fatalf("expected %v, got %v", e, i)
		}
	}
}

func TestPriorityQueueUpdate(t *testing.T) {
	a := TestItem{"A", 3, 0}
	b := TestItem{"B", 4, 0}

	pq := lib.New()
	pq.Insert(a, a.prio, a.second)
	pq.Insert(b, b.prio, b.second)
	pq.UpdatePriority(b, 2)

	item, err := pq.Pop()
	if err != nil {
		t.Fatal(err.Error())
	}

	if item.(TestItem) != b {
		t.Fatal("priority update failed")
	}
}

func TestPriorityQueueLen(t *testing.T) {
	a := TestItem{"A", 3, 0}
	b := TestItem{"B", 4, 0}

	pq := lib.New()
	if pq.Len() != 0 {
		t.Fatal("empty queue should have length of 0")
	}

	pq.Insert(a, a.prio, a.second)
	pq.Insert(b, b.prio, b.second)
	if pq.Len() != 2 {
		t.Fatal("queue should have lenght of 2 after 2 inserts")
	}
}

func TestDoubleAddition(t *testing.T) {
	a := TestItem{"A", 3, 0}
	b := TestItem{"B", 4, 0}

	pq := lib.New()
	pq.Insert(a, 2, a.second)
	pq.Insert(b, 3, b.second)
	pq.Insert(b, 1, b.second)

	if pq.Len() != 2 {
		t.Fatal("queue should ignore inserting the same element twice")
	}

	item, _ := pq.Pop()
	if item.(TestItem) != a {
		t.Fatal("queue should ignore duplicate insert, not update existing item")
	}
}

func TestPopEmptyQueue(t *testing.T) {
	pq := lib.New()
	_, err := pq.Pop()
	if err == nil {
		t.Fatal("should produce error when performing pop on empty queue")
	}
}

func TestUpdateNonExistingItem(t *testing.T) {
	a := TestItem{"A", 3, 0}
	b := TestItem{"B", 4, 0}

	pq := lib.New()

	pq.Insert(a, 4, a.second)
	pq.UpdatePriority(b, 5)

	if pq.Len() != 1 {
		t.Fatal("update should not add items")
	}

	item, _ := pq.Pop()
	if item.(TestItem) != a {
		t.Fatalf("update should not overwrite item, expected \"foo\", got \"%v\"", item.(TestItem))
	}
}

func TestMaxSize(t *testing.T) {
	a := TestItem{"A", 3, 0}
	b := TestItem{"B", 4, 0}
	c := TestItem{"C", 5, 0}
	d := TestItem{"D", 6, 0}

	pq := lib.NewWithMaxSize(3)
	pq.Insert(a, 1, a.second)
	pq.Insert(b, 2, b.second)
	pq.Insert(c, 3, c.second)
	pq.Insert(d, 4, d.second)

	if pq.Len() != 3 {
		t.Fatal("max size violated")
	}
}

func TestMaxSizeClean(t *testing.T) {
	counter = 0
	a := TestItem{"A", 3, 0}
	b := TestItem{"B", 4, 0}

	pq := lib.NewWithMaxSize(1)
	pq.Insert(a, 1, a.second)
	pq.Insert(b, 2, b.second)

	if counter == 0 {
		t.Fatalf("clean should be called when overflow item is popped")
	}
}

func TestSecondaryPrio(t *testing.T) {
	b := TestItem{"B", 4, 3}
	c := TestItem{"C", 4, 2}

	pq := lib.New()

	pq.Insert(b, b.prio, b.second)
	pq.Insert(c, c.prio, c.second)

	item, _ := pq.Pop()
	if item.(TestItem) != c {
		t.Fatal("queue should pop item with lower secondary prio")
	}
}
