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
	wantOutSrings = []string{"123\n...", "...\n456\n...", "...\n789\n"}
	i = 0
	for outString := range sliceMessage(inpString, 5) {
		// t.Logf("outString= %s, len(outString)=%d ", outString, len(outString))

		if outString != wantOutSrings[i] {

			t.Errorf("%s != %s\n", outString, wantOutSrings[i])
		}
		i += 1
	}

	inpString = "123Ф\n456\n789\n"
	wantOutSrings = []string{"123Ф\n...", "...\n456\n...", "...\n789\n"}
	i = 0
	for outString := range sliceMessage(inpString, 6) {
		// t.Logf("outString= %s, len(outString)=%d ", outString, len(outString))

		if outString != wantOutSrings[i] {

			t.Errorf("%s != %s\n", outString, wantOutSrings[i])
		}
		i += 1
	}

	inpString = "123\n456\n789"
	wantOutSrings = []string{"123\n...", "...\n456\n...", "...\n789"}
	i = 0
	for outString := range sliceMessage(inpString, 4) {
		// t.Logf("outString= %s, len(outString)=%d ", outString, len(outString))

		if outString != wantOutSrings[i] {

			t.Errorf("%s != %s\n", outString, wantOutSrings[i])
		}
		i += 1
	}

	inpString = "123Ф\n456\n789\n"
	wantOutSrings = []string{"123Ф\n456\n...", "...\n789\n"}
	i = 0
	for outString := range sliceMessage(inpString, 11) {
		// t.Logf("outString= %s, len(outString)=%d ", outString, len(outString))

		if outString != wantOutSrings[i] {

			t.Errorf("%s != %s\n", outString, wantOutSrings[i])
		}
		i += 1
	}

	inpString = `<b>new_Host name:</b> Микротик Панимба-СисАТ, <b>problems:</b>[Disk-131072: Disk space is critically low (used &gt; {$VFS.FS.PUSED.MAX.CRIT:&#34;Disk-131072&#34;}%)]
	<b>new_Host name:</b> Микротик ГореВая LTE_Megafon, <b>problems:</b>[Disk-131072: Disk space is critically low (used &gt; {$VFS.FS.PUSED.MAX.CRIT:&#34;Disk-131072&#34;}%)]
	<b>new_Host name:</b> Микротик Караган резервный, <b>problems:</b>[Нет пинга]
	<b>new_Host name:</b> Микротик Тюрепино ПК63, <b>problems:</b>[Нет пинга]
	<b>new_Host name:</b> Микротик Сисим Unlim, <b>problems:</b>[Нет пинга]
	<b>new_Host name:</b> Микротик Биза, <b>problems:</b>[Disk-131072: Disk space is critically low (used &gt; {$VFS.FS.PUSED.MAX.CRIT:&#34;Disk-131072&#34;}%)]
	<b>new_Host name:</b> Микротик Озёрка Резерв теле2, <b>problems:</b>[{HOST.NAME} ребутнулся (uptime &lt; 10m) Нет пинга]
	<b>new_Host name:</b> Микротик БК26, <b>problems:</b>[#1: High CPU utilization (over {$CPU.UTIL.CRIT}% for 5m)]
	<b>new_Host name:</b> Микротик Юхтахта, <b>problems:</b>[Нет пинга]
	<b>new_Host name:</b> Микротик Караган LTE, <b>problems:</b>[Нет пинга]
	<b>new_Host name:</b> Микротик Пескино Альтегра, <b>problems:</b>[Нет пинга]
	<b>new_Host name:</b> Микротик Тукша, <b>problems:</b>[Disk-131072: Disk space is critically low (used &gt; {$VFS.FS.PUSED.MAX.CRIT:&#34;Disk-131072&#34;}%)]
	<b>new_Host name:</b> Микротик Дражные Тайлы, <b>problems:</b>[Нет пинга]
	<b>new_Host name:</b> Микротик Н. Шумиха unlim, <b>problems:</b>[Disk-131072: Disk space is critically low (used &gt; {$VFS.FS.PUSED.MAX.CRIT:&#34;Disk-131072&#34;}%)]
	<b>new_Host name:</b> Микротик Мишкин, <b>problems:</b>[Нет пинга]
	<b>new_Host name:</b> Микротик Еленка, <b>problems:</b>[Disk-131072: Disk space is critically low (used &gt; {$VFS.FS.PUSED.MAX.CRIT:&#34;Disk-131072&#34;}%)]
	<b>new_Host name:</b> Микротик Илинский Unlim, <b>problems:</b>[Disk-131072: Disk space is critically low (used &gt; {$VFS.FS.PUSED.MAX.CRIT:&#34;Disk-131072&#34;}%)]
	<b>new_Host name:</b> Микротик Правая Безымянка Wi-Max, <b>problems:</b>[Disk-131072: Disk space is critically low (used &gt; {$VFS.FS.PUSED.MAX.CRIT:&#34;Disk-131072&#34;}%)]
	<b>new_Host name:</b> Микротик Ивановка, <b>problems:</b>[Нет пинга]
	<b>new_Host name:</b> Микротик Безобразовка, <b>problems:</b>[Disk-131072: Disk space is critically low (used &gt; {$VFS.FS.PUSED.MAX.CRIT:&#34;Disk-131072&#34;}%)]
	<b>new_Host name:</b> Микротик Безымянка-С, <b>problems:</b>[Нет пинга]
	<b>new_Host name:</b> Микротик ГореВая LTE - MTS, <b>problems:</b>[Нет пинга Disk-131072: Disk space is critically low (used &gt; {$VFS.FS.PUSED.MAX.CRIT:&#34;Disk-131072&#34;}%)]
	<b>new_Host name:</b> Микротик Абакан промышленная, <b>problems:</b>[Нет пинга]
	<b>new_Host name:</b> Микротик Американский, <b>problems:</b>[Disk-131072: Disk space is critically low (used &gt; {$VFS.FS.PUSED.MAX.CRIT:&#34;Disk-131072&#34;}%)]
	<b>new_Host name:</b> Микротик Горелая-Полигон, <b>problems:</b>[{HOST.NAME} ребутнулся (uptime &lt; 10m) Disk-131072: Disk space is critically low (used &gt; {$VFS.FS.PUSED.MAX.CRIT:&#34;Disk-131072&#34;}%)]
	<b>new_Host name:</b> Микротик Панимба Альтегра, <b>problems:</b>[Disk-131072: Disk space is critically low (used &gt; {$VFS.FS.PUSED.MAX.CRIT:&#34;Disk-131072&#34;}%)]
	<b>new_Host name:</b> Микротик Тюрепино ПК61, <b>problems:</b>[Нет пинга]
	<b>new_Host name:</b> Микротик Надежда MGF, <b>problems:</b>[Disk-131072: Disk space is critically low (used &gt; {$VFS.FS.PUSED.MAX.CRIT:&#34;Disk-131072&#34;}%)]
	<b>new_Host name:</b> Микротик Демидовский, <b>problems:</b>[Нет пинга]
	<b>new_Host name:</b> Микротик Кундат-2, <b>problems:</b>[Disk-131072: Disk space is critically low (used &gt; {$VFS.FS.PUSED.MAX.CRIT:&#34;Disk-131072&#34;}%)]
	<b>new_Host name:</b> Микротик Кувай Ростелеком, <b>problems:</b>[Disk-131072: Disk space is critically low (used &gt; {$VFS.FS.PUSED.MAX.CRIT:&#34;Disk-131072&#34;}%)]
	<b>new_Host name:</b> Микротик Караган Ямал, <b>problems:</b>[Disk-131072: Disk space is critically low (used &gt; {$VFS.FS.PUSED.MAX.CRIT:&#34;Disk-131072&#34;}%)]
	<b>new_Host name:</b> Микротик Тукша-unlim, <b>problems:</b>[Нет пинга Disk-131072: Disk space is critically low (used &gt; {$VFS.FS.PUSED.MAX.CRIT:&#34;Disk-131072&#34;}%)]
	<b>new_Host name:</b> Микротик Мамон R-телеком, <b>problems:</b>[Disk-131072: Disk space is critically low (used &gt; {$VFS.FS.PUSED.MAX.CRIT:&#34;Disk-131072&#34;}%)]
	<b>new_Host name:</b> Микротик Ковальков М.Н., <b>problems:</b>[Видеорегистратор отвалился]
	<b>new_Host name:</b> Микротик Н. Шумиха, <b>problems:</b>[Disk-131072: Disk space is critically low (used &gt; {$VFS.FS.PUSED.MAX.CRIT:&#34;Disk-131072&#34;}%)]
	<b>new_Host name:</b> Микротик Караган Алинк unlim, <b>problems:</b>[Disk-131072: Disk space is critically low (used &gt; {$VFS.FS.PUSED.MAX.CRIT:&#34;Disk-131072&#34;}%)]

	The number of problematic hosts is <b>38 (38 new, 0 changed)</b>`
	i = 0
	for outString := range sliceMessage(inpString, 1000) {
		t.Logf("outString= %s, len(outString)=%d ", outString, len(outString))

		// 	if outString != wantOutSrings[i] {

		//		t.Errorf("%s != %s\n", outString, wantOutSrings[i])
		//	}
		//
		// i += 1
	}
}
