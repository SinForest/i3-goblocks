package advancedclick

type ClickEvent struct {
	Name      string      `json:"name,omitempty"`
	Instance  string      `json:"instance,omitempty"`
	Button    MouseButton `json:"button,omitempty"`
	Modifiers []string    `json:"modifiers,omitempty"`
	X         int         `json:"x,omitempty"`
	Y         int         `json:"y,omitempty"`
	RelaiveX  int         `json:"relaive_x,omitempty"`
	RelativeY int         `json:"relative_y,omitempty"`
	OutputX   int         `json:"output_x,omitempty"`
	OutputY   int         `json:"output_y,omitempty"`
	Width     int         `json:"width,omitempty"`
	Height    int         `json:"height,omitempty"`
}

type MouseButton int

const (
	BtnLeft       MouseButton = 1 // left button
	BtnMiddle     MouseButton = 2 // middle button (pressing the scroll wheel)
	BtnRight      MouseButton = 3 // right button
	BtnWheelUp    MouseButton = 4 // turn scroll wheel up
	BtnWheelDown  MouseButton = 5 // turn scroll wheel down
	BtnWheelLeft  MouseButton = 6 // push scroll wheel left
	BtnWheelRight MouseButton = 7 // push scroll wheel right
	BtnBack       MouseButton = 8 // 4th button (aka browser backward button)
	BtnForward    MouseButton = 9 // 5th button (aka browser forward button)
)
