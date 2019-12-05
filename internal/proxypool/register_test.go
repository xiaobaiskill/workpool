package proxypool

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

type proxy1 struct{
	t *testing.T
}

func (p *proxy1)Pop() (HTTPClientMap, bool) {
	return HTTPClientMap{"127.0.0.1",&http.Client{}},true
}
func (p *proxy1) Push(client HTTPClientMap){
	assert.Equal(p.t,"127.0.0.1",client.Ip)
	return
}
func (p *proxy1) Del(ip string)  {
	assert.Equal(p.t,"127.0.0.1",ip)
	return
}
func (p *proxy1) Len()int{
	return 1
}

type proxy2 struct{
	t *testing.T
}


func (p *proxy2)Pop() (HTTPClientMap, bool) {
	return HTTPClientMap{"127.0.0.2",&http.Client{}},true
}
func (p *proxy2) Push(client HTTPClientMap){
	assert.Equal(p.t,"127.0.0.2",client.Ip)
	return
}
func (p *proxy2) Del(ip string)  {
	assert.Equal(p.t,"127.0.0.2",ip)
	return
}
func (p *proxy2) Len()int{
	return 1
}


func TestNewRegister(t *testing.T) {
	r := NewRegister()
	r.Add(&proxy1{t},2,"proxy1")
	r.Add(&proxy2{t},1,"proxy2")

	p := Newproxypools()
	client,index,ok := p.Pop()
	assert.Equal(t,"127.0.0.1",client.Ip)
	assert.Equal(t,0,index)
	assert.Equal(t,true,ok)
	p.Push(index,client)
	p.Del(index,client)

	client,index,ok = p.Pop()
	assert.Equal(t,"127.0.0.1",client.Ip)
	assert.Equal(t,0,index)
	assert.Equal(t,true,ok)
	p.Push(index,client)
	p.Del(index,client)

	client,index,ok = p.Pop()
	assert.Equal(t,"127.0.0.2",client.Ip)
	assert.Equal(t,1,index)
	assert.Equal(t,true,ok)
	p.Push(index,client)
	p.Del(index,client)

	_,_,ok = p.Pop()
	assert.Equal(t,false,ok)

}
