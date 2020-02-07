package core
import "fmt"
type ServerObserve interface{
	OnBooting()
	OnBooted(l []NetworkInet64)
	OnStoping()
	OnStoped()
}

type DefaultObserve struct{

}
func (d *DefaultObserve)OnBooting(){
	fmt.Println("On Booting")
}
func (d *DefaultObserve)OnBooted(l []NetworkInet64){
	fmt.Printf("On Booted %#v",l)
}
func (d *DefaultObserve)OnStoping(){
	fmt.Println("On Stoping")
}
func (d *DefaultObserve)OnStoped(){
	fmt.Println("On Stoped")

}
