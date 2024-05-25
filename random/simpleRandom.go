package random

// 非常简单的使用线性同余法的随机器
const (
	a int64 = 1664525
	c int64 = 1013904223
	m int64 = 1 << 32
)

type SimpleRandom struct {
	Previous int64
}

// XXX 使用uint32类型是因为如果是一个全长度的int64得不到想要的结果
// 使用种子初始化
func GetSimpleRandom(seed uint32) SimpleRandom {
	return SimpleRandom{Previous: int64(seed)}
}

// 强制设定前一次的结果，将便下一次结果与输入值相关
func (r *SimpleRandom) SetPrevious(previous uint32) {
	r.Previous = int64(previous)
}

// 获取下一个整型值
func (r *SimpleRandom) uintNext() uint32 {
	num := uint32((a*r.Previous + c) % m)
	r.Previous = int64(num)
	return num
}

// 获取下一个随机int值
func (r *SimpleRandom) Next(minValue int32, maxValue int32) int32 {
	if minValue == maxValue {
		return minValue
	} else if minValue < maxValue {
		return int32((float64(r.uintNext())/float64(m))*float64(maxValue-minValue+1)) + minValue
	} else {
		return int32((float64(r.uintNext())/float64(m))*float64(maxValue-minValue+1)) + maxValue
	}
}
