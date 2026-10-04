package entity

import "fmt"

type PersonaType string

const (
	TypeAdultCustomer    PersonaType = "adult_customer"
	TypeAdultNonCustomer PersonaType = "adult_non_customer"
	TypeChildCustomer    PersonaType = "child_customer"
	TypeChildNonCustomer PersonaType = "child_non_customer"
	TypeEmployee         PersonaType = "employee"
)

type Sign string

const (
	SignPos Sign = "pos"
	SignNeg Sign = "neg"
)

type Inevitability string

const (
	InevitabilityDaily          Inevitability = "daily"
	InevitabilityCriticalMoment Inevitability = "critical_moment"
	InevitabilityTransientNoise Inevitability = "transient_noise"
)

type PersonStatus string

const (
	StatusActive         PersonStatus = "active"
	StatusMonitoringOnly PersonStatus = "monitoring_only"
	StatusQuarantined    PersonStatus = "quarantined"
)

func ValidatePersonaType(t string) (PersonaType, error) {
	pt := PersonaType(t)
	switch pt {
	case TypeAdultCustomer, TypeAdultNonCustomer, TypeChildCustomer, TypeChildNonCustomer, TypeEmployee:
		return pt, nil
	default:
		return "", fmt.Errorf("invalid persona type: %s", t)
	}
}

func ValidateSign(s string) (Sign, error) {
	sn := Sign(s)
	switch sn {
	case SignPos, SignNeg:
		return sn, nil
	default:
		return "", fmt.Errorf("invalid sign: %s (must be pos or neg)", s)
	}
}

func ValidateInevitability(i string) (Inevitability, error) {
	in := Inevitability(i)
	switch in {
	case InevitabilityDaily, InevitabilityCriticalMoment, InevitabilityTransientNoise:
		return in, nil
	default:
		return "", fmt.Errorf("invalid inevitability: %s", i)
	}
}
