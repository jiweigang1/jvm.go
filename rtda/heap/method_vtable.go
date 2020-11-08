package heap
/**
* 返回方法存储的index值
* 通过方法名称和方法描述查找返回 
*/
func getVslot(class *Class, name, descriptor string) int {
	for i, m := range class.vtable {
		if m.Name == name && m.Descriptor == descriptor {
			return i
		}
	}
	// todo
	return -1
}
/**
* 初始化一个 class Vtable ，包含所有的继承的class 方法
* 
*/
func createVtable(class *Class) {
	class.vtable = copySuperVtable(class)

	for _, m := range class.Methods {
		//只复制实例方法
		if isVirtualMethod(m) {
			//如果已经存在，覆盖掉原有的方法
			if i := indexOf(class.vtable, m); i > -1 {
				class.vtable[i] = m // override
			//如果不存在添加到 vtable 中
			} else {
				class.vtable = append(class.vtable, m)
			}
		}
	}

	forEachInterfaceMethod(class, func(m *Method) {
		if i := indexOf(class.vtable, m); i < 0 {
			class.vtable = append(class.vtable, m)
		}
	})
}

func copySuperVtable(class *Class) []*Method {
	if class.SuperClass != nil {
		superVtable := class.SuperClass.vtable
		newVtable := make([]*Method, len(superVtable))
		copy(newVtable, superVtable)
		return newVtable
	} else {
		return nil
	}
}

func isVirtualMethod(method *Method) bool {
	return !method.IsStatic() &&
		//!method.IsFinal() &&
		!method.IsPrivate() &&
		method.Name != constructorName
}
/**
* 从已经存在的方法中，查找是否存在
* 如果存在返回 index 值，否则返回 -1
*/
func indexOf(vtable []*Method, m *Method) int {
	for i, vm := range vtable {
		if vm.Name == m.Name && vm.Descriptor == m.Descriptor {
			return i
		}
	}
	return -1
}

// visit all interface methods
func forEachInterfaceMethod(class *Class, f func(*Method)) {
	for _, iface := range class.Interfaces {
		forEachInterfaceMethod(iface, f)
		for _, m := range iface.Methods {
			f(m)
		}
	}
}
