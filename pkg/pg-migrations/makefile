BUILD_TARGET := migr

.PHONY: all
all: build run

.PHONY: build
build:
	go build -o $(BUILD_TARGET) .

.PHONY: run
run:
	./$(BUILD_TARGET)

.PHONY: clean
clean:
	rm -f $(BUILD_TARGET)