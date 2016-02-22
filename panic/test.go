package main

import "fmt"

func main() {
	//7. panic recover
	echoA()
	echoB()
	echoC()
	/**************************/
	/*运行结果:               */
	/*It is A                 */
	/*recover: It is B        */
	/*It is C                 */
	/**************************/
}

func echoA() {
	fmt.Println("It is A")
}

func echoB() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("recover: It is ", err)
		}
	}()

	panic("panic B") //注意顺序，要在 recover 后面
}

func echoC() {
	fmt.Println("It is C")
}
