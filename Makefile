shapes_profile:
	go build ./examples/shapes
	./shapes -profile cpu.pprof
	go tool pprof --pdf ./shapes ./cpu.pprof > cpu.pdf
	open cpu.pdf
	
mario_profile:
	go build ./examples/mario
	./mario -profile cpu.pprof
	go tool pprof --pdf ./mario ./cpu.pprof > cpu.pdf
	open cpu.pdf
	
sponza_profile:
	go build ./examples/sponza
	./sponza -profile cpu.pprof
	go tool pprof --pdf ./sponza ./cpu.pprof > cpu.pdf
	open cpu.pdf

mario-sculpture:
	go build ./cmd/pbr
	./pbr fixtures/models/mario/mario-sculpture.obj -width 500 -height 400 -from 200,200,200 -to 0,0,0 -ambient 40,0,90 -v --floor