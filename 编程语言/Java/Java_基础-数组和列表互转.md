## Java中数组和列表之间的转换

List和Array作为Java中经常使用到的集合类型。经常需要互相转换。下面我们将学习如何使用核心 Java 库、Guava 和 Apache Commons Collections 在 Array 和 List 之间进行转换。



### List To Array

#### 1. toArray方法

```java
List<String> list = List.of("A", "B", "C");
String[] array = list.toArray(new String[0]);

System.out.println(list);
System.out.println(Arrays.toString(array));

// 输出
// [A, B, C]
// [A, B, C]
```

> 注意：我们使用该方法的首选方式是 toArray(new T[0]) 与 toArray(new T[size])。
>
> toArray(new T[0]) 看起来更快、更安全，因此现在应该是默认选择。——[Collection.toArray(new T[0]) or Collection.toArray(new T[size]), that’s the question](https://shipilev.net/blog/2016/arrays-wisdom-ancients/#_conclusion)

#### 2. 使用Guava

```java
List<Integer> integers = List.of(1, 2, 3, 4, 5);
int[] ints = Ints.toArray(integers);

System.out.println(integers);
System.out.println(Arrays.toString(ints));

// 输出
// [1, 2, 3, 4, 5]
// [1, 2, 3, 4, 5]
```

增加了自动拆装箱操作，可以直接转成原始数据类型数组。

### Array To List

#### 1. Arrays.asList方法

```java
Integer[] integers1 = {1, 2, 3, 4, 5};
List<Integer> integerList = Arrays.asList(integers1);

System.out.println(integers1.length == integerList.size());

// 输出
// true
```

> 注意：`Arrays.asList`返回由指定数组支持的固定大小列表。所以不能调用`add`等方法，会报运行时异常：`java.lang.UnsupportedOperationException`。

如果需要一个标准的`ArrayList`,可以使用下面的方法。

```java
ist<Integer> notFixedList = new ArrayList<>(Arrays.asList(integers1));
notFixedList.add(6);
System.out.println(notFixedList);

// 输出
// [1, 2, 3, 4, 5, 6]
```

#### 2. 使用Guava

```java
Integer[] integers2 = {1, 2, 3, 4, 5};
List<Integer> integerList2 = Lists.newArrayList(integers2);

System.out.println(Arrays.toString(integers2));
System.out.println(integerList2);

integerList2.add(6);
System.out.println(integerList2);

// 输出
// [1, 2, 3, 4, 5]
// [1, 2, 3, 4, 5]
// [1, 2, 3, 4, 5, 6]
```

#### 3. 使用Collections

Guava的`Lists.newArrayList`内部调用的是

```java
ArrayList<E> list = new ArrayList(capacity);
Collections.addAll(list, elements);
```

所以这也是第三个方法，使用`Collections.addAll`。



