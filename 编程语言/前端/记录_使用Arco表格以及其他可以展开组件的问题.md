### 问题

在使用[arco design](https://arco.design/)的表格时，设置了默认展开所有节点的属性：**default-expand-all-rows = true**。但是在获取数据后，却没有效果。使用其他组件tree以及tree select都有同样的问题。

简答测试以后发现了问题，在acro design pro中如果使用mock数据，设置了default-expand-all-rows=true同样是无效的，但是不使用mock接口数据，而是使用直接定义的const数据就能生效。

### 问题猜想

应该是渲染的问题，第一次渲染使用了空数据，获取数据以后，由于已经组件已经渲染过，重新加载时，设置就不会生效了，默认展开所有无效了。

### 解决方案

#### 1. 数据没有加载之前不渲染

在组件上使用v-if。

```
<a-table
    v-if="renderData.length"
    row-key="id"
    :loading="loading"
    :pagination="pagination"
    :columns="(cloneColumns as TableColumnData[])"
    :data="renderData"
    :default-expand-all-rows="true"
    :bordered="false"
    :size="size"
    @page-change="onPageChange"
```

#### 2. 重新渲染

给组件上设置一个key，获取数据之后改变这个值，借助key改变自动变成新的组件。

#### 3. 使用expaned-其他属性

> expanded-keys（受控模式）、default-expanded-keys（非受控模式） 受控模式需要手动控制，如果设置了值，需要监听expand事件。不然组件点击展开或者收起就失效。