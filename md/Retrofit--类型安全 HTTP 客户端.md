![类型安全的HTTP客户端](https://p1-jj.byteimg.com/tos-cn-i-t2oaga2asx/gold-user-assets/2019/12/23/16f30577a5054c9a~tplv-t2oaga2asx-zoom-crop-mark:3024:3024:3024:1702.awebp)

适用于Java和Android的类型安全HTTP客户端，使用举例。

## 背景

Retrofit也是square开发的，和okHttp同属，而且乍看之下两者或许还有冲突的部分。Retrofit是一个RESTful的HTTP网络请求框架，可以使用OkHttp（也可以是其他的）作为HTTP client实现。

## 介绍

#### 功能

- 遵循RESTful API风格设计
- 通过注解配置网络请求参数
- 支持同步&异步网络请求
- 支持多种数据的解析&序列化格式
  - [Gson](https://github.com/google/gson): `com.squareup.retrofit2:converter-gson`
  - [Jackson](https://github.com/FasterXML/jackson): `com.squareup.retrofit2:converter-jackson`
  - [Moshi](https://github.com/square/moshi/): `com.squareup.retrofit2:converter-moshi`
  - [Protobuf](https://developers.google.com/protocol-buffers/): `com.squareup.retrofit2:converter-protobuf`
  - [Wire](https://github.com/square/wire): `com.squareup.retrofit2:converter-wire`
  - [Simple XML](http://simple.sourceforge.net/): `com.squareup.retrofit2:converter-simplexml`
  - [JAXB](https://docs.oracle.com/javase/tutorial/jaxb/intro/index.html): `com.squareup.retrofit2:converter-jaxb`
  - Scalars (primitives, boxed, and String): `com.squareup.retrofit2:converter-scalars`
- 提供对RxJava的支持

#### 优点

1. 功能强大：支持同步和异步
2. 支持多种数据解析&序列化格式
3. 支持RxJava
4. 简洁易用
5. 可拓展性好：功能模块高度封装、解耦度高。

#### 不完美

1. 构建Retrofit时必须指定baseURL，但是可以使用@URL动态修改URL
2. Delete不支持body
3. 执行机制严格（比如：定义了的map参数不能为空）

## 实例

- 添加Retrofit库的依赖

```xml
<dependency>
  <groupId>com.squareup.retrofit2</groupId>
  <artifactId>retrofit</artifactId>
  <version>(insert latest version)</version>
</dependency>
```

- 创建接收类

```
class AlbumRes {
        private Long userId;
        private Long id;
        private String title;

        public Long getUserId() {
            return userId;
        }

        public void setUserId(Long userId) {
            this.userId = userId;
        }

        public Long getId() {
            return id;
        }

        public void setId(Long id) {
            this.id = id;
        }

        public String getTitle() {
            return title;
        }

        public void setTitle(String title) {
            this.title = title;
        }

        @Override
        public String toString() {
            return "AlbumRes{" +
                    "userId=" + userId +
                    ", id=" + id +
                    ", title='" + title + '\'' +
                    '}';
        }
    }
```



- 定义网络请求接口

  ```
  public interface PostService {
    @GET("users/{user}/albums")
    Call<List<AlbumRes>> listalbums(@Path("user") Long user);
  }
  ```

- 创建Retrofit实例

  ```
  Retrofit retrofit = new Retrofit.Builder()
      .baseUrl("http://jsonplaceholder.typicode.com/")
      .addConverterFactory(GsonConverterFactory.create())
      .build();
  ```

- 创建网络请求接口实例

  ```
  PostService service = retrofit.create(PostService.clss);
  ```

- 发送网络请求

  异步：

  ```
  Call<List<AlbumRes>> albumsCall = service.listAlbums(1);
  albumsCall.enqueue(new Callback<>() {
              @Override
              public void onResponse(Call<List<GitHubService.AlbumRes>> call, Response<List<GitHubService.AlbumRes>> response) {
                  if (response.isSuccessful() && response.body() != null) {
                      System.out.println("请求结果：" + response.body().size());
                  }
              }
  
              @Override
              public void onFailure(Call<List<GitHubService.AlbumRes>> call, Throwable throwable) {
                  System.out.println("请求失败");
              }
          });
  ```

  同步：

  ```
  Call<List<AlbumRes>> albumsCall = service.listAlbums(1);
  try {
  	Response<List<>> response = albumsCall.execute();
  } cache(IOException){
  	///
  }
  
  ```

- 处理返回

  如果没有``` .addConverterFactory(GsonConverterFactory.create())``` 添加转换，默认情况下，Retrofit 只能将 HTTP 主体反序列化为 OkHttp 的 ResponseBody 类型，并且它只能接受其 RequestBody 类型作为@Body。