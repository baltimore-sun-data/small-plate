set -eux -o pipefail

go build -o "mac/Small Plate.app/Contents/Resources/small-plate"
cd mac
zip -r small-plate.zip "Small Plate.app" -x ".*"
