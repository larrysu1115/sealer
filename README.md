# sealer
archive files into one compressed file

```
go install github.com/larrysu1115/sealer && \
sealer --cmd=do --src="/Users/larrysu/Temp/testdir/src" --to="/Users/larrysu/Temp/testdir/to" --pre=nbcn

sealer --cmd=undo --src="/Users/larrysu/Temp/testdir/to" --to="/Users/larrysu/Temp/testdir/src" --file="nbcn_.*\.tgz"

rm ~/Temp/testdir/to/*.tgz   
```

Build for windows

```
env CGO_ENABLED=1 GOOS=windows GOARCH=amd64 \
  CC=/usr/local/Cellar/mingw-w64/5.0.2_2/bin/x86_64-w64-mingw32-gcc \
  go build -o ~/Downloads/Share_RDP/sealer.exe -v \
  github.com/larrysu1115/sealer
```
