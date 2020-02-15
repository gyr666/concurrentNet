package handle

type AcceptCallBack func(Connection) error
type AfterAcceptCallBack func(ChildChannel, ParentChannel) error
type DataCallBack func(ChildChannel) error
type ExceptionCallBack func(Throwable) error
type CloseCallback func(InetAddress64) error

type EventCallBack interface {
	Nc() AfterAcceptCallBack
	Dc() DataCallBack
	Ec() ExceptionCallBack
	Ac() AcceptCallBack
	Cc() CloseCallback
}
