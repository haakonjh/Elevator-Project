package driver

import(
	."time"
	."fmt"
	."globals"
)
var lamp_matrix = [N_FLOORS][N_BUTTONS] int {
{LIGHT_UP1, UNDEFINED, LIGHT_COMMAND1},
{LIGHT_UP2, LIGHT_DOWN2, LIGHT_COMMAND2},
{LIGHT_UP3, LIGHT_DOWN3, LIGHT_COMMAND3},
{UNDEFINED, LIGHT_DOWN4, LIGHT_COMMAND4},
}

var button_matrix = [N_FLOORS][N_BUTTONS] int {
{BUTTON_UP1, UNDEFINED, BUTTON_COMMAND1},
{BUTTON_UP2, BUTTON_DOWN2, BUTTON_COMMAND2},
{BUTTON_UP3, BUTTON_DOWN3, BUTTON_COMMAND3},
{UNDEFINED, BUTTON_DOWN4, BUTTON_COMMAND4},
}
//---------------------------------------------------------------------------------------------------
func Elev_init()bool{
	
	if(Init()){
		Println("ELEVATOR INITIATED SUCCESFULLY")
		Clear_all_lights()
		Motor_stop()	
	}	
	return false				
}
//---------------------------------------------------------------------------------------------------
func Return_lamp_matrix(floor int, button int)int{
	return lamp_matrix[floor][button]
}
//---------------------------------------------------------------------------------------------------
func drive_in_direction(){
 	switch Direction{
 		case DIRECTION_UP:
			Clear_bit(MOTORDIR)
			Write_analog(MOTOR,2800) 
			Direction = DIRECTION_UP
			Is_moving=true
		case DIRECTION_DOWN:
			Set_bit(MOTORDIR)
			Write_analog(MOTOR,2800);
			Direction = DIRECTION_DOWN
			Is_moving=true
 	
 	}
 
}
//---------------------------------------------------------------------------------------------------
func Drive_to_floor(floor int)bool{
	
	if floor > Current_floor{
		Direction = DIRECTION_UP		
	}else if floor < Current_floor{
		Direction = DIRECTION_DOWN
	}
	drive_in_direction()
	if Get_floor_sensor() == floor{
		Motor_stop()
		return true
	}
	return false
}
//---------------------------------------------------------------------------------------------------
func Motor_stop () {
	Write_analog(MOTOR,0)
	Is_moving=false
}
//---------------------------------------------------------------------------------------------------
func Set_door_open_light() {
	Set_bit(LIGHT_DOOR_OPEN)

}
//---------------------------------------------------------------------------------------------------
func Clear_door_open_light() {
	Clear_bit(LIGHT_DOOR_OPEN)
}
//---------------------------------------------------------------------------------------------------
func Get_door() int{
	return Read_bit(LIGHT_DOOR_OPEN)
}
//---------------------------------------------------------------------------------------------------
func Clear_light(floor int, button int){
	var lamp int = lamp_matrix[floor][button]
	Clear_bit(lamp)
}
//---------------------------------------------------------------------------------------------------
func Clear_all_lights(){
	for floor := 0;floor < N_FLOORS ; floor++ {
		for button :=0;button< N_BUTTONS; button++{
			if lamp_matrix[floor][button]!= UNDEFINED{
				Clear_light(floor, button)
			}		
		}	
	}
	Clear_door_open_light()	
}
//---------------------------------------------------------------------------------------------------
func Set_light(){
	
	var light_event Order_info
	
	for{	
		Sleep(Millisecond * 50)
		select{
			case light_event= <- Light_chan:
				switch light_event.Type{
					case BUTTON_CALL_DOWN, BUTTON_CALL_UP,BUTTON_COMMAND:
						Set_bit(lamp_matrix[light_event.Floor][light_event.Type])					
				}
		}
	
	}	

}
//---------------------------------------------------------------------------------------------------
func Read_input(){
	var last_floor int;
	temp_order := Order_info{Type: BUTTON_COMMAND, Floor: 0}	
	for{		
		Sleep(Millisecond * 50)	
		temp_floor:=Get_floor_sensor()
		if (temp_floor!=-1) && (temp_floor!=last_floor){
			last_floor=temp_floor
			set_floor_indicator(last_floor)
		}
		
		for floor:=0;floor < N_FLOORS; floor++{
			
			for Type := BUTTON_CALL_UP;Type <= BUTTON_COMMAND;Type++{
				if(Read_bit(button_matrix[floor][Type])==1){
					temp_order.Type=Type
					temp_order.Floor=floor
					Light_chan<-temp_order
					Input_chan<-temp_order	
												
				}
			}
		
		}
	
	}
}
//---------------------------------------------------------------------------------------------------
func set_floor_indicator(floor int){

	if floor != -1{	
		if ((floor & 0x02) != 0 ){
			Set_bit(LIGHT_FLOOR_IND1)
		}else{
			Clear_bit(LIGHT_FLOOR_IND1)
		}
		if ((floor & 0x01) != 0){
			Set_bit(LIGHT_FLOOR_IND2)
		} else {
			Clear_bit(LIGHT_FLOOR_IND2)
		}
	}	
}	
//---------------------------------------------------------------------------------------------------
func Get_floor_sensor() int{

	if Read_bit(SENSOR_FLOOR1)== 1{
		Current_floor=0
		return 0
	}
	if Read_bit(SENSOR_FLOOR2)== 1{
		Current_floor=1
		return 1
	}
	if Read_bit(SENSOR_FLOOR3)==1{
		Current_floor=2
		return 2
	}
	if Read_bit(SENSOR_FLOOR4)==1{
		Current_floor=3
		return 3
	}
	return -1
}
//---------------------------------------------------------------------------------------------------
func Print_direction_and_floor(){
	Println("Current Floor: ",Get_floor_sensor())
	if Direction == DIRECTION_UP{
		Println("Direction up")
	}else if Direction == DIRECTION_DOWN{
		Println("Direction Down")
	}else if Direction == DIRECTION_STOP{
		Println("Direction Stop")
	}
}
//---------------------------------------------------------------------------------------------------
