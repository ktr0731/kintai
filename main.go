package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

const (
	macWiFiCmd       = "/System/Library/PrivateFrameworks/Apple80211.framework/Versions/Current/Resources/airport"
	macWiFiCmdOpt    = "-I"
	macWiFiCmdPrefix = "SSID: "
)

func main() {
	os.Exit(run())
}

func run() int {
	logger := log.New(os.Stdout, "[kintai] ", log.LstdFlags)
	if runtime.GOOS != "darwin" {
		logger.Fatalln("OSs other than macOS are not supported yet")
	}

	if len(os.Args) != 3 {
		logger.Fatalln("Usage: kintai <SSID> <API KEY>")
	}

	client := NewClient(logger, os.Args[1], os.Args[2], 3*time.Second)
	if err := client.Start(); err != nil {
		panic(err)
	}

	return 0
}

func GetSSID() (string, error) {
	buf := new(bytes.Buffer)
	cmd := exec.Command(macWiFiCmd, macWiFiCmdOpt)
	cmd.Stdout = buf
	if err := cmd.Run(); err != nil {
		return "", err
	}

	s := bufio.NewScanner(buf)
	for s.Scan() {
		txt := strings.TrimSpace(s.Text())
		if strings.HasPrefix(txt, macWiFiCmdPrefix) {
			return strings.TrimPrefix(txt, macWiFiCmdPrefix), nil
		}
	}
	return "", fmt.Errorf("cannot resolve SSID by %s", macWiFiCmd)
}
