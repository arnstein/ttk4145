package driver

/*
#cgo LDFLAGS: -lcomedi -lm
#include "elev.h"
*/
import "C"

const (
	BUTTON_CALL_UP   = 0
	BUTTON_CALL_DOWN = 1
	BUTTON_COMMAND   = 2
	DOOROPEN         = 3

	DIR_DOWN = -1
	DIR_STOP = 0
	DIR_UP   = 1
)

func ElevInit() int {
	err := int(C.elev_init())
	SetMotorDir(DIR_DOWN)
	for GetFloorSensorSignal() != 0 {
	}
	SetMotorDir(DIR_STOP)

	return err

}

func SetDoorOpenLight(state int) {
	C.elev_set_door_open_lamp(C.int(state))
}

func SetFloorLight(floor int) {
	C.elev_set_floor_indicator(C.int(floor))
}

func SetMotorDir(dir int) {
	C.elev_set_motor_direction(C.elev_motor_direction_t(dir))
}

func GetFloorSensorSignal() int {
	return int(C.elev_get_floor_sensor_signal())
}

func SetFloorIndicator(floor int) {
	C.elev_set_floor_indicator(C.int(floor))
}

func GetButtonSignal(button int, floor int) int {
	return int(C.elev_get_button_signal(C.elev_button_type_t(button), C.int(floor)))
}

func SetButtonLamp(button int, floor int, value int) {
	C.elev_set_button_lamp(C.elev_button_type_t(button), C.int(floor), C.int(value))
}
