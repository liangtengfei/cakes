原文不知道在哪，记录备查。

**作者 | 百度小程序团队**

> 导读 本文收集一些使用 Go 开发过程中非常容易踩坑的 case，所有的 case 都有具体的代码示例，以及针对的代码修复方法，以避免大家再次踩坑。通常这些坑的特点就是代码正常能编译，但运行结果不及预期或是引入内存漏洞的风险。
>
> *全文 7866 字，预计阅读时间 20 分钟。*

# **01 参数传递误用**

## **1.1 误对指针计算 Sizeof**

对任何指针进行 unsafe.Sizeof 计算，返回的结果都是 8 (64 位平台下)。稍不注意就会引发错误。

错误示例：

```
func TestSizeofPtrBug(t *testing.T) {
    type CodeLocation struct {
        LineNo int64
        ColNo  int64
    }
    cl := &CodeLocation{10, 20}
    size := unsafe.Sizeof(cl)
    fmt.Println(size) // always return 8 for point size
}
```

建议使用示例：单独编写一个只处理值大小的函数 ValueSizeof。

```
func TestSizeofPtrWithoutBug(t *testing.T) {
    type CodeLocation struct {
        LineNo int64
        ColNo  int64
    }
    cl := &CodeLocation{10, 20}
    size := ValueSizeof(cl)
    fmt.Println(size) // 16
}

func ValueSizeof(v any) uintptr {
    typ := reflect.TypeOf(v)
    if typ.Kind() == reflect.Pointer {
        return typ.Elem().Size()
    }
    return typ.Size()

}
```

## **1.2 可变参数为 any 类型时，误传切片对象**

当参数的可变参数是 any 类型时，传入切片对象时一定要用展开方式。

```
    appendAnyF := func(t []any, toAppend ...any) []any {
        ret := append(t, toAppend...)
        return ret
    }

    emptySlice := []any{}
    slice2 := []any{"hello", "world"}

    // bug append slice as a element
    emptySlice = appendAnyF(emptySlice, slice2)
    fmt.Println(emptySlice) // only 1 element [[hello world]]

    emptySlice = []any{}
    emptySlice = appendAnyF(emptySlice, slice2...)
    fmt.Println(emptySlice) // [hello world]
```

## **1.3 数组是值传递**

数组在函数或方法中入参传递是值复制的方式，不能用入参的方式进函数或方法内修改数组内容进行返回的。

示例代码如下：

```
    arr := [3]int{0, 1, 2}
    f := func(v [3]int) {
        v[0] = 100
    }
    f(arr)           // no modify to arr
    fmt.Println(arr) // [0 1 2]
```

## **1.4 切片扩容后会新申请内存，不再与内存引用有任何关联**

这里坑在，如果从一个数组中引入一个切片，一旦这个切片引发扩容后，则与原来的引用内容没有任何关系。

```
    arr := []int{0, 1, 2}
    f := func(v []int) {
        v[0] = 100// can modify origin array
        v = append(v, 4) // new memory allocated
        v[0] = 50// no modify to origin array
    }
    f(arr)
    fmt.Println(arr) // [100 1 2]
```

上面的示例代码，扩容切片前对内容的修改可以影响到 arr 数组，说明是共享内存地址引用的，一旦扩容后，则是重新申请了内存，与数组不再是一个内存引用了。

## **1.5 返回参数尽量避免使用共享数据的切片对象，容易导致原始数据污染**

这种场景就是如果通过函数返回值方式从一个大数组获取部分内部，尽量不要用切片共享的方式，可以使用 copy 的方式来替换。

下面的代码，通过 ReadUnsafe 读取切片后，修改内容同步影响原始的内容。

```
type Queue struct {
    content []byte
    pos     int
}

func (q *Queue) ReadUnsafe(size int) []byte {
    if q.pos+size >= len(q.content) {
        return nil
    }
    pos := q.pos
    q.pos = q.pos + size
    return q.content[pos:q.pos]
}

func TestReadUnsafe(t *testing.T) {
    c := [200]byte{}
    q := &Queue{content: c[:]}
    v := q.ReadUnsafe(10)
    v[0] = 1

    fmt.Println(q.content[0]) // 1  q.content值已经被修改
}
```

正确的修改如下，使用 copy 创建一份新内存：

```
func (q *Queue) ReadSafe(size int) []byte {
    if q.pos+size >= len(q.content) {
        return nil
    }
    pos := q.pos
    q.pos = q.pos + size

    ret := make([]byte, size)
    copy(ret, q.content[pos:q.pos])
    return ret
}

func TestReadSafe(t *testing.T) {
    c := [200]byte{}
    q := &Queue{content: c[:]}
    v := q.ReadSafe(10)
    v[0] = 1

    fmt.Println(q.content[0]) // 0  q.content值安全
}
```

# **02 指针相关使用的坑**

## **2.1 误保存 uintptr 值**

uintptr 保存的当前地址的一个整型值，它一旦被获取后，是不会被编译器感知的，也就是它就是一个普通变量，不会追溯内存真实地址变化。

```
    slice := []int{0, 1, 2}
    ptr := unsafe.Pointer(&slice[0]) // get array element:0 pointer

    slice = append(slice, 3) // allocate new memory
    ptr2 := unsafe.Pointer(&slice[0])

    // ptr is 824633770392, ptr2 is 824633762896, ptr==ptr2 result is false
    fmt.Println(fmt.Sprintf("ptr is %d, ptr2 is %d, ptr==ptr2 result is %v", ptr, ptr2, ptr == ptr2))
```

## **2.2 len 与 cap 对空指针 nil 与空值返回相同**

针对切片， 用 len 与 cap 操作时，空值与 nil 都是返回 0， 针对 map， 用 len 操作时，空值与 nil 都是返回 0。

```
     var slice []int = nil
    fmt.Println(len(slice), cap(slice)) // 0 0

    var slice2 []int = []int{}
    fmt.Println(len(slice2), cap(slice2)) // 0 0

    var mp map[int]int = nil
    fmt.Println(len(mp)) // 0

    var mp2 map[int]int = map[int]int{}
    fmt.Println(len(mp2)) // 0
```

## **2.3 用 new 对 map 类型进行初始化**

用 new 对 map 进行创建，编译器不会报错，但是无法对 map 进行赋值操作的。正确应使用 make 进行内存分配。

```
        mp := new(map[int]int)
        f := func(m map[int]int) {
            m[10] = 10
        }
        f(*mp) // assignment to entry in nil map
```

## **2.4 空指针和空接口不等价**

对于接口类型是可以用 nil 赋值的，但如果对于接口指针类型，其值对应的并不一个空接口。Go 语言编译器似乎在这个处理，会特殊处理。

```
// MyErr just for demotype MyErr struct{}

func (e *MyErr) Error() string {
    return""
}

func TestInterfacePointBug(t *testing.T) {
    var e *MyErr = nil
    var e2 error = e // e2 will never be nil.
    fmt.Println(e2 == nil)
}
```

# **03 函数，方法与控制流相关**

## **3.1 循环中使用闭包错误引用同一个变量**

原因分析：闭包捕获外部变量，它不关心这些捕获的变量或常量是否超出作用域，只要闭包在使用，这些变量就会一直存在。

```
  type S struct {
        A string
        B string
        C string
    }
    typ := reflect.TypeOf(S{})
    funcArr := make([]func() string, typ.NumField())
    for i := 0; i < typ.NumField(); i++ {
        f := func() string {
            return typ.Field(i).Name
        }
        funcArr[i] = f
    }

    fmt.Println(funcArr[0]()) // error reflect: Field index out of bounds
```

所以上面的示例代码，在循环中闭包函数只记录了 i 变量的使用，当循环结束后，i 值变成了 3。当调用该匿名函数时，就会引用 i=3 的值 ，出现越界的异常。

正确处理的方式如下，只需要闭包前处理一下把 i 变量赋值给一个新变量。

```
  type S struct {
        A string
        B string
        C string
    }
    typ := reflect.TypeOf(S{})
    funcArr := make([]func() string, typ.NumField())
    for i := 0; i < typ.NumField(); i++ {
        index := i // assign to a new variable
        f := func() string {
            name := typ.Field(index).Name
            return name
        }
        funcArr[i] = f
    }

    fmt.Println(funcArr[0]()) // A
```

## **3.2 元素内容较大时，不要用 range 遍历**

用 range 来操作遍历使用上非常方便，但是它的遍历中是需要进行值赋值操作，遇到元素占用的内存比较大时，性能就会影响较大。

下面是针对两种方式做了一下基准测试。

```
func CreateABigSlice(count int) [][4096]int {
    ret := make([][4096]int, count)
    for i := 0; i < count; i++ {
        ret[i] = [4096]int{}
    }
    return ret
}

func BenchmarkRangeHiPerformance(b *testing.B) {
    v := CreateABigSlice(1 << 12)

    for i := 0; i < b.N; i++ {
        len := len(v)
        var tmp [4096]int
        for k := 0; k < len; k++ {
            tmp = v[k]
        }
        _ = tmp
    }
}

func BenchmarkRangeLowPerformance(b *testing.B) {
    v := CreateABigSlice(1 << 12)

    for i := 0; i < b.N; i++ {

        var tmp [4096]int
        for _, e := range v {
            tmp = e
        }
        _ = tmp
    }
}
```

测试结果如下：range 方式的性能较 for 方式相差了近 10000 倍。

```
cpu: 11th Gen Intel(R) Core(TM) i5-1145G7 @ 2.60GHz
BenchmarkRangeHiPerformance-8            9767457              1255 ns/op
BenchmarkRangeLowPerformance-8               975          11513216 ns/op
PASS
ok      withoutbug/avoidtofix   26.270s
```

## **3.3 循环内调用 defer 造成销毁处理延迟**

在很多场景，在循环内申请资源在循环完成后释放，但是使用 defer 语句处理，是需要在当前函数退出时才会执行，在循环中是不会触发的，导致资源延迟释放。

```
func main() {
    for i := 0; i < 5; i++ {
        f, err := os.Open("./mygo.go")
        if err != nil {
            log.Fatal(err)
        }
        defer f.Close()
    }
}
```

比较好的解决办法就是在 for 循环里不要使用 defer，直接进行销毁处理。

```
func main() {
    for i := 0; i < 5; i++ {
        f, err := os.Open("/path/to/file")
        if err != nil {
            log.Fatal(err)
        }
        f.Close()
    }
}
```

## **3.4 Goroutine 无法阻止主进程退出**

后台 Goroutine 无法保证在方法退出来执行完成。

```
func main() {
     gofunc() {
        time.Sleep(time.Second)
        fmt.Println("run")
    }()   
}
```

## **3.5 Goroutine 抛 panic 会导致进程退出**

后台 Goroutine 执行中，如果抛 panic 并不进行 recover 处理，会导致主进程退出。

下面的代码示例：

```
func main1() {
    go func() {
        panic("oh...")
    }()
    for i := 0; i < 3; i++ {
        fmt.Println(i)
        time.Sleep(time.Second)
    }
    fmt.Println("bye bye!")
}
```

修正代码如下：

```
func main2() {
    go func() {
        defer func() {
            recover() // should do some thing here
        }()

        panic("oh...")
    }()

    for i := 0; i < 3; i++ {
        fmt.Println(i)
        time.Sleep(time.Second)
    }

    fmt.Println("bye bye!")
}
```

**3.6 r\**\**ecover 函数 只在 defer 函数内生效**

需要注意：在非 defer 函数内，调用 recover 函数，是不会有任何的执行，也无法来处理 panic 错误。

下面的示例代码，是无法处理 panic 的错误：

```
func NoTestDeferBug(t *testing.T) {
    recover()
    panic(1) // could not catch
}

func NoTestDeferBug2(t *testing.T) {
    defer recover()
    panic(1) // could not catch
}
```

正确的代码如下：

```
func TestDeferFixed(t *testing.T) {
    defer func() {
        recover()
    }()
    panic("this is panic info") // could not catch
}
```

# **04 并发与内存同步相关**

## **4.1 跨 Goroutine 之间不支持顺序一致性内存模型**

在 Go 语言的内存模型设计中， 内存写入顺序性只能保障在单一 Goroutine 内一致，跨 Goroutine 之间无法保障监测变量操作顺序的一致性。

下面是官方的例子：

```
package main

var msg string
var done bool
func setup() {
    msg = "hello, world"
    done = true
}

func main() {
    go setup()
    for !done {
    }
    println(msg)
}
```

上面代码的问题是，不能保证在 main 中对 done 的写入的监测时， 会对变量 a 的写入也进行监测，因此该程序也可能会打印出一个空字符串。更糟的是，由于在两个线程之间没有同步事件，因此无法保证对 done 的写入总能被 main 监测到。main 中的循环不保证一定能结束。

解决办法就是使用显示同步方案，使用通道进行同步通信。

```
package main

var msg string 
var done = make(chan bool)

func setup() {
    msg = "hello, world"
    done <- true
}

func main() {
    go setup()
    <-done
    println(msg)
}
```

这样就可以保证代码执行过程中必定输出 hello，world。

更多内存同步阅读材料：[https://go-zh.org/ref/mem](https://www.oschina.net/action/GoToLink?url=https%3A%2F%2Fgo-zh.org%2Fref%2Fmem)

# **05 序列化相关**

## **5.1 基于指针参数方式传递的反序列功能，都不会初始化要反序列化的对象字段**

该问题经常发生的原因是基于指针参数方式传递的反序列函数其实做的只是值覆盖的功能，并不会把要反序化的对象的所有值进行初始化操作，这样就会导致未覆盖的值的保留。像 json.Unmarshal, xml.Unmarshal 函数等。

下面是基于 json 对 map 类型的变量进行 json.Unmarshal 的问题示例：

```
package main

import (
    "encoding/json"
    "fmt"
)

func main() {
    val := map[string]int{}
    s1 := `{"k1":1, "k2":2, "k3":3}`
    s2 := `{"k1":11, "k2":22, "k4":44}`
    json.Unmarshal([]byte(s1), &val)
    fmt.Println(s1, val)
    json.Unmarshal([]byte(s2), &val)
    fmt.Println(s2, val)
}
```

输出：

```
{"k1":1, "k2":2, "k3":3} map[k1:1 k2:2 k3:3]
{"k1":11, "k2":22, "k4":44} map[k1:11 k2:22 k3:3 k4:44]
```

由于 json.UnMarshal 方法只会新增和覆盖 map 中的 key，不会删除 key。虽然第二个 json 字符串中没有 k3 的内容，但输出结果中依然保留在了 k3 的内容。

要解决这个问题，每次 unmarshal 之前都重新声明变量即可。

# **06 其它杂项**

## **6.1 数字类型转换越界陷阱**

Go 语言中，任何操作符不会改变变量类型，下面示例引入一个坑，出现位移越界。

```
func TestOverFlowBug(t *testing.T) {
    var num int16 = 5000
    var result int64 = int64(num << 9)
    fmt.Println(result) // 4096 overflow
}
```

修正方式如下，需要操作前对类型转换：

```
func TestOverFlowFixed(t *testing.T) {
    var num int16 = 5000
    var result int64 = int64(num) << 9
    fmt.Println(result) // 2560000
}
```

## **6.2 map 遍历是顺序不固定**

map 的实现是通 hash 表进行分桶定位，同时 map 的遍历引入了随机实现，所以每次遍历的顺序都可能变化。

```
    mp := map[int]int{}
    for i := 0; i < 20; i++ {
        mp[i] = i
    }

    for k, v := range mp {
        fmt.Println(k, v)
    }
```

**——END——**