package utils

func TestFlags[L, R ~int32](lhs L, rhs R) bool {
	return int32(lhs) & int32(rhs) != 0
}
