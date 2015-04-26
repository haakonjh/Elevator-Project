package globals

//Channels--------------------------------------------------------
var Slave_order = make(chan Order_info ,1)
var Order_chan = make(chan int,1)
var Input_chan = make(chan Order_info,1)
var Light_chan = make(chan Order_info,1)
var Motor_chan = make(chan int,1)
var Event_chan = make(chan int,1)
var Quit_add_to_orders = make(chan bool,1)
var Send_chan = make(chan Order_info,1)

//Constants--------------------------------------------------------
const BUTTON_CALL_UP int = 0
const BUTTON_CALL_DOWN int = 1
const BUTTON_COMMAND int = 2
const N_BUTTONS int = 3
const N_FLOORS int = 4 
const UNDEFINED int = 0
const DIRECTION_UP int = 0
const DIRECTION_DOWN int = 1
const DIRECTION_STOP int = -1
const IDLE int = 9
const OPEN_DOORS int = 10
const MOVING int = 12
const INIT int = 13

//Variables--------------------------------------------------------
var Is_moving bool = false
var Just_arrived bool = false
var Is_master bool
var elevator int
var Current_floor int = -1
var Direction int

//Arrays--------------------------------------------------------
var Queue[] Order_info
var Slave_list[] string = []string{}
var Slave_last_message[] Udp_message=[]Udp_message{}


//Types--------------------------------------------------------
type Order_info struct{
	Floor int
	Type int 
}
type Udp_message struct {	
	Local_ip string
	Destination_ip string
	Is_order bool
	Order Order_info
	Direction int
	Current_floor int	
	Is_moving bool
	Slaves[] string
	Is_new_master bool
}
//---------------------------------------------------------------------------------------------------	
