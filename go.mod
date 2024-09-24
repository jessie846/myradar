module github.com/jessie846/myradar

go 1.23.1

// Replace the local paths to your custom modules
replace (
	github.com/jessie846/myradar/src/custom_map => ./src/custom_map
	github.com/jessie846/myradar/src/flight => ./src/flight
	github.com/jessie846/myradar/src/flight_list => ./src/flight_list
	github.com/jessie846/myradar/src/lat_long => ./src/lat_long
	github.com/jessie846/myradar/src/mca => ./src/mca
	github.com/jessie846/myradar/src/renderer => ./src/renderer
	github.com/jessie846/myradar/src/response_area => ./src/response_area
	github.com/jessie846/myradar/src/target_renderer => ./src/target_renderer
)
