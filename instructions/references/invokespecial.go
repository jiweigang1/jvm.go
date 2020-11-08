package references

import (
	"github.com/zxh0/jvm.go/instructions/base"
	"github.com/zxh0/jvm.go/rtda"
	"github.com/zxh0/jvm.go/rtda/heap"
)

// Invoke instance method;
// special handling for superclass, private, and instance initialization method invocations
// 调用父类，私有，构造方法
type InvokeSpecial struct{ base.Index16Instruction }

func (instr *InvokeSpecial) Execute(frame *rtda.Frame) {
	//获取常量池
	cp := frame.GetConstantPool()
	//获取方法的符号引用
	k := cp.GetConstant(instr.Index)
	//如果是类的方法
	if kMethodRef, ok := k.(*heap.ConstantMethodRef); ok {
		//直接获取当前类的方法，无需考虑方法的继承和覆盖
		method := kMethodRef.GetMethod(false)
		//执行方法
		frame.Thread.InvokeMethod(method)
	} else {
		//为什么可以调用到接口方法？
		method := k.(*heap.ConstantInterfaceMethodRef).GetMethod(false)
		frame.Thread.InvokeMethod(method)
	}
}
