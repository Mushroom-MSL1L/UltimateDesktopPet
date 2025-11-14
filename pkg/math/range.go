package math

func InRange(x, max, min int16) int16 {
	if x < min {
		return min 
	}
	if x > max {
		return max 
	}
	return x 
}