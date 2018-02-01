package geojson

// BoundingBox implementation as per https://tools.ietf.org/html/rfc7946
// BoundingBox syntax: "bbox": [west, south, east, north]
// BoundingBox defaults "bbox": [-180.0, -90.0, 180.0, 90.0]
func BoundingBox_Points(pts [][]float64) []float64 {
	// setting opposite default values
	west, south, east, north := 180.0, 90.0, -180.0, -90.0

	for _, pt := range pts {
		x, y := pt[0], pt[1]
		// can only be one condition
		// using else if reduces one comparison
		if x < west {
			west = x
		} else if x > east {
			east = x
		}

		if y < south {
			south = y
		} else if y > north {
			north = y
		}
	}
	return []float64{west, south, east, north}
}

func Push_Two_BoundingBoxs(bb1 []float64, bb2 []float64) []float64 {
	// setting opposite default values
	west, south, east, north := 180.0, 90.0, -180.0, -90.0

	// setting bb1 and bb2
	west1, south1, east1, north1 := bb1[0], bb1[1], bb1[2], bb1[3]
	west2, south2, east2, north2 := bb2[0], bb2[1], bb2[2], bb2[3]

	// handling west values: min
	if west1 < west2 {
		west = west1
	} else {
		west = west2
	}

	// handling south values: min
	if south1 < south2 {
		south = south1
	} else {
		south = south2
	}

	// handling east values: max
	if east1 > east2 {
		east = east1
	} else {
		east = east2
	}

	// handling north values: max
	if north1 > north2 {
		north = north1
	} else {
		north = north2
	}

	return []float64{west, south, east, north}
}

// this functions takes an array of bounding box objects and
// pushses them all out
func Expand_BoundingBoxs(bboxs [][]float64) []float64 {
	bbox := bboxs[0]
	for _, temp_bbox := range bboxs[1:] {
		bbox = Push_Two_BoundingBoxs(bbox, temp_bbox)
	}
	return bbox
}

// boudning box on a normal point geometry
// relatively useless
func BoundingBox_PointGeometry(pt []float64) []float64 {
	return []float64{pt[0], pt[1], pt[0], pt[1]}
}

// Returns BoundingBox for a MultiPoint
func BoundingBox_MultiPointGeometry(pts [][]float64) []float64 {
	return BoundingBox_Points(pts)
}

// Returns BoundingBox for a LineString
func BoundingBox_LineStringGeometry(line [][]float64) []float64 {
	return BoundingBox_Points(line)
}

// Returns BoundingBox for a MultiLineString
func BoundingBox_MultiLineStringGeometry(multiline [][][]float64) []float64 {
	bboxs := [][]float64{}
	for _, line := range multiline {
		bboxs = append(bboxs, BoundingBox_Points(line))
	}
	return Expand_BoundingBoxs(bboxs)
}

// Returns BoundingBox for a Polygon
func BoundingBox_PolygonGeometry(polygon [][][]float64) []float64 {
	bboxs := [][]float64{}
	for _, cont := range polygon {
		bboxs = append(bboxs, BoundingBox_Points(cont))
	}
	return Expand_BoundingBoxs(bboxs)
}

// Returns BoundingBox for a Polygon
func BoundingBox_MultiPolygonGeometry(multipolygon [][][][]float64) []float64 {
	bboxs := [][]float64{}
	for _, polygon := range multipolygon {
		for _, cont := range polygon {
			bboxs = append(bboxs, BoundingBox_Points(cont))
		}
	}
	return Expand_BoundingBoxs(bboxs)
}

// Returns a BoundingBox for a geometry collection
func BoundingBox_GeometryCollection(gs []*Geometry) []float64 {
	bboxs := [][]float64{}
	for _, g := range gs {
		bboxs = append(bboxs, g.Get_BoundingBox())
	}
	return Expand_BoundingBoxs(bboxs)
}

// retrieves a boundingbox for a given geometry
func (g *Geometry) Get_BoundingBox() []float64 {
	switch g.Type {
	case GeometryPoint:
		return BoundingBox_PointGeometry(g.Point)
	case GeometryMultiPoint:
		return BoundingBox_MultiPointGeometry(g.MultiPoint)
	case GeometryLineString:
		return BoundingBox_LineStringGeometry(g.LineString)
	case GeometryMultiLineString:
		return BoundingBox_MultiLineStringGeometry(g.MultiLineString)
	case GeometryPolygon:
		return BoundingBox_PolygonGeometry(g.Polygon)
	case GeometryMultiPolygon:
		return BoundingBox_MultiPolygonGeometry(g.MultiPolygon)
	case GeometryCollection:
		return BoundingBox_GeometryCollection(g.Geometries)
	}
	return []float64{}
}

func (feature *Feature) Get_BoundingBox() {
	bb := feature.Geometry.Get_BoundingBox()
	feature.BoundingBox = bb
}
