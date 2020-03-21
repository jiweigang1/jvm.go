package heap

import (
	"github.com/zxh0/jvm.go/classfile"
)

type ConstantMethodRef struct {
	ConstantMemberRef
	ParamSlotCount uint
	resolved       *Method
	vslot          int
}

func newConstantMethodRef(class *Class, cf *classfile.ClassFile,
	cfRef classfile.ConstantMethodRefInfo) *ConstantMethodRef {

	ref := &ConstantMethodRef{vslot: -1}
	ref.ConstantMemberRef = newConstantMemberRef(class, cf, cfRef.ClassIndex, cfRef.NameAndTypeIndex)
	ref.ParamSlotCount = calcParamSlotCount(ref.descriptor)
	return ref
}

func (ref *ConstantMethodRef) GetMethod(static bool) *Method {
	if ref.resolved == nil {
		if static {
			ref.resolveStaticMethod()
		} else {
			ref.resolveSpecialMethod()
		}
	}
	return ref.resolved
}

func (ref *ConstantMethodRef) resolveStaticMethod() {
	method := ref.findMethod(true)
	if method != nil {
		ref.resolved = method
	} else {
		// todo
		panic("static method not found!")
	}
}

func (ref *ConstantMethodRef) resolveSpecialMethod() {
	method := ref.findMethod(false)
	if method != nil {
		ref.resolved = method
		return
	}

	// todo
	// class := ref.cp.class.classLoader.LoadClass(ref.className)
	// if class.IsInterface() {
	// 	method = ref.findMethodInInterfaces(class)
	// 	if method != nil {
	// 		ref.method = method
	// 		return
	// 	}
	// }

	// todo
	panic("special method not found!")
}

func (ref *ConstantMethodRef) findMethod(isStatic bool) *Method {
	class := ref.getBootLoader().LoadClass(ref.className)
	return class.getMethod(ref.name, ref.descriptor, isStatic)
}

// todo
/*func (mr *ConstantMethodref) findMethodInInterfaces(iface *Class) *Method {
	for _, m := range iface.methods {
		if !m.IsAbstract() {
			if m.name == mr.name && m.descriptor == mr.descriptor {
				return m
			}
		}
	}

	for _, superIface := range iface.interfaces {
		if m := mr.findMethodInInterfaces(superIface); m != nil {
			return m
		}
	}

	return nil
}*/
/**
* 
*/
func (ref *ConstantMethodRef) GetVirtualMethod(obj *Object) *Method {
	if ref.vslot < 0 {
		//查找方法的存储的index
		ref.vslot = getVslot(obj.Class, ref.name, ref.descriptor)
	}
	// 如果index 存在，证明已经加载，直接返回方法
	if ref.vslot >= 0 {
		return obj.Class.vable[ref.vslot]
	}

	// TODO: invoking private method ?
	//println("GetVirtualMethod:", ref.className, ref.name, ref.descriptor)
	class := ref.getBootLoader().LoadClass(ref.className)
	//获取类声明的方法
	return class.getDeclaredMethod(ref.name, ref.descriptor, false)
}
