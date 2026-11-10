package config

import (
	"runtime/debug"
	"testing"
)

func TestResolveVersion(t *testing.T) {
	vcs := func(rev, modified string) *debug.BuildInfo {
		return &debug.BuildInfo{Settings: []debug.BuildSetting{
			{Key: "vcs.revision", Value: rev},
			{Key: "vcs.modified", Value: modified},
		}}
	}
	cases := []struct {
		name, ldflags, env string
		info               *debug.BuildInfo
		want               string
	}{
		{"ldflags menang", "abc", "env", vcs("rev", "false"), "abc"},
		{"env jika tanpa ldflags", "", " def ", vcs("rev", "false"), "def"},
		{"revisi git", "", "", vcs("123abc", "false"), "123abc"},
		{"revisi git dirty", "", "", vcs("123abc", "true"), "123abc-dirty"},
		{"tanpa info vcs", "", "", &debug.BuildInfo{}, "dev"},
		{"tanpa build info", "", "", nil, "dev"},
	}
	for _, c := range cases {
		if got := resolveVersion(c.ldflags, c.env, c.info); got != c.want {
			t.Errorf("%s: resolveVersion = %q, want %q", c.name, got, c.want)
		}
	}
}
