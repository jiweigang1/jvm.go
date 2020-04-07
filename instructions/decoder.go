package instructions

import (
	"github.com/zxh0/jvm.go/instructions/base"
)
// 从方法 code 二进制流中解析
func Decode(code []byte) []base.Instruction {
	reader := base.NewCodeReader(code)
	decoded := make([]base.Instruction, len(code))
        //循环生成指令
	for reader.Position() < len(code) {
		decoded[reader.Position()] = decodeInstruction(reader)
	}

	return decoded
}
//获取一个指令
func decodeInstruction(reader *base.CodeReader) base.Instruction {
	//获取指令码
	opcode := reader.ReadUint8()
	//生成一个指令
	instr := newInstruction(opcode)
	//获取指令的参数信息
	instr.FetchOperands(reader)
	return instr
}
