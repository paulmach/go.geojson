package geojson

import (
	"reflect"
	"testing"
)

func TestGetBoundingBox(t *testing.T) {
	rawJSON := ` { "type": "FeatureCollection",
	"features": [
	  { "type": "Feature",
		"geometry": {"type": "Point", "coordinates": [102.0, 0.5]},
		"properties": {"prop0": "value0"}
		},
	  { "type": "Feature",
		"geometry": {
		  "type": "LineString",
		  "coordinates": [
			[102.0, 0.0], [103.0, 1.0], [104.0, 0.0], [105.0, 1.0]
			]
		  },
		"properties": {
		  "prop0": "value0",
		  "prop1": 0.0
		  }
		},
	  { "type": "Feature",
		 "geometry": {
		   "type": "Polygon",
		   "coordinates": [
			 [ [100.5, 0.2], [101.5, 0.0], [101.0, 1.0],
			   [100.1, 1.0], [100.5, 0.2] ]
			 ]
		 },
		 "properties": {
		   "prop0": "value0",
		   "prop1": {"this": "that"}
		   }
		 }
	   ]
  	}`

	fc, err := UnmarshalFeatureCollection([]byte(rawJSON))
	if err != nil {
		t.Fatalf("should unmarshal feature collection without issue, err %v", err)
	}

	if fc.Type != "FeatureCollection" {
		t.Errorf("should have type of FeatureCollection, got %v", fc.Type)
	}

	if len(fc.Features) != 3 {
		t.Errorf("should have 3 features but got %d", len(fc.Features))
	}

	for i := 0; i < len(fc.Features); i++ {
		f := fc.Features[i]
		min, max := f.ComputeBoundingBox(true)
		switch f.Geometry.Type {
		case GeometryPoint:
			if !(reflect.DeepEqual(min, []float64{102.0, 0.5}) && reflect.DeepEqual(max, []float64{102.0, 0.5})) {
				t.Errorf("should have Point Boundingbox is [{102.0, 0.5},{102.0, 0.5}], got %v,%v", min, max)
			}
		case GeometryLineString:
			if !(reflect.DeepEqual(min, []float64{102.0, 0.0}) && reflect.DeepEqual(max, []float64{105.0, 1.0})) {
				t.Errorf("should have LineString Boundingbox is [{102.0, 0.0},{105.0, 1.0}], got %v,%v", min, max)
			}
		case GeometryPolygon:
			if !(reflect.DeepEqual(min, []float64{100.1, 0.0}) && reflect.DeepEqual(max, []float64{101.5, 1.0})) {
				t.Errorf("should have Point Boundingbox is [{100.1, 0.0},{101.5, 1.0}], got %v,%v", min, max)
			}
		}
	}

	min, max := fc.ComputeBoundingBox(true)
	if !(reflect.DeepEqual(min, []float64{100.1, 0.0}) && reflect.DeepEqual(max, []float64{105.0, 1.0})) {
		t.Errorf("should have FeatureCollection Boundingbox is [{100.1, 0.0},{105.0, 1.0}], got %v,%v", min, max)
	}
}
