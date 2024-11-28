package main

import (
	"fmt"
	"math"
	"math/big"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) != 3 && len(os.Args) != 4 {
		fmt.Println("Usage: binomial_distribution.exe <总次数> <事件发生数> [事件发生概率]")
		fmt.Println("<总次数> 期望为大于0的整数，大于等于<事件发生数>")
		fmt.Println("<事件发生数>期望为大于0的整数，小于等于总次数")
		fmt.Println("[事件发生概率] 默认为0.015，选填，范围为:（0,1）")
		fmt.Println("输出为：符合正态分布的随机事件发生<总次数>次，其中概率为[事件发生概率]的特定事件出现<事件发生数>次的概率，在正态分布中的位置")
		fmt.Println("人话:出货率为[事件发生概率]，抽<总次数>出货<事件发生数>的你，战胜了百分之几的人")
		return
	}

	precision := 500
	N, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	k, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	p_in := big.NewFloat(0).SetPrec(uint(precision))
	if len(os.Args) == 4 {
		p_in64,err := strconv.ParseFloat(os.Args[3],64)
		if err != nil || p_in64 >= 1 || p_in64 <=0 {
			fmt.Printf("请输入大于0小于1的概率\n", err)
			return
		}
		p_in = big.NewFloat(p_in64).SetPrec(uint(precision))
	}
	p := new(big.Float).SetPrec(uint(precision))
	if p_in.Sign() == 0 {
		p,_ = p.SetString("0.015") // 出现事件的概率
	} else {
		p,_ = p.SetString(os.Args[3])
	}
	M := new(big.Float).SetPrec(uint(precision)).Quo(big.NewFloat(float64(k)), big.NewFloat(float64(N))) // n/N的阈值
	M_float, _ := M.Float64()
	M_int := int(math.Ceil(M_float * float64(N))) // 将M转换为整数
	probability := cumulativeBinomialDistribution(N, p, M_int, precision) // 计算累积二项分布
	probabilityPercent := new(big.Float).SetPrec(uint(precision)).Mul(big.NewFloat(100),probability)
	fmt.Printf("出现事件次数/总数低于%s的概率为: %s%%\n", M.Text('f', 10), probabilityPercent.Text('f', 100))
}

// 计算二项分布的累积概率
func cumulativeBinomialDistribution(N int, p *big.Float, k int, precision int) *big.Float {
	sum := big.NewFloat(0).SetPrec(uint(precision))
	for i := 0; i <= k; i++ {
		binomialCoefficient := binomialCoefficient(N, i, precision)
		pPower := pow(p, int64(i), precision)
		q := new(big.Float).SetPrec(uint(precision)).Sub(big.NewFloat(1), p)
		qPower := pow(q, int64(N-i), precision)
		add := new(big.Float).SetPrec(uint(precision)).Mul(binomialCoefficient,pPower)
		add.Mul(add,qPower)
		probabilityPercent := new(big.Float).SetPrec(uint(precision)).Mul(add,big.NewFloat(100))
		fmt.Printf("%d次独立事件中出现%d次指定事件的概率为：%s%%\n", N, i, probabilityPercent.Text('f', 100))
		sum.Add(sum, add)
	}
	return sum
}

// 计算组合数
func binomialCoefficient(N int, k int, precision int) *big.Float {
	result := big.NewFloat(1).SetPrec(uint(precision))
	for i := 0; i < k; i++ {
		nMinusI := big.NewFloat(float64(N - i))
		iPlusOne := big.NewFloat(float64(i + 1))
		result.Mul(result, new(big.Float).SetPrec(uint(precision)).Quo(nMinusI, iPlusOne))
	}
	return result
}

// 计算幂
func pow(base *big.Float, exponent int64, precision int) *big.Float {
	result := big.NewFloat(1).SetPrec(uint(precision))
	for i := int64(0); i < exponent; i++ {
		result.Mul(result, base)
	}
	return result
}
