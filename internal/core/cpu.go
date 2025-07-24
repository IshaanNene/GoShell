package core

import (
	"github.com/shirou/gopsutil/v3/cpu"
	"time"
	"fmt"
	"github.com/spf13/cobra"
)

func GetCPUUsage() (float64, error) {		
	percentages, err := cpu.Percent(0, false)
	if err != nil {
		return 0, err
	}
	return percentages[0], nil
}

func GetCPUCount() (int, error) {
	return cpu.Counts(true)
}

func GetCPUFrequency() ([]cpu.InfoStat, error) {
	return cpu.Info()
}

func GetCPUPercentages(interval float64) ([]float64, error) {
	return cpu.Percent(time.Duration(interval * float64(time.Second)), true)
}

func GetCPUNice() (float64, error) {
	cpuStat, err := cpu.Times(false)
	if err != nil {
		return 0, err
	}
	if len(cpuStat) == 0 {
		return 0, nil
	}
	return cpuStat[0].Nice, nil
}

func GetCPUStealTime() (float64, error) {
	cpuStat, err := cpu.Times(false)
	if err != nil {
		return 0, err
	}
	if len(cpuStat) == 0 {
		return 0, nil
	}
	return cpuStat[0].Steal, nil
}

func GetCPUUserTime() (float64, error) {
	cpuStat, err := cpu.Times(false)
	if err != nil {
		return 0, err
	}
	if len(cpuStat) == 0 {
		return 0, nil
	}
	return cpuStat[0].User, nil
}

func GetCPUSystemTime() (float64, error) {
	cpuStat, err := cpu.Times(false)
	if err != nil {
		return 0, err
	}
	if len(cpuStat) == 0 {
		return 0, nil
	}
	return cpuStat[0].System, nil
}

func GetCPUTemperature()(float64) {
	return 0.0
}

func GetCPUIdle() (float64, error) {
	percentages, err := cpu.Percent(0, false)
	if err != nil {
		return 0, err
	}
	return 100 - percentages[0], nil
}

func GetClocksPerSecond() (uint64) {
	return uint64(cpu.ClocksPerSec)
}

var CPUCmd = &cobra.Command{
	Use:   "cpu",
	Short: "Display CPU statistics",
	Run: func(cmd *cobra.Command, args []string) {
		showUsage, _ := cmd.Flags().GetBool("cu")
		showCount, _ := cmd.Flags().GetBool("cc")
		showFreq, _ := cmd.Flags().GetBool("cfreq")
		showIdle, _ := cmd.Flags().GetBool("cidle")
		showNice, _ := cmd.Flags().GetBool("cnice")
		showSteal, _ := cmd.Flags().GetBool("csteal")
		showUser, _ := cmd.Flags().GetBool("cuser")
		showSystem, _ := cmd.Flags().GetBool("csys")
		showTemp, _ := cmd.Flags().GetBool("ctemp")
		showClocks, _ := cmd.Flags().GetBool("cclk")

		if showUsage {
			usage, err := GetCPUUsage()
			if err == nil {
				fmt.Printf("CPU Usage: %.2f%%\n", usage)
			}
		}
		if showCount {
			count, err := GetCPUCount()
			if err == nil {
				fmt.Println("Logical CPU Count:", count)
			}
		}
		if showFreq {
			freqs, err := GetCPUFrequency()
			if err == nil {
				for _, info := range freqs {
					fmt.Printf("CPU: %s, Cores: %d, Speed: %.2f MHz\n", info.ModelName, info.Cores, info.Mhz)
				}
			}
		}
		if showIdle {
			idle, err := GetCPUIdle()
			if err == nil {
				fmt.Printf("CPU Idle: %.2f%%\n", idle)
			}
		}
		if showNice {
			nice, err := GetCPUNice()
			if err == nil {
				fmt.Printf("CPU Nice Time: %.2f\n", nice)
			}
		}
		if showSteal {
			steal, err := GetCPUStealTime()
			if err == nil {
				fmt.Printf("CPU Steal Time: %.2f\n", steal)
			}
		}
		if showUser {
			user, err := GetCPUUserTime()
			if err == nil {
				fmt.Printf("CPU User Time: %.2f\n", user)
			}
		}
		if showSystem {
			system, err := GetCPUSystemTime()
			if err == nil {
				fmt.Printf("CPU System Time: %.2f\n", system)
			}
		}
		if showTemp {
			temp := GetCPUTemperature()
			fmt.Printf("CPU Temperature: %.2f°C (dummy)\n", temp)
		}
		if showClocks {
			clocks := GetClocksPerSecond()
			fmt.Println("Clocks Per Second:", clocks)
		}
	},
}

func init() {
	CPUCmd.Flags().Bool("cu", false, "Show CPU usage percent")
	CPUCmd.Flags().Bool("cc", false, "Show logical CPU count")
	CPUCmd.Flags().Bool("cfreq", false, "Show CPU frequency info")
	CPUCmd.Flags().Bool("cidle", false, "Show CPU idle percent")
	CPUCmd.Flags().Bool("cnice", false, "Show CPU nice time")
	CPUCmd.Flags().Bool("csteal", false, "Show CPU steal time")
	CPUCmd.Flags().Bool("cuser", false, "Show CPU user time")
	CPUCmd.Flags().Bool("csys", false, "Show CPU system time")
	CPUCmd.Flags().Bool("ctemp", false, "Show CPU temperature (dummy)")
	CPUCmd.Flags().Bool("cclk", false, "Show CPU clocks per second")
}
