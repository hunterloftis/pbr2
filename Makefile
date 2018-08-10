mario-sculpture:
	go build ./cmd/pbr
	./pbr fixtures/models/mario/mario-sculpture.obj -width 500 -height 400 -from 200,200,200 -to 0,0,0 -ambient 50,50,50 -v --floor

house:
	go build ./cmd/pbr
	./pbr 'fixtures/models/house/house interior.obj' -width 500 -height 400 -ambient 1000,1000,1000 -v

lambo:
	go build ./cmd/pbr
	./pbr fixtures/models/lambo/lambo.obj -width 256 -height 144 -env fixtures/envmaps/river.hdr -rad 1500 -to=-0.2,0.5,0.4 -from=-5,2,-5 -v --indirect -bounce 8

skull:
	go build ./cmd/pbr
	./pbr fixtures/models/simple/skull.obj -width 640 -height 480 -v

lucy:
	go build ./cmd/pbr
	./pbr fixtures/models/simple/lucy.obj -width 480 -height 640 -v

falcon:
	go build ./cmd/pbr
	./pbr fixtures/models/simple/falcon.obj -width 800 -height 400 -v -to=-86,-18,-2681 -from=800,200,-3000

moses:
	go build ./cmd/pbr
	./pbr fixtures/models/moses/model.obj -width 480 -height 640 -v

gopher:
	go build ./cmd/pbr
	./pbr fixtures/models/gopher/gopher.obj -width 480 -height 640 -v

cesar:
	go build ./cmd/pbr
	./pbr fixtures/models/simple/cesar.obj -width 500 -height 500 -v

chair:
	go build ./cmd/pbr
	./pbr fixtures/models/simple/chair.obj -width 480 -height 640 -v -from=40,300,-400

destroyer:
	go build ./cmd/pbr
	./pbr fixtures/models/simple/destroyer.obj -width 1000 -height 400 -v