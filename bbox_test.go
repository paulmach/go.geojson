package geojson	

import "testing"


func TestBoundingBoxPoint(t *testing.T) {
	g := NewPointGeometry([]float64{1, 2})
	expected := []float64{1, 2, 1, 2}
	value := g.Get_BoundingBox()

	if (expected[0] == value[0] && expected[1] == value[1] && expected[2] == value[2] && expected[3] == value[3]) == false {
		t.Errorf("BoundingBox Point expected: %v, got: %v",expected,value)
	}
}

func TestBoundingBoxMultiPoint(t *testing.T) {
	g := NewMultiPointGeometry([]float64{1, 2}, []float64{3, 4})
	expected := []float64{1, 2, 3, 4}
	value := g.Get_BoundingBox()
	
	if (expected[0] == value[0] && expected[1] == value[1] && expected[2] == value[2] && expected[3] == value[3]) == false {
		t.Errorf("BoundingBox MultiPoint expected: %v, got: %v",expected,value)
	}
}


func TestBoundingBoxLineString(t *testing.T) {
	g := NewLineStringGeometry([][]float64{{1, 2}, {3, 4}})
	expected := []float64{1, 2, 3, 4}
	value := g.Get_BoundingBox()

	if (expected[0] == value[0] && expected[1] == value[1] && expected[2] == value[2] && expected[3] == value[3]) == false {
		t.Errorf("BoundingBox LineString expected: %v, got: %v",expected,value)
	}
}


func TestBoundingBoxMultiLineString(t *testing.T) {
	g := NewMultiLineStringGeometry(
		[][]float64{{1, 2}, {3, 4}},
		[][]float64{{5, 6}, {7, 8}},
	)
	expected := []float64{1, 2, 7, 8}
	value := g.Get_BoundingBox()

	if (expected[0] == value[0] && expected[1] == value[1] && expected[2] == value[2] && expected[3] == value[3]) == false {
		t.Errorf("BoundingBox MultiLineString expected: %v, got: %v",expected,value)
	}
}


func TestBoundingBoxPolygon(t *testing.T) {
	g := NewPolygonGeometry([][][]float64{
		{{1, 2}, {3, 4}},
		{{5, 6}, {7, 8}},
	})
	expected := []float64{1, 2, 7, 8}
	value := g.Get_BoundingBox()

	if (expected[0] == value[0] && expected[1] == value[1] && expected[2] == value[2] && expected[3] == value[3]) == false {
		t.Errorf("BoundingBox Polygon expected: %v, got: %v",expected,value)
	}
}

func TestBoundingBoxMultiPolygon(t *testing.T) {
	g := NewMultiPolygonGeometry(
		[][][]float64{
			{{1, 2}, {3, 4}},
			{{5, 6}, {7, 8}},
		},
		[][][]float64{
			{{8, 7}, {6, 5}},
			{{4, 3}, {2, 1}},
		})
	expected := []float64{1, 1, 7, 8}
	value := g.Get_BoundingBox()

	if (expected[0] == value[0] && expected[1] == value[1] && expected[2] == value[2] && expected[3] == value[3]) == false {
		t.Errorf("BoundingBox MultiPolygon expected: %v, got: %v",expected,value)
	}
}


func TestCollectionGeometry(t *testing.T) {
	g := NewCollectionGeometry(
		NewPointGeometry([]float64{1, 2}),
		NewMultiPointGeometry([]float64{1, 2}, []float64{3, 4}),
	)
	expected := []float64{1, 2, 3, 4}
	value := g.Get_BoundingBox()

	if (expected[0] == value[0] && expected[1] == value[1] && expected[2] == value[2] && expected[3] == value[3]) == false {
		t.Errorf("BoundingBox Collection expected: %v, got: %v",expected,value)
	}
}
