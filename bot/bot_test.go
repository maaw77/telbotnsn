package bot

import (
	"testing"
)

// TestSliceMessage calls bot.sliceMessage with
// various arguments, checking the correctness of the output string.
func TestSliceMessage(t *testing.T) {

	inpString := "AABBCCDDEE"
	wantOutSrings := []string{""}
	var i int
	for outString := range sliceMessage(inpString, 0) {
		if outString != wantOutSrings[i] {
			t.Errorf("%s != %s\n", outString, wantOutSrings[i])
		}
		i += 1
	}

	inpString = "1234567Ы"
	wantOutSrings = []string{"1234567Ы"}
	i = 0
	for outString := range sliceMessage(inpString, 9) {
		// t.Logf("outString= %s", outString)
		if outString != wantOutSrings[i] {

			t.Errorf("%s != %s\n", outString, wantOutSrings[i])
		}
		i += 1
	}

	inpString = "123\n456\n789\n"
	wantOutSrings = []string{"123...", "...456...", "...789"}
	i = 0
	for outString := range sliceMessage(inpString, 5) {
		// t.Logf("outString= %s, len(outString)=%d ", outString, len(outString))

		if outString != wantOutSrings[i] {

			t.Errorf("%s != %s\n", outString, wantOutSrings[i])
		}
		i += 1
	}

	inpString = "123Ф\n456\n789\n"
	wantOutSrings = []string{"123Ф...", "...456...", "...789"}
	i = 0
	for outString := range sliceMessage(inpString, 6) {
		// t.Logf("outString= %s, len(outString)=%d ", outString, len(outString))

		if outString != wantOutSrings[i] {

			t.Errorf("%s != %s\n", outString, wantOutSrings[i])
		}
		i += 1
	}

	inpString = "123\n456\n789\n"
	wantOutSrings = []string{"123...", "...456...", "...789"}
	i = 0
	for outString := range sliceMessage(inpString, 4) {
		// t.Logf("outString= %s, len(outString)=%d ", outString, len(outString))

		if outString != wantOutSrings[i] {

			t.Errorf("%s != %s\n", outString, wantOutSrings[i])
		}
		i += 1
	}

	inpString = "123Ф\n456\n789\n"
	wantOutSrings = []string{"123Ф\n456...", "...789"}
	i = 0
	for outString := range sliceMessage(inpString, 11) {
		// t.Logf("outString= %s, len(outString)=%d ", outString, len(outString))

		if outString != wantOutSrings[i] {

			t.Errorf("%s != %s\n", outString, wantOutSrings[i])
		}
		i += 1
	}
}
