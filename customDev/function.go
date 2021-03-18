package customDev

import (
	"ehang.io/nps/lib/file"
	"math/rand"
	"time"
)

// Random one string with numbers and strings
func RandStr(length int) string {
	letterRunes := []rune("0123456789abcdefghijklmnopqrstuvwxyz")
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, length)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

// 随机选择tunnel
func RandChooseByNums(Data []*file.Tunnel, n int) []*file.Tunnel {
	RandomNums := uniqueRandomNum(0, len(Data)-1, n)

	var chosen []*file.Tunnel
	for _, j := range RandomNums {
		chosen = append(chosen, Data[j])
	}

	return chosen
}

//生成count个[start,end)结束的不重复的随机数切片
func uniqueRandomNum(start int, end int, count int) []int {
	//范围检查
	if end == start && end < count {
		return []int{start}
	} else if end < start || (end-start) < count {
		return nil
	}

	//存放结果的slice
	nums := make([]int, 0)
	//随机数生成器，加入时间戳保证每次生成的随机数不一样
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(nums) < count {
		//生成随机数
		num := r.Intn(end-start) + start

		//查重
		exist := false
		for _, v := range nums {
			if v == num {
				exist = true
				break
			}
		}

		if !exist {
			nums = append(nums, num)
		}
	}

	return nums
}
