## Comparator and Comparable in Java

### 1. 说明

当使用自定义类型，或尝试比较不能直接比较的对象时，我们需要使用比较策略。我们可以简单地通过使用 Comparator 或 Comparable 接口来构建一个。

### 2. 示例

```java
static class Teacher {
  private Long id;
  private String username;
  private Integer age;

  public Teacher(Long id, String username, Integer age) {
    this.id = id;
    this.username = username;
    this.age = age;
  }
}
```



```java
public static void main(String[] args) {
  List<Teacher> teachers = new ArrayList<>();

  Teacher t1 = new Teacher(1L, "Teacher-1", 30);
  Teacher t2 = new Teacher(2L, "Teacher-2", 19);
  Teacher t3 = new Teacher(3L, "Teacher-3", 83);
  Teacher t4 = new Teacher(4L, "Teacher-4", 50);
  teachers.add(t1);
  teachers.add(t2);
  teachers.add(t3);
  teachers.add(t4);

  
  Collections.sort(teachers);
}
```

上面的代码并不能运行，会出现编译时异常：reason: no instance(s) of type variable(s) T exist so that Teacher conforms to Comparable<? super T>。😭



### 3. Comparable

顾名思义，Comparable 是一个接口，定义了将一个对象与其他同类对象进行比较的策略。这被称为类的“自然排序”。

为了能够排序，我们必须通过实现 Comparable 接口将我们的 `Teacher` 对象定义为可比较的：

```java
static class Teacher implements Comparable<Teacher> {
	@Override
  public int compareTo(Teacher o) {
    return Long.compare(getId(), o.getId());
  }
}
```

实现`comparable`接口，并重写`compareTo`方法。问题解决。🌈

排序顺序由 `compareTo()` 方法的返回值决定。`Long.compare(x, y) `如果` x `小于` y` 则返回 -1，如果它们相等则返回 0，否则返回 1。

该方法返回一个数字，指示被比较的对象是小于、等于还是大于作为参数传递的对象。

排序结果：

```java
1 Teacher-1 年龄：30
2 Teacher-2 年龄：19
3 Teacher-3 年龄：83
4 Teacher-4 年龄：50
```

如果我们像按照年龄排序，需要修改`compareTo`方法，很不方便。

### 4. Comparator

*Comparator* 接口定义了一个 *compare(arg1, arg2)* 方法，它有两个参数，代表被比较的对象，其工作方式类似于 *Comparable.compareTo()* 方法。

#### 4.1实现接口

创建一个`Comparator`然后实现`compare`方法。

```java
static class TeacherAgeComparator implements Comparator<Teacher> {
  @Override
  public int compare(Teacher o1, Teacher o2) {
    return Integer.compare(o1.age, o2.age);
  }
}
```

使用：

```java
Collections.sort(teachers, new TeacherAgeComparator());
teachers.forEach(x -> System.out.println(x.getId() + " " + x.getUsername() + " 年龄：" + x.getAge()));

// 输出
// 2 Teacher-2 年龄：19
// 1 Teacher-1 年龄：30
// 4 Teacher-4 年龄：50
// 3 Teacher-3 年龄：83
```

更简单的方法：

```java1
teachers.sort(new TeacherAgeComparator());
```

#### 4.2 Java8 Comparators

```java
Comparator<Teacher> byName = Comparator.comparing(Teacher::getUsername);
Comparator<Teacher> byAge = Comparator.comparing(Teacher::getAge);

teachers.sort(byAge.thenComparing(byName).reversed());

        teachers.forEach(x -> System.out.println(x.getId() + " " + x.getUsername() + " 年龄：" + x.getAge()));

// 输出
// 3 Teacher-3 年龄：83
// 4 Teacher-4 年龄：50
// 1 Teacher-1 年龄：30
// 2 Teacher-2 年龄：19
```

Java 8 提供了使用 lambda 表达式和 *comparing()* 静态工厂方法定义*Comparators*的新方法。



## 结语

更多排序内容可以查看：[列表排序](https://www.jianshu.com/p/737c1038aab1)