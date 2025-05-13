APP_NAME=kviz
SRC_FILE=client.go
BIN_DIR=$(HOME)/.local/bin
ICON_DIR=$(HOME)/.local/share/icons
DESKTOP_DIR=$(HOME)/.local/share/applications

all: build install

build:
	go build -o $(APP_NAME) $(SRC_FILE)

install: build
	@echo "Instalacija aplikacije '$(APP_NAME)'..."

	# Креирање директоријума
	mkdir -p $(BIN_DIR)
	mkdir -p $(ICON_DIR)
	mkdir -p $(DESKTOP_DIR)

	# Копирање бинарног фајла
	cp $(APP_NAME) $(BIN_DIR)/
	chmod +x $(BIN_DIR)/$(APP_NAME)

	# Копирање иконе
	cp icon.png $(ICON_DIR)/$(APP_NAME).png

	# Креирање .desktop фајла
	sed "s|@HOME@|$(HOME)|g" $(APP_NAME).desktop.in > $(DESKTOP_DIR)/$(APP_NAME).desktop
	chmod +x $(DESKTOP_DIR)/$(APP_NAME).desktop

	@echo "✅ Instalacija završena! Pokreni aplikaciju iz menija."

clean:
	rm -f $(APP_NAME)

