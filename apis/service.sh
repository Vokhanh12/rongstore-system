#!/usr/bin/env bash
set -euo pipefail

### ====== CẤU HÌNH ĐƯỜNG DẪN (có thể chỉnh) ======
PROTO_DIR="proto"
GO_OUT="../server/api"
DART_OUT="../client/lib/features/api"
### ===============================================

# Màu sắc log (đẹp tí cho dễ nhìn)
BOLD="\033[1m"; RED="\033[31m"; GREEN="\033[32m"; YELLOW="\033[33m"; NC="\033[0m"

log() { echo -e "${BOLD}${GREEN}[OK]${NC} $*"; }
warn() { echo -e "${BOLD}${YELLOW}[WARN]${NC} $*"; }
err() { echo -e "${BOLD}${RED}[ERR]${NC} $*" 1>&2; }

need() {
  if ! command -v "$1" >/dev/null 2>&1; then
    err "Thiếu lệnh: $1"
    return 1
  fi
}

echo -e "${BOLD}=== Generate gRPC for Go & Dart ===${NC}"

### 1) Kiểm tra công cụ cần thiết
need protoc || {
  err "Bạn cần cài protoc (https://grpc.io/docs/protoc-installation/)."
  exit 1
}

# Go plugins
GO_PB_PLUGIN="protoc-gen-go"
GO_GRPC_PLUGIN="protoc-gen-go-grpc"
if ! command -v "${GO_PB_PLUGIN}" >/dev/null 2>&1 || ! command -v "${GO_GRPC_PLUGIN}" >/dev/null 2>&1; then
  warn "Chưa thấy ${GO_PB_PLUGIN} hoặc ${GO_GRPC_PLUGIN} trong PATH."
  echo "Cài đặt nhanh:"
  echo "  go install google.golang.org/protobuf/cmd/protoc-gen-go@latest"
  echo "  go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest"
  echo "  export PATH=\"\$(go env GOPATH)/bin:\$PATH\""
fi

# Dart plugin
if ! command -v protoc-gen-dart >/dev/null 2>&1; then
  warn "Chưa thấy protoc-gen-dart trong PATH."
  echo "Cài đặt nhanh:"
  echo "  dart pub global activate protoc_plugin"
  echo "  export PATH=\"\$HOME/.pub-cache/bin:\$PATH\""
fi

### 2) Kiểm tra thư mục
[[ -d "$PROTO_DIR" ]] || { err "Không tìm thấy thư mục $PROTO_DIR"; exit 1; }
mkdir -p "$GO_OUT" "$DART_OUT"

### 3) Lấy danh sách .proto (bỏ qua google/**)
# shellcheck disable=SC2207
PROTO_FILES=($(find "$PROTO_DIR" -type f -name '*.proto' ! -path "$PROTO_DIR/google/*" | sort))

if [[ ${#PROTO_FILES[@]} -eq 0 ]]; then
  err "Không tìm thấy file .proto nào trong $PROTO_DIR (ngoài google/*)."
  exit 1
fi

echo "Sẽ generate cho các file:"
for f in "${PROTO_FILES[@]}"; do echo "  - $f"; done

### 4) Generate cho Go
echo -e "\n${BOLD}---> Generating Go files vào: ${GO_OUT}${NC}"
protoc \
  -I="$PROTO_DIR" \
  --go_out="$GO_OUT" --go_opt=paths=source_relative \
  --go-grpc_out="$GO_OUT" --go-grpc_opt=paths=source_relative \
  "${PROTO_FILES[@]}"
log "Go: xong."

### 5) Generate cho Dart
echo -e "\n${BOLD}---> Generating Dart files vào: ${DART_OUT}${NC}"
protoc \
  -I="$PROTO_DIR" \
  --dart_out=grpc:"$DART_OUT" \
  "${PROTO_FILES[@]}"
log "Dart: xong."

### 6) Tóm tắt
echo -e "\n${BOLD}Done.${NC}"
echo "Go   => $GO_OUT"
echo "Dart => $DART_OUT"
echo -e "${BOLD}Gợi ý:${NC} nhớ add dependency:"
echo "  Go:   go get google.golang.org/grpc google.golang.org/protobuf"
echo "  Dart: flutter pub add grpc   (hoặc: dart pub add grpc)"
