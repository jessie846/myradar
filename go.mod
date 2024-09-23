module myradar

go 1.23.1

require (
    github.com/paulmach/go.geojson v1.5.0
    github.com/veandco/go-sdl2 v0.4.40
    github.com/streadway/amqp v1.0.0
)

replace github.com/veandco/go-sdl2/sdl/ttf v0.4.40 => github.com/veandco/go-sdl2 v0.4.40

// Replace the local paths to your custom modules
replace github.com/jessie846/myradar/src/custom_map v1.0.0 => ./src/custom_map
replace github.com/jessie846/myradar/src/flight v1.0.0 => ./src/flight
replace github.com/jessie846/myradar/src/flight_list v1.0.0 => ./src/flight_list
replace github.com/jessie846/myradar/src/lat_long v1.0.0 => ./src/lat_long
replace github.com/jessie846/myradar/src/renderer v1.0.0 => ./src/renderer
replace github.com/jessie846/myradar/src/target_renderer v1.0.0 => ./src/target_renderer
replace github.com/jessie846/myradar/src/response_area v1.0.0 => ./src/response_area
replace github.com/jessie846/myradar/src/mca v1.0.0 => ./src/mca