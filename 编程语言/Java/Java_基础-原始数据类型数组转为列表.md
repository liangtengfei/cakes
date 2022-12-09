## 将原始数组转换为列表

在这个简短的教程中，我们将展示如何将原始数组转换为相应类型的对象列表。通常，我们可以尝试在Java中使用**自动装箱**。然而，我们对自动装箱工作原理的直觉常常是错误的。

首先看一下下面的代码，有没有问题呢？

```java
int[] ints = new int[]{1, 2, 3, 4, 5};
List<Integer> integers = Arrays.asList(ints);
```

答案是肯定的，会出现编译错误：`no instance(s) of type variable(s) exist so that int[] conforms to Integer inference variable T has incompatible bounds: equality constraints: Integer lower bounds: int[]`。

**原因**：自动装箱只发生在单个元素上（例如从 int 到 Integer）。没有从基本类型数组到其装箱引用类型数组的自动转换（例如从 int[] 到 Integer[]）。

### 解决方案

#### 1. 迭代

```java
int[] ints = new int[]{1, 2, 3, 4, 5};
List<Integer> integers = new ArrayList<>();
for (int i : ints) {
  integers.add(i);
}
System.out.println(Arrays.toString(ints));
System.out.println(integers);
```

#### 2. Streams

从 Java 8 开始，我们可以使用 Stream API。

```java
int[] ints = new int[]{1, 2, 3, 4, 5};
List<Integer> integers = Arrays.stream(ints).boxed().toList();

// 也可以使用IntStream
List<Integer> integers2 = IntStream.of(ints).boxed().toList();
```

> 注意：上面两个方法返回的List都不可修改。这也是Stream的特点吧。

#### 3. Guava

```java
int[] ints = new int[]{1, 2, 3, 4, 5};
List<Integer> integers = Ints.asList(ints);
```

> 同样的，返回的List也不能修改。



## 结语

如果需要返回可以修改的列表，可以使用下面的方法。

```java
int[] ints = new int[]{1, 2, 3, 4, 5};
List<Integer> integers = new ArrayList<>(Ints.asList(ints));
integers.add(6);

System.out.println(Arrays.toString(ints));
System.out.println(integers);
```

`new ArrayList<>`新建实例。