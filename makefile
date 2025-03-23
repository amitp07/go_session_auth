APP_NAME = session_auth
SOURCE_ENTRY = ./cmd/api
OUTPUT_DIR = dist

build:
	mkdir -p $(OUTPUT_DIR)
	go build -o $(OUTPUT_DIR)/$(APP_NAME) $(SOURCE_ENTRY)

run: build
	$(OUTPUT_DIR)/$(APP_NAME)

clean:
	rm -rf $(OUTPUT_DIR)