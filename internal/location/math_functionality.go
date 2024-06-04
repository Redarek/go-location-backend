package location

import (
	"fmt"
	"math"
	. "math"
)

// const X, Y, Z int = 0, 1, 2

type XYcoordinate struct {
	x, y float64
}

type XYZcoordinate struct {
	x, y, z float64
}

type Vector struct {
	x, y, z float64
}

func main() {
	fmt.Print("Hello")
}

// ГОТОВО
// Аналог Math.hypot() в JS, но для трехмерного вектора
func Magnitude(vector Vector) float64 {
	var a, b, c float64

	a = vector.x * vector.x
	b = vector.y * vector.y
	c = vector.z * vector.z

	return Sqrt(a + b + c)
}

// ГОТОВО
/**
 * Returns True if segments` projections intersect.
 * @param range_1_s Start coordinate of first segment.
 * @param range_1_e End coordinate of first segment.
 * @param range_2_s Start coordinate of second segment.
 * @param range_2_e End coordinate of second segment.
 */
func _areProjectionsIntersect(range_1_s float64, range_1_e float64, range_2_s float64, range_2_e float64) bool {
	if range_1_s > range_1_e {
		range_1_s, range_1_e = range_1_e, range_1_s
	}
	if range_2_s > range_2_e {
		range_2_s, range_2_e = range_2_e, range_2_s
	}

	return Max(range_1_s, range_2_s) <= Min(range_1_e, range_2_e)
}

// ГОТОВО
/**
 * Returns 2D vectors` determinant.
 * @param startPoint1
 * @param int
 * @param int
 */
func _getVectorsDeterminant2D(startPoint1 XYZcoordinate, endPoint1 XYZcoordinate, startPoint2 XYZcoordinate, endPoint2 XYZcoordinate) float64 {
	var xComponent1 float64 = endPoint1.x - startPoint1.x
	var yComponent1 float64 = endPoint1.y - startPoint1.y
	var xComponent2 float64 = endPoint2.x - startPoint2.x
	var yComponent2 float64 = endPoint2.y - startPoint2.y

	return xComponent1*yComponent2 - xComponent2*yComponent1
}

// ГОТОВО
/**
 * Returns True in case of 2D segment intersection.
 * Z coordinate is unused!
 * @param startPoint1 first point of the first segment.
 * @param endPoint1 second point of the first segment.
 * @param startPoint2 first point of the second segment.
 * @param endPoint2 second point of the second segment.
 * @returns True in case of 2D segment intersection.
 */
func _areSegmentsIntersect2D(startPoint1 XYZcoordinate, endPoint1 XYZcoordinate, startPoint2 XYZcoordinate, endPoint2 XYZcoordinate) bool {
	if _areProjectionsIntersect(startPoint1.x, endPoint1.x, startPoint2.x, endPoint2.x) && _areProjectionsIntersect(startPoint1.y, endPoint1.y, startPoint2.y, endPoint2.y) {

		var d1 float64 = _getVectorsDeterminant2D(startPoint1, endPoint1, startPoint1, startPoint2)
		var d2 float64 = _getVectorsDeterminant2D(startPoint1, endPoint1, startPoint1, endPoint2)
		var d3 float64 = _getVectorsDeterminant2D(startPoint2, endPoint2, startPoint2, startPoint1)
		var d4 float64 = _getVectorsDeterminant2D(startPoint2, endPoint2, startPoint2, endPoint1)

		if ((d1 <= 0 && d2 >= 0) || (d1 >= 0 && d2 <= 0)) && ((d3 <= 0 && d4 >= 0) || (d3 >= 0 && d4 <= 0)) {
			return true
		}
	}
	return false
}

// ГОТОВО
/**
 * Returns vector from two points: (x1; y1; z1) -> (x2; y2; z2).
 * @param point1 [X; Y; z] coordinates for the first point.
 * @param point2 [X; Y; z] coordinates for the second point.
 * @returns Returns vector from two points: (x1; y1; z1) -> (x2; y2; z2).
 */
func _getVector(point1 XYZcoordinate, point2 XYZcoordinate) Vector {
	var x float64 = point2.x - point1.x
	var y float64 = point2.y - point1.y
	var z float64 = point2.z - point1.z
	return Vector{x: x, y: y, z: z}
}

// ГОТОВО
/**
 * Returns angle in radians between two vectors.
 * Returns ZeroDivisionError if any vector is zero.
 * @param vector1 First vector.
 * @param vector2 Second vector.
 * @returns Angle in radians between two vectors
 */
func _findVectorCosPhi(vector1 Vector, vector2 Vector) float64 {
	// const [a1, b1, c1]: Vector = vector1;
	// const [a2, b2, c2]: Vector = vector2;

	return (vector1.x*vector2.x + vector1.y*vector2.y + vector1.z*vector2.z) / (Magnitude(vector1) * Magnitude(vector2))
}

//////////////////OLD
///**
// * Returns determinant of vector.
// * Z coordinate is unused!
// * @param vector1
// * @param vector2
// * @returns
// */
//function _getDeterminant2D(vector1: Vector, vector2: Vector): number {
//    return vector1[X] * vector2[Y] - vector1[Y] * vector2[X];
//}
////////////////////OLD

// ГОТОВО
func _getSlope2D(point1 XYZcoordinate, point2 XYZcoordinate) float64 {
	return (point2.y - point1.y) / (point2.x - point1.x)
}

// ГОТОВО
func _getYIntercept2D(point XYZcoordinate, slope float64) float64 {
	return point.y - slope*point.x
}

// TODO
func _getLinesIntersection2D(line1Point1 XYZcoordinate, line1Point2 XYZcoordinate, line2Point1 XYZcoordinate, line2Point2 XYZcoordinate) (XYcoordinate, bool) {
	var slope1 float64 = _getSlope2D(line1Point1, line1Point2)
	var slope2 float64 = _getSlope2D(line2Point1, line2Point2)

	if slope1 == slope2 {
		return XYcoordinate{x: -1, y: -1}, false // lines are parallel
	}

	var yIntercept1 float64 = _getYIntercept2D(line1Point1, slope1)
	var yIntercept2 float64 = _getYIntercept2D(line2Point1, slope2)

	var x = (yIntercept2 - yIntercept1) / (slope1 - slope2)
	var y = slope1*x + yIntercept1

	return XYcoordinate{x: x, y: y}, true
}

/////////////////Старое
///**
// * Returns the lines cross point or null if parallel.
// * Z coordinate is unused!
// * @param a_point1
// * @param a_point2
// * @param b_point1
// * @param b_point2
// * @returns
// */
//function _getLinesIntersection(line1Point1: XYZcoordinate, line1Point2: XYZcoordinate, line2Point1: XYZcoordinate, line2Point2: XYZcoordinate,): XYcoordinate | null {
//    const div: number = _getDeterminant2D(line1Point1, line2Point1);
//    if (!div)
//        return null;

//    const X = (_getDeterminant2D(line1Point1, line2Point1) - line1Point1[1] * line2Point1[0]) / div;
//    const Y = (line1Point1[0] * line2Point1[1] - _getDeterminant2D(line1Point1, line2Point1)) / div;

//    return [X, Y];
//}
//////////////////////

// ГОТОВО
/**
 * Returns azimuth in degrees between(0; -1; 0) and (X; Y; 0) vectors.
 * Returns ZeroDivisionError if (X; Y; 0) vector is zero.
 * @param point1 Point 1.
 * @param point2 Point 2.
 * @param offset_deg ПО УМОЛЧАНИЮ 0!!!!
 * @returns Azimuth in degrees between(0; -1; 0) and (X; Y; 0) vectors.
 */
func getHorizontalAzimuthDeg(point1 XYZcoordinate, point2 XYZcoordinate, offset_deg int) int {
	var vector Vector = _getVector(point1, point2)
	return _getAzimuth(vector, offset_deg)
}

// ГОТОВО
/**
 * Returns azimuth in degrees between(0; 0; 1) and (0; Y; z) vectors.
 * Returns ZeroDivisionError if (0; Y; z) vector is zero.
 * @param point1 Point 1.
 * @param point2 Point 2.
 * @param offset_deg ПО УМОЛЧАНИЮ 0
 * @returns Azimuth in degrees between(0; 0; 1) and (0; Y; z) vectors.
 */
func getVerticalAzimuthDeg(point1 XYZcoordinate, point2 XYZcoordinate, offset_deg int) int {
	var vector Vector = _getVector(point1, point2)
	return _getPlunge(vector, offset_deg)
}

// ГОТОВО
/**
 * Returns azimuth in degrees between your vector and (0; -1; 0).
 * Returns ZeroDivisionError if (X;Y;0) vector is zero.
 * @param vector
 * @param zero_direction
 * @param offset_deg ПО УМОЛЧАНИЮ 0
 * @returns
 */
func _getAzimuth(vector Vector, offset_deg int) int {
	var angle float64 = Atan2(vector.x, -vector.y) * (180 / math.Pi)
	if angle < 0 {
		angle += 360
	}
	return (int(angle) + offset_deg) % 360
}

/**
 * Returns azimuth in degrees between your vector and (0; 0; 1).
 * @param vector
 * @param offset_deg ПО УМОЛЧАНИЮ 0!!!!!!!!!
 * @returns
 */
func _getPlunge(vector Vector, offset_deg int) int {
	// var X int = vector.X
	// var Y int = vector.Y
	var z float64 = vector.z
	var magnitude float64 = Magnitude(vector)

	var angle float64 = Asin(z/magnitude) * (180 / Pi)
	if angle < 0 {
		angle += 360
	}
	return (int(angle) + offset_deg) % 360
}

// ГОТОВО
func _getProjectionsOverlapLength(range_1_s float64, range_1_e float64, range_2_s float64, range_2_e float64) float64 {
	if range_1_s > range_1_e {
		range_1_s, range_1_e = range_1_e, range_1_s
	}
	if range_2_s > range_2_e {
		range_2_s, range_2_e = range_2_e, range_2_s
	}

	return Max(range_1_s, range_2_s) - Min(range_1_e, range_2_e)
}

///////OLD
//fi = (math.degrees(math.acos(_find_vectors_cos_fi(zero_direction, vector))) + offset_deg) % 360;
//if (vector[X] <= 0) {
//    return round(fi);
//}
//else {
//    return round(360 - fi);
//}
//}
//////////

// Что делать с NULL?
/**
 * Returns the path length(in meters) of line3D that goes throw the wall.
 * @param line1Point1 Client's coordinate.
 * @param line1Point2 Sensor;s coordinate.
 * @param line2Point1 Wall's start point.
 * @param line2Point2 Wall's end point.
 * @param thickness Wall's thickness.
 * @returns The path length(in meters) of line3D that goes throw the wall.
 */
func getWallPathLengthThrough(
	line1Point1 XYZcoordinate,
	line1Point2 XYZcoordinate,
	line2Point1 XYZcoordinate,
	line2Point2 XYZcoordinate,
	thickness float64,
	cell_size_meters float64) float64 {

	if _areSegmentsIntersect2D(line1Point1, line1Point2, line2Point1, line2Point2) {
		var vector1 Vector = _getVector(line1Point1, line1Point2)
		var vector2 Vector = _getVector(line2Point1, line2Point2)

		//if (vector1.every(el => el === 0) || vector2.every(el => el === 0)) // need to fix later
		//    return thickness;

		var cos_fi float64 = _findVectorCosPhi(vector1, vector2)
		if cos_fi == 0 || IsNaN(cos_fi) || IsInf(cos_fi, 0) {
			return thickness
		}

		var height_in_cells float64 = line1Point2.z - line1Point1.z
		var scaled_height_in_cells float64 = 0
		if height_in_cells != 0 {
			XY, tmp := _getLinesIntersection2D(line1Point1, line1Point2, line2Point1, line2Point2) // or null
			if tmp == true {
				var x float64 = XY.x
				var small_line_length float64 = Abs(x - line1Point1.x)
				var big_line_length float64 = Abs(line1Point2.x - line1Point1.x)
				var scale float64 = small_line_length / big_line_length
				if !IsNaN(scale) || !IsInf(scale, 0) {
					scaled_height_in_cells = scale * height_in_cells // здесь было округление до 2
				}
			}
		}
		var sin_fi float64 = Sqrt(1 - Pow(cos_fi, 2))
		var wallPathLengthThrough = Hypot(thickness/sin_fi, scaled_height_in_cells*cell_size_meters) // Здесь было округление до 2

		// if parallel
		if Abs(cos_fi) == 1 {
			var xProjection float64 = _getProjectionsOverlapLength(line1Point1.x, line1Point2.x, line2Point1.x, line2Point2.x)
			var yProjection float64 = _getProjectionsOverlapLength(line1Point1.y, line1Point2.y, line2Point1.y, line2Point2.y)
			//if (line1Point2[Y] >= line2Point2[Y])
			//    yProjection = line2Point2[Y] - line2Point1[Y];
			//else
			//    yProjection = line1Point2[Y] - line2Point1[Y];

			//if (line1Point2[X] >= line2Point2[X])
			//    xProjection = line2Point2[X] - line2Point1[X];
			//else
			//    xProjection = line1Point2[X] - line2Point1[X];
			var vector Vector = Vector{xProjection, yProjection, scaled_height_in_cells}
			wallPathLengthThrough = Magnitude(vector) * cell_size_meters // было округление до 2
			if wallPathLengthThrough == 0 {
				return Hypot(thickness, scaled_height_in_cells*cell_size_meters) // было округление до 2
			}
		}

		return wallPathLengthThrough
	}

	return 0
}

// export { _getVector, _areSegmentsIntersect2D, _findVectorCosPhi, _getLinesIntersection2D };
