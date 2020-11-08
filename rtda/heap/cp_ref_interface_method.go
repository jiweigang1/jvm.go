package heap

import (
	"github.com/zxh0/jvm.go/classfile"
)

type ConstantInterfaceMethodRef struct {
	//接口符号引用继承了 方法的接口引用
	ConstantMethodRef
}
// 创建一个接口的符号引用
func newConstantInterfaceMethodRef(class *Class, cf *classfile.ClassFile,
	cfRef classfile.ConstantInterfaceMethodRefInfo) *ConstantInterfaceMethodRef {

	ref := &ConstantInterfaceMethodRef{}
	ref.ConstantMemberRef = newConstantMemberRef(class, cf, cfRef.ClassIndex, cfRef.NameAndTypeIndex)
	ref.ParamSlotCount = calcParamSlotCount(ref.descriptor)
	return ref
}

// todo
// 查找对象的接口的方法，类实现的接口的方法是不能放入到 VTable 中的，所以需要遍历所有的父类来查找方法
func (ref *ConstantInterfaceMethodRef) FindInterfaceMethod(obj *Object) *Method {
	//查找所有的父类来获取方法
	for class := obj.Class; class != nil; class = class.SuperClass {
		method := class.getMethod(ref.name, ref.descriptor, false)
		if method != nil {
			return method
		}
	}
        // 注意 JDK8 以后 接口可以定义 default 方法，所以也需要从接口类背身进行查找
	if method := findInterfaceMethod(obj.Class.Interfaces, ref.name, ref.descriptor); method != nil {
		return method
	} else {
		//TODO
		panic("virtual method not found!")
	}
}
//从接口类中查找方法，处理 default 方法
func findInterfaceMethod(interfaces []*Class, name, descriptor string) *Method {
	for i := 0; i < len(interfaces); i++ {
		//递归查找所有的接口
		if method := findInterfaceMethod(interfaces[i].Interfaces, name, descriptor); method != nil {
			return method
		}
		method := interfaces[i].getMethod(name, descriptor, false)
		if method != nil {
			return method
		}
	}
	return nil
}
