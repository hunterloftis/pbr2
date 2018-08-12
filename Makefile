.PHONY: hello redblue shapes

hello:
	go run ./examples/hello/hello.go

redblue:
	go run ./examples/redblue/redblue.go

shapes:
	go run ./examples/shapes/shapes.go

sponza:
	go run ./examples/sponza/sponza.go

mario:
	go build ./cmd/pbr
	./pbr fixtures/models/mario/mario-sculpture.obj -width 500 -height 400 -from 200,200,200 -to 0,0,0 -ambient 50,50,50 -v --floor

house:
	go build ./cmd/pbr
	./pbr 'fixtures/models/house/house interior.obj' -width 500 -height 400 -ambient 1000,1000,1000 -v

lambo:
	go build ./cmd/pbr
	./pbr fixtures/models/lambo/lambo.obj -width 1280 -height 720 -env fixtures/envmaps/282.hdr -rad 2500 -to=-0.2,0.5,0.4 -from=-5,2,-5 -v --indirect -bounce 8

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

cesar:
	go build ./cmd/pbr
	./pbr fixtures/models/simple/cesar.obj -width 500 -height 500 -v

chair:
	go build ./cmd/pbr
	./pbr fixtures/models/simple/chair.obj -width 480 -height 640 -v -from=40,300,-400

destroyer:
	go build ./cmd/pbr
	./pbr fixtures/models/simple/destroyer.obj -width 1000 -height 400 -v

legobricks:
	go build ./cmd/pbr
	./pbr fixtures/models/legobricks/LegoBricks3.obj -from 12,4,12 -to 0,1.5,0

legoplane:
	go build ./cmd/pbr
	./pbr fixtures/models/legoplane/LEGO.Creator_Plane.obj -from 700,250,1100 -floor 1.1 -floorcolor 0.25,0.25,0.2 -floorrough 0.9 -ambient 1200,1200,1100

glassbowl:
	go build ./cmd/pbr
	./pbr fixtures/models/glassbowl/Glass\ Bowl\ with\ Cloth\ Towel.obj -from 6,4,6

glass:
	go build ./cmd/pbr
	./pbr fixtures/models/glass/glass-obj.obj -floor 1.5 -env fixtures/envmaps/ennis.hdr -from 840,120,600 -lens 80 -fstop 1.4 -focus 0.7

toilet:
	go build ./cmd/pbr
	./pbr fixtures/models/toilet/Toilet.obj -floor 10 -width 320 -height 640 -from 0,200,150

gopher:
	go build ./cmd/pbr
	./pbr fixtures/models/gopher/gopher.obj -floor
