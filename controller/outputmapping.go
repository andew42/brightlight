package controller

import "github.com/andew42/brightlight/framebuffer"

var mappingEnabled = true
var outputMappingTable = linear128

var linear128 = [256]byte{
	0, 0, 1, 1, 2, 2, 3, 3, 4, 4, 5, 5, 6, 6, 7, 7,
	8, 8, 9, 9, 10, 10, 11, 11, 12, 12, 13, 13, 14, 14, 15, 15,
	16, 16, 17, 17, 18, 18, 19, 19, 20, 20, 21, 21, 22, 22, 23, 23,
	24, 24, 25, 25, 26, 26, 27, 27, 28, 28, 29, 29, 30, 30, 31, 31,
	32, 32, 33, 33, 34, 34, 35, 35, 36, 36, 37, 37, 38, 38, 39, 39,
	40, 40, 41, 41, 42, 42, 43, 43, 44, 44, 45, 45, 46, 46, 47, 47,
	48, 48, 49, 49, 50, 50, 51, 51, 52, 52, 53, 53, 54, 54, 55, 55,
	56, 56, 57, 57, 58, 58, 59, 59, 60, 60, 61, 61, 62, 62, 63, 63,
	64, 64, 65, 65, 66, 66, 67, 67, 68, 68, 69, 69, 70, 70, 71, 71,
	72, 72, 73, 73, 74, 74, 75, 75, 76, 76, 77, 77, 78, 78, 79, 79,
	80, 80, 81, 81, 82, 82, 83, 83, 84, 84, 85, 85, 86, 86, 87, 87,
	88, 88, 89, 89, 90, 90, 91, 91, 92, 92, 93, 93, 94, 94, 95, 95,
	96, 96, 97, 97, 98, 98, 99, 99, 100, 100, 101, 101, 102, 102, 103, 103,
	104, 104, 105, 105, 106, 106, 107, 107, 108, 108, 109, 109, 110, 110, 111, 111,
	112, 112, 113, 113, 114, 114, 115, 115, 116, 116, 117, 117, 118, 118, 119, 119,
	120, 120, 121, 121, 122, 122, 123, 123, 124, 124, 125, 125, 126, 126, 127, 127,
}

var cie128 = [256]byte{
	0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 2, 2, 2, 2, 2,
	2, 2, 2, 2, 2, 2, 2, 2, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 4, 4, 4, 4, 4, 4, 4, 4, 5, 5, 5, 5, 5, 5, 5,
	6, 6, 6, 6, 6, 6, 7, 7, 7, 7, 7, 8, 8, 8, 8, 8,
	9, 9, 9, 9, 10, 10, 11, 11, 10, 11, 10, 10, 11, 12, 12, 12,
	13, 13, 13, 13, 14, 14, 14, 15, 15, 15, 16, 16, 16, 16, 17, 17,
	17, 18, 18, 19, 19, 19, 20, 20, 20, 21, 21, 22, 22, 22, 23, 23,
	24, 24, 24, 25, 25, 26, 26, 27, 27, 28, 28, 28, 29, 29, 30, 30,
	31, 31, 32, 32, 33, 34, 34, 35, 35, 36, 36, 37, 37, 38, 38, 39,
	40, 40, 41, 41, 42, 43, 43, 44, 45, 45, 46, 47, 47, 48, 49, 49,
	50, 51, 51, 52, 53, 53, 54, 55, 56, 56, 57, 58, 59, 59, 60, 61,
	62, 63, 63, 64, 65, 66, 67, 68, 68, 69, 70, 71, 72, 73, 74, 75,
	75, 76, 77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90,
	91, 92, 93, 94, 95, 96, 97, 98, 99, 101, 102, 103, 104, 105, 106, 107,
	108, 110, 111, 112, 113, 114, 115, 117, 118, 119, 120, 122, 123, 124, 125, 127,
}

var cie256 = [256]byte{
	0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 2, 2, 2,
	2, 2, 2, 2, 2, 2, 2, 3, 3, 3, 3, 3, 3, 3, 3, 4,
	4, 4, 4, 4, 4, 5, 5, 5, 5, 5, 6, 6, 6, 6, 6, 7,
	7, 7, 7, 8, 8, 8, 8, 9, 9, 9, 9, 10, 10, 10, 11, 11,
	11, 12, 12, 12, 13, 13, 13, 14, 14, 14, 15, 15, 16, 16, 16, 17,
	17, 18, 18, 19, 19, 20, 20, 20, 21, 21, 22, 22, 23, 24, 24, 25,
	25, 26, 26, 27, 27, 28, 29, 29, 30, 30, 31, 32, 32, 33, 34, 34,
	35, 36, 36, 37, 38, 39, 39, 40, 41, 42, 42, 43, 44, 45, 45, 46,
	47, 48, 49, 50, 51, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61,
	62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 72, 74, 75, 76, 77, 78,
	79, 81, 82, 83, 84, 85, 87, 88, 89, 90, 92, 93, 94, 96, 97, 98,
	100, 101, 103, 104, 105, 107, 108, 110, 111, 113, 114, 116, 117, 119, 120, 122,
	124, 125, 127, 128, 130, 132, 133, 135, 137, 138, 140, 142, 144, 145, 147, 149,
	151, 153, 155, 156, 158, 160, 162, 164, 166, 168, 170, 172, 174, 176, 178, 180,
	182, 184, 186, 188, 190, 192, 194, 197, 199, 201, 203, 205, 208, 210, 212, 215,
	217, 219, 221, 224, 226, 229, 231, 233, 236, 238, 241, 243, 246, 248, 251, 253,
}

func SetOutputMapping(mapping string) {

	switch mapping {
	case "Linear 128":
		outputMappingTable = linear128
		mappingEnabled = true
	case "Cie 128":
		outputMappingTable = cie128
		mappingEnabled = true
	case "Cie 256":
		outputMappingTable = cie256
		mappingEnabled = true
	default:
		mappingEnabled = false
	}
}

func mapOutput(in framebuffer.Rgb) framebuffer.Rgb {

	if !mappingEnabled {
		return in
	} else {
		return framebuffer.Rgb{
			Red:   outputMappingTable[in.Red],
			Green: outputMappingTable[in.Green],
			Blue:  outputMappingTable[in.Blue]}
	}
}
