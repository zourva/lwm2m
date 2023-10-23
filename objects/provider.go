package objects

func GetOMAObjectDescriptors() []string {
	return []string{
		SecurityDescriptor,
		ServerDescriptor,
		AccessControlDescriptor,
		DeviceDescriptor,
		ConnMonitorDescriptor,
		FirmwareUpdateDescriptor,
		LocationDescriptor,
		ConnStatsDescriptor,
	}
}
