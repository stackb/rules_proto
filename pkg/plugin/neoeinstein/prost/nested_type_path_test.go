package prost

import "testing"

func TestProtoTypePathToRustPath(t *testing.T) {
	cases := []struct{ in, want string }{
		{"MyMessage", "MyMessage"},
		{"PIKInfo", "PikInfo"},
		// Nested: every non-final segment becomes snake_cased module name,
		// the leaf stays UpperCamel.
		{"UstOrderStatus.Confirmed", "ust_order_status::Confirmed"},
		{"Outer.Inner.Leaf", "outer::inner::Leaf"},
		{"PIKInfo.SubType", "pik_info::SubType"},
	}
	for _, c := range cases {
		if got := protoTypePathToRustPath(c.in); got != c.want {
			t.Errorf("protoTypePathToRustPath(%q): got %q, want %q", c.in, got, c.want)
		}
	}
}

func TestToSnakeFromCamel(t *testing.T) {
	cases := []struct{ in, want string }{
		{"MyType", "my_type"},
		{"PIKInfo", "pik_info"},
		{"URLLoader", "url_loader"},
		{"already_snake", "already_snake"},
		{"X", "x"},
	}
	for _, c := range cases {
		if got := toSnakeFromCamel(c.in); got != c.want {
			t.Errorf("toSnakeFromCamel(%q): got %q, want %q", c.in, got, c.want)
		}
	}
}
