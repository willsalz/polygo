package main

import (
	"github.com/golang/geo/s2"
	geojson "github.com/paulmach/go.geojson"
	"go.uber.org/zap"
)

func main() {
	logger := zap.NewExample().Sugar()
	defer logger.Sync()

	rawFeatureCollection := []byte(`
	{
       "type": "FeatureCollection",
       "features": [{
           "type": "Feature",
           "geometry": {
               "type": "Polygon",
               "coordinates": [
                   [
                       [100.0, 0.0],
                       [101.0, 0.0],
                       [101.0, 1.0],
                       [100.0, 1.0],
                       [100.0, 0.0]
                   ]
               ]
           },
           "properties": {
               "prop0": "value0",
               "prop1": {
                   "this": "that"
               }
           }
       }]
   }
	`)

	fc, err := geojson.UnmarshalFeatureCollection(rawFeatureCollection)
	if err != nil {
		logger.Fatal(err)
	}

	toLoop := func(points [][]float64) *s2.Loop {
		var pts []s2.Point
		for _, pt := range points {
			pts = append(pts, s2.PointFromLatLng(s2.LatLngFromDegrees(pt[1], pt[0])))
		}
		return s2.LoopFromPoints(pts)
	}
	polygon := s2.PolygonFromOrientedLoops([]*s2.Loop{toLoop(fc.Features[0].Geometry.Polygon[0])})
	point := s2.PointFromLatLng(s2.LatLngFromDegrees(0.1, 100.1))

	logger.Infof("Combined area: %.7f\n", polygon.Area())
	logger.Infof("%v contains %v: %v", polygon, point, polygon.ContainsPoint(point))

}
