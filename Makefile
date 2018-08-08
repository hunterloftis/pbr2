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
	./pbr fixtures/models/mario/mario-sculpture.obj -width 500 -height 400 -from 200,200,200 -to 0,0,0 -ambient 50,50,50 -v --floor

sponza:
	go build ./cmd/pbr
	./pbr fixtures/models/sponza/sponza.obj -width 500 -height 400 -ambient 1000,1000,1000 -v

house:
	go build ./cmd/pbr
	./pbr 'fixtures/models/house/house interior.obj' -width 500 -height 400 -ambient 1000,1000,1000 -v

lambo:
	go build ./cmd/pbr
	./pbr fixtures/models/lambo/lambo.obj -width 640 -height 480 -ambient 1000,1000,1000 -to=-0.2,0.5,0.4 -from=-5,2,-5 -v