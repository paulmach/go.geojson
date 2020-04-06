package geojson

import (
	"fmt"
	"math"
)

func decodeBoundingBox(bb interface{}) ([]float64, error) {
	if bb == nil {
		return nil, nil
	}

	switch f := bb.(type) {
	case []float64:
		return f, nil
	case []interface{}:
		bb := make([]float64, 0, 4)
		for _, v := range f {
			switch c := v.(type) {
			case float64:
				bb = append(bb, c)
			default:
				return nil, fmt.Errorf("bounding box coordinate not usable, got %T", v)
			}

		}
		return bb, nil
	default:
		return nil, fmt.Errorf("bounding box property not usable, got %T", bb)
	}
}

// checkBoundingBox Check if BoundingBox is valid
func checkBoundingBox(boundingbox []float64) bool {
	if boundingbox == nil || len(boundingbox)%2 != 0 {
		return false
	}
	// x0,y0,... x1,y1,...
	dimension := len(boundingbox) / 2
	for i := 0; i < dimension; i++ {
		if boundingbox[i] > boundingbox[i+dimension] {
			return false
		}
	}
	return true
}

// ComputeBoundingBox  Get Geometry geo outsourcing box
func (g *Geometry) ComputeBoundingBox(force bool) (min, max []float64) {
	if !force && checkBoundingBox(g.BoundingBox) {
		dimension := len(g.BoundingBox) / 2
		return g.BoundingBox[0:dimension], g.BoundingBox[dimension:]
	}

	bmin := make([]float64, 0, 16)
	bmax := make([]float64, 0, 16)
	switch g.Type {
	case GeometryPoint:
		bmin = append(bmin, g.Point...)
		bmax = append(bmax, g.Point...)
	case GeometryMultiPoint:
		dimension := len(g.MultiPoint[0])
		bmin = append(bmin, g.MultiPoint[0]...)
		bmax = append(bmax, g.MultiPoint[0]...)
		for i := 1; i < len(g.MultiPoint); i++ {
			for d := 0; d < dimension; d++ {
				bmin[d] = math.Min(bmin[d], g.MultiPoint[i][d])
				bmax[d] = math.Max(bmax[d], g.MultiPoint[i][d])
			}
		}
	case GeometryLineString:
		dimension := len(g.LineString[0])
		bmin = append(bmin, g.LineString[0]...)
		bmax = append(bmax, g.LineString[0]...)
		for i := 1; i < len(g.LineString); i++ {
			for d := 0; d < dimension; d++ {
				bmin[d] = math.Min(bmin[d], g.LineString[i][d])
				bmax[d] = math.Max(bmax[d], g.LineString[i][d])
			}
		}
	case GeometryMultiLineString:
		dimension := len(g.MultiLineString[0][0])
		bmin = append(bmin, g.MultiLineString[0][0]...)
		bmax = append(bmax, g.MultiLineString[0][0]...)
		for line := 0; line < len(g.MultiLineString); line++ {
			linestring := g.MultiLineString[line]
			for i := 0; i < len(linestring); i++ {
				for d := 0; d < dimension; d++ {
					bmin[d] = math.Min(bmin[d], linestring[i][d])
					bmax[d] = math.Max(bmax[d], linestring[i][d])
				}
			}
		}
	case GeometryPolygon:
		dimension := len(g.Polygon[0][0])
		bmin = append(bmin, g.Polygon[0][0]...)
		bmax = append(bmax, g.Polygon[0][0]...)
		for line := 0; line < len(g.Polygon); line++ {
			linestring := g.Polygon[line]
			for i := 0; i < len(linestring); i++ {
				for d := 0; d < dimension; d++ {
					bmin[d] = math.Min(bmin[d], linestring[i][d])
					bmax[d] = math.Max(bmax[d], linestring[i][d])
				}
			}
		}
	case GeometryMultiPolygon:
		dimension := len(g.MultiPolygon[0][0][0])
		bmin = append(bmin, g.MultiPolygon[0][0][0]...)
		bmax = append(bmax, g.MultiPolygon[0][0][0]...)
		for poly := 0; poly < len(g.MultiPolygon); poly++ {
			for line := 0; line < len(g.MultiPolygon[poly]); line++ {
				linestring := g.MultiPolygon[poly][line]
				for i := 0; i < len(linestring); i++ {
					for d := 0; d < dimension; d++ {
						bmin[d] = math.Min(bmin[d], linestring[i][d])
						bmax[d] = math.Max(bmax[d], linestring[i][d])
					}
				}
			}
		}
	case GeometryCollection:
		if g.Geometries == nil || len(g.Geometries) == 0 {
			return nil, nil
		}
		if bmin, bmax = g.Geometries[0].ComputeBoundingBox(force); bmin == nil {
			return nil, nil
		}
		dimension := len(bmin)
		for i := 1; i < len(g.Geometries); i++ {
			tmin, tmax := g.Geometries[i].ComputeBoundingBox(force)
			if tmin == nil || len(tmin) != dimension {
				return nil, nil
			}
			for d := 0; d < dimension; d++ {
				bmin[d] = math.Min(bmin[d], tmin[d])
				bmax[d] = math.Max(bmax[d], tmax[d])
			}
		}
	}
	return bmin, bmax
}

// ComputeBoundingBox  Get Feature geo outsourcing box
func (f *Feature) ComputeBoundingBox(force bool) (bmin, bmax []float64) {
	if !force && checkBoundingBox(f.BoundingBox) {
		dimension := len(f.BoundingBox) / 2
		bmin = f.BoundingBox[:dimension]
		bmax = f.BoundingBox[dimension:]
		return bmin, bmax
	}
	return f.Geometry.ComputeBoundingBox(force)
}

// ComputeBoundingBox  Get FeatureCollection geo outsourcing box
func (fc *FeatureCollection) ComputeBoundingBox(force bool) (bmin, bmax []float64) {
	if !force && checkBoundingBox(fc.BoundingBox) {
		dimension := len(fc.BoundingBox) / 2
		bmin = fc.BoundingBox[:dimension]
		bmax = fc.BoundingBox[dimension:]
		return bmin, bmax
	}
	if fc.Features == nil || len(fc.Features) == 0 {
		return nil, nil
	}
	if bmin, bmax = fc.Features[0].ComputeBoundingBox(force); bmin == nil {
		return nil, nil
	}
	dimension := len(bmin)
	for i := 1; i < len(fc.Features); i++ {
		tmin, tmax := fc.Features[i].ComputeBoundingBox(force)
		if tmin == nil || len(tmin) != dimension {
			return nil, nil
		}
		for d := 0; d < dimension; d++ {
			bmin[d] = math.Min(bmin[d], tmin[d])
			bmax[d] = math.Max(bmax[d], tmax[d])
		}
	}
	return bmin, bmax
}
