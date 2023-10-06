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

func TestCalcFilterParams(t *testing.T) {
	t.Run("Valid with K 3", func(t *testing.T) {
		m, k := CalcFilterParams(5000, 0.1)
		expectedM := uint64(24000)
		expectedK := 3

		if m != expectedM || k != expectedK {
			t.Errorf("got (%d, %d), expected (%d, %d)", m, k, expectedM, expectedK)
		}
	})
	t.Run("Valid with K 7", func(t *testing.T) {
		m, k := CalcFilterParams(1000, 0.01)
		expectedM := uint64(9600)
		expectedK := 7

		if m != expectedM || k != expectedK {
			t.Errorf("got (%d, %d), expected (%d, %d)", m, k, expectedM, expectedK)
		}
	})

	t.Run("Zero items", func(t *testing.T) {
		m, k := CalcFilterParams(0, 0.01)
		if m != 0 || k != 0 {
			t.Errorf("got (%d, %d), expected (0, 0)", m, k)
		}
	})

	t.Run("Invalid probability low", func(t *testing.T) {
		m, k := CalcFilterParams(1000, -0.01)
		if m != 0 || k != 0 {
			t.Errorf("got (%d, %d), expected (0, 0)", m, k)
		}
	})

	t.Run("Invalid probability high", func(t *testing.T) {
		m, k := CalcFilterParams(1000, 1.5)
		if m != 0 || k != 0 {
			t.Errorf("got (%d, %d), expected (0, 0)", m, k)
		}
	})
}
