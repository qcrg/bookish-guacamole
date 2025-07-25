package utils

type Bitwise interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

func Has[T Bitwise](flags T, flag T) bool {
	return flags&flag != 0
}

func HasOneOf[T Bitwise](flags T, vflag ...T) bool {
	for _, flag := range vflag {
		if Has(flags, flag) {
			return true
		}
	}
	return false
}
