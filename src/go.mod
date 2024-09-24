module myradar

go 1.23.1

require (
	github.com/paulmach/go.geojson v1.5.0
	github.com/streadway/amqp v1.1.0
	github.com/veandco/go-sdl2 v0.4.40
)

replace (
  github.com/jessie846/myradar/src/file_list => ./src/file_list
)
