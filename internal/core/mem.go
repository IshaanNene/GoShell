package core

import (
	"github.com/shirou/gopsutil/v3/mem"
	"fmt"
	"github.com/spf13/cobra"
)

func GetMemoryUsage() (float64, uint64, error) {
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		return 0, 0, err
	}
	return vmStat.UsedPercent, vmStat.Used, nil
}

func GetMemoryTotal() (uint64, error) {
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		return 0, err
	}
	return vmStat.Total, nil
}

func GetMemorySwapTotal() (uint64, error) {
	swapStat, err := mem.SwapMemory()
	if err != nil {
		return 0, err
	}
	return swapStat.Total, nil
}

func GetMemorySwapUsed() (uint64, error) {
	swapStat, err := mem.SwapMemory()
	if err != nil {
		return 0, err
	}
	return swapStat.Used, nil
}

func GetMemoryFree() (uint64, error) {
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		return 0, err
	}
	return vmStat.Free, nil
}

func GetSwapTotal() (uint64, error) {
	swapStat, err := mem.SwapMemory()
	if err != nil {
		return 0, err
	}
	return swapStat.Total, nil
}

func GetSwapUsed() (uint64, error) {
	swapStat, err := mem.SwapMemory()
	if err != nil {
		return 0, err
	}
	return swapStat.Used, nil
}

func GetSwapFree() (uint64, error) {
	swapStat, err := mem.SwapMemory()
	if err != nil {
		return 0, err
	}
	return swapStat.Free, nil
}

func GetMemoryCached() (uint64, error) {
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		return 0, err
	}
	return vmStat.Cached, nil
}

func GetMemoryBuffers() (uint64, error) {
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		return 0, err
	}
	return vmStat.Buffers, nil
}

func GetMemorySwapFree() (uint64, error) {
	swapStat, err := mem.SwapMemory()
	if err != nil {
		return 0, err
	}
	return swapStat.Free, nil
}

func GetMemoryTotalSwap() (uint64, error) {
	swapStat, err := mem.SwapMemory()
	if err != nil {
		return 0, err
	}
	return swapStat.Total, nil
}

func GetMemoryTotalUsed() (uint64, error) {
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		return 0, err
	}
	return vmStat.Used, nil
}

func GetMemoryAvailable() (uint64, error) {
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		return 0, err
	}
	return vmStat.Available, nil
}


var MemCmd = &cobra.Command{
	Use:   "mem",
	Short: "Display memory statistics",
	Run: func(cmd *cobra.Command, args []string) {
		showUsage, _ := cmd.Flags().GetBool("mu")
		showTotal, _ := cmd.Flags().GetBool("mt")
		showFree, _ := cmd.Flags().GetBool("mf")
		showAvailable, _ := cmd.Flags().GetBool("ma")
		showCached, _ := cmd.Flags().GetBool("mc")
		showBuffers, _ := cmd.Flags().GetBool("mb")
		showUsed, _ := cmd.Flags().GetBool("mused")

		showSwapTotal, _ := cmd.Flags().GetBool("st")
		showSwapUsed, _ := cmd.Flags().GetBool("su")
		showSwapFree, _ := cmd.Flags().GetBool("sf")

		if showUsage {
			usedPercent, used, err := GetMemoryUsage()
			if err == nil {
				fmt.Printf("Memory Usage: %.2f%% (%d bytes)\n", usedPercent, used)
			}
		}
		if showTotal {
			total, err := GetMemoryTotal()
			if err == nil {
				fmt.Println("Total Memory:", total)
			}
		}
		if showFree {
			free, err := GetMemoryFree()
			if err == nil {
				fmt.Println("Free Memory:", free)
			}
		}
		if showAvailable {
			avail, err := GetMemoryAvailable()
			if err == nil {
				fmt.Println("Available Memory:", avail)
			}
		}
		if showCached {
			cached, err := GetMemoryCached()
			if err == nil {
				fmt.Println("Cached Memory:", cached)
			}
		}
		if showBuffers {
			buffers, err := GetMemoryBuffers()
			if err == nil {
				fmt.Println("Buffer Memory:", buffers)
			}
		}
		if showUsed {
			used, err := GetMemoryTotalUsed()
			if err == nil {
				fmt.Println("Used Memory:", used)
			}
		}
		if showSwapTotal {
			swapTotal, err := GetSwapTotal()
			if err == nil {
				fmt.Println("Swap Total:", swapTotal)
			}
		}
		if showSwapUsed {
			swapUsed, err := GetSwapUsed()
			if err == nil {
				fmt.Println("Swap Used:", swapUsed)
			}
		}
		if showSwapFree {
			swapFree, err := GetSwapFree()
			if err == nil {
				fmt.Println("Swap Free:", swapFree)
			}
		}
	},
}

func init() {
	
	MemCmd.Flags().Bool("mu", false, "Show memory usage (percent + used)")
	MemCmd.Flags().Bool("mt", false, "Show total memory")
	MemCmd.Flags().Bool("mf", false, "Show free memory")
	MemCmd.Flags().Bool("ma", false, "Show available memory")
	MemCmd.Flags().Bool("mc", false, "Show cached memory")
	MemCmd.Flags().Bool("mb", false, "Show buffer memory")
	MemCmd.Flags().Bool("mused", false, "Show total used memory")
	MemCmd.Flags().Bool("st", false, "Show swap total")
	MemCmd.Flags().Bool("su", false, "Show swap used")
	MemCmd.Flags().Bool("sf", false, "Show swap free")
}