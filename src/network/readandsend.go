package network

import(
	."time"
	"encoding/json"
	."globals"
)
//---------------------------------------------------------------------------------------------------
func read(message* Udp_message){
	msg := make([]byte,1024)
	
	for{
		Sleep(10*Millisecond)
		msgSize,_,_ := read_conn.ReadFromUDP(msg)
		var temp Udp_message
		json.Unmarshal(msg[:msgSize],&temp)
		*message=temp
	}	
}
//---------------------------------------------------------------------------------------------------
func send(message Udp_message){
	msg,_:=json.Marshal(message)
	for{
		Sleep(10*Millisecond)
		Broadcast_conn.Write([]byte(msg))
	}
}
//---------------------------------------------------------------------------------------------------
