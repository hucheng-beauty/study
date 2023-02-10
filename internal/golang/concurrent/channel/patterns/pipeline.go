package patterns

import "fmt"

// PipelineWithBuy 工序1采购
func PipelineWithBuy(n int) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		for i := 1; i <= n; i++ {
			out <- fmt.Sprint("零件", i)
		}
	}()
	return out
}

// PipelineWithBuild 工序2组装
func PipelineWithBuild(in <-chan string) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		for c := range in {
			out <- "组装(" + c + ")"
		}
	}()
	return out
}

// PipelineWithPack 工序3打包
func PipelineWithPack(in <-chan string) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		for c := range in {
			out <- "打包(" + c + ")"
		}
	}()
	return out
}

// Pipeline 模式
func Pipeline() {
	packs := /*打包它们以便售卖*/ PipelineWithPack( /*组装10部手机*/ PipelineWithBuild( /*采购10套零件*/ PipelineWithBuy(10)))

	for p := range packs {
		fmt.Println(p)
	}
}
