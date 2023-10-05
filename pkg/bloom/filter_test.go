package bloom

import "testing"

func TestFNV64a(t *testing.T) {
	t.Run("Valid test", func(t *testing.T) {
		if hash := hashFNV64aFunc([]byte("hello")); hash == 0 {
			t.Errorf("Expected non-zero hash")
		}
	})

	t.Run("Deterministic test", func(t *testing.T) {
		if hash1, hash2 := hashFNV64aFunc([]byte("hello")), hashFNV64aFunc([]byte("hello")); hash1 != hash2 {
			t.Errorf("Hash should be deterministic")
		}
	})

	t.Run("Different input test", func(t *testing.T) {
		if hash1, hash2 := hashFNV64aFunc([]byte("hello")), hashFNV64aFunc([]byte("world")); hash1 == hash2 {
			t.Errorf("Different input should not produce the same hash")
		}
	})
}

func TestFNV32(t *testing.T) {
	t.Run("Valid test", func(t *testing.T) {
		if hash := hashFNV32Func([]byte("hello")); hash == 0 {
			t.Errorf("Expected non-zero hash")
		}
	})

	t.Run("Deterministic test", func(t *testing.T) {
		if hash1, hash2 := hashFNV32Func([]byte("hello")), hashFNV32Func([]byte("hello")); hash1 != hash2 {
			t.Errorf("Hash should be deterministic")
		}
	})

	t.Run("Different input test", func(t *testing.T) {
		if hash1, hash2 := hashFNV32Func([]byte("hello")), hashFNV32Func([]byte("world")); hash1 == hash2 {
			t.Errorf("Different input should not produce the same hash")
		}
	})
}

func TestCRC32(t *testing.T) {
	t.Run("Valid test", func(t *testing.T) {
		if hash := hashCRC32Func([]byte("hello")); hash == 0 {
			t.Errorf("Expected non-zero hash")
		}
	})

	t.Run("Deterministic test", func(t *testing.T) {
		if hash1, hash2 := hashCRC32Func([]byte("hello")), hashCRC32Func([]byte("hello")); hash1 != hash2 {
			t.Errorf("Hash should be deterministic")
		}
	})

	t.Run("Different input test", func(t *testing.T) {
		if hash1, hash2 := hashCRC32Func([]byte("hello")), hashCRC32Func([]byte("world")); hash1 == hash2 {
			t.Errorf("Different input should not produce the same hash")
		}
	})
}
