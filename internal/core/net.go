package core

import (
	"github.com/shirou/gopsutil/v3/net"		
	"github.com/spf13/cobra"
)

func GetNetworkStats() ([]net.IOCountersStat, error) {
	stats, err := net.IOCounters(true)
	if err != nil {
		return nil, err
	}
	return stats, nil
}

func GetNetworkInterfaces() ([]net.InterfaceStat, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	return interfaces, nil
}

func GetNetworkSentBytes() (uint64, error) {
	stats, err := net.IOCounters(true)
	if err != nil {
		return 0, err
	}
	var totalSent uint64
	for _, stat := range stats {
		totalSent += stat.BytesSent
	}
	return totalSent, nil
}

func GetNetworkReceivedBytes() (uint64, error) {
	stats, err := net.IOCounters(true)
	if err != nil {
		return 0, err
	}
	var totalReceived uint64
	for _, stat := range stats {
		totalReceived += stat.BytesRecv
	}
	return totalReceived, nil
}

func GetNetworkPacketsSent() (uint64, error) {
	stats, err := net.IOCounters(true)
	if err != nil {
		return 0, err
	}
	var totalPacketsSent uint64
	for _, stat := range stats {
		totalPacketsSent += stat.PacketsSent
	}
	return totalPacketsSent, nil
}

func GetNetworkPacketsReceived() (uint64, error) {
	stats, err := net.IOCounters(true)
	if err != nil {
		return 0, err
	}
	var totalPacketsReceived uint64
	for _, stat := range stats {
		totalPacketsReceived += stat.PacketsRecv
	}
	return totalPacketsReceived, nil
}

func GetNetworkErrorStats() ([]net.IOCountersStat, error) {
	stats, err := net.IOCounters(true)
	if err != nil {
		return nil, err
	}
	var errorStats []net.IOCountersStat
	for _, stat := range stats {
		if stat.Errin > 0 || stat.Errout > 0 {
			errorStats = append(errorStats, stat)
		}
	}
	return errorStats, nil
}

func GetNetworkInterfaceStats(name string) (*net.InterfaceStat, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, iface := range interfaces {
		if iface.Name == name {
			return &iface, nil
		}
	}
	return nil, nil
}

func GetNetworkPacketsDropped() (uint64, error) {
	stats, err := net.IOCounters(true)
	if err != nil {
		return 0, err
	}
	var totalDropped uint64
	for _, stat := range stats {
		totalDropped += stat.Dropin
		totalDropped += stat.Dropout
	}
	return totalDropped, nil
}

func GetNetworkTotalErrors() (uint64, error) {
	stats, err := net.IOCounters(true)
	if err != nil {
		return 0, err
	}
	var totalErrors uint64
	for _, stat := range stats {
		totalErrors += stat.Errin + stat.Errout
	}
	return totalErrors, nil
}
var NetCmd = &cobra.Command{
	Use:   "net",
	Short: "Display network statistics",
	Run: func(cmd *cobra.Command, args []string) {
		getNetworkTotalErrors, _ := cmd.Flags().GetBool("te")
		getNetworkPacketsDropped, _ := cmd.Flags().GetBool("pd")
		getNetworkInterfaceStats, _ := cmd.Flags().GetBool("is")
		getNetworkErrorStats, _ := cmd.Flags().GetBool("es")
		getNetworkPacketsReceived, _ := cmd.Flags().GetBool("pr")
		getNetworkPacketsSent, _ := cmd.Flags().GetBool("ps")
		getNetworkReceivedBytes, _ := cmd.Flags().GetBool("rb")
		getNetworkSentBytes, _ := cmd.Flags().GetBool("sb")
		getNetworkInterfaces, _ := cmd.Flags().GetBool("ni")
		getNetworkStats, _ := cmd.Flags().GetBool("ns")

		
		if getNetworkTotalErrors {
			errors, err := GetNetworkTotalErrors()
			if err == nil {
				println("Total Network Errors:", errors)
			}
		}
		if getNetworkPacketsDropped {
			dropped, err := GetNetworkPacketsDropped()
			if err == nil {
				println("Total Packets Dropped:", dropped)
			}
		}
		if getNetworkInterfaceStats {
			ifaces, err := GetNetworkInterfaces()
			if err == nil {
				for _, iface := range ifaces {
					println("Interface:", iface.Name)
				}
			}
		}
		if getNetworkErrorStats {
			stats, err := GetNetworkErrorStats()
			if err == nil {
				for _, s := range stats {
					println("Error Interface:", s.Name)
				}
			}
		}
		if getNetworkPacketsReceived {
			pr, err := GetNetworkPacketsReceived()
			if err == nil {
				println("Packets Received:", pr)
			}
		}
		if getNetworkPacketsSent {
			ps, err := GetNetworkPacketsSent()
			if err == nil {
				println("Packets Sent:", ps)
			}
		}
		if getNetworkReceivedBytes {
			rb, err := GetNetworkReceivedBytes()
			if err == nil {
				println("Bytes Received:", rb)
			}
		}
		if getNetworkSentBytes {
			sb, err := GetNetworkSentBytes()
			if err == nil {
				println("Bytes Sent:", sb)
			}
		}
		if getNetworkInterfaces {
			ifaces, err := GetNetworkInterfaces()
			if err == nil {
				for _, iface := range ifaces {
					println("Interface:", iface.Name)
				}
			}
		}
		if getNetworkStats {
			stats, err := GetNetworkStats()
			if err == nil {
				for _, s := range stats {
					println("Interface:", s.Name, "Bytes Sent:", s.BytesSent, "Bytes Recv:", s.BytesRecv)
				}
			}
		}
	},
}

func init() {
	NetCmd.Flags().Bool("te", false, "Show total network errors")
	NetCmd.Flags().Bool("pd", false, "Show packets dropped")
	NetCmd.Flags().Bool("is", false, "Show interface stats")
	NetCmd.Flags().Bool("es", false, "Show error stats per interface")
	NetCmd.Flags().Bool("pr", false, "Show packets received")
	NetCmd.Flags().Bool("ps", false, "Show packets sent")
	NetCmd.Flags().Bool("rb", false, "Show received bytes")
	NetCmd.Flags().Bool("sb", false, "Show sent bytes")
	NetCmd.Flags().Bool("ni", false, "Show network interfaces")
	NetCmd.Flags().Bool("ns", false, "Show raw network stats")
}
