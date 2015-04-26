package Orders

import(
	."driver"
	."fmt"
	."time"
	."network"
	."globals"	
)
var orders[N_FLOORS][N_BUTTONS] bool
//---------------------------------------------------------------------------------------------------
func orders_queue_check(){
	for{
		Sleep(300*Millisecond)
		for floor := 0; floor < N_FLOORS; floor++ {
			for button := 0; button < N_BUTTONS; button++ {
				if orders[floor][button]==false && Read_bit(Return_lamp_matrix(floor,button))==1{
					Clear_bit(Return_lamp_matrix(floor,button))
				}
				if orders[floor][button] == true && Is_in_Queue(floor) == false{
					Clear_bit(Return_lamp_matrix(floor,button))
					orders[floor][button]=false
				}
			
			}
		}
	}

}
//---------------------------------------------------------------------------------------------------
func Is_in_Queue(floor int)bool{
	for i:=0;i<len(Queue);i++{
		if Queue[i].Floor == floor{
			return true
		}
	}
	return false
}
//---------------------------------------------------------------------------------------------------
func Init_orders(){
	orders = [N_FLOORS][N_BUTTONS]bool{
		{false,false,false},
		{false,false,false},
		{false,false,false},
		{false,false,false},
	}
	Queue=[]Order_info{}
	go orders_queue_check()
}
//---------------------------------------------------------------------------------------------------
func Add_to_orders(){
	if Is_master== true{
		for{
			Sleep(Millisecond*10)
				select{
					case input_event := <-Input_chan:
						switch input_event.Type{
							case BUTTON_CALL_UP,BUTTON_CALL_DOWN:
								var pos int = Pos_best_elev(input_event)
								if pos == -1{
									Sleep(Millisecond*50)
									if Is_in_orders(input_event){
										Print_orders()
										Print_queue()
										Print_direction_and_floor()
										continue;
									}else{
										Add_to_queue(input_event)
										orders[input_event.Floor][input_event.Type] = true;
										Print_orders()
										Print_queue()
										Print_direction_and_floor()
									}
								}else if pos>-1{
									Send_chan <- input_event
									Sleep(Millisecond*50)
								}
								
							case BUTTON_COMMAND:
								if Is_in_orders(input_event){
										Print_orders()
										Print_queue()
										Print_direction_and_floor()
										continue;
								}else{
										Add_to_queue(input_event)
										orders[input_event.Floor][input_event.Type] = true;
										Print_orders()
										Print_queue()
										Print_direction_and_floor()
								}			
						}
			}
		}
	}else if Is_master == false{
		for{
				select{
					case input_event := <-Slave_order:
						switch input_event.Type{
							case BUTTON_CALL_UP,BUTTON_CALL_DOWN,BUTTON_COMMAND:
								if Is_in_orders(input_event){
									Print_orders()
									Print_queue()
									Print_direction_and_floor()
									continue;
								}else{
									Add_to_queue(input_event)
									orders[input_event.Floor][input_event.Type] = true;
									Print_orders()
									Print_queue()
									Print_direction_and_floor()
								}		
						}
					case <- Quit_add_to_orders:
						go Add_to_orders()
						return	
				}
		}
	}	
}
//---------------------------------------------------------------------------------------------------
func Is_in_orders(order Order_info)bool{
	return orders[order.Floor][order.Type]
}
//---------------------------------------------------------------------------------------------------
func Queue_insert_to_pos(order Order_info, position int,){
	temp := Order_info{Floor:0,Type:order.Type}
	Queue = append(Queue, temp)
	copy(Queue[position+1:], Queue[position:])
	Queue[position] = order
}
//---------------------------------------------------------------------------------------------------
func Insert_by_floorsize(order Order_info, choice string){
	if len(Queue) == 1 {
		if choice == "up"{
			Queue_insert_to_pos(order,0)
			return
		}else{
			
			Queue = append(Queue,order)
			return
		
		}
	}
	var prev int = 0
	for i := 1; i < (len(Queue)); i++{
		if choice == "up"{
			if Queue[prev].Floor > Queue[i].Floor{
				Queue_insert_to_pos(order,i)
				break		
			}else if (i == (len(Queue)-1)){
				Queue = append(Queue,order)
				break
			}
		}else if choice == "down"{
			if Queue[prev].Floor < Queue[i].Floor{
				Queue_insert_to_pos(order,i)
				break		
			}else if (i == (len(Queue)-1)){
				Queue = append(Queue,order)
				break
			}
		}
		prev += 1
	}
}
//---------------------------------------------------------------------------------------------------
func Insert_higher_order(order Order_info){
	if order.Type == BUTTON_CALL_UP{
		if order.Floor < Queue[0].Floor{
			Queue_insert_to_pos(order,0)
		}else if order.Floor > Queue[0].Floor{
			Insert_by_floorsize(order,"down")
		}
	}else if order.Type == BUTTON_CALL_DOWN{
		if Is_order_above(Current_floor){
			for i := 0;i<len(Queue);i++{
				if Queue[i].Floor > Current_floor{
					if Queue[i].Type == BUTTON_CALL_DOWN {
						if order.Floor < Queue[i].Floor {
							Queue_insert_to_pos(order,(i+1))
							break
						}else if order.Floor > Queue[i].Floor{
							Queue_insert_to_pos(order,(i))
							break
						}
					}else if i == len(Queue)-1{
						Queue = append(Queue,order)
						break
					}
				}					
			}

		}else{
			Queue = append(Queue,order)
		}
	}else if order.Type == BUTTON_COMMAND{
		if order.Floor < Queue[0].Floor{
			Queue_insert_to_pos(order,0)
		}else if order.Floor > Queue[0].Floor{
			if Queue[0].Type != BUTTON_CALL_DOWN{
				Insert_by_floorsize(order,"down")
			}else{
				Queue_insert_to_pos(order,0)
			}
		}
	}	
}
//---------------------------------------------------------------------------------------------------
func Insert_lower_order(order Order_info){
	if order.Type == BUTTON_CALL_UP{
		if Is_order_below(Current_floor){
			for i := 0;i<len(Queue);i++{
				if Queue[i].Floor < Current_floor{
					if Queue[i].Type == BUTTON_CALL_UP {
						if order.Floor > Queue[i].Floor {
							Queue_insert_to_pos(order,(i+1))
							break
						}else if order.Floor < Queue[i].Floor{
							Queue_insert_to_pos(order,(i))
							break
						}
					}
				}else if i == len(Queue)-1{
					Queue = append(Queue,order)
					break
				}
			}

		}else{
			Queue = append(Queue,order)
		}		
	
	}else if order.Type == BUTTON_CALL_DOWN{
		if order.Floor > Queue[0].Floor{
			Queue_insert_to_pos(order,0)			
		}else if order.Floor < Queue[0].Floor{
			Insert_by_floorsize(order,"down")
		}else if order.Floor == Queue[0].Floor{
			if Queue[0].Type != BUTTON_CALL_DOWN{
				Insert_by_floorsize(order,"up")
			}else{
				Queue_insert_to_pos(order,0)
			}
		}
	}else if order.Type == BUTTON_COMMAND{
			if order.Floor > Queue[0].Floor{
				Queue_insert_to_pos(order,0)
			}else if order.Floor < Queue[0].Floor{
				if Queue[0].Type != BUTTON_CALL_DOWN{
					Insert_by_floorsize(order,"down")
				}else{
					Queue_insert_to_pos(order,0)
				}

			}
	}

}
//---------------------------------------------------------------------------------------------------
func Insert_order_opposite_direction(order Order_info){
	for i := 0; i < len(Queue); i++ {
		if order.Type == BUTTON_CALL_UP{
			if Queue[i].Type == BUTTON_CALL_UP{
				if order.Floor < Queue[i].Floor{
					Queue_insert_to_pos(order,i)
					return
				}
			}
			
		}else if order.Type == BUTTON_CALL_DOWN{
			if Queue[i].Type == BUTTON_CALL_DOWN{
				if order.Floor > Queue[i].Floor{
					Queue_insert_to_pos(order,i)
					return
				}
			}
		}
	}
	Queue = append(Queue,order)
}
//---------------------------------------------------------------------------------------------------
func Add_to_queue(order Order_info){
	if Is_empty_queue(){
		if order.Floor == Current_floor && Direction == order.Type{
			Remove_order(order)
			Clear_light(order.Floor, order.Type)
			return
		
		}else{
			Queue = append(Queue,order)
			return
		}
	}	
	switch order.Type{
		case BUTTON_CALL_UP, BUTTON_CALL_DOWN, BUTTON_COMMAND:				
			if order.Type != Direction && order.Type != BUTTON_COMMAND && order.Floor!=Current_floor{
				Insert_order_opposite_direction(order)
			}else{	
					if order.Floor > Current_floor{
						Insert_higher_order(order)

					}else if order.Floor < Current_floor{
						Insert_lower_order(order)

					}else if order.Floor == Current_floor{
						if Direction != order.Type{
							Insert_order_opposite_direction(order)	
						}else{
							Remove_order(order)
							Clear_light(order.Floor, order.Type)
						
						}
					}
			}	
	}	
}
//---------------------------------------------------------------------------------------------------
func Is_order_above(floor int)bool{
	for f := floor+1; f < N_FLOORS; f++ {
		for b := 0; b < N_BUTTONS; b++ {
			if orders[f][b] {
				return true
			}
		}
	}
	return false
}
//---------------------------------------------------------------------------------------------------
func Is_order_below(floor int)bool{
	for f := 0; f < floor; f++ {
		for b := 0; b < N_BUTTONS; b++ {
			if orders[f][b] {
				return true
			}
		}
	}
	return false
}
//---------------------------------------------------------------------------------------------------
func Any_order()bool{
	for f := 0; f < N_FLOORS; f++ {
		for b := 0; b < N_BUTTONS; b++ {
			if orders[f][b] {
				return true
			}
		}
	}
	return false
}
//---------------------------------------------------------------------------------------------------
func Remove_order(order Order_info){
	for i := 0; i < N_FLOORS; i++{
		for j := 0; j < N_BUTTONS; j++ {
			orders[order.Floor][order.Type] = false
		}
	}
	Print_orders()
}
//---------------------------------------------------------------------------------------------------
func Pop_queue(){
	Queue = Queue[1:]
	Print_queue()
}
//---------------------------------------------------------------------------------------------------
func Remove_all_orders(){
	for i := 0; i < N_FLOORS; i++ {
		for j := 0; j < N_BUTTONS; j++ {
			orders[i][j] = false
		}
	}
}
//---------------------------------------------------------------------------------------------------
func Is_empty_queue()bool{
	if len(Queue)==0{
		return true
	}
	return false
}
//---------------------------------------------------------------------------------------------------
func Print_orders(){
	Println("\n")
	Println("UP-----DOWN--COMMAND-------------------")
	for Floor:=N_FLOORS-1;Floor >= 0; Floor--{
		Println(orders[Floor],"\n")		
	}
	Println("----------------------------------------")
	Println("\n")	
}
//---------------------------------------------------------------------------------------------------
func Print_queue(){
	Println("\n")
	Println("Queue")
	Println("----------------------------------------")
	for i:=0 ;i < len(Queue); i++{
		switch Queue[i].Type{
			case BUTTON_CALL_UP:
				Println("Floor: ",Queue[i].Floor,"      |")
				Println("Type: CALL-UP   |")
			case BUTTON_CALL_DOWN:
				Println("Floor: ",Queue[i].Floor,"      |")
				Println("Type: CALL-DOWN |")
			case BUTTON_COMMAND:
				Println("Floor: ",Queue[i].Floor,"         |")
				Println("Type: COMMAND      |")
		}
			
		Println("")	
	}
	Println("----------------------------------------")
}
//---------------------------------------------------------------------------------------------------
