package main

import (
    "fmt"
)

/*
   抖音活跃用户行为分析
   表A：用户浏览视频日志 user_behavior: date, user_id,
       video_id, start_time, end_time
   表B：视频信息 video_info: video_id, video_duration
   表C：用户信息 user_info: user_id, gender
   问题：
   （1）2020年12月11日，观看不同视频个数的前 5 名 user_id ；
   （2）2020年12月11日，观看超过 50 个不同视频的女性用户中，
       完整观看率最高的 10 个 user_id。
*/

/*

   select user_id, count(distinct video_id) as video_count
   from user_behavior
   where date = "2020-12-11"
   group by video_id
   order by video_count desc
   limit 5;


    // having: 对分组后的结果进行过滤,类似于 where
    // timestampdiff: 返回两个时间点之间的时间差,单位为秒。
    // case when: 判断一个条件,并返回一个结果:
    CASE
       WHEN condition1 THEN result1
       WHEN condition2 THEN result2
       ...
       ELSE default_result
    END
    // 每个 WHEN 子句都会判断一个条件。
    // 如果某个条件成立，则返回对应的 THEN 后的结果。
    // 所有条件都不满足时，返回 ELSE 后的结果（可选）。
    // 最后必须以 END 结束。

   select
       ub.user_id,
       count(distinct ub.video_id) as total_videos,
       count(distinct case
           when timestampdiff(second, ub.start_time, ub.end_time) >= vi.video_duration
           then ub.video_id
           end) as full_watch_count,
       count(distinct case
           when timestampdiff(second, ub.start_time, ub.end_time) >= vi.video_duration
           then ub.video_id
           end) / count(distinct ub.video_id) as full_watch_rate
   from user_behavior ub
   join user_info ui on ub.user_id = ui.user_id
   join video_info vi on ub.video_id = vi.video_id
   where ub.date = '2020-12-11' and ui.gender = 'female'
   group by ub.user_id
   having count(distinct ub.video_id) > 50
   order by full_watch_rate desc
   limit 10;

*/

func main() {
    fmt.Println("hello world!")
}
