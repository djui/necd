package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"text/template"

	"github.com/kardianos/osext"
)

// Daemonize uses Launchd to keep this program running in the background as
// daemon/agent.
//
// See: http://www.goinggo.net/2013/06/running-go-programs-as-background.html
func Daemonize(name string) error {
	// See: http://stackoverflow.com/questions/1023306/finding-current-executables-path-without-proc-self-exe/1024937#1024937
	execPath, err := osext.Executable()
	if err != nil {
		return fmt.Errorf("Can't get executable path: %v", err)
	}

	usr, err := user.Current()
	if err != nil {
		return fmt.Errorf("Can't get user: %v", err)
	}

	logPath := filepath.Join(usr.HomeDir, "Library", "Logs", name+".log")

	agent := LaunchdAgent{
		Name:     name,
		ExecPath: execPath,
		LogPath:  logPath,
	}

	if err := agent.Install(); err != nil {
		return err
	}

	return agent.Load()
}

// LaunchdAgent describes a Launchd agent.
type LaunchdAgent struct {
	ExecPath string
	LogPath  string
	Name     string
	path     string
}

// Install write the Launchd agent configuration file to ~/Library/LaunchAgents.
func (l *LaunchdAgent) Install() error {
	usr, err := user.Current()
	if err != nil {
		return fmt.Errorf("Can't get user: %v", err)
	}

	l.path = filepath.Join(usr.HomeDir, "Library", "LaunchAgents", l.Name+".plist")
	f, err := os.Create(l.path)
	if err != nil {
		return fmt.Errorf("Can't create launchd agent config file: %v", err)
	}
	defer f.Close()

	tmpl, err := template.New("launchdAgent").Parse(launchdAgentTemplate)
	if err != nil {
		return fmt.Errorf("Can't parse launchd agent template: %v", err)
	}

	if err := tmpl.Execute(f, l); err != nil {
		return fmt.Errorf("Can't write launchd agent template: %v", err)
	}

	return nil
}

// Load loads a Launchd agent.
func (l *LaunchdAgent) Load() error {
	cmd := exec.Command("launchctl", "load", l.path)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("Can't launch agent: %v", err)
	}

	return nil
}

// Unload unloads a Launchd agent.
func (l *LaunchdAgent) Unload() error {
	cmd := exec.Command("launchctl", "unload", l.path)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("Can't launch agent: %v", err)
	}

	return nil
}

const launchdAgentTemplate = `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple Computer//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd" >
<plist version="1.0">
  <dict>
    <key>Label</key>
    <string>{{.Name}}</string>
    <key>ProgramArguments</key>
    <array>
      <string>{{.ExecPath}}</string>
    </array>
    <key>KeepAlive</key>
    <true/>
    <key>RunAtLoad</key>
    <false/>
    <key>Disabled</key>
    <false/>
    <key>StandardOutPath</key>
    <string>{{.LogPath}}</string>
    <key>StandardErrorPath</key>
    <string>{{.LogPath}}</string>
  </dict>
</plist>
`
