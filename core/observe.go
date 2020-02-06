package core
import "fmt"
type ServerObserve interface{
	OnBooting()
	OnBooted()
	OnStoping()
	OnStoped()
}

type DefaultObserve struct{

}
func	(d *DefaultObserve)OnBooting(){
	fmt.Println("On Booting")
}
func	(d *DefaultObserve)OnBooted(){
	fmt.Println("On Booted")
}
func	(d *DefaultObserve)OnStoping(){

	fmt.Println("On Stoping")

}
func	(d *DefaultObserve)OnStoped(){
	fmt.Println("On Stoped")

}
