package core

import "strconv"

type ResourceID = uint16

// LwM2MSecurity resources
const (
	LwM2MSecurityLwM2MServerURI                   ResourceID = 0
	LwM2MSecurityBootstrapServer                  ResourceID = 1
	LwM2MSecuritySecurityMode                     ResourceID = 2
	LwM2MSecurityPublicKeyOrIdentity              ResourceID = 3
	LwM2MSecurityServerPublicKeyOrIdentity        ResourceID = 4
	LwM2MSecuritySecretKey                        ResourceID = 5
	LwM2MSecuritySMSSecurityMode                  ResourceID = 6
	LwM2MSecuritySMSBindingKeyParameters          ResourceID = 7
	LwM2MSecuritySMSBindingSecretKeys             ResourceID = 8
	LwM2MSecurityLwM2MServerSMSNumber             ResourceID = 9
	LwM2MSecurityShortServerID                    ResourceID = 10
	LwM2MSecurityClientHoldOffTime                ResourceID = 11
	LwM2mSecurityBootstrapServerAccountTimeout    ResourceID = 12
	LwM2mSecurityMatchingType                     ResourceID = 13
	LwM2mSecuritySNI                              ResourceID = 14
	LwM2mSecurityCertificateUsage                 ResourceID = 15
	LwM2mSecurityDTLSTLSCipherSuite               ResourceID = 16
	LwM2mSecurityOSCORESecurityMode               ResourceID = 17
	LwM2mSecurityGroupsToUse                      ResourceID = 18
	LwM2mSecuritySignatureAlgorithmSupported      ResourceID = 19
	LwM2mSecuritySignatureAlgorithmToUse          ResourceID = 20
	LwM2mSecuritySignatureAlgorithmCertsSupported ResourceID = 21
	LwM2mSecurityTLS13FeaturesToUse               ResourceID = 22
	LwM2mSecurityTLSExtensionsSupported           ResourceID = 23
	LwM2mSecurityTLSExtensionsToUse               ResourceID = 24
	LwM2mSecuritySecondaryLwM2MServerURI          ResourceID = 25
	LwM2mSecurityMQTTServer                       ResourceID = 26
	LwM2mSecurityLwM2MCOSESecurity                ResourceID = 27
	LwM2mSecurityRDSDestinationPort               ResourceID = 28
	LwM2mSecurityRDSSourcePort                    ResourceID = 29
	LwM2mSecurityRDSApplicationID                 ResourceID = 30
)

const (
	LwM2MServerShortServerID                            ResourceID = 0
	LwM2MServerLifetime                                 ResourceID = 1
	LwM2MServerDefaultMinimumPeriod                     ResourceID = 2
	LwM2MServerDefaultMaximumPeriod                     ResourceID = 3
	LwM2MServerDisable                                  ResourceID = 4
	LwM2MServerDisableTimeout                           ResourceID = 5
	LwM2MServerNotificationStoringWhenDisabledOrOffline ResourceID = 6
	LwM2MServerBinding                                  ResourceID = 7
	LwM2MServerRegistrationUpdateTrigger                ResourceID = 8
	LwM2MServerBootstrapRequestTrigger                  ResourceID = 9
	LwM2MServerAPNLink                                  ResourceID = 10
	LwM2MServerTLSDTLSAlertCode                         ResourceID = 11
	LwM2MServerLastBootstrapped                         ResourceID = 12
	LwM2MServerRegistrationPriorityOrder                ResourceID = 13
	LwM2MServerInitialRegistrationDelayTimer            ResourceID = 14
	LwM2MServerRegistrationFailureBlock                 ResourceID = 15
	LwM2MServerBootstrapOnRegistrationFailure           ResourceID = 16
	LwM2MServerCommunicationRetryCount                  ResourceID = 17
	LwM2MServerCommunicationRetryTimer                  ResourceID = 18
	LwM2MServerCommunicationSequenceDelayTimer          ResourceID = 19
	LwM2MServerCommunicationSequenceRetryCount          ResourceID = 20
	LwM2MServerTrigger                                  ResourceID = 21
	LwM2MServerPreferredTransport                       ResourceID = 22
	LwM2MServerMuteSend                                 ResourceID = 23
	LwM2MAlternateAPNLinks                              ResourceID = 24
	LwM2MSupportedServerVersions                        ResourceID = 25
	LwM2MDefaultNotificationMode                        ResourceID = 26
	LwM2MProfileIDHashAlgorithm                         ResourceID = 27
)

// AccessControl resources
const (
	LwM2MAccessControlObjectID           ResourceID = 0
	LwM2MAccessControlObjectInstanceID   ResourceID = 1
	LwM2MAccessControlACL                ResourceID = 2
	LwM2MAccessControlAccessControlOwner ResourceID = 3
)

// Device resources
const (
	DeviceManufacturer          ResourceID = 0
	DeviceModelNumber           ResourceID = 1
	DeviceSerialNumber          ResourceID = 2
	DeviceFirmwareVersion       ResourceID = 3
	DeviceReboot                ResourceID = 4
	DeviceFactoryReset          ResourceID = 5
	DeviceAvailablePowerSources ResourceID = 6
	DevicePowerSourceVoltage    ResourceID = 7
	DevicePowerSourceCurrent    ResourceID = 8
	DeviceBatteryLevel          ResourceID = 9
	DeviceMemoryFree            ResourceID = 10
	DeviceErrorCode             ResourceID = 11
	DeviceResetErrorCode        ResourceID = 12
	DeviceCurrentTime           ResourceID = 13
	DeviceUTCOffset             ResourceID = 14
	DeviceTimezone              ResourceID = 15
	DeviceSupportedBindingModes ResourceID = 16
)

// ConnectivityMonitoring resources
const (
	ConnectivityMonitoringNetworkBearer            ResourceID = 0
	ConnectivityMonitoringAvailableNetworkBearer   ResourceID = 1
	ConnectivityMonitoringRadioSignalStrength      ResourceID = 2
	ConnectivityMonitoringLinkQuality              ResourceID = 3
	ConnectivityMonitoringIPAddresses              ResourceID = 4
	ConnectivityMonitoringRouterIPAddresses        ResourceID = 5
	ConnectivityMonitoringLinkUtilization          ResourceID = 6
	ConnectivityMonitoringAPN                      ResourceID = 7
	ConnectivityMonitoringCellID                   ResourceID = 8
	ConnectivityMonitoringSMNC                     ResourceID = 9
	ConnectivityMonitoringSMCC                     ResourceID = 10
	ConnectivityMonitoringSignalSNR                ResourceID = 11
	ConnectivityMonitoringLAC                      ResourceID = 12
	ConnectivityMonitoringCoverageEnhancementLevel ResourceID = 13
)

// FirmwareUpdate resources
const (
	FirmwareUpdatePackage                ResourceID = 0
	FirmwareUpdatePackageURI             ResourceID = 1
	FirmwareUpdateUpdate                 ResourceID = 2
	FirmwareUpdateState                  ResourceID = 3
	FirmwareUpdateUpdateSupportedObjects ResourceID = 4
	FirmwareUpdateUpdateResult           ResourceID = 5
	FirmwareUpdatePkgName                ResourceID = 6
	FirmwareUpdatePkgVersion             ResourceID = 7
	FirmwareUpdateProtocolSupport        ResourceID = 8
	FirmwareUpdateDeliveryMethod         ResourceID = 9
	FirmwareUpdateCancel                 ResourceID = 10
	FirmwareUpdateSeverity               ResourceID = 11
	FirmwareUpdateLastStateChangeTime    ResourceID = 12
	FirmwareUpdateMaximumDeferPeriod     ResourceID = 13
)

// Location resources
const (
	LocationLatitude  ResourceID = 0 //float, WGS-84
	LocationLongitude ResourceID = 1 //float, WGS-84
	LocationAltitude  ResourceID = 2 //float, m
	LocationRadius    ResourceID = 3 //float, m
	LocationVelocity  ResourceID = 4 //opaque
	LocationTimestamp ResourceID = 5 //time
	LocationSpeed     ResourceID = 6 //float, m/s
)

// ConnectivityStatistics resources
const (
	ConnectivityStatisticsSMSTxCounter       ResourceID = 0
	ConnectivityStatisticsSMSRxCounter       ResourceID = 1
	ConnectivityStatisticsTxData             ResourceID = 2
	ConnectivityStatisticsRxData             ResourceID = 3
	ConnectivityStatisticsMaxMessageSize     ResourceID = 4
	ConnectivityStatisticsAverageMessageSize ResourceID = 5
	ConnectivityStatisticsStart              ResourceID = 6
	ConnectivityStatisticsStop               ResourceID = 7
	ConnectivityStatisticsCollectionPeriod   ResourceID = 8
)

const (
	SecurityModePreSharedKey = 0
	SecurityModeRawPublicKey = 1
	SecurityModeCertificate  = 2
	SecurityModeNoSec        = 3
)

const (
	FirmwareUpdateStateIdle        = 1
	FirmwareUpdateStateDownloading = 2
	FirmwareUpdateStateDownloaded  = 3

	FirmwareUpdateResultDefault                = 0
	FirmwareUpdateResultSuccessful             = 1
	FirmwareUpdateResultNotEnoughStorage       = 2
	FirmwareUpdateResultOutOfMemory            = 3
	FirmwareUpdateResultConnectionLost         = 4
	FirmwareUpdateResultCrcCheck               = 5
	FirmwareUpdateResultUnsupportedPackageType = 6
	FirmwareUpdateResultInvalidUri             = 7
)

// Battery Status enum
const (
	BatteryStatusNormal         = 0 //operating normally and not on power
	BatteryStatusCharging       = 1 //currently charging
	BatteryStatusChargeComplete = 2 //fully charged and still on power
	BatteryStatusDamaged        = 3 //has some problem
	BatteryStatusLowBattery     = 4 //low on charge
	BatteryStatusNotInstalled   = 5 //not installed
	BatteryStatusUnknown        = 6 //information is not available
)

// Resource maps to LwM2M Resource depicted in
// OMA-TS-LightweightM2M_Core-V1_2_1-20221209-A
// Appendix D.1
type Resource interface {
	Id() ResourceID
	Name() string
	Type() ValueType
	Operations() OpCode //combination of Read (R), Write (W), and Execute (E)
	Mandatory() bool
	Multiple() bool
	Units() string
	RangeOrEnums() string
	Description() string

	SetId(ResourceID)
	SetName(string)
	SetType(ValueType)
	SetOperations(OpCode)
	SetMandatory(bool)
	SetMultiple(bool)
	SetUnits(string)
	SetRangeOrEnums(string)
	SetDescription(string)

	MarshalJSON() ([]byte, error)
}

type ResourceImpl struct {
	id           ResourceID
	name         string
	kind         ValueType
	operations   OpCode
	multiple     bool
	mandatory    bool
	units        string
	rangeOrEnums string
	description  string
}

func (r *ResourceImpl) MarshalJSON() ([]byte, error) {
	buf := []byte(`{`)
	buf = append(buf, `"id":`+strconv.Itoa(int(r.id))...)
	buf = append(buf, `,"name":"`+r.name+`"`...)
	buf = append(buf, `,"kind":`+strconv.Itoa(int(r.kind))...)
	buf = append(buf, `,"operations":`+strconv.Itoa(int(r.operations))...)
	buf = append(buf, `,"multiple":`...)
	buf = strconv.AppendBool(buf, r.multiple)
	buf = append(buf, `,"mandatory":`...)
	buf = strconv.AppendBool(buf, r.mandatory)
	buf = append(buf, `,"units":"`+r.units+`"`...)
	buf = append(buf, `,"rangeOrEnums":"`+r.rangeOrEnums+`"`...)
	buf = append(buf, `,"description":"`+r.description+`"`...)
	buf = append(buf, `}`...)

	return buf, nil
}

func (r *ResourceImpl) Id() ResourceID {
	return r.id
}

func (r *ResourceImpl) SetId(id ResourceID) {
	r.id = id
}

func (r *ResourceImpl) Name() string {
	return r.name
}

func (r *ResourceImpl) SetName(name string) {
	r.name = name
}

func (r *ResourceImpl) Type() ValueType {
	return r.kind
}

func (r *ResourceImpl) SetType(kind ValueType) {
	r.kind = kind
}

func (r *ResourceImpl) Operations() OpCode {
	return r.operations
}

func (r *ResourceImpl) SetOperations(operations OpCode) {
	r.operations = operations
}

func (r *ResourceImpl) Executable() bool {
	return r.operations == OpExecute
}

func (r *ResourceImpl) Readable() bool {
	return r.operations == OpRead || r.operations == OpReadWrite
}

func (r *ResourceImpl) Writable() bool {
	return r.operations == OpWrite || r.operations == OpReadWrite
}

func (r *ResourceImpl) Multiple() bool {
	return r.multiple
}

func (r *ResourceImpl) SetMultiple(multiple bool) {
	r.multiple = multiple
}

func (r *ResourceImpl) Mandatory() bool {
	return r.mandatory
}

func (r *ResourceImpl) SetMandatory(mandatory bool) {
	r.mandatory = mandatory
}

func (r *ResourceImpl) Units() string {
	return r.units
}

func (r *ResourceImpl) SetUnits(units string) {
	r.units = units
}

func (r *ResourceImpl) RangeOrEnums() string {
	return r.rangeOrEnums
}

func (r *ResourceImpl) SetRangeOrEnums(rangeOrEnums string) {
	r.rangeOrEnums = rangeOrEnums
}

func (r *ResourceImpl) Description() string {
	return r.description
}

func (r *ResourceImpl) SetDescription(description string) {
	r.description = description
}
