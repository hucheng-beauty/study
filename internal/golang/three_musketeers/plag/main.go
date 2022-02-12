package main

import (
	flag "github.com/spf13/pflag"
)

var (
	name string
	port int
)

func main() {

	// 支持长选项、默认值和使用文本，并将标志的值存储在指针中
	flag.String("name", "hello", "this is string about name!")

	// 支持长选项、短选项、默认值和使用文本，并将标志的值存储在指针中。
	flag.StringP("version", "v", "v0.3.9", "this is app version!")

	// 支持长选项、默认值和使用文本，并将标志的值绑定到变量。
	flag.StringVar(&name, "nam", "hello", "this si string about name!")

	// 支持长选项、短选项、默认值和使用文本，并将标志的值绑定到变量。
	flag.StringVarP(&name, "nam", "n", "hello", "this si string about name!")

	/*
		// 使用Get获取参数的值。
		flagSet := flag.NewFlagSet("test", flag.ContinueOnError)
		flagSet.Int("age", 19, "age")

		age1, err := flagSet.GetInt("age")
		if err != nil {
			return
		}
		fmt.Println(age1)
	*/

	/*
		// 获取非选项参数
		// go run ./internal/plag/main.go name age
		fmt.Printf("argument number is: %v\n", flag.NArg())
		fmt.Printf("argument list is: %v\n", flag.Args())
		fmt.Printf("the first argument is: %v\n", flag.Arg(0))
	*/

	/*
		// 指定了选项但是没有指定选项值时的默认值。
		// go run ./internal/plag/main.go --age  ==> 4321
		// go run ./internal/plag/main.go --age=1312  ==> 1312
		// go run ./internal/plag/main.go  ==> 1234
		var ip = flag.IntP("age", "a", 1234, "help message")
		flag.Lookup("age").NoOptDefVal = "4321"
		fmt.Println("++", *ip)
	*/

	/*
		// 弃用标志或者标志的简写
		// 弃用的标志或标志简写在帮助文本中会被隐藏，并在使用不推荐的标志或简写时打印正确的用法提示。
		//例如，弃用名为 logmode 的标志，并告知用户应该使用哪个标志代替：
		fmt.Println(*flag.String("logmode", "logmode string", "the usage of log mode!"))
		err := flag.CommandLine.MarkDeprecated("logmode",
			"please use --log-mode instead")
		if err != nil {
			fmt.Println(err)
			return
		}
	*/

	/*
		// 保留名为 port 的标志，但是弃用它的简写形式
		// 这样隐藏了帮助文本中的简写 P，并且当使用简写 P 时，打印了
		// Flag shorthand -P has been deprecated, please use --port only。
		// usage message 在此处必不可少，并且不应为空。
		flag.IntVarP(&port, "port", "P", 3306, "MySQL service host port!")
		flag.CommandLine.MarkShorthandDeprecated("port", "please use --port only")
	*/

	/*
		// 隐藏标志
		// 可以将 Flag 标记为隐藏的，这意味着它仍将正常运行，但不会显示在 usage/help 文本中。
		//例如：隐藏名为 secretFlag 的标志，只在内部使用，并且不希望它显示在帮助文本或者使用文本中
		flag.CommandLine.MarkHidden("secretFlag")
	*/
	//flag.Parse()
}
