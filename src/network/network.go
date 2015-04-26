package network

import(
	."fmt"
	"net"
	"strconv"	
	."time"
	."os"
	."globals"
	
)
//---------------------------------------------------------------------------------------------------
func get_local_broadcast_ip() (string,string) {
	var Local_ip string = ""
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		Stderr.WriteString("Oops: " + err.Error() + "\n")
		Exit(1)
	}
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				Local_ip+=(ipnet.IP.String())
			}
		}
	}
	slice := []byte(Local_ip)
	slice[len(slice)-3] = 50
	slice[len(slice)-2] = 53
	slice[len(slice)-1] = 53
	broadcast_ip := string(slice)
	return Local_ip,broadcast_ip
}
//---------------------------------------------------------------------------------------------------
func get_local_broadcast_ips() (string) {
	var Local_ip string = ""
	
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		Stderr.WriteString("Oops: " + err.Error() + "\n")
		Exit(1)
	}
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				Local_ip+=(ipnet.IP.String())
			}
		}
	}
	return Local_ip
}
//---------------------------------------------------------------------------------------------------
func should_be_master(){              
        msg := make([]byte,1024)
		read_conn.SetReadDeadline(Now().Add(1000*Millisecond))
		_,_,err := read_conn.ReadFromUDP(msg)
		if err != nil {
			Is_master = true		
			return 
		}	
		Is_master = false
		
		return      
}
//---------------------------------------------------------------------------------------------------
func read_disconnect_and_reconnect(){
	_=read_conn.Close()
	read_conn,_ = net.ListenUDP("udp",read_addr)
}
//---------------------------------------------------------------------------------------------------
var work_space int=25
var broadcast_port string=strconv.Itoa(20000+work_space)
var Local_ip,broadcast_ip string=get_local_broadcast_ip()
var broadcast_addr,_ = net.ResolveUDPAddr("udp",broadcast_ip+ ":" +broadcast_port)		
var read_addr,_ = net.ResolveUDPAddr("udp", ":" + broadcast_port)		
var read_conn,_ = net.ListenUDP("udp",read_addr)
var Broadcast_conn,_ = net.DialUDP("udp",nil,broadcast_addr)

func Network(){
	should_be_master()	
	Sleep(1000*Millisecond)
	read_disconnect_and_reconnect()
	if Is_master == true{
		go master()
		Println("Master")
	}else{
		go slave()
		Println("Slave")
	}
	return	
}
//---------------------------------------------------------------------------------------------------

