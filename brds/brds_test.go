package brds

import (
	"slices"
	"testing"
)

func TestAddHost(t *testing.T) {
	DbDef = 1 // use 1 for the test database
	rdb, ctx := InitClient()
	// t.Log(ctx, rdb)
	if _, err := AddHost(rdb, nil, ZabbixHost{}); err.Error() != "the input data is empty" {
		t.Fatalf("AddHost(rdb, ctx, ZabbixHost{})=%v, want `the input data is empty`", err)
	}

	if _, err := AddHost(rdb, nil, ZabbixHost{}); err.Error() != "the input data is empty" {
		t.Fatalf("ddHost(rdb, ctx, ZabbixHost{}))=%v, want `the input data is empty`", err)
	}

	if _, err := AddHost(rdb, ctx, ZabbixHost{}); err.Error() != "the input data is empty" {
		t.Fatalf("AddHost(rdb, ctx, ZabbixHost{}))=%v, want `the input data is empty`", err)
	}

	if _, err := AddHost(rdb, ctx, ZabbixHost{HostIdZ: "foo"}); err.Error() != "the input data is empty" {
		t.Fatalf("AddHost(rdb, ctx, ZabbixHost{HostIdZ: `foo`}))=%v, want `the input data is empty`", err)
	}

	if _, err := AddHost(rdb, ctx, ZabbixHost{ProblemZ: []string{"foo1", "foo2"}}); err.Error() != "the input data is empty" {
		t.Fatalf("AddHost(rdb, ctx,ZabbixHost{ZabbixHost{ProblemZ: []string{`foo1`, `foo2}}}))=%v, want `the input data is empty`", err)
	}

	if r, err := AddHost(rdb, ctx, ZabbixHost{HostIdZ: "foo", ProblemZ: []string{"foo1", "foo2"}}); err != nil || r != 7 {
		t.Fatalf("AddHost(rdb, ctx, ZabbixHost{ZabbixHost{ProblemZ: []string{`foo1`, `foo2}}}))=%v, %v , want 7, nil", r, err)
	}

	if r, err := AddHost(rdb, ctx, ZabbixHost{HostIdZ: "foo_1",
		HostZ:    "foo_1_23",
		NameZ:    "name_foo_1",
		StatusZ:  "status_foo_1",
		ItNew:    true,
		ProblemZ: []string{"pfoo_1_1", "pfoo_2_1"}}); err != nil || r != 7 {
		t.Fatalf("AddHost(rdb, ctx, ZabbixHost{%v)=%v, %v , want 7, nil", ZabbixHost{HostIdZ: "foo_1",
			HostZ:    "foo_1_23",
			NameZ:    "name_foo_1",
			StatusZ:  "status_foo_1",
			ItNew:    true,
			ProblemZ: []string{"pfoo_1_1", "pfoo_2_1"}},
			r, err)
	}

}

func TestGetHost(t *testing.T) {
	DbDef = 1 // use 1 for the test database
	rdb, ctx := InitClient()
	want := ZabbixHost{HostIdZ: "foo_1",
		HostZ:    "foo_1_23",
		NameZ:    "name_foo_1",
		StatusZ:  "status_foo_1",
		ItNew:    true,
		ProblemZ: []string{"pfoo_1_1", "pfoo_2_1"}}
	res, err := GetHost(rdb, ctx, "foo_1")
	if err != nil {
		t.Fatalf("Get err = %v, want err = nil", err)
	}
	if want.HostIdZ != res.HostIdZ ||
		want.HostZ != res.HostZ ||
		want.NameZ != res.NameZ ||
		want.StatusZ != res.StatusZ ||
		want.ItChanged != res.ItChanged ||
		want.ItNew != res.ItNew ||
		slices.Compare(want.ProblemZ, res.ProblemZ) != 0 {
		t.Fatalf("Get %v, want %v", res, want)
	}
}

func TestDelHost(t *testing.T) {
	DbDef = 1 // use 1 for the test database
	rdb, ctx := InitClient()

	if res, err := DelHost(rdb, ctx, "foo_1"); err != nil || res != 1 {
		t.Fatalf("res, err= %v, %v, want 1, nil", res, err)
	}

	if res, err := DelHost(rdb, ctx, "foo"); err != nil || res != 1 {
		t.Fatalf("res, err= %v, %v, want 1, nil", res, err)
	}

}
