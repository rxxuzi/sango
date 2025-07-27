package ast

// ReceiverType represents the type of method receiver
type ReceiverType int

const (
	ValueReceiver ReceiverType = iota
	PointerReceiver
	ReferenceReceiver
)

// ReceiverInfo holds information about a method receiver
type ReceiverInfo struct {
	Type     ReceiverType
	TypeName string // e.g., "Point"
}

// String returns the string representation of the receiver
func (r *ReceiverInfo) String() string {
	switch r.Type {
	case PointerReceiver:
		return "*" + r.TypeName
	case ReferenceReceiver:
		return "&" + r.TypeName
	default:
		return r.TypeName
	}
}

// TypeString returns the receiver type as a string
func (rt ReceiverType) String() string {
	switch rt {
	case PointerReceiver:
		return "PointerReceiver"
	case ReferenceReceiver:
		return "ReferenceReceiver"
	default:
		return "ValueReceiver"
	}
}

// ParseReceiverType parses a receiver type string into ReceiverInfo
func ParseReceiverType(typeStr string) *ReceiverInfo {
	if len(typeStr) == 0 {
		return nil
	}
	
	if typeStr[0] == '*' {
		return &ReceiverInfo{
			Type:     PointerReceiver,
			TypeName: typeStr[1:],
		}
	}
	
	if typeStr[0] == '&' {
		return &ReceiverInfo{
			Type:     ReferenceReceiver,
			TypeName: typeStr[1:],
		}
	}
	
	return &ReceiverInfo{
		Type:     ValueReceiver,
		TypeName: typeStr,
	}
}