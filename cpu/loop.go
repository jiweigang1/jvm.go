package cpu

import (
	"fmt"

	"github.com/zxh0/jvm.go/instructions"
	"github.com/zxh0/jvm.go/instructions/base"
	"github.com/zxh0/jvm.go/rtda"
	"github.com/zxh0/jvm.go/rtda/heap"
	"github.com/zxh0/jvm.go/vm"
)

// 执行方法调用
func ExecMethod(thread *rtda.Thread, method *heap.Method, args []heap.Slot) heap.Slot {
	// 生成新的栈桢
	shimFrame := rtda.NewShimFrame(thread, args)
	//栈桢压入栈顶
	thread.PushFrame(shimFrame)
	thread.InvokeMethod(method)

	debug := thread.VMOptions.XDebugInstr
	defer _catchErr(thread) // todo

	for {
		frame := thread.CurrentFrame()
		if frame == shimFrame {
			thread.PopFrame()
			if frame.IsStackEmpty() {
				return heap.EmptySlot
			} else {
				return frame.Pop()
			}
		}

		pc := frame.NextPC
		thread.PC = pc

		// fetch instruction
		instr, nextPC := fetchInstruction(frame.Method, pc)
		frame.NextPC = nextPC

		// execute instruction
		instr.Execute(frame)
		if debug {
			_logInstruction(frame, instr)
		}
	}
}

func Loop(thread *rtda.Thread) {
	threadObj := thread.JThread()
	isDaemon := threadObj != nil && threadObj.GetFieldValue("daemon", "Z").IntValue() == 1
	if !isDaemon {
		nonDaemonThreadStart()
	}

	_loop(thread)

	// terminate thread
	threadObj = thread.JThread()
	threadObj.Monitor.NotifyAll()
	if !isDaemon {
		nonDaemonThreadStop()
	}
}
/**
循环执行一个线程中的指令
*/
func _loop(thread *rtda.Thread) {
	debug := thread.VMOptions.XDebugInstr
	defer _catchErr(thread) // todo

	for {
		//获取当前的栈桢
		frame := thread.CurrentFrame()
		//把当前栈桢的计数器设置到当前的线程的计数器上
		pc := frame.NextPC
		thread.PC = pc

		// fetch instruction 获取当前的指令，同时返回下一个计数器的值
		instr, nextPC := fetchInstruction(frame.Method, pc)
		//栈桢上设置下一个指令的位置
		frame.NextPC = nextPC

		// execute instruction 执行指令
		instr.Execute(frame)
		if debug {
			_logInstruction(frame, instr)
		}
		// 如果当前线程的栈为空，证明当前的线程的指令（方法）已经执行完毕，结束这个线程
		if thread.IsStackEmpty() {
			break
		}
	}
}
/**
* 根据计数器获指定线程的 指令
*/
func fetchInstruction(method *heap.Method, pc int) (base.Instruction, int) {
	//如果方法的指令还没有转换，从二进制的字节进行转换
	if method.Instructions == nil {
		//执行指令的解析的时候也会解析 指令自身的参数信息 比如 对应的常量池的 Index。
		method.Instructions = instructions.Decode(method.Code)
	}

	instrs := method.Instructions.([]base.Instruction)
	//获取指令
	instr := instrs[pc]

	// calc nextPC
	pc++
	//计数器，指向下个不为空的指令
	for pc < len(instrs) && instrs[pc] == nil {
		pc++
	}

	return instr, pc
}

// todo
func _catchErr(thread *rtda.Thread) {
	if r := recover(); r != nil {
		if err, ok := r.(vm.ClassNotFoundError); ok {
			thread.ThrowClassNotFoundException(err.Error())
			_loop(thread)
			return
		}

		_logFrames(thread)

		err, ok := r.(error)
		if !ok {
			err = fmt.Errorf("%v", r)
			panic(err.Error())
		} else {
			panic(err.Error())
		}
	}
}

func _logFrames(thread *rtda.Thread) {
	for !thread.IsStackEmpty() {
		frame := thread.PopFrame()
		method := frame.Method
		className := method.Class.Name
		lineNum := method.GetLineNumber(frame.NextPC)
		fmt.Printf(">> line:%4d pc:%4d %v.%v%v \n",
			lineNum, frame.NextPC, className, method.Name, method.Descriptor)
	}
}

func _logInstruction(frame *rtda.Frame, instr base.Instruction) {
	thread := frame.Thread
	method := frame.Method
	className := method.Class.Name
	pc := thread.PC

	if method.IsStatic() {
		fmt.Printf("[instruction] thread:%p %v.%v() #%v %T %v\n",
			thread, className, method.Name, pc, instr, instr)
	} else {
		fmt.Printf("[instruction] thread:%p %v#%v() #%v %T %v\n",
			thread, className, method.Name, pc, instr, instr)
	}
}
