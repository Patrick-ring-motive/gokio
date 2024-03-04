package submodules

import (
	"fmt"
	"sync"
	. "unsafe"
  "reflect"
)

var initPromiseSync = Pass(sync.Mutex{})

var p func()


func GenericFunction(i int,h int) int {
  fmt.Println("GenericFunction ")
  fmt.Println(i)
  return i + h
}

/* Generic Structures */
/* 
  INT acts as a generic type definition
  so that it can be passed around as a type
  and be used to define the generic return type 
  without setting rigid input types
*/
func INT(int){}

/******************Promise Structures*************************/



type Func[A func()Z|func(B)Z|func(B,C)Z|func(B,C,D)Z,B any,C any,D any,Z any] struct{
  Function A
  Output Z
}



type Promise[T any] struct {
	PromiseChannel chan (T)
	Error          error
	Result         T
	Resolved       bool
	Rejected       bool
}



/*
PromiseBuilder and FinalizePromise exist to allow initizing a promise
with an empty result. Otherwise it would be next to impossible to instantiate
a promise generically
*/

type PromiseBuilder[T any] struct {
  PromiseChannel chan (T)
  Error          error
  Result         any
  Resolved       bool
  Rejected       bool
}

func FinalizePromise[T any](prom PromiseBuilder[T]) Promise[T] {
	promise := *(*Promise[T])(Pointer(&prom))
	return promise
}

func Go[T any](fn any,t func(T),payload ...any)Promise[T]{
  return  GoGeneric(fn,t,payload)
}


func GoGeneric[T any](f interface{},t func(T), payload []any) Promise[T] {
	promiseChannel := make(chan T)
	promiseBuilder := PromiseBuilder[T]{PromiseChannel: promiseChannel, Error: nil, Result: nil, Resolved: false, Rejected: false}
	promise := FinalizePromise(promiseBuilder)
	go Async(f,func(T){}, payload, promise)
	return promise
}



func Async[T any](f interface{},t func(T), payload []any, promise Promise[T]) Promise[T] {
	promise.Result = InvokeAnyFunc(f,t, payload)
	promise.Resolved = true
	promise.PromiseChannel <- promise.Result
	return promise
}

func Await[T any](promise Promise[T]) T {
	if promise.Resolved || promise.Rejected {
		return promise.Result
	} else {
		promise.Result = <-promise.PromiseChannel
		promise.Resolved = true
		close(promise.PromiseChannel)
		return promise.Result
	}
}

func InvokeAnyFunc[T any](fn interface{},t func(T), args interface{}) (T) {
  fnVal := reflect.ValueOf(fn)
  fnType := fnVal.Type()
  numIn := fnType.NumIn()
  in := make([]reflect.Value, numIn)
  for i, arg := range args.([]any) {
    argVal := reflect.ValueOf(arg)
    in[i] = argVal
  }
  out := fnVal.Call(in)
  result := make([]interface{}, len(out))
  for i, o := range out {
    result[i] = o.Interface()
  }

  return result[0].(T)
}