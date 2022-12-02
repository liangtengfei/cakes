> **Rust 的集合可以分为四个主要类别**：
>
> - 序列: Vec、VecDeque、LinkedList
> - Maps: HashMap, BTreeMap
> - 集合: HashSet、BTreeSet
> - 混杂: BinaryHeap

## 什么时候应该使用哪个集合？

#### 在以下情况下，请使用 Vec：
- 您想要收集项以供以后处理或发送到其他地方，而不必关心所存储的实际值的任何属性。
- 您需要一个按特定顺序排列的元素序列，并且只会追加到末尾 (或接近末尾)。
- 您想要一个栈。
- 您需要一个可调整大小的数组。
- 您需要一个堆分配的数组。
#### 在以下情况下，请使用 VecDeque：
- 您需要 Vec 支持序列两端的有效插入。
- 您想要一个队列。
- 您需要一个双端队列 (deque)。
- #### 在以下情况下，请使用 LinkedList：
- 您需要一个未知大小的 Vec 或 VecDeque，并且不能容忍摊销。
- 您想有效地拆分和追加列表。
- 您确实，想要一个双向链表。


#### 在以下情况下，请使用 HashMap：
- 您想要将任意键与任意值相关联。
- 您需要一个缓存。
- 您需要一个没有额外功能的 map。

#### 在以下情况下，请使用 BTreeMap：
- 您需要一个按其键排序的 map。
- 您希望能够按需获取一系列条目。
- 您对最小或最大键值对是什么感兴趣。
- 您想要找到小于或大于某项的最大或最小键。

#### 在以下情况下，请使用任何这些 Map 的 Set 变体：
- 您只想记住您所看到的键。
- 与您的键关联没有任何有意义的值。
- 您只想要一个 set。


#### 在以下情况下，请使用 BinaryHeap：
- 您想存储一堆元素，但只想在任何给定时间处理 “biggest” 或 “most important”。
- 您需要一个优先级队列。

## 性能比较
### Sequences

|            | get(i)         | insert(i)       | remove(i)      | append | split_off(i)   |
| ---------- | -------------- | --------------- | -------------- | ------ | -------------- |
| Vec        | O(1)           | O(n-i)*         | O(n-i)         | O(m)*  | O(n-i)         |
| VecDeque   | O(1)           | O(min(i, n-i))* | O(min(i, n-i)) | O(m)*  | O(min(i, n-i)) |
| LinkedList | O(min(i, n-i)) | O(min(i, n-i))  | O(min(i, n-i)) | O(1)   | O(min(i, n-i)) |

### Maps

|          | get       | insert    | range     | remove    | append |
| -------- | --------- | --------- | --------- | --------- | ------ |
| HashMap  | O(1)~     | O(1)~*    | O(1)~     | N/A       | N/A    |
| BTreeMap | O(log(n)) | O(log(n)) | O(log(n)) | O(log(n)) | O(n+m) |