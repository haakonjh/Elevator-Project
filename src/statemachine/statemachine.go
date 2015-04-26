package statemachine

import(
	."fmt"
	."driver"
	."orders"
	."time"
	."globals"
) 
//---------------------------------------------------------------------------------------------------
func Statemachine(){
	for{
		Sleep(Millisecond * 10)	
		select {
			case event_signal := <-Event_chan:
				switch event_signal{
					case INIT:
						Initialize()
						Event_chan<-IDLE
						break	
					case IDLE:
						Println("--IDLE--\n")
						idle()
						Event_chan<-MOVING
						break
					case MOVING:
						Println("--MOVING--\n")
						moving()
						Event_chan<-OPEN_DOORS
						break
					case OPEN_DOORS:
						Println("--OPENDOORS--\n")
						open_doors()
						Event_chan<- IDLE
						break
				}
		}
	}
}
//---------------------------------------------------------------------------------------------------
func Initialize(){
	Elev_init()
	Init_orders()
	go Add_to_orders()
	go Read_input()
	go Set_light()
}	
//---------------------------------------------------------------------------------------------------
func open_doors(){
	Set_door_open_light()
	Sleep(3000*Millisecond)
	Clear_door_open_light()
	Just_arrived = true
}
//---------------------------------------------------------------------------------------------------
func idle(){
	if Just_arrived == true{
		Clear_light(Queue[0].Floor,Queue[0].Type)
		Remove_order(Queue[0])
		Pop_queue()
	}
	for{
		Sleep(10*Millisecond)
		if len(Queue) > 0{
			return
		}
	}
}
//---------------------------------------------------------------------------------------------------
func moving(){
	
	for{
		Sleep(10*Millisecond)
		var floor = (Queue[0].Floor)
		
		if Drive_to_floor(floor) == true{
			return
		}
	}	
}
//---------------------------------------------------------------------------------------------------
