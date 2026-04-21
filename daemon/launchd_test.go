//go:build darwin

package daemon

import (
	"fmt"
	"strings"
	"testing"
)

func TestBuildPlist_KeepAliveDoesNotRestartOnCleanExit(t *testing.T) {
	cfg := Config{
		BinaryPath: "/opt/cc-connect/cc-connect",
		WorkDir:    "/tmp/wd",
		LogFile:    "/tmp/log",
		LogMaxSize: 10485760,
		EnvPATH:    "/usr/bin",
	}
	xml := buildPlist(cfg)
	if !strings.Contains(xml, "<key>SuccessfulExit</key>") {
		t.Fatal("plist should use KeepAlive dict with SuccessfulExit so exit 0 does not respawn")
	}
	// Boolean KeepAlive causes launchd to restart after every exit, including SIGTERM shutdown.
	if strings.Contains(xml, "<key>KeepAlive</key>\n\t<true/>") {
		t.Fatal("plist must not use boolean KeepAlive true")
	}
	if !strings.Contains(xml, "<key>LimitLoadToSessionType</key>") ||
		!strings.Contains(xml, "<string>Aqua</string>") ||
		!strings.Contains(xml, "<string>Background</string>") {
		t.Fatal("plist should allow both Aqua and Background sessions")
	}
}

func TestPreferredLaunchdDomainFallsBackToUserWhenGUIDomainUnavailable(t *testing.T) {
	orig := runLaunchctl
	t.Cleanup(func() { runLaunchctl = orig })

	guiDomain := launchdGUIDomain()
	userDomain := launchdUserDomain()
	runLaunchctl = func(args ...string) (string, error) {
		if len(args) >= 2 && args[0] == "print" && args[1] == guiDomain {
			return "Bootstrap failed: 125: Domain does not support specified action", fmt.Errorf("exit status 125")
		}
		if len(args) >= 2 && args[0] == "print" && args[1] == userDomain {
			return "subsystem", nil
		}
		return "", nil
	}

	if got := preferredLaunchdDomain(); got != userDomain {
		t.Fatalf("preferredLaunchdDomain() = %q, want %q", got, userDomain)
	}
}

func TestLaunchdStatusUsesUserDomainWhenGUIDomainUnavailable(t *testing.T) {
	orig := runLaunchctl
	t.Cleanup(func() { runLaunchctl = orig })

	guiDomain := launchdGUIDomain()
	userDomain := launchdUserDomain()
	guiTarget := launchdTarget(guiDomain)
	userTarget := launchdTarget(userDomain)
	runLaunchctl = func(args ...string) (string, error) {
		if len(args) < 2 || args[0] != "print" {
			return "", nil
		}
		switch args[1] {
		case guiDomain, guiTarget:
			return "Bootstrap failed: 125: Domain does not support specified action", fmt.Errorf("exit status 125")
		case userDomain:
			return "subsystem", nil
		case userTarget:
			return "pid = 4321\nstate = running", nil
		default:
			return "", fmt.Errorf("unexpected target %q", args[1])
		}
	}

	mgr := &launchdManager{}
	st, err := mgr.Status()
	if err != nil {
		t.Fatalf("Status() error = %v", err)
	}
	if !st.Running {
		t.Fatal("Status().Running = false, want true")
	}
	if st.PID != 4321 {
		t.Fatalf("Status().PID = %d, want 4321", st.PID)
	}
}
