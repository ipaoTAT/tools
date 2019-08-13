package tstat

//function to calculate axis mark
// num: count of marks
// max: max value
// min: min value
// minMark: min mark
// step: step of marks
func CalculateAxisMark(num int, max, min int) (minMark, step int) {
	step = 1
	if max-min > num {
		step = axisAdjustUp((max - min) / num)
	}
	minMark = axisAdjustDown(min)
	return
}

func axisAdjustUp(i int) int {
	std := int(1)
	if i < 0 {
		return -1 * axisAdjustDown(-1*i)
	}
	for ; i > std; std *= 10 {
	}
	if std/5 >= i {
		std = std / 5
	} else if std/2 >= i {
		std = std / 2
	}
	return std
}
func axisAdjustDown(i int) int {
	std := int(1)
	if i < 0 {
		return -1 * axisAdjustUp(-1*i)
	}
	for ; i > std; std *= 10 {
	}
	std = std / 10
	if std*5 <= i {
		std = std * 5
	} else if std*2 <= i {
		std = std * 2
	}
	return std
}
