package msgmngr

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/maaw77/telbotnsn/brds"
)

// TestFormatProblemHostZbxEmpty  calls msgmngr.formatProblemHostZbx with an empty input data,
// checking for an error.
func TestFormatProblemHostZbxEmpty(t *testing.T) {
	var currentHosts *brds.SavedHosts

	if v, err := formatProblemHostZbx(currentHosts); err == nil {
		t.Fatalf(`formatProblemHosts(nil)= %s, %v, want "",error`, v, err)
	}

	currentHosts = &brds.SavedHosts{}
	if v, err := formatProblemHostZbx(currentHosts); err == nil {
		t.Fatalf(`formatProblemHosts(nil)= %s, %v, want "", error`, v, err)
	}
}

// TestFormatProblemHostZbx calls msgmngr.formatProblemHostZbx with
// various arguments, checking the correctness of the output string.
func TestFormatProblemHostZbx(t *testing.T) {
	currentHosts := &brds.SavedHosts{Hosts: map[string]brds.ZabbixHost{"1": {HostIdZ: "1",
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
		"9": {HostIdZ: "9",
			HostZ:     "host_9",
			NameZ:     "name_host_9",
			ProblemZ:  []string{"problem_3_host_9"},
			ItNew:     false,
			ItChanged: true},
		"10": {HostIdZ: "10",
			HostZ:     "host_10",
			NameZ:     "name_host_10",
			ProblemZ:  []string{"problem_111_host_10", "problem_32_host_10"},
			ItNew:     true,
			ItChanged: false},
		"11": {HostIdZ: "11",
			HostZ:     "host_11",
			NameZ:     "name_host_11",
			ProblemZ:  []string{"problem_111_host_11", "problem_32_host_11"},
			ItNew:     true,
			ItChanged: false},
		"12": {HostIdZ: "12",
			HostZ:     "host_12",
			NameZ:     "name_host_12",
			ProblemZ:  []string{"problem_111_host_12", "problem_32_host_12"},
			ItNew:     false,
			ItChanged: false},
	}}

	var numName, numChName, numNewName int
	currentHosts.RWD.RLock()
	for _, host := range currentHosts.Hosts {
		switch {
		case host.ItNew:
			numNewName += 1
		case host.ItChanged:
			numChName += 1
		default:
			numName += 1
		}
	}
	currentHosts.RWD.RUnlock()

	outString, err := formatProblemHostZbx(currentHosts)
	if err != nil {
		t.Fatalf(`formatProblemHosts(notNil)= %s, %v, want "outString", nil`, outString, err)
	}

	currentHosts.RWD.RLock()
	parStr := fmt.Sprintf(`%d \(%d new, %d changed\)</b>`, len(currentHosts.Hosts), numNewName, numChName)
	currentHosts.RWD.RUnlock()

	want := regexp.MustCompile(parStr)
	num := len(want.FindAllString(outString, -1))
	t.Log(num)
	if num != 1 {
		t.Errorf(`The string "%s" was not found`, parStr)
	}

	want = regexp.MustCompile(`<b>Host name`)
	num = len(want.FindAllString(outString, -1))
	if num != numName {
		t.Errorf(`The string "<b>ch_Host name" was not found in the required number(%d), want %d`, num, numName)
	}

	want = regexp.MustCompile(`<b>ch_Host name`)
	num = len(want.FindAllString(outString, -1))
	if num != numChName {
		t.Errorf(`The string "<b>ch_Host name" was not found in the required number(%d), want %d`, num, numChName)
	}

	want = regexp.MustCompile(`<b>new_Host name`)
	num = len(want.FindAllString(outString, -1))
	if num != numNewName {
		t.Errorf(`The string "<b>new_Host name" was not found in the required number(%d), want %d`, num, numNewName)
	}
}

// TestFormatResotoredHostZbxEmpty  calls msgmngr.formatRestoredHostZbx with an empty input data,
// checking for an error.
func TestFormatRestoredHostZbxEmpty(t *testing.T) {
	var currentHosts *brds.SavedHosts

	if v, err := formatRestoredHostZbx(currentHosts); err == nil {
		t.Fatalf(`formatProblemHosts(nil)= %s, %v, want "",error`, v, err)
	}

	currentHosts = &brds.SavedHosts{}
	if v, err := formatRestoredHostZbx(currentHosts); err == nil {
		t.Fatalf(`formatProblemHosts(nil)= %s, %v, want "", error`, v, err)
	}
}

// TestFormatRestoredHostZbx calls msgmngr.formatProblemHostZbx with
// various arguments, checking the correctness of the output string.
func TestFormatRestoredHostZbx(t *testing.T) {
	currentHosts := &brds.SavedHosts{Hosts: map[string]brds.ZabbixHost{"1": {HostIdZ: "1",
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
		"5": {HostIdZ: "5",
			HostZ:     "host_5",
			NameZ:     "name_host_5",
			ProblemZ:  []string{"problem_3_host_8"},
			ItNew:     false,
			ItChanged: true},
	}}

	outString, err := formatRestoredHostZbx(currentHosts)
	if err != nil {
		t.Fatalf(`formatProblemHosts(notNil)= %s, %v, want "outString", nil`, outString, err)
	}
	// t.Log(outString)

	currentHosts.RWD.RLock()
	defer currentHosts.RWD.RUnlock()
	parStr := fmt.Sprintf(`of restored hosts is %d.</b>`, len(currentHosts.Hosts))

	want := regexp.MustCompile(parStr)
	num := len(want.FindAllString(outString, -1))
	if num != 1 {
		t.Errorf(`The string "%s" was not found`, parStr)
	}

	for i := 1; i <= len(currentHosts.Hosts); i++ {

		want = regexp.MustCompile(`<b>Host name:</b> name_host_` + strconv.Itoa(i))
		if len(want.FindAllString(outString, -1)) != 1 {
			t.Errorf(`The string "<b>Host name:<b> name_host_%s" was not found`, strconv.Itoa(i))
		}
	}
}
