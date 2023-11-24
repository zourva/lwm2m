package server

// DeviceControlDelegator implements the server-side operations
// for Device Management and Service Enablement Interface.
type DeviceControlDelegator struct {
	server *LwM2MServer //server context
}

func NewDeviceControlServerDelegator(server *LwM2MServer) *DeviceControlDelegator {
	return &DeviceControlDelegator{
		server: server,
	}
}
