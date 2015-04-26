package main

import(	
	."statemachine"
	."network"
	."globals"
)
//---------------------------------------------------------------------------------------------------
func main(){
	Network()
	Event_chan<-INIT
	Statemachine()
	
	dead_chan :=make(chan bool,1)
	<- dead_chan
}
//---------------------------------------------------------------------------------------------------
