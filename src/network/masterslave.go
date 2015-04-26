package network

import(
	."time"
	."fmt"
	"encoding/json"
	."globals"
	
)
//---------------------------------------------------------------------------------------------------
func in_array(arr[] string,item string)bool{
	for i:=0;i<len(arr);i++{
		if arr[i] == item{
			return true
		}
	}
	return false
}
//---------------------------------------------------------------------------------------------------
func Pos_best_elev (order Order_info)(int){ 
	var best_cost int = 99999
	var slave_number int
	for slave:=0;slave<len(Slave_last_message);slave++{
		var slave_cost int =0
		if Slave_last_message[slave].Is_moving==true{
			slave_cost+=5
		}
		if Slave_last_message[slave].Current_floor < order.Floor && Slave_last_message[slave].Direction == DIRECTION_DOWN{
			
			slave_cost += order.Floor- Slave_last_message[slave].Current_floor
		}
		if Slave_last_message[slave].Current_floor > order.Floor && Slave_last_message[slave].Direction == DIRECTION_UP{
			slave_cost += Slave_last_message[slave].Current_floor-order.Floor
		}
		if slave_cost<best_cost{
			best_cost=slave_cost
			slave_number=slave
		}		
	}
	var master_cost int =0
	if Is_moving==true{
		master_cost+=5
	}
	if Current_floor < order.Floor && Direction == DIRECTION_DOWN{
			
		master_cost+=order.Floor-Current_floor
	}
	if Current_floor > order.Floor && Direction==DIRECTION_UP{
		master_cost += Current_floor-order.Floor
	}
	if master_cost<=best_cost{
		best_cost=master_cost
		slave_number=-1
	}
	
	return slave_number
}
//---------------------------------------------------------------------------------------------------
func master(){
	var empty_order Order_info = Order_info{1,1}
	var received_message Udp_message
	go read(&received_message)
	
	for{
		Sleep(20*Millisecond)
		select{
			case send_event := <- Send_chan:
			var pos int = Pos_best_elev (send_event)
			var slave_message Udp_message= Udp_message{Local_ip:Local_ip,Destination_ip:Slave_list[pos],Is_order:true,Order:send_event,Direction:-5,Current_floor:-5,Is_moving:false,Slaves:Slave_list,Is_new_master:false}	
						msg,_:=json.Marshal(slave_message)
						Broadcast_conn.Write([]byte(msg))	

		default:
			if received_message.Is_order==true {
				Input_chan<-received_message.Order
				if in_array(Slave_list,received_message.Local_ip)==false && received_message.Local_ip!=""  && received_message.Local_ip!=Local_ip {
					Slave_list=append(Slave_list,received_message.Local_ip)
					Slave_last_message=append(Slave_last_message,received_message)
				}
			}

			if in_array(Slave_list,received_message.Local_ip)==false && received_message.Local_ip!=""  && received_message.Local_ip!=Local_ip{
				Slave_list=append(Slave_list,received_message.Local_ip)
				Slave_last_message=append(Slave_last_message,received_message)
				
			}

			if in_array(Slave_list,Local_ip){
				var slave_number int=0
				for slave:=0;slave<len(Slave_list);slave++{
					if Local_ip==Slave_list[slave]{
						slave_number=slave
					}
				}
				Slave_list = append(Slave_list[:slave_number], Slave_list[slave_number+1:]...)
				Slave_last_message = append(Slave_last_message[:slave_number], Slave_last_message[slave_number+1:]...)
				
			}

			var message Udp_message= Udp_message{Local_ip:Local_ip,Destination_ip:"",Is_order:false,Order:empty_order,Direction:-5,Current_floor:-5,Is_moving:false,Slaves:Slave_list,Is_new_master:false}
			for i:=0;i<len(Slave_last_message);i++{
				if received_message.Local_ip==Slave_last_message[i].Local_ip{
					Slave_last_message[i]=received_message
				}
			}
			msg,_:=json.Marshal(message)
			Broadcast_conn.Write([]byte(msg))
		}
	}
}
//---------------------------------------------------------------------------------------------------
func slave(){
	var empty_order Order_info = Order_info{1,1}
	var received_message Udp_message
	go read(&received_message)
	time_last_received_master :=Now()
	for{
		Sleep(50*Millisecond)
		Slave_list=received_message.Slaves
		
		if received_message.Is_order==true && received_message.Destination_ip==Local_ip{
			Slave_order<-received_message.Order
			Light_chan<-received_message.Order

		}

		if received_message.Direction==-5 && received_message.Current_floor==-5{
			time_last_received_master =Now()	
		}

		if received_message.Is_new_master==true && received_message.Destination_ip==Local_ip{
			for slave:=0;slave<len(Slave_list);slave++{
				Slave_last_message=append(Slave_last_message,Udp_message{Local_ip:Slave_list[slave],Destination_ip:"",Is_order:false,Order:empty_order,Direction:Direction,Current_floor:Current_floor,Is_moving:Is_moving,Slaves:Slave_list,Is_new_master:false})
			}
			go master()
			Is_master=true
			Quit_add_to_orders <- true
			Sleep(100*Millisecond)
			return
		}

		if Since(time_last_received_master)>(3000*Millisecond){
			var message = 
				Udp_message{Local_ip:Local_ip,Destination_ip:Slave_list[0],Is_order:false,Order:empty_order,Direction:Direction,Current_floor:Current_floor,Is_moving:Is_moving,Slaves:Slave_list,Is_new_master:true}
				msg,_:=json.Marshal(message)
				Broadcast_conn.Write([]byte(msg))
				Println("IP:" + Slave_list[0]+" Is the new master!")
				continue				
		}

		select{
			case order_event := <- Input_chan: 
				if order_event.Type!= BUTTON_COMMAND{
					var message = 
					Udp_message		{Local_ip:Local_ip,Destination_ip:"",Is_order:true,Order:order_event,Direction:Direction,Current_floor:Current_floor,Is_moving:Is_moving,Slaves:Slave_list,Is_new_master:false}
					msg,_:=json.Marshal(message)
					Broadcast_conn.Write([]byte(msg))					
				}else{
					Slave_order <- order_event
				
				}									
			default:
				var message = 
				Udp_message{Local_ip:Local_ip,Destination_ip:"",Is_order:false,Order:empty_order,Direction:Direction,Current_floor:Current_floor,Is_moving:Is_moving,Slaves:Slave_list,Is_new_master:false}
				msg,_:=json.Marshal(message)
				Broadcast_conn.Write([]byte(msg))
		}
	}	
}
//---------------------------------------------------------------------------------------------------
