## Java 中的数字格式化

学习Java中不同的数字格式化方法，以及如何实现。

### 1. 基础方法：String.format

String.format 方法对于格式化数字非常有用。并且不仅仅在格式化数字上有用，更详细的使用可以查看[Java基础-格式化输出](https://www.jianshu.com/p/09e17876ab24)。

String.format[官方语法说明](https://docs.oracle.com/en/java/javase/17/docs/api/java.base/java/util/Formatter.html#syntax)。一般、字符和数字类型的格式说明符具有以下语法：

```java
%[argument_index$][flags][width][.precision]conversion
```

言归正传，下面看看在数字格式化方面能有什么作用吧。可以用来格式化数字的格式化符号有以下几个：

| 格式化符号 | 是否常用 | 用途                                                         |
| ---------- | -------- | ------------------------------------------------------------ |
| d          | 是       | 结果被格式化为十进制整数，用于整数                           |
| o          | 否       | 结果被格式化为八进制整数，用于整数                           |
| x, X       | 否       | 结果被格式化为十六进制整数，用于整数                         |
| e, E       | 否       | 结果被格式化为计算机科学记数法中的十进制数，用于浮点数       |
| f          | 是       | 结果被格式化为十进制数，用于浮点数                           |
| g, G       | 否       | 根据四舍五入后的精度和值，使用计算机化的科学记数法或十进制格式对结果进行格式化，用于浮点数 |
| a, A       | 否       | 结果被格式化为带有尾数和指数的十六进制浮点数。 BigDecimal 类型不支持此转换，尽管后者属于浮点参数类别，用于浮点数 |
| %          | 是       | 输出'%'                                                      |

实例1：格式化正数，左对齐输出

```java
int num = 2022;
int num2 = 22;
System.out.printf("这个一个整数：%4d\n", num);
System.out.printf("这个一个整数：%04d\n", num2);

// 输出
// 这个一个整数：2022
// 这个一个整数：0022
```

实例2：格式化浮点数

```java
double pi = 3.1415d;
double pi2 = 3.1415926d;
// 如果不指定精度 默认保留6位小数，不足6位的后面补0
System.out.printf("原样输出：%f\n", pi);
System.out.printf("原样输出：%f\n", pi2);

// 四舍五入 保留5位小数
System.out.printf("原样输出：%.5f\n", pi);
System.out.printf("原样输出：%.5f\n", pi2);

// 输出
// 默认输出：3.141500
// 默认输出：3.141593
// 指定精度：3.14150
// 指定精度：3.14159
```

### 2. 使用BigDecimal

BigDecimal 类提供用于算术、比例操作、舍入、比较、散列和格式转换的操作。 toString() 方法提供了 BigDecimal 的规范表示。

BigDecimal[官方语法说明](https://docs.oracle.com/en/java/javase/17/docs/api/java.base/java/math/BigDecimal.html)。

> 在阿里云的开发手册中写道：
>
> 【强制】禁止使用构造方法 BigDecimal(double) 的方式把 double 值转化为 BigDecimal 对象。 说明：BigDecimal(double) 存在精度损失风险，在精确计算或值比较的场景中可能会导致业务逻辑异常。

```java
float f1 = 3.141592678F;

// 新建一个BigDecimal对象，位了避免丢失精度，使用根据浮点数的toString()方法
BigDecimal bigDecimal1 = new BigDecimal(Float.toString(f1));
// 四舍五入 保留5位小数 第一格参数指定精度，第二个参数指定舍入模式
bigDecimal1 = bigDecimal1.setScale(5, RoundingMode.HALF_UP);

System.out.println(3.14159 == bigDecimal1.doubleValue());

// 输出
// true

double D = 4.2352989244D;
assertThat(withBigDecimal(D, 2)).isEqualTo(4.24);
assertThat(withBigDecimal(D, 3)).isEqualTo(4.235);

public static double withBigDecimal(double value, int places) {
    BigDecimal bigDecimal = new BigDecimal(value);
    bigDecimal = bigDecimal.setScale(places, RoundingMode.HALF_UP);
    return bigDecimal.doubleValue();
}
```

### 3. 使用Math.round

Math.round[官方语法说明](https://docs.oracle.com/en/java/javase/17/docs/api/java.base/java/lang/Math.html#round(double))。

使用Math静态方法去格式化浮点数，不是直接的而是通过“乘除法”。

```java
double scale1 = Math.pow(10, 5);
double newPi1 = Math.round(pi2 * scale1) / scale1;
System.out.println("四舍五入 保留5位小数：" + newPi1);

double scale2 = Math.pow(10, 3);
double newPi2 = Math.round(pi2 * scale2) / scale2;
System.out.println("四舍五入 保留3位小数：" + newPi2);

// 输出
// 四舍五入 保留5位小数：3.14159
// 四舍五入 保留3位小数：3.142
```

> 这种方法要慎用，如果数值本身就比较大，再进行乘法的话，很容易就溢出了，造成结果错误。
>
> 再有就是round，是用来取整的，会对数据进行截断，也会造成数据不是预期。

所以请注意，列出此方法仅用于学习目的。

### 4. 使用DecimalFormat

DecimalFormat 是 NumberFormat 的具体子类，用于格式化十进制数。它具有多种功能，旨在使在任何区域设置中解析和格式化数字成为可能，包括对西方、阿拉伯和印度数字的支持。它还支持不同类型的数字，包括整数 (123)、定点数 (123.4)、科学记数法 (1.23E4)、百分比 (12%) 和货币金额 ($123)。所有这些都可以本地化。

功能非常强大，DecimalFormat[官方语法说明](https://docs.oracle.com/en/java/javase/17/docs/api/java.base/java/text/DecimalFormat.html)。

#### 4.1 用逗号格式化大整数

```java
// 使用默认语言环境的模式。
DecimalFormat decimalFormatDefault = new DecimalFormat();
System.out.println("默认输出：" + decimalFormatDefault.format(123456789L));

// 输出
// 默认输出：123,456,789
```

#### 4.2 在小数点后面填充“0”

```java
DecimalFormat decimalFormatPattern2 = new DecimalFormat("#.00");
System.out.println("自定义格式：" + decimalFormatPattern2.format(11D));
System.out.println("自定义格式：" + decimalFormatPattern2.format(12.0));

// 输出
// 自定义格式：11.00
// 自定义格式：12.00
```

### 5. 百分比格式化

实例1: 使用NumberFormat

NumberFormat 是所有数字格式的抽象基类。此类提供用于格式化和解析数字的接口。 NumberFormat 还提供了一些方法来确定哪些语言环境具有数字格式，以及它们的名称是什么。

```java
// 返回当前默认 FORMAT 语言环境的百分比格式。
NumberFormat numberFormatPercent = NumberFormat.getPercentInstance();
System.out.println("默认环境输出：" + numberFormatPercent.format(0.15));
System.out.println("默认环境输出：" + numberFormatPercent.format(25F / 100F));

// 输出
// 默认环境输出：15%
// 默认环境输出：25%
```

实例2: 使用String.format

```java
System.out.printf("增加百分号：%d%%", 25);
// 输出
// 增加百分号：25%
```

> 可以看到两种方法的区别，NumberFormat格式化之前会先进行计算，String.format只会格式化输出，不会对数字进行计算。

### 6. 货币数字格式化

```java
// 返回当前默认 FORMAT 语言环境的货币格式。
NumberFormat numberFormatCurrency = NumberFormat.getCurrencyInstance();
System.out.println("默认环境输出：" + numberFormatCurrency.format(123.45D));

NumberFormat numberFormatCurrencyUS = NumberFormat.getCurrencyInstance(Locale.US);
System.out.println("美国环境输出：" + numberFormatCurrencyUS.format(123.45D));

// 输出
// 默认环境输出：¥123.45
// 美国环境输出：$123.45
```

如果数字比较大，会输出位科学记数法格式。

```java
// 返回当前默认 FORMAT 语言环境的货币格式。
NumberFormat numberFormatCurrency = NumberFormat.getCurrencyInstance();
System.out.println("默认环境输出：" + numberFormatCurrency.format(123456.78D));

NumberFormat numberFormatCurrencyUS = NumberFormat.getCurrencyInstance(Locale.US);
System.out.println("美国环境输出：" + numberFormatCurrencyUS.format(123456.78D));

// 输出
// 默认环境输出：¥123,456.78
// 美国环境输出：$123,456.78
```



## 结尾

Java 中数字格式化的不同方式。正如我们所看到的，没有一种最好的方法可以做到这一点。可以使用多种方法，因为每种方法都有自己的特点。
