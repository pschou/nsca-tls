package main

import (
	"time"

	"github.com/coreos/go-systemd/daemon"
)

func systemdStarted()  { daemon.SdNotify(false, daemon.SdNotifyReady) }
func systemdStopping() { daemon.SdNotify(false, daemon.SdNotifyStopping) }

func systemdWatchdog() {
	interval, err := daemon.SdWatchdogEnabled(false)
	if err != nil || interval == 0 {
		return
	}
	for {
		daemon.SdNotify(false, daemon.SdNotifyWatchdog)
		time.Sleep(interval / 3)
	}
}
