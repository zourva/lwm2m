package objects

func GetOMAObjectDescriptors() []string {
	return []string{
		SecurityDescriptor,
		AccessControlDescriptor,
		DeviceDescriptor,
		ConnMonitorDescriptor,
		FirmwareUpdateDescriptor,
		LocationDescriptor,
		ConnStatsDescriptor,
	}
}
