SOURCES=$(wildcard */go.mod)
TARGETS-BIN=$(patsubst %/go.mod,%-bin,$(SOURCES))
TARGETS=$(patsubst %/go.mod,%-exec,$(SOURCES))
TARGETS_EXEC_1=$(patsubst %/go.mod,%-exec-1,$(SOURCES))
TARGETS_EXEC_2=$(patsubst %/go.mod,%-exec-2,$(SOURCES))
.PHONY: $(TARGETS)
all: $(TARGETS) day03-exec

# C code
day03/part1 day03/part2: $(wildcard day03/*.c day03/*.txt)
	cd day03 && make -j

.PHONY: day03-exec-1
day03-exec-1: day03/part2
	cd day03 && time ./part1

.PHONY: day03-exec-2
day03-exec-2: day03/part2
	cd day03 && time ./part2

day03-exec: day03-exec-1 day03-exec-2

# Go code
$(TARGETS-BIN): %-bin: % %/go.mod
	cd $< && go build

.PHONY: %-exec-1
$(TARGETS_EXEC_1): %-exec-1: % %-bin
	cd $</ && time ./$< -p false 

.PHONY: %-exec-2
$(TARGETS_EXEC_2): %-exec-2: % %-bin
	cd $</ && time ./$< -p true

.PHONY: %-exec
$(TARGETS): %-exec: %-exec-1 %-exec-2
