package brds

import "testing"

func TestAddHost(t *testing.T) {
	DbDef = 1 // use 1 for the test database
	ctx, rdb := InitClient()
	// t.Log(ctx, rdb)
	if _, err := AddHost(nil, rdb, ZabbixHost{}); err.Error() != "the input data is nil" {
		t.Fatalf("AddHost(nil, rdb, ZabbixHost{})=%v, want `the input data is nil`", err)
	}

	if _, err := AddHost(ctx, nil, ZabbixHost{}); err.Error() != "the input data is nil" {
		t.Fatalf("ddHost(ctx, nil, ZabbixHost{}))=%v, want `the input data is nil`", err)
	}

	if _, err := AddHost(ctx, rdb, ZabbixHost{}); err.Error() != "the input data is nil" {
		t.Fatalf("AddHost(ctx, rdb, ZabbixHost{}))=%v, want `the input data is nil`", err)
	}

	if _, err := AddHost(ctx, rdb, ZabbixHost{HostIdZ: "foo"}); err.Error() != "the input data is nil" {
		t.Fatalf("AddHost(ctx, rdb, ZabbixHost{HostIdZ: `foo`}))=%v, want `the input data is nil`", err)
	}

	if _, err := AddHost(ctx, rdb, ZabbixHost{ProblemZ: []string{"foo1", "foo2"}}); err.Error() != "the input data is nil" {
		t.Fatalf("AddHost(ctx, rdb, ZabbixHost{ZabbixHost{ProblemZ: []string{`foo1`, `foo2}}}))=%v, want `the input data is nil`", err)
	}

	if _, err := AddHost(ctx, rdb, ZabbixHost{HostIdZ: "foo", ProblemZ: []string{"foo1", "foo2"}}); err != nil {
		t.Fatalf("AddHost(ctx, rdb, ZabbixHost{ZabbixHost{ProblemZ: []string{`foo1`, `foo2}}}))=%v, want nil", err)
	}

}
