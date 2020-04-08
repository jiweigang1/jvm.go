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
//调用静态方法
func (instr *InvokeStatic) Execute(frame *rtda.Frame) {
	//如果关联的方法为空的化进行关联，静态方法可以直接从符号应用进行关联的，静态方法不存在覆盖的问题可以直接获取到指定的方法
	if instr.method == nil {
		//获取常量池
		cp := frame.GetConstantPool()
		k := cp.GetConstant(instr.Index)
		//如果是类的静态方法
		if kMethodRef, ok := k.(*heap.ConstantMethodRef); ok {
			instr.method = kMethodRef.GetMethod(true)
		//如果是接口的构造方法
		} else {
			instr.method = k.(*heap.ConstantInterfaceMethodRef).GetMethod(true)
		}
	}

	// init class
	//如果class没有初始化，初始化class
	class := instr.method.Class
	if class.InitializationNotStarted() {
		frame.RevertNextPC()
		frame.Thread.InitClass(class)
		return
	}

	frame.Thread.InvokeMethod(instr.method)
}
