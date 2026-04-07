package main

type JoinMessage struct {
	Type     string `json:"type"`
	RoomCode string `json:"room_code"`
}

type InputMessage struct {
	Type string  `json:"type"`
	X    float64 `json:"x"`
	Z    float64 `json:"z"`
}

type StateMessage struct {
	Type       string  `json:"type"`
	PosX       float64 `json:"pos_x"`
	PosY       float64 `json:"pos_y"`
	PosZ       float64 `json:"pos_z"`
	RotX       float64 `json:"rot_x"`
	RotY       float64 `json:"rot_y"`
	RotZ       float64 `json:"rot_z"`
	RotW       float64 `json:"rot_w"`
	VelX       float64 `json:"vel_x"`
	VelY       float64 `json:"vel_y"`
	VelZ       float64 `json:"vel_z"`
	HostSteer  float64 `json:"host_steer"`
	HostMove   float64 `json:"host_move"`
	GuestSteer float64 `json:"guest_steer"`
	GuestMove  float64 `json:"guest_move"`
}

type InfoMessage struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}
