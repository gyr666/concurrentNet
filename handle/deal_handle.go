package handle

type AcceptCallBack func(Connection) error
type AfterAcceptCallBack func(ChildChannel,ParentChannel) error
type DataCallBack func(ChildChannel) error
type ExceptionCallBack(Throwable) error
type CloseCallback(InetAddress64) error

type EventCallBack interface{
	Nc() AfterAcceptCallBack
	Dc() DataCallBack
	Ec() ExceptionCallBack
	Ac() AcceptCallBack
	Cc() CloseCallback
}



