package sets

import "testing"

func TestBasicStuff(t *testing.T) {
	s := NewSet([]string{"asd", "bds", "uvw", "asd"})
	t.Log(s)

	s = Remove(s, "bss")
	s = Remove(s, "bds")
	t.Log(s)

	t.Log(Exists(s, "asd"))
	t.Log(Exists(s, "bds"))

	//s := NewSet()
}
