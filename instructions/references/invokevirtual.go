package references

import (
	"github.com/zxh0/jvm.go/instructions/base"
	"github.com/zxh0/jvm.go/rtda"
	"github.com/zxh0/jvm.go/rtda/heap"
)

// Invoke instance method; dispatch based on class
type InvokeVirtual struct {
	base.Index16Instruction
	kMethodRef   *heap.ConstantMethodRef
	argSlotCount uint
}
/**
* 执行实例方法的调用
*/
func (instr *InvokeVirtual) Execute(frame *rtda.Frame) {
	if instr.kMethodRef == nil {
		//获取常量池
		cp := frame.GetConstantPool()
		//从常量池中获取方法信息
		instr.kMethodRef = cp.GetConstant(instr.Index).(*heap.ConstantMethodRef)
		instr.argSlotCount = instr.kMethodRef.ParamSlotCount
	}
	//获取调用方法的对象
	ref := frame.TopRef(instr.argSlotCount)
	if ref == nil {
		frame.Thread.ThrowNPE()
		return
	}
	//获取实例的方法
	method := instr.kMethodRef.GetVirtualMethod(ref)
	frame.Thread.InvokeMethod(method)
}
