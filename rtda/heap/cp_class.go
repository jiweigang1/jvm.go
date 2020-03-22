package heap

import (
	"github.com/zxh0/jvm.go/classfile"
)

type ConstantClass struct {
	class    *Class
	name     string
	resolved *Class
}

func newConstantClass(class *Class, cf *classfile.ClassFile,
	cfc classfile.ConstantClassInfo) *ConstantClass {

	return &ConstantClass{
		class: class,
		name:  cf.GetUTF8(cfc.NameIndex),
	}
}
/**
* 获取class常量对应 class内容
* 如果class内容还没有关联，执行关联
*/
func (ref *ConstantClass) GetClass() *Class {
	if ref.resolved == nil {
		ref.resolve()
	}
	return ref.resolved
}

// todo
func (ref *ConstantClass) resolve() {
	// load class
	ref.resolved = ref.class.bootLoader.LoadClass(ref.name)
}
