package location

func MetersToPixels(meters, scale float64) (px float64) {
	return meters * 1000 / scale
}

// func MetersToCells(meters float64) (cells int) {
// 	return
// }

func PixelsToMeters(pixels, scale float64) (meters float64) {
	return pixels * scale / 1000
}

func PixelsToCells(pixels int, scale, cellSizeMeter float64) (cells int) {
	//TODO Убедиться в правильном округлении
	return int(PixelsToMeters(float64(pixels), scale) / cellSizeMeter)
}
