package references

import (
	"github.com/zxh0/jvm.go/instructions/base"
	"github.com/zxh0/jvm.go/rtda"
	"github.com/zxh0/jvm.go/rtda/heap"
)

// Invoke a class (static) method
type InvokeStatic struct {
	base.Index16Instruction
	method *heap.Method
}
/**
* 执行静态方法的调用
* 当前接口中也可以定义静态方法。
*/
func (instr *InvokeStatic) Execute(frame *rtda.Frame) {
	if instr.method == nil {
		cp := frame.GetConstantPool()
		k := cp.GetConstant(instr.Index)
		
		if kMethodRef, ok := k.(*heap.ConstantMethodRef); ok {
			instr.method = kMethodRef.GetMethod(true)
		//如果是接口方法
		} else {
			instr.method = k.(*heap.ConstantInterfaceMethodRef).GetMethod(true)
		}
	}

	// init class
	class := instr.method.Class
	if class.InitializationNotStarted() {
		frame.RevertNextPC()
		frame.Thread.InitClass(class)
		return
	}

	frame.Thread.InvokeMethod(instr.method)
}
