package union

import "testing"

func TestUnion(t *testing.T) {
	u := NewUnion([]string{"alice", "bob", "charlie"})
	t.Log(Defined(u))

	err := Set(&u, "alice")
	t.Log(err)
	t.Log(u)

	v, exists := Get(u)
	t.Log(*v, exists)

}
