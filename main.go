package main

import (
	"brightnescli/brightnes"
	_ "embed"
	"fmt"
	"github.com/spf13/cobra"
	"log/slog"
	"os"
	"os/exec"
	"os/user"
	"strconv"
)

//go:embed 10-backlight.rules
var brightnessRules []byte

var ErrLogger = slog.New(slog.NewTextHandler(os.Stderr, nil))
var rootCmd = &cobra.Command{
	Use:   "brightnesscli",
	Short: "A CLI for controlling your computer's brightness",
	Long: `A CLI for controlling your computer's brightness. 
Requires root privileges to run, unless you install the udev rules.
To Install udev rules, run: 'sudo brightnesscli install-rules'`,
}

func main() {
	rootCmd.AddCommand(upCmd)
	rootCmd.AddCommand(downCmd)
	rootCmd.AddCommand(setBrightnessCmd)
	rootCmd.AddCommand(getCmd)
	installRulesCmd.Flags().StringP("user", "u", "", "The user to add to the 'video' group (default is the current user)")
	rootCmd.AddCommand(installRulesCmd)
	rootCmd.AddGroup(&cobra.Group{
		ID:    "global",
		Title: "Brightness Commands",
	})

	if err := rootCmd.Execute(); err != nil {
		ErrLogger.Error("Error executing command", "Error", err)
		os.Exit(1)
	}
}

var upCmd = &cobra.Command{
	Use:     "up",
	Short:   "Increase brightness by 5%",
	GroupID: "global",
	Run: func(cmd *cobra.Command, args []string) {
		brightness, err := brightnes.NewDevice()
		if err != nil {
			ErrLogger.Error("Error getting brightness device", "Error", err)
			os.Exit(1)
		}
		level, err := brightness.GetBrightness()
		if err != nil {
			ErrLogger.Error("Error getting brightness level", "Error", err)
			os.Exit(1)
		}
		if err := brightness.SetBrightness(level + 5); err != nil {
			ErrLogger.Error("Error setting brightness level", "Error", err)
			os.Exit(1)
		}
		fmt.Print(level + 5)
	},
}

var downCmd = &cobra.Command{
	Use:     "down",
	Short:   "Decrease brightness by 5%",
	GroupID: "global",
	Run: func(cmd *cobra.Command, args []string) {
		brightness, err := brightnes.NewDevice()
		if err != nil {
			ErrLogger.Error("Error getting brightness device", "Error", err)
			os.Exit(1)
		}
		level, err := brightness.GetBrightness()
		if err != nil {
			ErrLogger.Error("Error getting brightness level", "Error", err)
			os.Exit(1)
		}
		if err := brightness.SetBrightness(level - 5); err != nil {
			ErrLogger.Error("Error setting brightness level", "Error", err)
			os.Exit(1)
		}
		fmt.Print(level - 5)
	},
}

var setBrightnessCmd = &cobra.Command{
	Use:     "set [value]",
	Short:   "Set brightness to a specific value",
	GroupID: "global",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		value, err := strconv.Atoi(args[0])
		if err != nil {
			ErrLogger.Error("Error parsing brightness value", "Error", err)
			os.Exit(1)
		}
		brightness, err := brightnes.NewDevice()
		if err != nil {
			ErrLogger.Error("Error getting brightness device", "Error", err)
			os.Exit(1)
		}
		err = brightness.SetBrightness(value)
		if err != nil {
			ErrLogger.Error("Error setting brightness level", "Error", err)
			os.Exit(1)
		}
		fmt.Print(value)

	},
}

var getCmd = &cobra.Command{
	Use:     "get",
	Short:   "Get the current brightness",
	GroupID: "global",
	Run: func(cmd *cobra.Command, args []string) {
		brightness, err := brightnes.NewDevice()
		if err != nil {
			fmt.Println("Error getting brightness device")
			os.Exit(1)
		}
		level, err := brightness.GetBrightness()
		if err != nil {
			fmt.Println("Error getting brightness level")
			os.Exit(1)
		}
		fmt.Print(level)
	},
}
var installRulesCmd = &cobra.Command{
	Use:   "install-rules",
	Short: "Install brightness rules",
	Run: func(cmd *cobra.Command, args []string) {
		// Get the user flag
		userFlag, _ := cmd.Flags().GetString("user")

		// If no user flag is provided, get the current user, and run the command again with sudo
		if userFlag == "" {
			currentUser, err := user.Current()
			if err != nil {
				ErrLogger.Error("Error getting current user", "Error", err)
				os.Exit(1)
			}
			elevateCommand := exec.Command("sudo", os.Args[0], "install-rules", "--user", currentUser.Username)
			elevateCommand.Stdout = os.Stdout
			elevateCommand.Stderr = os.Stderr
			elevateCommand.Stdin = os.Stdin
			if err := elevateCommand.Run(); err != nil {
				ErrLogger.Error("Error elevating command", "Error", err)
				os.Exit(1)
			}
			return
		}

		rulesPath := "/etc/udev/rules.d/10-backlight.rules"

		// Write the rules to the file
		if err := os.WriteFile(rulesPath, brightnessRules, 0644); err != nil {
			ErrLogger.Error("Error writing rules file, make sure you are running with sudo", "Error", err)
			os.Exit(1)
		}

		// Reload the udev rules
		reloadCommand := exec.Command("udevadm", "control", "--reload-rules")
		if err := reloadCommand.Run(); err != nil {
			ErrLogger.Error("Error reloading udev rules", "Error", err)
			os.Exit(1)
		}
		triggerCommand := exec.Command("sudo", "udevadm", "trigger")
		if err := triggerCommand.Run(); err != nil {
			ErrLogger.Error("Error triggering udev rules", "Error", err)
			os.Exit(1)
		}

		// Prepare the command
		usermodCommand := exec.Command("sudo", "usermod", "-a", "-G", "video", userFlag)
		// Run the command
		if err := usermodCommand.Run(); err != nil {
			ErrLogger.Error("Error adding user to video group", "Error", err)
			os.Exit(1)
		}

		fmt.Println("Brightness rules installed successfully, Path:", rulesPath)
		fmt.Println("You've been added to the 'video' group. Please log out and back in, or start a new shell with 'newgrp video' or 'su - $USER', to use this group.")
	},
}
