package brds

import (
	"errors"
	"slices"
	"testing"
)

func TestAddHost(t *testing.T) {

	DbDef = 1 // use 1 for the test database
	AddrDef = "localhost:6380"
	rdb, ctx := InitClient()
	defer rdb.Close()

	t.Log("Cleanin up!")
	DelAllHosts(rdb, ctx)
	// t.Log(ctx, rdb)

	if _, err := AddHost(rdb, nil, ZabbixHost{}); !errors.Is(err, ErrEmptyInputData) {
		t.Fatalf("AddHost(rdb, ctx, ZabbixHost{})=%v, want %v", err, ErrEmptyInputData)
	}

	if _, err := AddHost(rdb, nil, ZabbixHost{}); !errors.Is(err, ErrEmptyInputData) {
		t.Fatalf("ddHost(rdb, ctx, ZabbixHost{}))=%v, want %v", err, ErrEmptyInputData)
	}

	if _, err := AddHost(rdb, ctx, ZabbixHost{}); !errors.Is(err, ErrEmptyInputData) {
		t.Fatalf("AddHost(rdb, ctx, ZabbixHost{}))=%v, want %v", err, ErrEmptyInputData)
	}

	if _, err := AddHost(rdb, ctx, ZabbixHost{HostIdZ: "foo"}); !errors.Is(err, ErrEmptyInputData) {
		t.Fatalf("AddHost(rdb, ctx, ZabbixHost{HostIdZ: `foo`}))=%v, want %v", err, ErrEmptyInputData)
	}

	if _, err := AddHost(rdb, ctx, ZabbixHost{ProblemZ: []string{"foo1", "foo2"}}); !errors.Is(err, ErrEmptyInputData) {
		t.Fatalf("AddHost(rdb, ctx,ZabbixHost{ZabbixHost{ProblemZ: []string{`foo1`, `foo2}}}))=%v, want %v", err, ErrEmptyInputData)
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
	AddrDef = "localhost:6380"
	rdb, ctx := InitClient()
	defer rdb.Close()

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
	if !compareHosts(want, res) {
		t.Fatalf("Get %v, want %v", res, want)
	}

}

func TestGetAllHosts(t *testing.T) {
	DbDef = 1 // use 1 for the test database
	AddrDef = "localhost:6380"
	rdb, ctx := InitClient()
	defer rdb.Close()

	want := map[string]ZabbixHost{"foo_1": {HostIdZ: "foo_1",
		HostZ:    "foo_1_23",
		NameZ:    "name_foo_1",
		StatusZ:  "status_foo_1",
		ItNew:    true,
		ProblemZ: []string{"pfoo_1_1", "pfoo_2_1"}},
		"foo": {HostIdZ: "foo",
			ProblemZ: []string{"foo1", "foo2"}}}

	hosts, _ := GetAllHosts(rdb, ctx)
	for k, v := range want {
		if !compareHosts(v, hosts[k]) {
			t.Errorf("%v != %v", v, hosts[k])
		}
	}

}

func TestDelHost(t *testing.T) {
	DbDef = 1 // use 1 for the test database
	AddrDef = "localhost:6380"
	rdb, ctx := InitClient()
	defer rdb.Close()

	if res, err := DelHost(rdb, ctx, "foo_1"); err != nil || res != 1 {
		t.Fatalf("res, err= %v, %v, want 1, nil", res, err)
	}

	if res, err := DelHost(rdb, ctx, "foo"); err != nil || res != 1 {
		t.Fatalf("res, err= %v, %v, want 1, nil", res, err)
	}

	if res, err := DelHost(rdb, ctx, "foo"); err != nil || res != 0 {
		t.Fatalf("res, err= %v, %v, want 1, nil", res, err)
	}

}

func TestAddMultHosts(t *testing.T) {
	DbDef = 1 // use 1 for the test database
	AddrDef = "localhost:6380"
	rdb, ctx := InitClient()
	defer rdb.Close()

	hosts := map[string]ZabbixHost{"foo_1": {HostIdZ: "foo_1",
		HostZ:    "foo_1_23",
		NameZ:    "name_foo_1",
		StatusZ:  "status_foo_1",
		ItNew:    true,
		ProblemZ: []string{"pfoo_1_1", "pfoo_2_1"}},
		"foo": {HostIdZ: "foo",
			NameZ:    "name_foo",
			StatusZ:  "status_foo",
			ProblemZ: []string{"foo1", "foo2"},
		}}
	if res, err := AddMultHosts(rdb, ctx, hosts); res != 2 && err != nil {
		t.Fatalf("AddMultHosts(rdb, ctx) = %d, %v, want 2, nil", res, err)
	}

	want, _ := GetAllHosts(rdb, ctx)
	for k, v := range want {
		if !compareHosts(v, hosts[k]) {
			t.Errorf("%v != %v", v, hosts[k])
		}
	}

}

func TestDelAllHossts(t *testing.T) {
	DbDef = 1 // use 1 for the test database
	AddrDef = "localhost:6380"

	rdb, ctx := InitClient()
	defer rdb.Close()

	if res, err := DelAllHosts(nil, ctx); res != 0 || !errors.Is(err, ErrEmptyInputData) {
		t.Fatalf("DelAllHosts(rdb, ctx) = %d, %v, want 0, err = %v", res, err, errors.New("the input data is empty"))
	}

	if res, err := DelAllHosts(rdb, nil); res != 0 || !errors.Is(err, ErrEmptyInputData) {
		t.Fatalf("DelAllHosts(rdb, ctx) = %d, %v, want 0, err = %v", res, err, errors.New("the input data is empty"))
	}

	want := map[string]ZabbixHost{"foo_1": {HostIdZ: "foo_1",
		HostZ:    "foo_1_23",
		NameZ:    "name_foo_1",
		StatusZ:  "status_foo_1",
		ItNew:    true,
		ProblemZ: []string{"pfoo_1_1", "pfoo_2_1"}},
		"foo": {HostIdZ: "foo",
			NameZ:    "name_foo",
			StatusZ:  "status_foo",
			ProblemZ: []string{"foo1", "foo2"},
		}}
	hosts, _ := GetAllHosts(rdb, ctx)

	for k, v := range want {
		if !compareHosts(v, hosts[k]) {
			t.Errorf("%v != %v", v, hosts[k])
		}
	}

	if res, err := DelAllHosts(rdb, ctx); res != int64(len(want)) && err != nil {
		t.Fatalf("DelAllHosts(rdb, ctx) = %d, %v, want %d, nil", res, err, len(want))
	}

	if hosts, err := GetAllHosts(rdb, ctx); len(hosts) != 0 && err != nil {
		t.Fatalf("GetAllHosts(rdb, ctx) = %v, %v, want map[], nil", hosts, err)
	}
}

func TestUpdateZabixHosts(t *testing.T) {
	DbDef = 1 // use 1 for the test database
	AddrDef = "localhost:6380"

	rdb, ctx := InitClient()
	defer rdb.Close()

	svdHosts := SavedHosts{}
	// svdHosts.RWD.Lock()
	// defer svdHosts.RWD.Unlock()

	if err := UpdateZabixHosts(rdb, ctx, &svdHosts); err != nil || len(svdHosts.Hosts) != 0 || svdHosts.Hosts == nil {
		t.Fatalf("err = %v, len(svdHosts.Hosts) = %d, want err = nil, len(svdHosts.Hosts) = 0 ", err, len(svdHosts.Hosts))
	}

	svdHosts.RWD.Lock()
	svdHosts.Hosts = map[string]ZabbixHost{"foo_1": {HostIdZ: "foo_1",
		HostZ:    "foo_1_23",
		NameZ:    "name_foo_1",
		StatusZ:  "status_foo_1",
		ItNew:    true,
		ProblemZ: []string{"pfoo_1_1", "pfoo_2_1"}},
		"foo": {HostIdZ: "foo",
			NameZ:    "name_foo",
			StatusZ:  "status_foo",
			ProblemZ: []string{"foo1", "foo2"},
		}}
	svdHosts.RWD.Unlock()

	if err := UpdateZabixHosts(rdb, ctx, &svdHosts); err != nil || len(svdHosts.Hosts) != 2 {
		t.Fatalf("err = %v, len(svdHosts.Hosts) = %d, want err = nil, len(svdHosts.Hosts) = 2 ", err, len(svdHosts.Hosts))
	}

	svdHosts.RWD.Lock()
	svdHosts.Hosts = nil
	t.Logf(" len(svdHosts.Hosts) = %d", len(svdHosts.Hosts))
	svdHosts.RWD.Unlock()

	if err := UpdateZabixHosts(rdb, ctx, &svdHosts); err != nil || len(svdHosts.Hosts) != 2 {
		t.Fatalf("err = %v, len(svdHosts.Hosts) = %d, want err = nil, len(svdHosts.Hosts) = 2 ", err, len(svdHosts.Hosts))
	}

	want := map[string]ZabbixHost{"foo_1": {HostIdZ: "foo_1",
		HostZ:    "foo_1_23",
		NameZ:    "name_foo_1",
		StatusZ:  "status_foo_1",
		ItNew:    true,
		ProblemZ: []string{"pfoo_1_1", "pfoo_2_1"}},
		"foo": {HostIdZ: "foo",
			NameZ:    "name_foo",
			StatusZ:  "status_foo",
			ProblemZ: []string{"foo1", "foo2"},
		}}

	svdHosts.RWD.RLock()
	for k, v := range want {
		if !compareHosts(v, svdHosts.Hosts[k]) {
			t.Errorf("%v != %v", v, svdHosts.Hosts[k])
		}
	}
	svdHosts.RWD.RUnlock()

	svdHosts.RWD.Lock()
	svdHosts.Hosts = map[string]ZabbixHost{"foo_100": {HostIdZ: "foo_100",
		HostZ:    "foo_1_2300",
		NameZ:    "name_foo_100",
		StatusZ:  "status_foo_100",
		ItNew:    true,
		ProblemZ: []string{"pfoo_1_100", "pfoo_2_10"}},
		"foo_200": {HostIdZ: "foo_200",
			HostZ:    "foo_2_2300",
			NameZ:    "name_foo_200",
			StatusZ:  "status_foo_200",
			ItNew:    true,
			ProblemZ: []string{"pfoo_1_200", "pfoo_2_20"}},

		"foo": {HostIdZ: "foo_00",
			NameZ:    "name_foo_00",
			StatusZ:  "status_foo_00",
			ProblemZ: []string{"foo1_00", "foo2_00"},
		}}
	svdHosts.RWD.Unlock()

	if err := UpdateZabixHosts(rdb, ctx, &svdHosts); err != nil || len(svdHosts.Hosts) != 3 {
		t.Fatalf("err = %v, len(svdHosts.Hosts) = %d, want err = nil, len(svdHosts.Hosts) = 3 ", err, len(svdHosts.Hosts))
	}

	svdHosts.RWD.RLock()
	for k, v := range want {
		if compareHosts(v, svdHosts.Hosts[k]) {
			t.Errorf("%v = %v", v, svdHosts.Hosts[k])
		}
	}
	svdHosts.RWD.RUnlock()

	want = map[string]ZabbixHost{"foo_100": {HostIdZ: "foo_100",
		HostZ:    "foo_1_2300",
		NameZ:    "name_foo_100",
		StatusZ:  "status_foo_100",
		ItNew:    true,
		ProblemZ: []string{"pfoo_1_100", "pfoo_2_10"}},
		"foo_200": {HostIdZ: "foo_200",
			HostZ:    "foo_2_2300",
			NameZ:    "name_foo_200",
			StatusZ:  "status_foo_200",
			ItNew:    true,
			ProblemZ: []string{"pfoo_1_200", "pfoo_2_20"}},

		"foo": {HostIdZ: "foo_00",
			NameZ:    "name_foo_00",
			StatusZ:  "status_foo_00",
			ProblemZ: []string{"foo1_00", "foo2_00"},
		}}

	svdHosts.RWD.RLock()
	for k, v := range want {
		if !compareHosts(v, svdHosts.Hosts[k]) {
			t.Errorf("%v != %v", v, svdHosts.Hosts[k])
		}
	}
	svdHosts.RWD.RUnlock()

	// t.Cleanup(func() { DelAllHosts(rdb, ctx) })
	t.Log("Cleanin up!")
	DelAllHosts(rdb, ctx)
}

func compareHosts(host1, host2 ZabbixHost) bool {
	if host1.HostIdZ != host2.HostIdZ ||
		host1.HostZ != host2.HostZ ||
		host1.NameZ != host2.NameZ ||
		host1.StatusZ != host2.StatusZ ||
		host1.ItChanged != host2.ItChanged ||
		host1.ItNew != host2.ItNew ||
		slices.Compare(host1.ProblemZ, host2.ProblemZ) != 0 {
		return false
	}

	return true
}
