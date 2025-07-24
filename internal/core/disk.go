package core

import (
	"github.com/shirou/gopsutil/v3/disk"
	"time"
	"fmt"
	"github.com/spf13/cobra"
)

func GetDiskUsage() (float64, uint64, error) {
	diskStat, err := disk.Usage("/")
	if err != nil {
		return 0, 0, err
	}
	return diskStat.UsedPercent, diskStat.Used, nil
}

func GetDiskTotal() (uint64, error) {
	diskStat, err := disk.Usage("/")
	if err != nil {
		return 0, err
	}
	return diskStat.Total, nil
}

func GetDiskReadWriteStats() (map[string]disk.IOCountersStat, error) {
	stats, err := disk.IOCounters()
	if err != nil {
		return nil, err
	}
	return stats, nil
}

func GetDiskInodes() (uint64, uint64, error) {
	inodeStat, err := disk.Usage("/")
	if err != nil {
		return 0, 0, err
	}
	return inodeStat.InodesUsed, inodeStat.InodesFree, nil
}

func GetDiskReadBytes() (uint64, error) {
	stats, err := disk.IOCounters()
	if err != nil {
		return 0, err
	}
	var totalRead uint64
	for _, stat := range stats {
		totalRead += stat.ReadBytes
	}
	return totalRead, nil
}

func GetDiskWriteBytes() (uint64, error) {
	stats, err := disk.IOCounters()
	if err != nil {
		return 0, err
	}
	var totalWrite uint64
	for _, stat := range stats {
		totalWrite += stat.WriteBytes
	}
	return totalWrite, nil
}

func GetDiskFree() (uint64, error) {
	diskStat, err := disk.Usage("/")
	if err != nil {
		return 0, err
	}
	return diskStat.Free, nil
}

func GetDiskInodesUsed() (uint64, error) {
	inodeStat, err := disk.Usage("/")
	if err != nil {
		return 0, err
	}
	return inodeStat.InodesUsed, nil
}

func GetDiskInodesFree() (uint64, error) {
	inodeStat, err := disk.Usage("/")
	if err != nil {
		return 0, err
	}
	return inodeStat.InodesFree, nil
}

func GetDiskIOStats() (map[string]disk.IOCountersStat, error) {
	return disk.IOCounters()
}

func GetDiskWriteBytesPerSecond() (uint64, error) {
	stats, err := disk.IOCounters()
	if err != nil {
		return 0, err
	}
	var totalWrite uint64
	for _, stat := range stats {
		totalWrite += stat.WriteBytes
	}
	return totalWrite / uint64(time.Now().Unix()), nil
}

func GetDiskReadBytesPerSecond() (uint64, error) {
	stats, err := disk.IOCounters()
	if err != nil {
		return 0, err
	}
	var totalRead uint64
	for _, stat := range stats {
		totalRead += stat.ReadBytes
	}
	return totalRead / uint64(time.Now().Unix()), nil 
}

func GetDiskUsageByPath(path string) (float64, uint64, error) {
	diskStat, err := disk.Usage(path)
	if err != nil {
		return 0, 0, err
	}
	return diskStat.UsedPercent, diskStat.Used, nil
}

func GetDiskInodeUsage() (uint64, uint64, error) {
	inodeStat, err := disk.Usage("/")
	if err != nil {
		return 0, 0, err
	}
	return inodeStat.InodesUsed, inodeStat.InodesFree, nil
}

func GetDiskTotalInodes() (uint64, error) {
	inodeStat, err := disk.Usage("/")
	if err != nil {
		return 0, err
	}
	return inodeStat.InodesTotal, nil
}

func GetDiskTotalUsed() (uint64, error) {
	diskStat, err := disk.Usage("/")
	if err != nil {
		return 0, err
	}
	return diskStat.Used, nil
}

func GetDiskSerialNumber() (string, error) {
	stats, err := disk.IOCounters()
	if err != nil {
		return "", err
	}
	for _, stat := range stats {
		if stat.SerialNumber != "" {
			return stat.SerialNumber, nil
		}
	}
	return "", nil
}

var DiskCmd = &cobra.Command{
	Use:   "disk",
	Short: "Display disk statistics",
	Run: func(cmd *cobra.Command, args []string) {
		showUsage, _ := cmd.Flags().GetBool("du")
		showTotal, _ := cmd.Flags().GetBool("dt")
		showFree, _ := cmd.Flags().GetBool("df")
		showUsed, _ := cmd.Flags().GetBool("dused")

		showReadBytes, _ := cmd.Flags().GetBool("drb")
		showWriteBytes, _ := cmd.Flags().GetBool("dwb")
		showReadBps, _ := cmd.Flags().GetBool("drbps")
		showWriteBps, _ := cmd.Flags().GetBool("dwbps")

		showInodesUsed, _ := cmd.Flags().GetBool("diu")
		showInodesFree, _ := cmd.Flags().GetBool("dif")
		showInodesTotal, _ := cmd.Flags().GetBool("dit")
		showSerial, _ := cmd.Flags().GetBool("dsn")

		if showUsage {
			usedPercent, used, err := GetDiskUsage()
			if err == nil {
				fmt.Printf("Disk Usage: %.2f%% (%d bytes)\n", usedPercent, used)
			}
		}
		if showTotal {
			total, err := GetDiskTotal()
			if err == nil {
				fmt.Println("Disk Total:", total)
			}
		}
		if showFree {
			free, err := GetDiskFree()
			if err == nil {
				fmt.Println("Disk Free:", free)
			}
		}
		if showUsed {
			used, err := GetDiskTotalUsed()
			if err == nil {
				fmt.Println("Disk Used:", used)
			}
		}
		if showReadBytes {
			rb, err := GetDiskReadBytes()
			if err == nil {
				fmt.Println("Disk Read Bytes:", rb)
			}
		}
		if showWriteBytes {
			wb, err := GetDiskWriteBytes()
			if err == nil {
				fmt.Println("Disk Write Bytes:", wb)
			}
		}
		if showReadBps {
			rbps, err := GetDiskReadBytesPerSecond()
			if err == nil {
				fmt.Println("Disk Read Bytes/sec:", rbps)
			}
		}
		if showWriteBps {
			wbps, err := GetDiskWriteBytesPerSecond()
			if err == nil {
				fmt.Println("Disk Write Bytes/sec:", wbps)
			}
		}
		if showInodesUsed {
			iu, err := GetDiskInodesUsed()
			if err == nil {
				fmt.Println("Disk Inodes Used:", iu)
			}
		}
		if showInodesFree {
			ifree, err := GetDiskInodesFree()
			if err == nil {
				fmt.Println("Disk Inodes Free:", ifree)
			}
		}
		if showInodesTotal {
			itotal, err := GetDiskTotalInodes()
			if err == nil {
				fmt.Println("Disk Inodes Total:", itotal)
			}
		}
		if showSerial {
			serial, err := GetDiskSerialNumber()
			if err == nil && serial != "" {
				fmt.Println("Disk Serial Number:", serial)
			} else {
				fmt.Println("Disk Serial Number: Not available")
			}
		}
	},
}

func init() {
	
	DiskCmd.Flags().Bool("du", false, "Show disk usage (percent + used)")
	DiskCmd.Flags().Bool("dt", false, "Show total disk space")
	DiskCmd.Flags().Bool("df", false, "Show free disk space")
	DiskCmd.Flags().Bool("dused", false, "Show total used disk space")

	
	DiskCmd.Flags().Bool("drb", false, "Show total read bytes")
	DiskCmd.Flags().Bool("dwb", false, "Show total write bytes")
	DiskCmd.Flags().Bool("drbps", false, "Show read bytes per second")
	DiskCmd.Flags().Bool("dwbps", false, "Show write bytes per second")

	
	DiskCmd.Flags().Bool("diu", false, "Show inodes used")
	DiskCmd.Flags().Bool("dif", false, "Show inodes free")
	DiskCmd.Flags().Bool("dit", false, "Show total inodes")

	
	DiskCmd.Flags().Bool("dsn", false, "Show disk serial number")
}
