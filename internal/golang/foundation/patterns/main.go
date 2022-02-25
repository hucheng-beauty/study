package patterns

// worker
/*
func worker(i int) {
	fmt.Printf("worker %d\n", i)
}
*/

// for select 无限循环模式
/*func main() {
	ch := make(chan int, 1)
	for i := 0; i < 10; i++ {
		go worker(i)
		if i == 9 {
			ch <- i
		}

		//go func(i int) {
		//	if i == 9 {
		//		ch <- i
		//	}
		//}(i)
	}

	time.Sleep(time.Millisecond)
	for {
		select {
		case <-ch:
			fmt.Println("over")
			return
		}
	}
}*/

// for range select 有限循环模式
/*
func main() {
	ch := make(chan struct{})
	for _, i := range []int{1, 2, 3} {
		fmt.Println(i)
		go func(i int) {
			if i == 3 {
				ch <- struct{}{}
			}
		}(i)
	}

	time.Sleep(time.Millisecond)
	for {
		select {
		case <-ch:
			fmt.Println("over")
			return
		}
	}
}
*/

// select timeout 模式
/*
func main() {
	result := make(chan string)
	timeout := time.After(3 * time.Second) //
	go func() {
		//模拟网络访问
		time.Sleep(5 * time.Second)
		result <- "服务端结果"
	}()
	for {
		select {
		case v := <-result:
			fmt.Println(v)
		case <-timeout:
			fmt.Println("网络访问超时了")
			return
		default:
			fmt.Println("等待...")
			time.Sleep(1 * time.Second)
		}
	}
}*/

// Context 的 WithTimeout 函数超时取消
/*
func main() {
	// 创建一个子节点的context,3秒后自动超时
	ctx, stop := context.WithTimeout(context.Background(), 3*time.Second)

	go func() {
		worker(ctx, "worker bee first")
	}()
	go func() {
		worker(ctx, "worker bee second")
	}()

	fn(stop)
}

func fn(cancelFunc context.CancelFunc) {
	time.Sleep(5 * time.Second) // 工作5秒后休息
	cancelFunc()                // 5秒后发出停止指令
	fmt.Println("结束了")
}

func worker(ctx context.Context, name string) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("下班咯~~~")
			return
		default:
			fmt.Println(name, "认真摸鱼中，请勿打扰...")
		}
		time.Sleep(1 * time.Second)
	}
}
*/

// Pipeline 模式
/*
func main() {
	coms := buy(10)       // 采购10套零件
	phones := build(coms) // 组装10部手机
	packs := pack(phones) // 打包它们以便售卖

	//输出测试，看看效果
	for p := range packs {
		fmt.Println(p)
	}
}

//工序1采购
func buy(n int) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		for i := 1; i <= n; i++ {
			out <- fmt.Sprint("零件", i)
		}
	}()
	return out
}

//工序2组装
func build(in <-chan string) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		for c := range in {
			out <- "组装(" + c + ")"
		}
	}()
	return out
}

//工序3打包
func pack(in <-chan string) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		for c := range in {
			out <- "打包(" + c + ")"
		}
	}()
	return out
}
*/

// 扇入扇出模式
/*
func main() {
	coms := buy(10) //采购10套配件

	//三班人同时组装100部手机
	phones1 := build(coms)
	phones2 := build(coms)
	phones3 := build(coms)

	//汇聚三个channel成一个
	phones := fanIn(phones1, phones2, phones3)

	packs := pack(phones) //打包它们以便售卖

	//输出测试，看看效果
	for p := range packs {
		fmt.Println(p)
	}
}

//工序1采购
func buy(n int) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		for i := 1; i <= n; i++ {
			out <- fmt.Sprint("零件", i)
		}
	}()
	return out
}

//工序2组装
func build(in <-chan string) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		for c := range in {
			out <- "组装(" + c + ")"
		}
	}()
	return out
}

//工序3打包
func pack(in <-chan string) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		for c := range in {
			out <- "打包(" + c + ")"
		}
	}()
	return out
}

//扇入函数（组件）,把多个 channel 中的数据发送到一个 channel 中
func fanIn(ins ...<-chan string) <-chan string {
	var wg sync.WaitGroup
	out := make(chan string)

	//把一个channel中的数据发送到out中
	p := func(in <-chan string) {
		defer wg.Done()
		for c := range in {
			out <- c
		}
	}
	wg.Add(len(ins))

	//扇入，需要启动多个goroutine用于处于多个channel中的数据
	for _, cs := range ins {
		go p(cs)
	}

	//等待所有输入的数据ins处理完，再关闭输出out
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}
*/

// Futures 模式
/*
func main() {
	vegetablesCh := wash() // 洗菜
	waterCh := getWater()  // 烧水

	fmt.Println("已经安排好洗菜和烧水了，我先开一局")
	time.Sleep(2 * time.Second)

	fmt.Println("要做火锅了，看看菜和水好了吗")
	vegetables := <-vegetablesCh
	water := <-waterCh

	fmt.Println("准备好了，可以做火锅了:", vegetables, water)

}

// 洗菜
func wash() <-chan string {
	out := make(chan string)
	go func() {
		time.Sleep(5 * time.Second)
		out <- "洗好的菜"
	}()
	return out
}

//烧水
func getWater() <-chan string {
	out := make(chan string)
	go func() {
		time.Sleep(5 * time.Second)
		out <- "烧开的水"
	}()
	return out
}
*/
