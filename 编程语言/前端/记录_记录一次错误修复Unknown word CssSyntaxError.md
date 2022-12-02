在使用[Arco Pro](https://arco.design/vue/docs/pro/start)遇到的问题：Unknown word CssSyntaxError。

在编写系统代码的时候没发现，提交代码时，自动触发了代码规范检查的钩子，造成代码一直无法提交，把解决办法记录一下。

## Unknown word CssSyntaxError



> 可能还会出现这样的描述信息：
>
> When linting something other than CSS, you should install an appropriate syntax, e.g. "postcss-html", and use the "customSyntax" option
>

然后就找到了*stylelint*的官网https://stylelint.io/migration-guide/to-14/，这个页面列出了使用各种语言的语法配置。

为每种语言提供了其他共享配置：

- [stylelint-config-recommended-vue](https://www.npmjs.com/package/stylelint-config-recommended-vue) ... Shared config for Vue.
- [stylelint-config-html](https://www.npmjs.com/package/stylelint-config-html) ... Shared config that enables parsing for HTML, XML, Vue, Svelte, and PHP.

选择第一个针对vue的配置，进入到https://www.npmjs.com/package/stylelint-config-recommended-vue插件页面。

**但是先不要安装**,往下看：

> 但是，某些文体规则可能不适用于 Vue。我们建议您改为安装 stylelint-config-standard-vue。

**所以直接安装**stylelint-config-standard-vue即可。

``` shell
npm install --save-dev postcss-html stylelint-config-standard-vue
```

修改配置文件：stylelint，可能是js、json后缀的文件。

```
"extends": "stylelint-config-standard-vue"
"customSyntax": "postcss-html",
```

到这里就解决了上面的问题。

另外，对于其他语言和嵌入式样式，我们建议使用以下 PostCSS 语法：

- Less language (`.less`) use [postcss-less](https://www.npmjs.com/package/postcss-less)
- Sass language (`.sass`) use [postcss-sass](https://www.npmjs.com/package/postcss-sass)
- CSS-in-JS embeds (`.js`, `.jsx`, `.ts` etc.) use [@stylelint/postcss-css-in-js](https://www.npmjs.com/package/@stylelint/postcss-css-in-js)
- HTML, XML and HTML-like embeds (`.html`, `.xml`, `.svelte`, `.vue` etc.) use [postcss-html](https://www.npmjs.com/package/postcss-html)
- Markdown embeds (`.md`, `.markdown` etc.) use [postcss-markdown](https://www.npmjs.com/package/postcss-markdown)

##### PS：提交的时候又遇到的问题😭。

- subject may not be empty [subject-empty] 

- type may not be empty [type-empty]

也是提交的时候设置的钩子，出现这个问题的原因是，提交的commit message不规范。😵‍💫

参考 [commitlint](https://commitlint.js.org/#/)即可解决。



*PSPS：出现问题最好先去官网找解决办法，不然胡乱google一通，浪费时间不说，可能还解决不了问题。*
