version=1.0.0b3

rm -rf ./build
mkdir build

go build -o build/pvm.exe .

cp -r ./scripts ./build/scripts

cd ./build
tar.exe -a -c -f pvm-$version.zip Scripts pvm.exe