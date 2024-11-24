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
	// for outString := range sliceMessage(inpString, 2) {
	// 	t.Log(outString)
	// }

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
	inpString = "А3А\nВ3ВС3СД3Сс\n"
	wantOutSrings = []string{"А3...", "...А\n...", "...В3...", "...В...", "...С3...", "...С...", "...Д3...", "...С...", "...с\n"}
	i = 0
	for outString := range sliceMessage(inpString, 3) {
		if outString != wantOutSrings[i] {
			t.Errorf("%s != %s\n", outString, wantOutSrings[i])
		}
		i += 1
		// t.Log(outString)
	}

	inpString = ""
	wantOutSrings = []string{""}
	i = 0
	for outString := range sliceMessage(inpString, -3) {
		if outString != wantOutSrings[i] {
			t.Errorf("%s != %s\n", outString, wantOutSrings[i])
		}
		i += 1
		// t.Log(outString)
	}

	// HTMl tags
	inpString = "The number of problematic hosts is <b>%d (%d new, %d changed)</b>.\nThe number of restored hosts is <b>%d</b>."
	// wantOutSrings = []string{""}
	i = 0
	for outString := range sliceMessage(inpString, 40) {
		// if outString != wantOutSrings[i] {
		// 	t.Errorf("%s != %s\n", outString, wantOutSrings[i])
		// }
		// i += 1
		t.Log(outString)
	}

}
