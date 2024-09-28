task:
    build

# 出力先のディレクトリ
BINDIR:=bin

# ルートパッケージ名の取得
ROOT_PACKAGE:=$(shell go list .)

# 出力先バイナリファイル名(bin/server など)
BINARIES:=$(COMMAND_PACKAGES:$(ROOT_PACKAGE)/cmd/%=$(BINDIR)/%)

# ビルド時にチェックする .go ファイル
GO_FILES:=$(shell find . -type f -name '*.go' -print)

# ビルドタスク
.PHONY: build
build: $(BINARIES)

# 実ビルドタスク
$(BINARIES): $(GO_FILES)
	@go build -o $@ $(@:$(BINDIR)/%=$(ROOT_PACKAGE)/cmd/%)
