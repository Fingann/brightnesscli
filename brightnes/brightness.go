package brightnes

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	path    = "/sys/class"
	classes = "backlight"
)

type Device struct {
	class          string
	id             string
	currBrightness int
	maxBrightness  int
	brightnessPath string
}

func New(iface string) (*Device, error) {
	maxBrightnessPath := filepath.Join(path, classes, iface, "max_brightness")
	maxBrightnessStr, err := os.ReadFile(maxBrightnessPath)
	if err != nil {
		return nil, fmt.Errorf("error reading max brightness file: %w", err)
	}

	maxBrightness, err := strconv.Atoi(strings.TrimSpace(string(maxBrightnessStr)))
	if err != nil {
		return nil, fmt.Errorf("error converting max brightness to int: %w", err)
	}
	return &Device{
		class:          classes,
		id:             iface,
		brightnessPath: filepath.Join(path, classes, iface, "brightness"),
		maxBrightness:  maxBrightness,
	}, nil
}

func NewDevice() (*Device, error) {
	// List of possible backlight interfaces
	interfaces := []string{"intel_backlight", "acpi_video0", "acpi_video1", "radeon_bl0"}

	for _, iface := range interfaces {
		brightnessPath := filepath.Join(path, classes, iface, "brightness")
		if _, err := os.Stat(brightnessPath); err == nil {
			return New(iface)
		}
	}
	return nil, fmt.Errorf("no backlight interface found")
}

func (d *Device) SetBrightness(level int) error {
	level = max(0, min(100, level)) // clamp between 0 and 100

	brightness := (level * d.maxBrightness) / 100
	brightnessStr := strconv.Itoa(brightness)

	return os.WriteFile(d.brightnessPath, []byte(brightnessStr), 0644)
}

func (d *Device) IncreaseBrightness(level int) error {

	level = max(0, min(100, level)) // clamp between 0 and 100

	brightness := (level * d.maxBrightness) / 100
	brightnessStr := strconv.Itoa(brightness)

	return os.WriteFile(d.brightnessPath, []byte(brightnessStr), 0644)
}

func (d *Device) GetBrightness() (int, error) {
	brightnessStr, err := os.ReadFile(d.brightnessPath)
	if err != nil {
		return -1, fmt.Errorf("error reading brightness file: %w", err)
	}

	currBrightness, err := strconv.Atoi(strings.TrimSpace(string(brightnessStr)))
	if err != nil {
		return -1, fmt.Errorf("error converting brightness to int: %w", err)
	}

	return (currBrightness * 100) / d.maxBrightness, nil

}
