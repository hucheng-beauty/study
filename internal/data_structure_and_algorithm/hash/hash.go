package hash

/*
	原理:
		key => hash函数 => 数组下标 => value(数组元素)

	冲突解决:
		链表法
		开放寻止法(线性探测、二次探测、双重哈希)

	扩容:
		类似数组扩容,分摊移数据

	应用场景:
		位图: 用一个二进制位(bit)表示真假
		布隆过滤器: key => hash函数 => 位图
	优缺点:
		优点:
		缺点:
*/

const BitType int = 8

// Bitmap implement a bitmap.
type Bitmap struct {
	s     []byte
	nBits int
}

func (b *Bitmap) Get(v int) bool {
	if v > b.nBits {
		return false
	}

	sIdx := v / BitType   // bitmap数组下标
	bitIdx := v % BitType // 数组数据的所在的二进制位

	return (b.s[sIdx] & (1 << bitIdx)) != 0
}

func (b *Bitmap) Set(v int) {
	if v > b.nBits {
		return
	}

	sIdx := v / BitType
	bitIdx := v % BitType

	b.s[sIdx] |= 1 << bitIdx
}

func NewBitmap(nBits int) *Bitmap {
	return &Bitmap{
		s:     make([]byte, (nBits-1)/BitType+1),
		nBits: nBits,
	}
}
