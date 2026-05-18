package domain

func ScorePrediction(predictedHome int, predictedAway int, finalHome int, finalAway int) int {
	if predictedHome == finalHome && predictedAway == finalAway {
		return 10
	}

	if sign(predictedHome-predictedAway) == sign(finalHome-finalAway) {
		return 5
	}

	return 0
}

func sign(value int) int {
	switch {
	case value > 0:
		return 1
	case value < 0:
		return -1
	default:
		return 0
	}
}
