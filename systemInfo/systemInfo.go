package systemInfo

import (
	"errors"
	"fmt"
	"log"
	"os/exec"
	"time"
)

func GetUptime(service string) (time.Duration, error) {
	// get uptime of the server or container
	switch service {
	case "server":
		return getServerUptime(), nil
	case "container":
		return getContainerUptime(), nil
	default:
		return 0, errors.New("service not found")
	}
}

func getServerUptime() time.Duration {
	// get the start time of the server and current time, then calculate the uptime
	start := GetServerStart()
	current := time.Now()
	uptime := current.Sub(start)
	return uptime
}

func getContainerUptime() time.Duration {
	// get the start time of the container and current time, then calculate the uptime
	start := GetContainerStart()
	current := time.Now()
	uptime := current.Sub(start)
	return uptime
}

func GetServerStart() time.Time {
	// get uptime of the server by executing uptime command
	out, err := exec.Command("sh", "./systemInfo/shell-scripts/uptimeServer.sh").Output()
	if err != nil {
		log.Fatal(err)
	}

	// remove the newline character from the output
	out = out[:len(out)-1]

	// Parse time looking like this: 2024-05-22 20:09:01 (the result of the stat command)
	tt, err := time.Parse("2006-01-02 15:04:05", string(out))
	return tt
}

func GetContainerStart() time.Time {
	// get uptime of the container by executing uptime command
	out, err := exec.Command("sh", "./systemInfo/shell-scripts/uptimeContainer.sh").Output()
	if err != nil {
		log.Fatal(err)
	}

	// remove the newline character from the output
	out = out[:len(out)-1]

	// Parse time looking like this: 2024-05-22 20:47:25.471086076 +0000 (the result of the docker command)
	tt, err := time.Parse("2006-01-02 15:04:05.999999999 -0700", string(out))
	return tt
}

func FormatDurationHumanReadable(d time.Duration) string {
	// format the duration to a human-readable format
	days := int(d.Hours()) / 24
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60
	return fmt.Sprintf("%d days, %d hours, %d minutes, %d seconds", days, hours, minutes, seconds)
}
