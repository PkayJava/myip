package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

func init() {
	log.SetPrefix("[myip] : ")
	log.SetOutput(os.Stdout)
}
func main() {
	var _mode string
	var _port int

	if runtime.GOOS != "linux" {
		log.Println(runtime.GOOS + " is not supported")
	}

	flag.StringVar(&_mode, "mode", "", "Required, server, systemd-install, systemd-remove")
	flag.IntVar(&_port, "port", 80, "Optional, default is port 80")

	flag.Parse()

	cli, _ := filepath.Abs(os.Args[0])

	if _mode == "systemd-install" {
		lines := []string{
			"[Unit]",
			"Description=myip",
			"After=network-online.target",
			"",
			"[Service]",
			"Type=simple",
			"Restart=always",
			"RestartSec=1",
			"User=root",
			"WorkingDirectory=/tmp",
			fmt.Sprintf("ExecStart=%s -mode server -port %d", cli, _port),
			"StartLimitInterval=0",
			"",
			"[Install]",
			"WantedBy=multi-user.target",
		}

		file, errorOpenFile := os.OpenFile("/etc/systemd/system/myip.service", os.O_CREATE|os.O_WRONLY, 0755)

		if errorOpenFile != nil {
			log.Println(errorOpenFile.Error())
			return
		}

		datawriter := bufio.NewWriter(file)

		for _, data := range lines {
			_, _ = datawriter.WriteString(fmt.Sprintln(data))
		}

		_ = datawriter.Flush()
		_ = file.Close()

		{
			cmd := exec.Command("/usr/bin/systemctl", "daemon-reload")
			_, _ = cmd.CombinedOutput()
		}
		{
			cmd := exec.Command("/usr/bin/systemctl", "enable", "myip")
			_, _ = cmd.CombinedOutput()
		}
		{
			cmd := exec.Command("/usr/sbin/service", "myip", "start")
			_, _ = cmd.CombinedOutput()
		}
	} else if _mode == "server" {
		app := fiber.New(fiber.Config{
			Prefork:       true,
			StrictRouting: true,
			CaseSensitive: false,
			ServerHeader:  "Fiber",
			AppName:       "MyIp",
		})

		app.Get("/", MyIp)

		_ = app.Listen(fmt.Sprintf(":%d", _port))
	} else if _mode == "systemd-remove" {
		{
			cmd := exec.Command("/usr/sbin/service", "myip", "stop")
			_, _ = cmd.CombinedOutput()
		}
		{
			cmd := exec.Command("/usr/bin/systemctl", "disable", "myip")
			_, _ = cmd.CombinedOutput()
		}
		_ = os.Remove("/etc/systemd/system/myip.service")
		{
			cmd := exec.Command("/usr/bin/systemctl", "daemon-reload")
			_, _ = cmd.CombinedOutput()
		}
	} else {
		flag.PrintDefaults()
		return
	}
}

func MyIp(c *fiber.Ctx) error {
	client := &http.Client{}
	request, _ := http.NewRequest("GET", "https://ip.me", nil)
	request.Header.Set("User-Agent", "curl/7.54")

	response, errorGet := client.Do(request)
	if errorGet != nil {
		log.Println(errorGet.Error())
		return c.SendString(errorGet.Error())
	}

	body, errorReadAll := io.ReadAll(response.Body)
	if errorReadAll != nil {
		log.Println(errorReadAll.Error())
		return c.SendString(errorGet.Error())
	}

	ip := string(body)

	return c.SendString(ip)
}
