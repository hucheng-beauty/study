package fourth_0914

/*
    Hercy 想要为购买第一辆车存钱。他每天都往银行里存钱。
    最开始，他在周一的时候存入 1 块钱。
    从周二到周日，他每天都比前一天多存入 1 块钱。
    在接下来每一个周一，他都会比前一个周一多存入 1 块钱。
    给你 n ，请你返回在第 n 天结束的时候他在银行总共存了多少块钱。

   例子1： 输入：4 输出：10

   例子2： 输入：10 输出：37

   例子3： 输入：20 输出：96
*/

// 1 2 3 4 5 6 7
// 2 3 4 5 6 7 8

func CountMoney(n int) int {
    weeks := n / 7
    days := n % 7
    total := 0

    for i := 0; i < weeks; i++ {
        total += (i * 7) + 28
    }

    for i := 0; i < days; i++ {
        total += weeks + 1 + i
    }

    return total
}
