package mytrade

type subscription[T any] struct {
	id         string        //订阅ID
	resultChan chan T        //接收订阅结果的通道
	errChan    chan error    //接收订阅错误的通道
	closeChan  chan struct{} //接收订阅关闭的通道
}

// 获取订阅结果
func (sub *subscription[T]) ResultChan() chan T {
	return sub.resultChan
}

// 获取错误订阅
func (sub *subscription[T]) ErrChan() chan error {
	return sub.errChan
}

// 获取关闭订阅信号
func (sub *subscription[T]) CloseChan() chan struct{} {
	return sub.closeChan
}
