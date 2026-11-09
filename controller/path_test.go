package controller

import "testing"

func TestPathParam(t *testing.T) {
	cases := []struct {
		path, prefix, want string
		ok                 bool
	}{
		{"/data/user/123", "/data/user/", "123", true},
		{"/data/user/abc-DEF_9", "/data/user/", "abc-DEF_9", true},
		{"/data/user/", "/data/user/", "", false},
		{"/data/user", "/data/user/", "", false},
		{"/data/user/123/", "/data/user/", "", false},
		{"/data/user/123/order", "/data/user/", "", false},
		{"/data/item/123", "/data/user/", "", false},
		{"/data/users/123", "/data/user/", "", false},
		{"", "/data/user/", "", false},
	}
	for _, c := range cases {
		got, ok := PathParam(c.path, c.prefix)
		if got != c.want || ok != c.ok {
			t.Errorf("PathParam(%q, %q) = %q, %v; want %q, %v", c.path, c.prefix, got, ok, c.want, c.ok)
		}
	}
}

func TestPathParamZeroAlloc(t *testing.T) {
	allocs := testing.AllocsPerRun(1000, func() {
		PathParam("/data/user/6543210fedcba9876543210f", "/data/user/")
	})
	if allocs != 0 {
		t.Errorf("PathParam mengalokasikan %v kali per panggilan, want 0", allocs)
	}
}

func BenchmarkPathParam(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		PathParam("/data/user/6543210fedcba9876543210f", "/data/user/")
	}
}
