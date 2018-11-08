package filefaker

import (
	"testing"
)

func TestObserve(t *testing.T) {
	expected := "are you sane?"
	w := New()
	w.Write([]byte(expected))

	if actual := w.Observe(); actual != expected {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}

func TestBytesWritten(t *testing.T) {
	expected := 5
	w := New()
	bw, _ := w.Write([]byte("xxxxx"))

	if bw != expected {
		t.Errorf("Expected %d, got %d", expected, bw)
	}

	if w.N != expected {
		t.Errorf("Expected %d, got %d", expected, w.N)
	}
}

func TestClose(t *testing.T) {
	w := New()
	if w.Closed != false {
		t.Errorf("Expected %t, got %t", false, w.Closed)
	}

	w.Close()
	if w.Closed != true {
		t.Errorf("Expected %t, got %t", true, w.Closed)
	}
}
