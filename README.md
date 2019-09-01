# Amasd
Amasd是一款使用golang编写的基于scrapyd的scrapy爬虫部署工具，并且具有后台管理的功能。 Amasd抛弃了传统工具使用配置文件配置的方式，在运行时可无需任务配置即可运行，且无需任何依赖。Amasd所有的调度任务均基于golang的协程来完成，占用服务器资源极少。
同时还具有爬虫定时调度、爬虫任务状态监控、服务器性能监控、通知等功能。
- 运行工具无需任何依赖
- 可运行在Linux、Mac、Windows等系统上
- 可一次性执行爬虫任务，也可添加定时任务来对爬虫进行定时的启动
- 可视化监控服务器性能指标(需要单独下载agent),目前有cpu负载(windows下不支持)、cpu使用情况、内存使用情况、网络上下行速度。
- 可对目标服务器Scrapyd运行状态、爬虫任务运行结束、爬虫任务运行异常进行监控并通知
- 有邮件通知、钉钉群机器人通知、企业微信群机器人通知

#### 运行参数
- 后台登录默认用户名、密码均为admin
- -p 监听端口，默认 8000
- -e 运行环境，可使用dev(对应trace日志级别)、testing(对应debug日志级别)、production(对应info日志级别)中任意一种，默认使用dev
- -log 日志存放路径，默认Amasd运行目录
- -db 数据库存放路径，默认Amasd运行目录

#### 作者联系方式
- QQ 376202990
- QQ交流群 3059362
- 微信 hxz_lhq

#### 部分功能界面截图
- 服务器性能概览
![image](https://github.com/ltachi1/amasd/raw/dev/images/monitor.jpg)

- 服务器近一小时性能数据
![image](https://github.com/ltachi1/amasd/raw/dev/images/monitor_detail.jpg)

- 通知设置
![image](https://github.com/ltachi1/amasd/raw/dev/images/notice.png)

- 任务完成通知设置
![image](https://github.com/ltachi1/amasd/raw/dev/images/notice_task_finished.png)

#### Scrapy 课程推荐(本人已看完此教程所有视频,并且还在不断更新中)
- 课程内容：Python基础+脚本爬虫+Scrapy框架+实战训练
- 涵盖功能：数据提取、数据入库、模拟登录、反爬和反反爬、代理IP管理
- 课程优势：一对一专业答疑，远程调配环境，爬虫思路指导
- 课程链接：https://study.163.com/course/courseMain.htm?courseId=1003729016&share=2&shareId=3749780
- ![image](https://github.com/ltachi1/amasd/raw/dev/images/scrapy.jpg)