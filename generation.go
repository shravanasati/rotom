package main

var generationRanges = map[int][2]int{
	1: {1, 151},
	2: {151, 251},
	3: {252, 386},
	4: {387, 493},
	5: {494, 649},
	6: {650, 721},
	7: {722, 809},
	8: {810, 905},
	9: {906, 1025},
}

func generationFromDex(dex int) int {
	for gen, rng := range generationRanges {
		if dex >= rng[0] && dex <= rng[1] {
			return gen
		}
	}
	return 0
}

func formatGeneration(g int) string {
	switch g {
	case 1:
		return "gen I"
	case 2:
		return "gen II"
	case 3:
		return "gen III"
	case 4:
		return "gen IV"
	case 5:
		return "gen V"
	case 6:
		return "gen VI"
	case 7:
		return "gen VII"
	case 8:
		return "gen VIII"
	case 9:
		return "gen IX"
	default:
		return "gen unknown"
	}
}
