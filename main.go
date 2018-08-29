package main

import (
	"github.com/labstack/echo"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/encoding/mvt"
	"github.com/paulmach/orb/geojson"
	"github.com/paulmach/orb/maptile"
	"github.com/paulmach/orb/simplify"
	"net/http"
	"strconv"
)

func main() {
	e := echo.New()
	e.Static("/", "assets")
	e.GET("/neighborhood/:z/:x/:y", neighborhood)
	e.Logger.Fatal(e.Start(":3000"))
}

func contextToMaptile(c echo.Context) maptile.Tile {
	x, _ := strconv.ParseUint(c.Param("x"), 10, 32)
	y, _ := strconv.ParseUint(c.Param("y"), 10, 32)
	z, _ := strconv.ParseUint(c.Param("z"), 10, 32)
	return maptile.New(uint32(x), uint32(y), maptile.Zoom(z))
}

func neighborhood(c echo.Context) error {
	fc := geojson.NewFeatureCollection()
	p := orb.Point{-77.032353, 38.905511}
	f := geojson.NewFeature(p)
	f.Properties["icon"] = "marker"
	fc.Append(f)
	collections := map[string]*geojson.FeatureCollection{
		"test": fc,
	}
	layers := mvt.NewLayers(collections)
	layers.ProjectToTile(contextToMaptile(c))
	layers.Clip(mvt.MapboxGLDefaultExtentBound)
	layers.Simplify(simplify.DouglasPeucker(1.0))
	layers.RemoveEmpty(1.0, 1.0)
	data, _ := mvt.Marshal(layers)
	return c.Blob(http.StatusOK, "application/vnd.mapbox-vector-tile", data)
}
