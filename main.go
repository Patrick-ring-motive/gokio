package main

import (
	. "fmt"
	. "handler/api/main/src"
)

type FuncStruct[T any,A any,B any,FN func()T|func(A)T|func(...A)T|func(A,B)T|func(A,...B)T ] struct{
  Func FN
}

func test(i int,b int, c int)int{
  return i + b + c
}
func testInt(a func(any,any,any)int){
  a(1,2,3)
}

func main() {
//  testInt(test)

    i := GenericFunction(1,2)
    prom := Go(GenericFunction,INT,1,2);
    x := Await(prom);
    Println(i)
    Println(x)


}



