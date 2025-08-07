package main

import (
    "context"
    "time"
)

/*
   products（商品表）：
        product_id（主键）、product_name（商品名称）、
        category_id（分类 ID）、price（价格）、
        stock（库存）、create_time（创建时间）
   orders（订单表）：
       order_id（主键）、user_id（用户 ID）、
       product_id（关联商品 ID ）、order_num（订单数量）、
       order_time（下单时间）、
       status（订单状态：0 - 未支付 1 - 已支付 2 - 已完成 3 - 已取消 ）
   users（用户表）：
       user_id（主键）、user_name（用户名）、
       phone（手机号）、address（收货地址）

    销售总金额（金额 = 数量 × 商品价格 ）
    问题：营销活动需要统计近 7 天，
    每个商品分类（category_id 关联）下，“已完成” 订单的
    商品销售总数量和销售总金额,并按销售总金额降序排序，
    取 Top5 的分类。

    select
        p.category_id,
        sum(o.order_num) as total_quantity,
        sum(o.order_num * p.price) as total_amount
    from
        orders o
    join
        products p on o.product_id = p.product_id
    where
        o.status = 2  -- 已完成的订单
        and o.order_time >= now() - interval 7 day  -- 近 7 天的数据
    group by
        p.category_id
    order by
        total_amount desc
    limit 5;

*/

/*
   如何用 Go 快速实现一个 GET/POST 接口（要求写关键代码片段），
   需体现 goroutine 或 channel 在接口优化中的应用（如异步处理逻辑）。
*/

type handler interface {
    Handler(param string, url string) (interface{}, error)
}

var Mapper = map[string]handler{
    "Get":  nil,
    "Post": nil,
}

var Handler = make([]handler, 0)

var Q = make(chan handler, 100)

func transfer(ctx context.Context) {
    for {
        select {
        case h, ok := <-Q:
            if ok {
                go h.Handler("", "")
            }
        case <-ctx.Done():
            return
        default:
        }
    }
}

func Submit(handler handler) {
    Q <- handler
}

func GetHandler(param string, url string) (interface{}, error) {
    h, ok := Mapper[param]
    if ok {
        return h.Handler(param, url)
    } else {
        return Mapper["Get"].Handler(param, url)
    }
}

type Get struct{}

func (g Get) Handler(param string, url string) (interface{}, error) {
    // TODO implement me
    panic("implement me")
}

func main() {
    ctx, cancelFunc := context.WithCancel(context.Background())

    Submit(nil)

    go transfer(ctx)

    time.Sleep(time.Second)
    cancelFunc()

}
