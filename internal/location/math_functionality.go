package location

import "math";

const X, Y, Z int = 0, 1, 2;

type XYcoordinate struct {
    x int
    y int
}

type XYZcoordinate struct {
    x int
    y int
    z int
}

type Vector struct {
    x int
    y int
    z int
};



/**
 * Returns True if segments` projections intersect.
 * @param range_1_s Start coordinate of first segment.
 * @param range_1_e End coordinate of first segment.
 * @param range_2_s Start coordinate of second segment.
 * @param range_2_e End coordinate of second segment.
 */
// function _areProjectionsIntersect(range_1_s: number, range_1_e: number, range_2_s: number, range_2_e: number): boolean {
//     if (range_1_s > range_1_e) {
//         [range_1_s, range_1_e] = [range_1_e, range_1_s];
//     }
//     if (range_2_s > range_2_e) {
//         [range_2_s, range_2_e] = [range_2_e, range_2_s];
//     }

//     return Math.max(range_1_s, range_2_s) <= Math.min(range_1_e, range_2_e);
// }


/**
 * Returns 2D vectors` determinant.
 * @param startPoint1
 * @param int
 * @param int
 */
func _getVectorsDeterminant2D(startPoint1 XYZcoordinate, endPoint1 XYZcoordinate, startPoint2 XYZcoordinate, endPoint2 XYZcoordinate) float32 {
    const xComponent1 int = endPoint1.x - startPoint1.x;
    const yComponent1 int = endPoint1.y - startPoint1.y;
    const xComponent2 int = endPoint2.x - startPoint2.x;
    const yComponent2 int = endPoint2.y - startPoint2.y;

    return xComponent1 * yComponent2 - xComponent2 * yComponent1;
}


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
    if (_areProjectionsIntersect(startPoint1.x, endPoint1.x, startPoint2.x, endPoint2.x) && _areProjectionsIntersect(startPoint1.y, endPoint1.y, startPoint2.y, endPoint2.y)) {

        const d1 float32 = _getVectorsDeterminant2D(startPoint1, endPoint1, startPoint1, startPoint2);
        const d2 float32 = _getVectorsDeterminant2D(startPoint1, endPoint1, startPoint1, endPoint2);
        const d3 float32 = _getVectorsDeterminant2D(startPoint2, endPoint2, startPoint2, startPoint1);
        const d4 float32 = _getVectorsDeterminant2D(startPoint2, endPoint2, startPoint2, endPoint1);

        if (((d1 <= 0 && d2 >= 0) || (d1 >= 0 && d2 <= 0)) && ((d3 <= 0 && d4 >= 0) || (d3 >= 0 && d4 <= 0))) {
            return true;
        }
    }
    return false;
}


/**
 * Returns vector from two points: (x1; y1; z1) -> (x2; y2; z2).
 * @param point1 [x; y; z] coordinates for the first point.
 * @param point2 [x; y; z] coordinates for the second point.
 * @returns Returns vector from two points: (x1; y1; z1) -> (x2; y2; z2).
 */
func _getVector(point1 XYZcoordinate, point2 XYZcoordinate) Vector { 
    const x int = point2.x - point1.x;
    const y int = point2.y - point1.y;
    const z int = point2.z - point1.z;
    return Vector {x: x, y: y, z: z};
}


/**
 * Returns angle in radians between two vectors.
 * Returns ZeroDivisionError if any vector is zero.
 * @param vector1 First vector.
 * @param vector2 Second vector.
 * @returns Angle in radians between two vectors
 */
func _findVectorCosPhi(vector1 Vector, vector2 Vector) float32 {
    // const [a1, b1, c1]: Vector = vector1;
    // const [a2, b2, c2]: Vector = vector2;

    return (vector1.x * vector2.x + vector1.y * vector2.y + vector1.z * vector2.z) / (math.Hypot(vector1.x, vector1.y, vector1.z) * Math.Hypot(vector2.x, vector2.y, vector2.z));
}


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

function _getSlope2D(point1: XYZcoordinate, point2: XYZcoordinate): number {
    return (point2[Y] - point1[Y]) / (point2[X] - point1[X]);
}

function _getYIntercept2D(point: XYZcoordinate, slope: number): number {
    return point[Y] - slope * point[X];
}

function _getLinesIntersection2D(line1Point1: XYZcoordinate, line1Point2: XYZcoordinate, line2Point1: XYZcoordinate, line2Point2: XYZcoordinate): XYcoordinate | null {
    const slope1 = _getSlope2D(line1Point1, line1Point2);
    const slope2 = _getSlope2D(line2Point1, line2Point2);

    if (slope1 === slope2) {
        return null; // lines are parallel
    }

    const yIntercept1 = _getYIntercept2D(line1Point1, slope1);
    const yIntercept2 = _getYIntercept2D(line2Point1, slope2);

    const x = (yIntercept2 - yIntercept1) / (slope1 - slope2);
    const y = slope1 * x + yIntercept1;

    return [x, y];
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

//    const x = (_getDeterminant2D(line1Point1, line2Point1) - line1Point1[1] * line2Point1[0]) / div;
//    const y = (line1Point1[0] * line2Point1[1] - _getDeterminant2D(line1Point1, line2Point1)) / div;

//    return [x, y];
//}
//////////////////////

/**
 * Returns azimuth in degrees between(0; -1; 0) and (x; y; 0) vectors.
 * Returns ZeroDivisionError if (x; y; 0) vector is zero.
 * @param point1 Point 1.
 * @param point2 Point 2.
 * @param offset_deg 
 * @returns Azimuth in degrees between(0; -1; 0) and (x; y; 0) vectors.
 */
function getHorizontalAzimuthDeg(point1: XYZcoordinate, point2: XYZcoordinate, offset_deg = 0): number { 
    const vector: Vector = _getVector(point1, point2);
    return _getAzimuth(vector, offset_deg);
}


/**
 * Returns azimuth in degrees between(0; 0; 1) and (0; y; z) vectors.
 * Returns ZeroDivisionError if (0; y; z) vector is zero.
 * @param point1 Point 1.
 * @param point2 Point 2.
 * @param offset_deg
 * @returns Azimuth in degrees between(0; 0; 1) and (0; y; z) vectors.
 */
function getVerticalAzimuthDeg(point1: XYZcoordinate, point2: XYZcoordinate, offset_deg = 0): number {
    const vector: Vector = _getVector(point1, point2);
    return _getPlunge(vector, offset_deg);
}


/**
 * Returns azimuth in degrees between your vector and (0; -1; 0).
 * Returns ZeroDivisionError if (x;y;0) vector is zero.
 * @param vector
 * @param zero_direction
 * @param offset_deg
 * @returns
 */
function _getAzimuth(vector: Vector, offset_deg = 0): number {
    let angle: number = Math.atan2(vector[X], -vector[Y]) * (180 / Math.PI);
    if (angle < 0) {
        angle += 360;
    }
    return (angle + offset_deg) % 360;
}


/**
 * Returns azimuth in degrees between your vector and (0; 0; 1).
 * @param vector
 * @param offset_deg
 * @returns
 */
function _getPlunge(vector: Vector, offset_deg = 0): number {
    const x: number = vector[X];
    const y: number = vector[Y];
    const z: number = vector[Z];
    const magnitude = Math.hypot(x, y, z);

    let angle: number = Math.asin(z / magnitude) * (180 / Math.PI);
    if (angle < 0) {
        angle += 360;
    }
    return (angle + offset_deg) % 360;
}


function _getProjectionsOverlapLength(range_1_s: number, range_1_e: number, range_2_s: number, range_2_e: number): number {
    if (range_1_s > range_1_e) {
        [range_1_s, range_1_e] = [range_1_e, range_1_s];
    }
    if (range_2_s > range_2_e) {
        [range_2_s, range_2_e] = [range_2_e, range_2_s];
    }

    return Math.max(range_1_s, range_2_s) - Math.min(range_1_e, range_2_e);
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

/**
 * Returns the path length(in meters) of line3D that goes throw the wall.
 * @param line1Point1 Client's coordinate. 
 * @param line1Point2 Sensor;s coordinate.
 * @param line2Point1 Wall's start point.
 * @param line2Point2 Wall's end point.
 * @param thickness Wall's thickness.
 * @returns The path length(in meters) of line3D that goes throw the wall.
 */
function getWallPathLengthThrough(
    line1Point1: XYZcoordinate,
    line1Point2: XYZcoordinate,
    line2Point1: XYZcoordinate,
    line2Point2: XYZcoordinate,
    thickness: number,
    cell_size_meters: number): number {

    if (_areSegmentsIntersect2D(line1Point1, line1Point2, line2Point1, line2Point2)) {
        const vector1: Vector = _getVector(line1Point1, line1Point2)
        const vector2: Vector = _getVector(line2Point1, line2Point2)

        //if (vector1.every(el => el === 0) || vector2.every(el => el === 0)) // need to fix later
        //    return thickness;

        const cos_fi: number = _findVectorCosPhi(vector1, vector2);
        if (!cos_fi || !isFinite(cos_fi))
            return thickness;

        const height_in_cells: number = line1Point2[Z] - line1Point1[Z];
        let scaled_height_in_cells = 0;
        if (height_in_cells) {
            const XY: XYcoordinate | null = _getLinesIntersection2D(line1Point1, line1Point2, line2Point1, line2Point2);
            if (XY !== null) {
                const [x,]: XYcoordinate = XY;
                const small_line_length: number = Math.abs(x - line1Point1[X]);
                const big_line_length: number = Math.abs(line1Point2[X] - line1Point1[X]);
                const scale: number = small_line_length / big_line_length;
                if (Number.isFinite(scale))
                    scaled_height_in_cells = Number((scale * height_in_cells).toFixed(2));
            }
        }
        const sin_fi: number = Math.sqrt(1 - cos_fi ** 2);
        let wallPathLengthThrough = Number(Math.hypot(thickness / sin_fi, scaled_height_in_cells * cell_size_meters).toFixed(2));

        // if parallel
        if (Math.abs(cos_fi) == 1) {
            const xProjection = _getProjectionsOverlapLength(line1Point1[X], line1Point2[X], line2Point1[X], line2Point2[X]);
            const yProjection = _getProjectionsOverlapLength(line1Point1[Y], line1Point2[Y], line2Point1[Y], line2Point2[Y]);
            //if (line1Point2[Y] >= line2Point2[Y])
            //    yProjection = line2Point2[Y] - line2Point1[Y];
            //else
            //    yProjection = line1Point2[Y] - line2Point1[Y];

            //if (line1Point2[X] >= line2Point2[X])
            //    xProjection = line2Point2[X] - line2Point1[X];
            //else
            //    xProjection = line1Point2[X] - line2Point1[X];

            wallPathLengthThrough = Number((Math.hypot(xProjection, yProjection, scaled_height_in_cells) * cell_size_meters).toFixed(2));
            if (!wallPathLengthThrough)
                return Number(Math.hypot(thickness, scaled_height_in_cells * cell_size_meters).toFixed(2));
        }

        return wallPathLengthThrough;
    }

    return 0;       
}


// export { _getVector, _areSegmentsIntersect2D, _findVectorCosPhi, _getLinesIntersection2D };