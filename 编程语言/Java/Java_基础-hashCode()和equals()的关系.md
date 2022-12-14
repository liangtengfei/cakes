## 复习java基础

### hashCode()和equals()的关系

hashCode()用于获取哈希码（散列码），eauqls()用于比较两个对象是否相等。

- 如果两个对象相等，则它们的hashcode必须相同；
- 如果两个对象的hashcode相同，则它们不一定相等。

**这也就是为什么有时候需要重写hashCode和equals方法。因为Object类提供的equals()方法默认是用==来进行比较的，也就是说只有两个对象是同一个对象时，才能返回相等的结果。但是在实际的业务中，有时候需要根据对象的内容是否相同，来判断对象是否相等，所以需要重写这两个方法。**

### “==”和“equals()”的区别

==运算符：

- 作用域基本数据类型时，是比较两个数值是否相等。
- 作用于引用数据类型时，是比较两个对象的内存地址是否相同，即判断它们是否为同一个对象。

equals方法：

- 没有重写时，Object默认以“==”来实现，即比较两个对象的内存地址是否相同；
- 重写后，具体按照重写方法进行比较。



### 两个字符串相加的底层是如何实现的

- 如果拼接的都是字符串直接量，则在编译时编译器会将其直接优化为一个完整的字符串，和你直接写一个完整的字符串是一样的。

- 如果拼接的字符串中包含变量，则在编译时编译器采用StringBuilder对其进行优化，即自动创建StringBuilder实例并调用其append()方法，将这些字符串拼接在一起。