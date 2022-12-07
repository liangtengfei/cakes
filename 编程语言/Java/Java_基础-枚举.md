## 枚举
在本教程中，我们将了解什么是 Java 枚举、它们解决的问题以及它们的一些设计模式如何在实践中使用。

### 1. 概述

Java 5 首先引入了 enum 关键字。它表示一种特殊类型的类，它总是扩展 java.lang.Enum 类。有关使用的官方文档，我们可以转到[文档](https://docs.oracle.com/en/java/javase/17/docs/api/java.base/java/lang/Enum.html)。 以这种方式定义的常量使代码更具可读性，允许进行编译时检查，预先记录可接受值的列表，并避免由于传入无效值而导致的意外行为。

### 2. 定义

最简单的定义一个枚举类。
```java
public enum WeekDay {    
    MONDAY, 
    TUESDAY, 
    WEDNESDAY,
    THURSDAY, 
    FRIDAY, 
    SATURDAY,
    SUNDAY
}
```

