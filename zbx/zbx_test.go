package zbx

import (
	"fmt"
	"slices"
	"testing"

	"github.com/maaw77/telbotnsn/brds"
	"github.com/maaw77/telbotnsn/msgmngr"
)

// TestCompareHosts calls zbx.compareHosts with nil arguments, checking
// for an error.
func TestCompareHostsNil(t *testing.T) {
	var lastHosts, fixHosts, currentHosts *brds.SavedHosts
	commandQueueFromZbx := make(chan msgmngr.CommandFromZbx, 5)

	if err := compareHosts(lastHosts, fixHosts, currentHosts, commandQueueFromZbx); err == nil {
		t.Fatalf("compareHosts(lastHost = Nil, currentHosts = nil) = %v, want error", err)
	} else if err.Error() != "the input data is nil" {
		t.Errorf("compareHosts(lastHost.Hosts = Nil, currentHosts.Hosts = nil) = \"%v\", want error.Error()=\"the input data is nil\"", err)
	}

	lastHosts = &brds.SavedHosts{}
	currentHosts = &brds.SavedHosts{}
	fixHosts = &brds.SavedHosts{}

	if err := compareHosts(lastHosts, fixHosts, currentHosts, commandQueueFromZbx); err == nil {
		t.Fatalf("compareHosts(lastHost.Hosts = Nil, currentHost.Hosts = nil) = %v, want error", err)
	} else if err.Error() != "hosts are nil" {
		t.Fatalf("compareHosts(lastHost.Hosts = Nil, currentHosts.Hosts = nil) = \"%v\", want error.Error()=\"hosts are nil\"", err)
	}
}

// TestCompareHosts calls zbx.compareHosts with arguments, checking
// fixHosts and currentHosts for a valid value in the brds.ZabbixHost fields.
func TestCompareHosts(t *testing.T) {

	commandQueueFromZbx := make(chan msgmngr.CommandFromZbx, 5)

	fixHosts := &brds.SavedHosts{Hosts: map[string]brds.ZabbixHost{"22": {HostIdZ: "22",
		HostZ:    "host_22",
		NameZ:    "name_host_22",
		ProblemZ: []string{"problem_1_host_22", "problem_2_host_22"}},
		"33": {HostIdZ: "33",
			HostZ:    "host_33",
			NameZ:    "name_host_33",
			ProblemZ: []string{"problem_1_host_22"}},
	}}

	lastHosts := &brds.SavedHosts{Hosts: map[string]brds.ZabbixHost{"1": {HostIdZ: "1",
		HostZ:    "host_1",
		NameZ:    "name_host_1",
		ProblemZ: []string{"problem_1_host_1", "problem_2_host_1"}},
		"2": {HostIdZ: "2",
			HostZ:    "host_2",
			NameZ:    "name_host_2",
			ProblemZ: []string{"problem_1_host_2", "problem_2_host_2"}},
		"3": {HostIdZ: "3",
			HostZ:    "host_3",
			NameZ:    "name_host_3",
			ProblemZ: []string{"problem_1_host_3", "problem_2_host_3"}},
		"5": {HostIdZ: "5",
			HostZ:    "host_5",
			NameZ:    "name_host_5",
			ProblemZ: []string{"problem_1_host_5", "problem_2_host_5"}},
		"6": {HostIdZ: "6",
			HostZ:    "host_6",
			NameZ:    "name_host_6",
			ProblemZ: []string{"problem_8_host_6", "problem_9_host_6"}},
		"7": {HostIdZ: "7",
			HostZ:    "host_7",
			NameZ:    "name_host_7",
			ProblemZ: []string{"problem_8_host_7"}},
		"8": {HostIdZ: "8",
			HostZ:    "host_8",
			NameZ:    "name_host_8",
			ProblemZ: []string{"problem_33_host_88"}},
	}}

	currentHosts := &brds.SavedHosts{Hosts: map[string]brds.ZabbixHost{"1": {HostIdZ: "1",
		HostZ:    "host_1",
		NameZ:    "name_host_1",
		ProblemZ: []string{"problem_1_host_1", "problem_2_host_1"}},
		"2": {HostIdZ: "2",
			HostZ:    "host_2",
			NameZ:    "name_host_2",
			ProblemZ: []string{"problem_1_host_2", "problem_2_host_2", "problem_3_host_2"}},
		"3": {HostIdZ: "3",
			HostZ:    "host_3",
			NameZ:    "name_host_3",
			ProblemZ: []string{"problem_1_host_3", "problem_2_host_3"}},
		"4": {HostIdZ: "4",
			HostZ:    "host_4",
			NameZ:    "name_host_4",
			ProblemZ: []string{"problem_11_host_44", "problem_2_host_4"}},
		"8": {HostIdZ: "8",
			HostZ:    "host_8",
			NameZ:    "name_host_8",
			ProblemZ: []string{"problem_3_host_8"}},
	}}

	wantFixHosts := map[string]brds.ZabbixHost{"5": {HostIdZ: "5",
		HostZ:    "host_5",
		NameZ:    "name_host_5",
		ProblemZ: []string{"problem_1_host_5", "problem_2_host_5"}},
		"6": {HostIdZ: "6",
			HostZ:    "host_6",
			NameZ:    "name_host_6",
			ProblemZ: []string{"problem_8_host_6", "problem_9_host_6"}},
		"7": {HostIdZ: "7",
			HostZ:    "host_7",
			NameZ:    "name_host_7",
			ProblemZ: []string{"problem_8_host_7"}},
	}

	wantCurrentHosts := map[string]brds.ZabbixHost{"1": {HostIdZ: "1",
		HostZ:     "host_1",
		NameZ:     "name_host_1",
		ProblemZ:  []string{"problem_1_host_1", "problem_2_host_1"},
		ItNew:     false,
		ItChanged: false},
		"2": {HostIdZ: "2",
			HostZ:     "host_2",
			NameZ:     "name_host_2",
			ProblemZ:  []string{"problem_1_host_2", "problem_2_host_2", "problem_3_host_2"},
			ItNew:     false,
			ItChanged: true},
		"3": {HostIdZ: "3",
			HostZ:     "host_3",
			NameZ:     "name_host_3",
			ProblemZ:  []string{"problem_1_host_3", "problem_2_host_3"},
			ItNew:     false,
			ItChanged: false},
		"4": {HostIdZ: "4",
			HostZ:     "host_4",
			NameZ:     "name_host_4",
			ProblemZ:  []string{"problem_11_host_44", "problem_2_host_4"},
			ItNew:     true,
			ItChanged: false},
		"8": {HostIdZ: "8",
			HostZ:     "host_8",
			NameZ:     "name_host_8",
			ProblemZ:  []string{"problem_3_host_8"},
			ItNew:     false,
			ItChanged: true},
	}

	if err := compareHosts(lastHosts, fixHosts, currentHosts, commandQueueFromZbx); err != nil {
		t.Errorf("compareHosts() = \"%v\", want nil", err)
	}

	// Checking fixHost.Hosts
	fixHosts.RWD.RLock()
	defer fixHosts.RWD.RUnlock()

	if len(fixHosts.Hosts) != len(wantFixHosts) {
		t.Fatal("len(fixHost.Hosts) != len(wantFixHosts)")
	}

	for k := range wantFixHosts {
		_, ok := fixHosts.Hosts[k]
		if !ok {
			t.Errorf("wantFixHosts[%s] is not in fixHost.Hosts", k)
		}

	}

	// Checking currentHosts
	currentHosts.RWD.RLock()
	defer currentHosts.RWD.RUnlock()

	if len(currentHosts.Hosts) != len(wantCurrentHosts) {
		t.Fatal("len(currentHosts.Hosts) != len(wantCurretHosts)")
	}

	var counterNewHosts, counterChangedHosts int
	for kw, vw := range wantCurrentHosts {
		v, ok := currentHosts.Hosts[kw]
		if !ok {
			t.Fatalf("wantCurrentHosts[%s] is not in currentHosts.Hosts", kw)
		} else {
			if vw.ItNew {
				counterNewHosts += 1
			}
			if vw.ItChanged {
				counterChangedHosts += 1
			}
			if vw.HostIdZ != v.HostIdZ {
				t.Errorf("wantCurrentHosts[%s].HostIdZ=%s != currentHosts.Hosts[%s].HostIdZ=%s", kw, vw.HostIdZ, kw, v.HostIdZ)
			}
			if vw.HostZ != v.HostZ {
				t.Errorf("wantCurrentHosts[%s].HostZ=%s != currentHosts.Hosts[%s].HostZ=%s", kw, vw.HostZ, kw, v.HostZ)
			}
			if vw.NameZ != v.NameZ {
				t.Errorf("wantCurrentHosts[%s].NameZ=%s != currentHosts.Hosts[%s].NameZ=%s", kw, vw.NameZ, kw, v.NameZ)
			}
			if slices.Compare(vw.ProblemZ, currentHosts.Hosts[kw].ProblemZ) != 0 {
				t.Errorf("wantCurrentHosts[%s]ProblemZ=%v != currentHosts.Hosts[%s].ProblemZ=%v", kw, vw.ProblemZ, kw, v.ProblemZ)
			}
			if vw.ItNew != v.ItNew {
				t.Errorf("wantCurrentHosts[%s].ItNew=%t != currentHosts.Hosts[%s].ItNew=%t", kw, vw.ItNew, kw, v.ItNew)
			}
			if vw.ItChanged != v.ItChanged {
				t.Errorf("wantCurrentHosts[%s].ItChanged=%t != currentHosts.Hosts[%s].ItChanged=%t", kw, vw.ItChanged, kw, v.ItChanged)
			}
		}
	}
	info := fmt.Sprintf("The number of problematic hosts is <b>%d (%d new, %d changed)</b>.\nThe number of restored hosts is <b>%d</b>.",
		len(currentHosts.Hosts), counterNewHosts, counterChangedHosts, len(fixHosts.Hosts))

	mess := <-commandQueueFromZbx
	if info != mess.TextMessage {
		t.Error(info)

	}
}
