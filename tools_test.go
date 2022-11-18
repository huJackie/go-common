package pkg

import "testing"

func TestExpireTime(t *testing.T) {
	t.Log(ExpireTime(0).Seconds())
}

func TestNow(t *testing.T) {
	t.Log(Now())
}

func TestHashCode(t *testing.T) {
	var tables = []struct {
		raw  string
		md5  string
		want uint32
	}{
		{
			"fff62b4a2b4a57bf",
			"dbe8f0d666a5278593469833ec9cf812",
			4,
		},
	}

	for i, v := range tables {
		m := Md5(v.raw)
		if m != v.md5 {
			t.Errorf("md5 not equal want [%s] real [%s]", v.raw, m)
		}

		if d := HashCode(m) % 10; d != v.want {
			t.Errorf("No1.%d want %d real %d", i, v.want, d)
		}

		if d := HashCode(v.md5) % 10; d != v.want {
			t.Errorf("No2.%d want %d real %d", i, v.want, d)
		}
	}

}
