package main

import (
	"errors"
	"os"
	"testing"
	"time"
)

type fakeFileInfo struct {
	name  string
	isDir bool
}

func (info fakeFileInfo) Name() string       { return info.name }
func (info fakeFileInfo) Size() int64        { return 0 }
func (info fakeFileInfo) Mode() os.FileMode  { return 0755 }
func (info fakeFileInfo) ModTime() time.Time { return time.Time{} }
func (info fakeFileInfo) IsDir() bool        { return info.isDir }
func (info fakeFileInfo) Sys() interface{}   { return nil }

func fakeStat(existing map[string]bool) func(string) (os.FileInfo, error) {
	return func(path string) (os.FileInfo, error) {
		if existing[path] {
			return fakeFileInfo{name: path}, nil
		}

		return nil, os.ErrNotExist
	}
}

func TestDefaultUnixShellPrefersShellEnvironment(t *testing.T) {
	got := defaultUnixShell(
		func(string) string { return "/usr/local/bin/fish" },
		fakeStat(map[string]bool{"/usr/local/bin/fish": true}),
		func(string) ([]byte, error) { return nil, errors.New("not read") },
		501,
	)

	if got != "/usr/local/bin/fish" {
		t.Fatalf("defaultUnixShell() = %q, want %q", got, "/usr/local/bin/fish")
	}
}

func TestEffectiveCommandArgsKeepsExplicitCommandReadOnly(t *testing.T) {
	got, usingDefaultCommand := effectiveCommandArgs([]string{"top", "-u"})
	if usingDefaultCommand {
		t.Fatal("effectiveCommandArgs() usingDefaultCommand = true, want false")
	}
	if len(got) != 2 || got[0] != "top" || got[1] != "-u" {
		t.Fatalf("effectiveCommandArgs() = %#v, want %#v", got, []string{"top", "-u"})
	}
}

func TestEffectiveCommandArgsUsesDefaultForEmptyCommand(t *testing.T) {
	got, usingDefaultCommand := effectiveCommandArgs(nil)
	if !usingDefaultCommand {
		t.Fatal("effectiveCommandArgs() usingDefaultCommand = false, want true")
	}
	if len(got) != 1 || got[0] == "" {
		t.Fatalf("effectiveCommandArgs() = %#v, want one non-empty command", got)
	}
}

func TestDefaultUnixShellUsesPasswdWhenShellEnvironmentMissing(t *testing.T) {
	got := defaultUnixShell(
		func(string) string { return "" },
		fakeStat(map[string]bool{"/bin/zsh": true}),
		func(path string) ([]byte, error) {
			if path != passwdPath {
				t.Fatalf("readFile path = %q, want %q", path, passwdPath)
			}

			return []byte("root:x:0:0:root:/root:/bin/sh\nmoloch:x:501:20::/Users/moloch:/bin/zsh\n"), nil
		},
		501,
	)

	if got != "/bin/zsh" {
		t.Fatalf("defaultUnixShell() = %q, want %q", got, "/bin/zsh")
	}
}

func TestDefaultUnixShellFallsBackToBashThenSh(t *testing.T) {
	got := defaultUnixShell(
		func(string) string { return "" },
		fakeStat(map[string]bool{"/bin/bash": true}),
		func(string) ([]byte, error) { return nil, errors.New("not found") },
		501,
	)

	if got != "/bin/bash" {
		t.Fatalf("defaultUnixShell() = %q, want %q", got, "/bin/bash")
	}

	got = defaultUnixShell(
		func(string) string { return "" },
		fakeStat(map[string]bool{}),
		func(string) ([]byte, error) { return nil, errors.New("not found") },
		501,
	)

	if got != "/bin/sh" {
		t.Fatalf("defaultUnixShell() = %q, want %q", got, "/bin/sh")
	}
}

func TestDefaultWindowsShellPrefersPowerShell(t *testing.T) {
	got := defaultWindowsShell(func(path string) (string, error) {
		if path == "powershell.exe" {
			return path, nil
		}

		return "", os.ErrNotExist
	})

	if got != "powershell.exe" {
		t.Fatalf("defaultWindowsShell() = %q, want %q", got, "powershell.exe")
	}
}

func TestDefaultWindowsShellFallsBackToCmd(t *testing.T) {
	got := defaultWindowsShell(func(string) (string, error) {
		return "", os.ErrNotExist
	})

	if got != "cmd.exe" {
		t.Fatalf("defaultWindowsShell() = %q, want %q", got, "cmd.exe")
	}
}
