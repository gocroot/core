package controller

import "strings"

// PathParam mengambil satu segmen path setelah prefix, untuk route dengan parameter.
//
//	PathParam("/data/user/123", "/data/user/")     // "123", true
//	PathParam("/data/user/", "/data/user/")        // "", false (kosong)
//	PathParam("/data/user/123/x", "/data/user/")   // "", false (lebih dari satu segmen)
//	PathParam("/data/item/123", "/data/user/")     // "", false (prefix beda)
//
// Bisa dipakai di routes untuk mencocokkan route sekaligus di controller untuk membaca nilainya.
// Tidak mengalokasikan memori: hasilnya adalah potongan dari path asli.
func PathParam(path, prefix string) (string, bool) {
	param, ok := strings.CutPrefix(path, prefix)
	if !ok || param == "" || strings.IndexByte(param, '/') >= 0 {
		return "", false
	}
	return param, true
}
