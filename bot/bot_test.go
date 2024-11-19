package bot

import (
	"testing"
)

// TestSliceMessage calls bot.sliceMessage with
// various arguments, checking the correctness of the output string.
func TestSliceMessage(t *testing.T) {
	inpString := "AABBCCDDEE"
	wantOutSrings := []string{"AA...", "...BB...", "...CC...", "...DD...", "...EE"}
	var i int
	for outString := range sliceMessage(inpString, 2) {
		if outString != wantOutSrings[i] {
			t.Errorf("%s != %s\n", outString, wantOutSrings[i])
		}
		i += 1
	}

	inpString = "AABBCCDDEEC"
	wantOutSrings = []string{"AA...", "...BB...", "...CC...", "...DD...", "...EE...", "...C"}
	i = 0
	for outString := range sliceMessage(inpString, 2) {
		if outString != wantOutSrings[i] {
			t.Errorf("%s != %s\n", outString, wantOutSrings[i])
		}
		i += 1
	}

	inpString = "A3AB3BC3CD3DE3ECc"
	wantOutSrings = []string{"A3A...", "...B3B...", "...C3C...", "...D3D...", "...E3E...", "...Cc"}
	i = 0
	for outString := range sliceMessage(inpString, 3) {
		if outString != wantOutSrings[i] {
			t.Errorf("%s != %s\n", outString, wantOutSrings[i])
		}
		i += 1
	}
	inpString = ""
	wantOutSrings = []string{""}
	i = 0
	for outString := range sliceMessage(inpString, 3) {
		if outString != wantOutSrings[i] {
			t.Errorf("%s != %s\n", outString, wantOutSrings[i])
		}
		i += 1
	}
	inpString = "123"
	wantOutSrings = []string{"123"}
	i = 0
	for outString := range sliceMessage(inpString, 3) {
		if outString != wantOutSrings[i] {
			t.Errorf("%s != %s\n", outString, wantOutSrings[i])
		}
		i += 1
	}

	inpString = "1234"
	wantOutSrings = []string{"123...", "...4"}
	i = 0
	for outString := range sliceMessage(inpString, 3) {
		if outString != wantOutSrings[i] {
			t.Errorf("%s != %s\n", outString, wantOutSrings[i])
		}
		i += 1
	}

	// Russian
	inpString = "А3АВ3ВС3СД3ДЕ3ЕСс"
	wantOutSrings = []string{"А3...", "...АВ...", "...3В...", "...С3...", "...СД...", "...3Д...", "...Е3...", "...ЕС...", "...с"}
	i = 0
	for outString := range sliceMessage(inpString, 3) {
		if outString != wantOutSrings[i] {
			t.Errorf("%s != %s\n", outString, wantOutSrings[i])
		}
		i += 1
		// t.Log(outString)
	}
}
