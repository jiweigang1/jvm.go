package references

import (
	"github.com/zxh0/jvm.go/instructions/base"
	"github.com/zxh0/jvm.go/rtda"
	"github.com/zxh0/jvm.go/rtda/heap"
)

// Create new object
type New struct {
	base.Index16Instruction
	class *heap.Class
}

func (instr *New) Execute(frame *rtda.Frame) {
	//如果指令还没有关联class
	if instr.class == nil {
		//获取常量池
		cp := frame.GetConstantPool()
		//获取常量池中的 class 
		kClass := cp.GetConstantClass(instr.Index)
		//获取常量池中class信息 真实关联的class信息
		instr.class = kClass.GetClass()
	}

	// init class
	if instr.class.InitializationNotStarted() {
		frame.RevertNextPC() // undo new
		frame.Thread.InitClass(instr.class)
		return
	}
    //创建对象
	ref := instr.class.NewObj()
	//对象引用压入操作数栈
	frame.PushRef(ref)
}
