package main

func int32Min(n, v int32) int32 {
	if n < v {
		return n
	}
	return v
}

func int32Max(n, v int32) int32 {
	if n > v {
		return n
	}
	return v
}

func int32Abs(n int32) int32 {
	if n < 0 {
		return -n
	}
	return n
}

func intAbs(n int) int {
	if n < 0 {
		return int(-n)
	}
	return int(n)
}

func intMax(n, j int) int {
	if n > j {
		return n
	}
	return j
}

func intMin(n, j int) int {
	if n < j {
		return n
	}
	return j
}
