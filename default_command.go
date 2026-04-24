package main

import (
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

const passwdPath = "/etc/passwd"

func defaultCommandArgs() []string {
	return []string{defaultCommand()}
}

func effectiveCommandArgs(args []string) ([]string, bool) {
	if len(args) > 0 {
		return args, false
	}

	return defaultCommandArgs(), true
}

func defaultCommand() string {
	if runtime.GOOS == "windows" {
		return defaultWindowsShell(exec.LookPath)
	}

	return defaultUnixShell(os.Getenv, os.Stat, os.ReadFile, os.Getuid())
}

func defaultWindowsShell(lookPath func(string) (string, error)) string {
	for _, shell := range []string{"powershell.exe", "powershell"} {
		if _, err := lookPath(shell); err == nil {
			return shell
		}
	}

	return "cmd.exe"
}

func defaultUnixShell(
	getenv func(string) string,
	stat func(string) (os.FileInfo, error),
	readFile func(string) ([]byte, error),
	uid int,
) string {
	if shell := getenv("SHELL"); shellExists(stat, shell) {
		return shell
	}

	if shell := passwdShell(readFile, uid); shellExists(stat, shell) {
		return shell
	}

	if shellExists(stat, "/bin/bash") {
		return "/bin/bash"
	}

	return "/bin/sh"
}

func shellExists(stat func(string) (os.FileInfo, error), shell string) bool {
	if shell == "" || !strings.HasPrefix(shell, "/") {
		return false
	}

	info, err := stat(shell)
	return err == nil && !info.IsDir()
}

func passwdShell(readFile func(string) ([]byte, error), uid int) string {
	data, err := readFile(passwdPath)
	if err != nil {
		return ""
	}

	for _, line := range strings.Split(string(data), "\n") {
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		fields := strings.Split(line, ":")
		if len(fields) < 7 {
			continue
		}

		entryUID, err := strconv.Atoi(fields[2])
		if err == nil && entryUID == uid {
			return fields[6]
		}
	}

	return ""
}
