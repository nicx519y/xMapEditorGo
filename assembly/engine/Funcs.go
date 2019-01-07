package main

import "math"

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

func intMax(first int, args ...int) int {
	for _, v := range args {
		if first < v {
			first = v
		}
	}
	return first
}

func intMin(first int, args ...int) int {
	for _, v := range args {
		if first > v {
			first = v
		}
	}
	return first
}

func round(n float64) int {
	return int(math.Floor(n + 0.5))
}

// 逆时针旋转 一个点绕另一个点旋转后的坐标
func nrotate(angle float64, valuex, valuey, pointx, pointy int) (int, int) {
	nRotatex := float64(valuex-pointx)*math.Cos(angle) - float64(valuey-pointy)*math.Sin(angle) + float64(pointx)
	nRotatey := float64(valuey-pointy)*math.Sin(angle) + float64(valuey-pointy)*math.Cos(angle) + float64(pointy)
	return round(nRotatex), round(nRotatey)
}

// 顺时针旋转 一个点绕另一个点旋转后的坐标
func srotate(angle float64, valuex, valuey, pointx, pointy int) (int, int) {
	sRotatex := float64(valuex-pointx)*math.Cos(angle) + float64(valuey-pointy)*math.Sin(angle) + float64(pointx)
	sRotatey := float64(valuey-pointy)*math.Cos(angle) - float64(valuex-pointx)*math.Sin(angle) + float64(pointy)
	return round(sRotatex), round(sRotatey)
}
