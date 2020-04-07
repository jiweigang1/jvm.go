package references

import (
	"github.com/zxh0/jvm.go/instructions/base"
	"github.com/zxh0/jvm.go/rtda"
	"github.com/zxh0/jvm.go/rtda/heap"
)

// Invoke instance method; dispatch based on class
type InvokeVirtual struct {
	base.Index16Instruction
	//方法的符合引用，包含方法的 描述等信息
	kMethodRef   *heap.ConstantMethodRef
	//方法参数的 slot数量，用于获取从操作数栈中获取参数
	argSlotCount uint
}
/**
* 执行实例方法的调用
*/
func (instr *InvokeVirtual) Execute(frame *rtda.Frame) {
        //如果指令还没有执行关联常量池中的内容，进行初始化。
	if instr.kMethodRef == nil {
		//获取常量池
		cp := frame.GetConstantPool()
		//获取方法信息
		instr.kMethodRef = cp.GetConstant(instr.Index).(*heap.ConstantMethodRef)
		//获取方法的参数的 slot数量
		instr.argSlotCount = instr.kMethodRef.ParamSlotCount
	}
	//从调用的栈帧中获被取调用方法的对象，被调用的对象是在栈顶，根据方法参数 slot的数量就可以从操作数栈中获取出来参数信息
	ref := frame.TopRef(instr.argSlotCount)
	// 如果对象为空，直接抛出异常
	if ref == nil {
		frame.Thread.ThrowNPE()
		return
	}
	// 这里执行 方法的关联，返回的 mehod 是包含执行的 code 的字节码的，这里也会处理方法的覆盖等内容。
	method := instr.kMethodRef.GetVirtualMethod(ref)
	// 执行方法
	frame.Thread.InvokeMethod(method)
}
