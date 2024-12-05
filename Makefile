all: day01-exec day02-exec day03-exec day04-exec

day01/day01: $(wildcard day01/*.go day01/*.txt)
	cd day01 && go build 

day01-exec: day01 day01/day01
	cd day01 && time ./day01 -p false
	cd day01 && time ./day01 -p true

day02/day02: $(wildcard day02/*.go day02/*.txt)
	cd day02 && go build 

day02-exec: day02 day02/day02
	cd day02 && time ./day02 -p false
	cd day02 && time ./day02 -p true

day03/part1 day03/part2: $(wildcard day03/*.c day03/*.txt)
	cd day03 && make -j

day03-exec: day03/part1 day03/part2
	cd day03 && time ./part1
	cd day03 && time ./part2

day04/day04: $(wildcard day04/*.go day04/*.txt)
	cd day04 && go build 

day04-exec: day04 day04/day04
	cd day04 && time ./day04 -p false
	cd day04 && time ./day04 -p true

day05/day05: $(wildcard day05/*.go day05/*.txt)
	cd day05 && go build 

day05-exec: day05 day05/day05
	cd day05 && time ./day05 -p false
	cd day05 && time ./day05 -p true
