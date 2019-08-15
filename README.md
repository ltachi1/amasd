# scrapyd-admin
方便、快捷、动态部署scrapyd的管理工具。
#### 特点
1. 0配置
2. 无任何依赖
3. 可在任意系统运行(如果需要可联系作者进行打包)
4. 可对scrapyd服务状态、任务运行结束、任务运行异常进行监控，可自由配置
5. 通知方式有邮件通知、钉钉群机器人通知、企业微信群机器人通知，可自由配置

#### 运行参数
- -p 监听端口，默认 8000
- -e 运行环境，可使用dev、testing、production中任意一种
- -log 日志存放路径，默认当前运行的目录
- -db 数据库存放路径，默认当前运行的目录

#### 作者联系方式
- QQ: 376202990
- 微信: hxz_lhq

#### 部分功能截图
- 默认用户名、密码均为admin

- 登录
![image](https://github.com/ltachi1/scrapyd-admin/raw/master/images/login.jpg)

- 添加服务器
![image](https://github.com/ltachi1/scrapyd-admin/raw/master/images/server_add.jpg)

- 添加项目
![image](https://github.com/ltachi1/scrapyd-admin/raw/master/images/project_add.jpg)

- 项目列表
![image](https://github.com/ltachi1/scrapyd-admin/raw/master/images/project_list.png)

- 更新项目文件
![image](https://github.com/ltachi1/scrapyd-admin/raw/master/images/update_version.png)

- 更新关联服务器
![image](https://github.com/ltachi1/scrapyd-admin/raw/master/images/update_servers.png)

- 添加任务
![image](https://github.com/ltachi1/scrapyd-admin/raw/master/images/add_task.jpg)

- 添加计划任务
![image](https://github.com/ltachi1/scrapyd-admin/raw/master/images/add_schedule.png)

- 计划任务列表
![image](https://github.com/ltachi1/scrapyd-admin/raw/master/images/schedule_list.jpg)

- 正在运行的任务列表
![image](https://github.com/ltachi1/scrapyd-admin/raw/master/images/task_running_list.jpg)

- 已完成任务列表
![image](https://github.com/ltachi1/scrapyd-admin/raw/master/images/task_finished_list.png)

- 通知设置
![image](https://github.com/ltachi1/scrapyd-admin/raw/master/images/notice.png)

- 任务完成通知设置
![image](https://github.com/ltachi1/scrapyd-admin/raw/master/images/notice_task_finished.png)

#### Scrapy 课程推荐(本人已看完此教程所有视频)
- 课程内容：Python基础+脚本爬虫+Scrapy框架+实战训练
- 涵盖功能：数据提取、数据入库、模拟登录、反爬和反反爬、代理IP管理
- 课程优势：一对一专业答疑，远程调配环境，爬虫思路指导
- 课程链接：https://study.163.com/course/courseMain.htm?courseId=1003729016&share=2&shareId=3749780
- ![image](https://github.com/ltachi1/scrapyd-admin/raw/master/images/scrapy.jpg)