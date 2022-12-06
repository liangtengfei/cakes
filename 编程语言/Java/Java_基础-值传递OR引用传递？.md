## Java 中的按值传递作为参数传递机制

### 1. 介绍

将参数传递给方法的两种最普遍的模式是**“按值传递”**和**“按引用传递”**。不同的编程语言以不同的方式使用这些概念。就 Java 而言，***一切都是严格的值传递***。下面，我们将说明 Java 如何为各种类型传递参数。

### 2. 按值传递与按引用传递

现代编程语言中最常见的两种机制是“按值传递”和“按引用传递”。在我们继续之前，让我们先讨论这些：

#### 2.1 值传递

当参数是按值传递时，调用方和被调用方方法对两个不同的变量进行操作，这两个变量是彼此的副本。对一个变量的任何更改都不会修改另一个变量。 

这意味着在调用方法时，**传递给被调用方方法的参数将是原始参数的克隆**。在被调用方方法中所做的任何修改都不会影响调用方方法中的原始参数。

#### 2.2 引用传递

当参数是按引用传递时，调用者和被调用者对同一个对象进行操作。

 这意味着当一个变量被引用传递时，**对象的唯一标识符被发送到方法**。对参数实例成员的任何更改都将导致对原始值进行更改。

### 3. Java中的参数传递

在 Java 中，**原始变量存储实际值，而非原始变量存储引用变量，这些引用变量指向它们所引用的对象的地址。**值和引用都存储在堆栈内存中。

Java 中的参数始终按值传递。在方法调用期间，每个参数的副本，无论是值还是引用，都会在堆栈内存中创建，然后传递给方法。

如果是原始的，则简单地将其复制到堆栈内存中，然后将其传递给被调用方方法。如果是非原始类型，则堆栈内存中的应用指向堆中的实际数据。当我们通过对象时，堆栈内存中的引用将复制，并将新引用传递给该方法。

#### 3.1 原始类型传递

Java 编程语言具有八种原始数据类型。**原始变量直接存储在堆栈内存中**。每当任何原始数据类型的变量作为参数传递时，实际参数都会被复制到形参中，并且这些形参会在堆栈内存中累积自己的空间。

只要该方法正在运行，这些形式参数的生命周期就会持续，在返回时，这些形式参数将从堆栈中清除并被丢弃。

```java
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertEquals;

public class JavaPassByValueTests {
    @Test
    public void originalValuesNotModified() {
        int a = 1, b = 2;

        assertEquals(a, 1);
        assertEquals(b, 2);

        numberAddOne(a, b);

        assertEquals(a, 1);
        assertEquals(b, 2);
    }

    private static void numberAddOne(int num1, int num2) {
        num1 += 1;
        num2 += 1;
    }
}
```

让我们通过分析这些值在内存中的存储方式来尝试理解上述程序中的断言：

1. "a"和“b”变量是原始类型数据，它们的值存储在堆栈内存中。
2. 当我们调用**numberAddOne()**方法时，每个变量的精确副本被创建并存储在堆栈内存中的不同位置。
3. 对这些副本的任何修改只会影响变量的副本，而不会改变原始变量

#### 3.2 对象传递

在 Java 中，所有对象都动态存储在堆空间中。这些对象是从称为引用变量的引用中引用的。😵😵😵

与 原始类型数据相比，Java 对象分两个阶段存储。引用变量存储在堆栈内存中，它们所引用的对象存储在堆内存中。

每当将对象作为参数传递时，都会创建引用变量的精确副本，该副本**指向对象在堆内存中与原始引用变量相同的位置**。

因此，每当我们在方法中对同一对象进行任何更改时，该更改都会反映在原始对象中。但是，如果我们为传递的引用变量分配一个新对象，那么它就不会反映在原始对象中。🌈

```java
@Test
public void passingObject() {
    Person p1 = new Person(1L, "唐僧");
    Person p2 = new Person(2L, "孙悟空");

    assertEquals(p1.username, "唐僧");
    assertEquals(p2.username, "孙悟空");

    modifyPerson(p1, p2);

    assertEquals(p1.username, "唐三藏");
    assertEquals(p2.username, "齐天大圣");
    
    modifyPersonNo(p1, p2);

    assertEquals(p1.username, "唐三藏");
    assertNotEquals(p2.username, "六耳猕猴");
}

private static void modifyPerson(Person p1, Person p2) {
    p1.setUsername("唐三藏");
    p2.setUsername("齐天大圣");
}

private static void modifyPersonNo(Person p1, Person p2) {
    p1.setUsername("唐三藏");
    p2 = new Person(p2.id, "六耳猕猴");
    p2.id++;
}
```

我们来分析一下上面程序中的断言。

调用第一个方法修改*Person*对象时，会创建*p1*和*p2*这两个引用变量的副本，它们指向相同的原始对象。

调用第二个方法时，也会创建引用变量的副本，但是，*p2*变量，分配了一个新的对象，和变量副本指向的已经不是同一个原始对象，所以对*p2*的修改不会改变原始对象。

## 结语

我们了解到，Java 中的参数传递始终是按值传递。但是，上下文会根据我们处理的是原始数据还是对象而变化：

1. 对于 原始类型，参数是按值传递的
2. 对于对象类型，对象引用是按值传递的

🎉🎉🎉🎉🎉🎉🎉🎉🎉🎉🎉🎉🎉🎉🎉🎉🎉

参考资料：https://www.baeldung.com/java-pass-by-value-or-pass-by-reference 以及 https://docs.oracle.com/javase/tutorial/java/javaOO/arguments.html