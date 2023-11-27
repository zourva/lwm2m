package core

type PropertiesAttributeName = string

const (
	Dimension      PropertiesAttributeName = "dim"
	ShortServerID  PropertiesAttributeName = "ssid"
	ServerURI      PropertiesAttributeName = "uri"
	ObjectVersion  PropertiesAttributeName = "ver"
	EnablerVersion PropertiesAttributeName = "lwm2m"
)

type NotificationAttributeName = string

const (
	MinimumPeriod           NotificationAttributeName = "pmin"  //Readable
	MaximumPeriod           NotificationAttributeName = "pmax"  //Readable
	GreaterThan             NotificationAttributeName = "gt"    //Numerical & Readable
	LesserThan              NotificationAttributeName = "lt"    //Numerical & Readable
	Step                    NotificationAttributeName = "st"    //Numerical & Readable
	MinimumEvaluationPeriod NotificationAttributeName = "epmin" //Readable
	MaximumEvaluationPeriod NotificationAttributeName = "epmax" //Readable
	Edge                    NotificationAttributeName = "edge"  //Boolean & Readable
	ConfirmableNotification NotificationAttributeName = "con"   //Boolean & Readable
	MaximumHistoricalQueue  NotificationAttributeName = "hqmax" //Readable
)

type NotificationAttrs = map[string]any
