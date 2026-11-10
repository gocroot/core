package config

import (
	"os"
	"runtime/debug"
	"strings"
)

// Version adalah versi kode yang sedang berjalan, ditampilkan di GET /version.
// Workflow deploy membandingkannya dengan commit yang baru di-push untuk
// memastikan proses baru benar-benar melayani, bukan proses lama yang belum
// berhenti. Isinya hanya SHA commit, jangan pernah menaruh secret di sini.
//
// Sumber, dari yang paling diutamakan:
//  1. -ldflags "-X <modul>/config.Version=..." saat build
//  2. env VERSION (Cloud Functions: Google membangun tanpa folder .git)
//  3. revisi git yang ditanam otomatis oleh go build (vcs.revision),
//     ditambah "-dirty" jika ada perubahan yang belum di-commit
//  4. "dev"
var Version string

func resolveVersion(ldflags, env string, info *debug.BuildInfo) string {
	if v := strings.TrimSpace(ldflags); v != "" {
		return v
	}
	if v := strings.TrimSpace(env); v != "" {
		return v
	}
	if info != nil {
		var rev string
		var dirty bool
		for _, s := range info.Settings {
			switch s.Key {
			case "vcs.revision":
				rev = s.Value
			case "vcs.modified":
				dirty = s.Value == "true"
			}
		}
		if rev != "" {
			if dirty {
				return rev + "-dirty"
			}
			return rev
		}
	}
	return "dev"
}

func setVersion() {
	info, _ := debug.ReadBuildInfo()
	Version = resolveVersion(Version, os.Getenv("VERSION"), info)
}
