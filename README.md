# kube-trashman

自动清理K8Sjob的程序

tag: v1.0版本   [release-v1.0](https://github.com/fangxx-github/kube-trashman/releases/tag/v1.0)

注释：这个版本回删除哪些job：

ExpiredJobs：

1.首先他会识别job中是否带有kube.trashman.io/expiration表示的注释，如果有记录这个过期时间，他首先将job的创建时间取时间戳和job注释的过期时间作比较，如果大于过期时间返回true，再通过运行时间和传入的过期时间作比较，如果还大于传入的过期时间再返回true。这是函数ExpiredJobs的逻辑；

2.然后在执行TargetJobs函数，这个函数的确认最终哪些job需要删除，其中一个判断，是判断job中是否带有annotation/ignore注释的job，如果失败job中有这个注释则从这个job列表中去除这个job
