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
	