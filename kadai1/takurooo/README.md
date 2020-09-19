# 概要
静止画像のフォーマット変換を行うツール  
指定されたディレクトリ下にある静止画像を全てフォーマット変換する。  
オプションで変換前のフォーマットと変換後のフォーマットを指定することができる。

# 対応フォーマット
- jpeg
- png


# ビルド
``go build -o imgtrans cmd/imgtrans/imgtrans.go``

# 使い方

### Usage
```
Usage of imgtrans:
  -d string
    	directory path
  -i string
    	input format(jpg or png) (default "jpg")
  -o string
    	output formatjpg or png (default "png")
```

### Examples
例1)  
testディレクトリ以下にある全てのjpg or jpegファイルをpngフォーマットに変換  
``imgtrans -d ./test``  

例2)  
testディレクトリ以下にある全てのpngファイルをjpgフォーマットに変換  
``imgtrans -i png -o jpg -d ./test``

