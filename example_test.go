package geojson_test

import (
	geojson "."
	"fmt"
)

func ExampleUnmarshalFeatureCollection() {
	rawFeatureJSON := []byte(`
		  { "type": "FeatureCollection",
		    "features": [
		      { "type": "Feature",
		        "geometry": {"type": "Point", "coordinates": [102.0, 0.5]},
		        "properties": {"prop0": "value0"}
		      }
		    ]
		  }`)

	fc, err := geojson.UnmarshalFeatureCollection(rawFeatureJSON)
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}

	fmt.Printf("%s", fc.Features[0].Properties["prop0"])
	// Output: value0
}

func ExampleUnmarshalGeometry() {
	rawGeometryJSON := []byte(`{"type": "Point", "coordinates": [102.0, 0.5]}`)
	g, err := geojson.UnmarshalGeometry(rawGeometryJSON)
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}

	fmt.Printf("%s", g.Type)
	// Output: Point
}

func ExampleMarshalJSON() {
	fc := geojson.NewFeatureCollection()
	fc.AddFeature(geojson.NewPointFeature([]float64{1, 2}))
	rawJSON, err := fc.MarshalJSON()
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}

	fmt.Printf("%s", string(rawJSON))
	// Output: {"type":"FeatureCollection","features":[{"type":"Feature","geometry":{"type":"Point","coordinates":[1,2]},"properties":null}]}
}
