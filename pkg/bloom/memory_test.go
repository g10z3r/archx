package bloom

import (
	"crypto/rand"
	"testing"
)

func TestPutAndMightContain(t *testing.T) {
	filter := NewBloomFilter(1000)

	t.Run("Put valid item", func(t *testing.T) {
		err := filter.Put([]byte("valid_item"))
		if err != nil {
			t.Fatalf("Put failed: %v", err)
		}
	})

	t.Run("MightContain with existing item", func(t *testing.T) {
		contains, err := filter.MightContain([]byte("valid_item"))
		if err != nil {
			t.Fatalf("MightContain failed: %v", err)
		}
		if !contains {
			t.Fatalf("Item should be in the filter")
		}
	})

	t.Run("MightContain with non-existing item", func(t *testing.T) {
		contains, err := filter.MightContain([]byte("non_existent_item"))
		if err != nil {
			t.Fatalf("MightContain failed: %v", err)
		}
		if contains {
			t.Fatalf("Item should not be in the filter")
		}
	})

	t.Run("MightContain with empty item", func(t *testing.T) {
		contains, err := filter.MightContain([]byte(""))
		if err != nil {
			t.Fatalf("MightContain failed: %v", err)
		}
		if contains {
			t.Fatalf("Empty item should not be in the filter by default")
		}
	})

	t.Run("MightContain with nil item", func(t *testing.T) {
		contains, err := filter.MightContain(nil)
		if err != nil {
			t.Fatalf("MightContain failed: %v", err)
		}
		if contains {
			t.Fatalf("Nil item should not be in the filter by default")
		}
	})
}

func TestConcurrentAccess(t *testing.T) {
	filter := NewBloomFilter(10000)
	n := 1000

	t.Run("Put valid item", func(t *testing.T) {
		c := make(chan bool, n)

		for i := 0; i < n; i++ {
			go func() {
				b := make([]byte, 10)
				rand.Read(b)
				err := filter.Put(b)
				if err != nil {
					t.Errorf("Put failed: %v", err)
				}
				c <- true
			}()
		}

		for i := 0; i < n; i++ {
			<-c
		}
	})

	t.Run("Concurrent might contain valid", func(t *testing.T) {
		c := make(chan bool, n)

		for i := 0; i < n; i++ {
			go func() {
				b := make([]byte, 10)
				rand.Read(b)
				filter.Put(b)
				contains, err := filter.MightContain(b)
				if err != nil {
					t.Errorf("MightContain failed: %v", err)
				}
				if !contains {
					t.Errorf("Item should be in filter")
				}
				c <- true
			}()
		}

		for i := 0; i < n; i++ {
			<-c
		}
	})

	t.Run("Concurrent might contain invalid", func(t *testing.T) {
		c := make(chan bool, n)

		for i := 0; i < n; i++ {
			go func() {
				b := make([]byte, 10)
				rand.Read(b)
				contains, err := filter.MightContain(b)
				if err != nil {
					t.Errorf("MightContain failed: %v", err)
				}
				if contains {
					t.Logf("False positive: Item was detected in filter but should not be")
				}
				c <- true
			}()
		}

		for i := 0; i < n; i++ {
			<-c
		}
	})
}
