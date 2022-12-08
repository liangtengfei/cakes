## 结论：在 Golang 中所有的类型传递都是通过值传递实现的

#### 举例说明：

### 普通参数

```go
func main() {
  x, y := int64(1), int64(2)
  Swap(x, y)
  fmt.Println(x, y) //1 2
  p := Person{1, "boy"}
  attachPersonName(&p)
  fmt.Println(p) //{1 boy}
}

func Swap(x, y int64) {
	x, y = y, x
}

type Person struct {
    Id int
    Name string
}

func attachPersonName(p Person) {
	p.Name = strings.ToUpper(p.Name)
}
```

### Slice

```
func main(){
	data := []int64{1, 2, 3}
	SortSlice(data)
	fmt.Println(data) //[3 2 1]
}

func SortSlice(data []int64) {
	sort.Slice(data, func(i, j int) bool {
		return data[i] > data[j]
	})
}
```

> **因为slice值包含指向第一个slice元素的指针，因此向函数传递slice将允许在函数内部修改底层数组的元素。****换句话说，复制一个slice只是对底层的数组创建了一个新的slice别名**

> **另外：map,chan和slice是类似的**

