package enum

type Color string

const (
    Unknown = "Unknown"

    ColorRed   Color = "Red"
    ColorGreen Color = "Green"
    ColorBlue  Color = "Blue"
)

func (c Color) String() string {
    if !c.IsValid() {
        return Unknown
    }

    return string(c)
}

func (c Color) IsValid() bool {
    switch c {
    case ColorRed, ColorGreen, ColorBlue:
        return true
    default:
        return false
    }
}

type Status int

const (
    StatusActive Status = iota
    StatusInactive
)

func (s Status) String() string {
    if !s.IsValid() {
        return Unknown
    }

    if s == StatusInactive {
        return "Inactive"
    } else {
        return "Active"
    }
}

func (s Status) IsValid() bool {
    switch s {
    case StatusActive, StatusInactive:
        return true
    default:
        return false
    }
}

type SwitchState bool

const (
    SwitchOn  SwitchState = true
    SwitchOff SwitchState = false
)

func (s SwitchState) String() string {
    if !s.IsValid() {
        return Unknown
    }

    if s == SwitchOn {
        return "On"
    } else {
        return "Off"
    }
}

func (s SwitchState) IsValid() bool {
    switch s {
    case SwitchOn, SwitchOff:
        return true
    default:
        return false
    }
}
